package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/varnit-ta/PlacementLog/internal/db"
)

type PostsRepo struct {
	db *sql.DB
}

func NewPostsRepo(db *sql.DB) *PostsRepo {
	return &PostsRepo{
		db: db,
	}
}

func (repo PostsRepo) GetAllPosts() ([]db.Post, error) {
	query := `SELECT id, user_id, post_body FROM placement_log_posts;`

	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("failed to get all the posts: %v", err)
	}

	defer rows.Close()

	var posts []db.Post

	for rows.Next() {
		var p db.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.PostBody); err != nil {
			return nil, fmt.Errorf("failed to scan posts: %v", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (repo PostsRepo) GetPostsByUserId(userId string) ([]db.Post, error) {
	if userId == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	query := `SELECT id, user_id, post_body FROM placement_log_posts WHERE user_id=$1;`

	rows, err := repo.db.Query(query, userId)

	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %v", err)
	}

	defer rows.Close()

	var posts []db.Post

	for rows.Next() {
		var p db.Post

		if err := rows.Scan(&p.ID, &p.UserID, &p.PostBody); err != nil {
			return nil, fmt.Errorf("failed to scan posts: %v", err)
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func (repo PostsRepo) AddPost(userId string, postBody json.RawMessage) (*db.Post, error) {
	if userId == "" || postBody == nil {
		return nil, fmt.Errorf("all fields are required")
	}

	query := `
		INSERT INTO placement_log_posts (user_id, post_body)
		VALUES ($1, $2)
		RETURNING id;
	`

	var postId string

	err := repo.db.QueryRow(query, userId, postBody).Scan(&postId)

	if err != nil {
		return nil, fmt.Errorf("failed to add post: %v", err)
	}

	return &db.Post{
		ID:       postId,
		UserID:   userId,
		PostBody: postBody,
	}, nil
}

func (repo PostsRepo) DeletePost(postId string) error {
	if postId == "" {
		return fmt.Errorf("post ID is required")
	}

	query := `
		DELETE FROM placement_log_posts
		WHERE id = $1;
	`

	result, err := repo.db.Exec(query, postId)

	if err != nil {
		return fmt.Errorf("failed to delete post: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not confirm deletion: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with given ID")
	}

	return nil
}
