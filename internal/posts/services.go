package posts

import (
	"encoding/json"
	"fmt"

	"github.com/varnit-ta/PlacementLog/internal/db"
)

type PostsService struct {
	repo *PostsRepo
}

func NewPostsService(repo *PostsRepo) *PostsService {
	return &PostsService{repo: repo}
}

func (s *PostsService) AddPost(userId string, postBody map[string]any) (*db.Post, error) {
	if userId == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	bytes, err := json.Marshal(postBody)

	if err != nil {
		return nil, fmt.Errorf("error marshalling post bytes: %v", err)
	}

	return s.repo.AddPost(userId, json.RawMessage(bytes))
}

func (s *PostsService) DeletePost(postId string) error {
	return s.repo.DeletePost(postId)
}

func (s *PostsService) GetAll() ([]db.Post, error) {
	return s.repo.GetAllPosts()
}

func (s *PostsService) GetByUser(userId string) ([]db.Post, error) {
	return s.repo.GetPostsByUserId(userId)
}
