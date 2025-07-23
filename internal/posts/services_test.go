package posts

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/varnit-ta/PlacementLog/internal/db"
)

type mockPostsRepo struct {
	AddPostFunc             func(userId string, postBody json.RawMessage) (*db.Post, error)
	UpdatePostFunc          func(postId, userId string, postBody json.RawMessage) (*db.Post, error)
	DeletePostFunc          func(postId, userId string) error
	DeletePostAsAdminFunc   func(postId string) error
	GetAllPostsFunc         func() ([]db.Post, error)
	GetAllPostsForAdminFunc func() ([]db.Post, error)
	GetPostsByUserIdFunc    func(userId string) ([]db.Post, error)
	ReviewPostFunc          func(postId, action string) error
}

func (m *mockPostsRepo) AddPost(userId string, postBody json.RawMessage) (*db.Post, error) {
	return m.AddPostFunc(userId, postBody)
}
func (m *mockPostsRepo) UpdatePost(postId, userId string, postBody json.RawMessage) (*db.Post, error) {
	return m.UpdatePostFunc(postId, userId, postBody)
}
func (m *mockPostsRepo) DeletePost(postId, userId string) error {
	return m.DeletePostFunc(postId, userId)
}
func (m *mockPostsRepo) DeletePostAsAdmin(postId string) error {
	return m.DeletePostAsAdminFunc(postId)
}
func (m *mockPostsRepo) GetAllPosts() ([]db.Post, error) {
	return m.GetAllPostsFunc()
}
func (m *mockPostsRepo) GetAllPostsForAdmin() ([]db.Post, error) {
	return m.GetAllPostsForAdminFunc()
}
func (m *mockPostsRepo) GetPostsByUserId(userId string) ([]db.Post, error) {
	return m.GetPostsByUserIdFunc(userId)
}
func (m *mockPostsRepo) ReviewPost(postId, action string) error {
	return m.ReviewPostFunc(postId, action)
}

func TestPostsService_AddPost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPostsRepo{
			AddPostFunc: func(userId string, postBody json.RawMessage) (*db.Post, error) {
				return &db.Post{ID: "1", UserID: userId, PostBody: postBody, Reviewed: false}, nil
			},
		}
		s := NewPostsService(repo)
		postBody := map[string]any{"company": "TestCo", "role": "Engineer"}
		post, err := s.AddPost("user1", postBody)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if post.UserID != "user1" {
			t.Errorf("expected user1, got %s", post.UserID)
		}
	})
	t.Run("missing userId", func(t *testing.T) {
		s := NewPostsService(&mockPostsRepo{})
		_, err := s.AddPost("", map[string]any{"company": "TestCo"})
		if err == nil || err.Error() != "user ID is required" {
			t.Errorf("expected user ID is required error, got %v", err)
		}
	})
	t.Run("marshal error", func(t *testing.T) {
		s := NewPostsService(&mockPostsRepo{})
		_, err := s.AddPost("user1", map[string]any{"bad": func() {}})
		if err == nil || !strings.Contains(err.Error(), "error marshalling post bytes") {
			t.Errorf("expected marshalling error, got %v", err)
		}
	})
}

func TestPostsService_UpdatePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPostsRepo{
			UpdatePostFunc: func(postId, userId string, postBody json.RawMessage) (*db.Post, error) {
				return &db.Post{ID: postId, UserID: userId, PostBody: postBody, Reviewed: false}, nil
			},
		}
		s := NewPostsService(repo)
		postBody := map[string]any{"company": "TestCo"}
		post, err := s.UpdatePost("p1", "u1", postBody)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if post.ID != "p1" || post.UserID != "u1" {
			t.Errorf("unexpected post: %+v", post)
		}
	})
	t.Run("missing postId or userId", func(t *testing.T) {
		s := NewPostsService(&mockPostsRepo{})
		_, err := s.UpdatePost("", "u1", map[string]any{})
		if err == nil || err.Error() != "post ID and user ID are required" {
			t.Errorf("expected post ID and user ID are required error, got %v", err)
		}
		_, err = s.UpdatePost("p1", "", map[string]any{})
		if err == nil || err.Error() != "post ID and user ID are required" {
			t.Errorf("expected post ID and user ID are required error, got %v", err)
		}
	})
	t.Run("marshal error", func(t *testing.T) {
		s := NewPostsService(&mockPostsRepo{})
		_, err := s.UpdatePost("p1", "u1", map[string]any{"bad": func() {}})
		if err == nil || !strings.Contains(err.Error(), "error marshalling post bytes") {
			t.Errorf("expected marshalling error, got %v", err)
		}
	})
}

func TestPostsService_DeletePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPostsRepo{
			DeletePostFunc: func(postId, userId string) error { return nil },
		}
		s := NewPostsService(repo)
		err := s.DeletePost("p1", "u1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
	t.Run("missing postId or userId", func(t *testing.T) {
		s := NewPostsService(&mockPostsRepo{})
		err := s.DeletePost("", "u1")
		if err == nil || err.Error() != "post ID and user ID are required" {
			t.Errorf("expected post ID and user ID are required error, got %v", err)
		}
		err = s.DeletePost("p1", "")
		if err == nil || err.Error() != "post ID and user ID are required" {
			t.Errorf("expected post ID and user ID are required error, got %v", err)
		}
	})
	// Edge: repo returns error
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPostsRepo{
			DeletePostFunc: func(postId, userId string) error { return errors.New("db error") },
		}
		s := NewPostsService(repo)
		err := s.DeletePost("p1", "u1")
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

func TestPostsService_DeletePostAsAdmin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPostsRepo{
			DeletePostAsAdminFunc: func(postId string) error { return nil },
		}
		s := NewPostsService(repo)
		err := s.DeletePostAsAdmin("p1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
	t.Run("missing postId", func(t *testing.T) {
		s := NewPostsService(&mockPostsRepo{})
		err := s.DeletePostAsAdmin("")
		if err == nil || err.Error() != "post ID is required" {
			t.Errorf("expected post ID is required error, got %v", err)
		}
	})
	// Edge: repo returns error
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPostsRepo{
			DeletePostAsAdminFunc: func(postId string) error { return errors.New("db error") },
		}
		s := NewPostsService(repo)
		err := s.DeletePostAsAdmin("p1")
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

func TestPostsService_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		posts := []db.Post{{ID: "1"}, {ID: "2"}}
		repo := &mockPostsRepo{
			GetAllPostsFunc: func() ([]db.Post, error) { return posts, nil },
		}
		s := NewPostsService(repo)
		got, err := s.GetAll()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, posts) {
			t.Errorf("expected %v, got %v", posts, got)
		}
	})
	// Edge: repo returns error
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPostsRepo{
			GetAllPostsFunc: func() ([]db.Post, error) { return nil, errors.New("db error") },
		}
		s := NewPostsService(repo)
		_, err := s.GetAll()
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

func TestPostsService_GetAllPostsForAdmin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		posts := []db.Post{{ID: "1"}, {ID: "2"}}
		repo := &mockPostsRepo{
			GetAllPostsForAdminFunc: func() ([]db.Post, error) { return posts, nil },
		}
		s := NewPostsService(repo)
		got, err := s.GetAllPostsForAdmin()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, posts) {
			t.Errorf("expected %v, got %v", posts, got)
		}
	})
	// Edge: repo returns error
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPostsRepo{
			GetAllPostsForAdminFunc: func() ([]db.Post, error) { return nil, errors.New("db error") },
		}
		s := NewPostsService(repo)
		_, err := s.GetAllPostsForAdmin()
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

func TestPostsService_GetByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		posts := []db.Post{{ID: "1"}, {ID: "2"}}
		repo := &mockPostsRepo{
			GetPostsByUserIdFunc: func(userId string) ([]db.Post, error) { return posts, nil },
		}
		s := NewPostsService(repo)
		got, err := s.GetByUser("u1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, posts) {
			t.Errorf("expected %v, got %v", posts, got)
		}
	})
	// Edge: repo returns error
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPostsRepo{
			GetPostsByUserIdFunc: func(userId string) ([]db.Post, error) { return nil, errors.New("db error") },
		}
		s := NewPostsService(repo)
		_, err := s.GetByUser("u1")
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

func TestPostsService_ReviewPost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPostsRepo{
			ReviewPostFunc: func(postId, action string) error { return nil },
		}
		s := NewPostsService(repo)
		err := s.ReviewPost("p1", "approve")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
	t.Run("missing postId or action", func(t *testing.T) {
		s := NewPostsService(&mockPostsRepo{})
		err := s.ReviewPost("", "approve")
		if err == nil || err.Error() != "post ID and action are required" {
			t.Errorf("expected post ID and action are required error, got %v", err)
		}
		err = s.ReviewPost("p1", "")
		if err == nil || err.Error() != "post ID and action are required" {
			t.Errorf("expected post ID and action are required error, got %v", err)
		}
	})
	// Edge: repo returns error
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPostsRepo{
			ReviewPostFunc: func(postId, action string) error { return errors.New("db error") },
		}
		s := NewPostsService(repo)
		err := s.ReviewPost("p1", "approve")
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}
