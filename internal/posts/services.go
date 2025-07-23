package posts

import (
	"encoding/json"
	"fmt"

	"github.com/varnit-ta/PlacementLog/internal/db"
)

// Define PostsRepository interface for testability
//go:generate mockgen -destination=mock_posts_repo.go -package=posts . PostsRepository

type PostsRepository interface {
	AddPost(userId string, postBody json.RawMessage) (*db.Post, error)
	UpdatePost(postId, userId string, postBody json.RawMessage) (*db.Post, error)
	DeletePost(postId, userId string) error
	DeletePostAsAdmin(postId string) error
	GetAllPosts() ([]db.Post, error)
	GetAllPostsForAdmin() ([]db.Post, error)
	GetPostsByUserId(userId string) ([]db.Post, error)
	ReviewPost(postId, action string) error
}

/*
PostsService handles post-related business logic.
Provides methods for creating, reading, updating, and deleting posts.
Includes both user and admin-specific operations.
*/
type PostsService struct {
	repo PostsRepository
}

/*
NewPostsService creates a new PostsService instance with the provided repository.

Parameters:
- repo: The posts repository

Returns:
- *PostsService: A new service instance
*/
func NewPostsService(repo PostsRepository) *PostsService {
	return &PostsService{repo: repo}
}

/*
AddPost creates a new post for a user.
The post is created with reviewed=false, requiring admin approval.

Parameters:
- userId: The ID of the user creating the post
- postBody: The post content as a map

Returns:
- *db.Post: The created post
- error: Any error that occurred during creation

The function:
1. Validates that user ID is provided
2. Marshals the post body to JSON
3. Creates the post in the database with reviewed=false
4. Returns the created post information
*/
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

/*
UpdatePost updates an existing post.
Users can only update their own posts.

Parameters:
- postId: The ID of the post to update
- userId: The ID of the user updating the post
- postBody: The updated post content as a map

Returns:
- *db.Post: The updated post
- error: Any error that occurred during update

The function:
1. Validates that post ID and user ID are provided
2. Marshals the post body to JSON
3. Updates the post in the database (sets reviewed=false)
4. Returns the updated post information

Note: When a post is updated, it needs to be reviewed again by an admin.
*/
func (s *PostsService) UpdatePost(postId string, userId string, postBody map[string]any) (*db.Post, error) {
	if postId == "" || userId == "" {
		return nil, fmt.Errorf("post ID and user ID are required")
	}

	bytes, err := json.Marshal(postBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling post bytes: %v", err)
	}

	return s.repo.UpdatePost(postId, userId, json.RawMessage(bytes))
}

/*
DeletePost deletes a post.
Users can only delete their own posts.

Parameters:
- postId: The ID of the post to delete
- userId: The ID of the user deleting the post

Returns:
- error: Any error that occurred during deletion

The function:
1. Validates that post ID and user ID are provided
2. Deletes the post from the database (only if owned by the user)
3. Returns any error that occurred
*/
func (s *PostsService) DeletePost(postId string, userId string) error {
	if postId == "" || userId == "" {
		return fmt.Errorf("post ID and user ID are required")
	}
	return s.repo.DeletePost(postId, userId)
}

/*
DeletePostAsAdmin deletes any post (admin operation).
Admins can delete any post, regardless of ownership.

Parameters:
- postId: The ID of the post to delete

Returns:
- error: Any error that occurred during deletion

The function:
1. Validates that post ID is provided
2. Deletes the post from the database
3. Returns any error that occurred
*/
func (s *PostsService) DeletePostAsAdmin(postId string) error {
	if postId == "" {
		return fmt.Errorf("post ID is required")
	}
	return s.repo.DeletePostAsAdmin(postId)
}

/*
GetAll retrieves all approved posts.
This is a public operation that doesn't require authentication.

Returns:
- []db.Post: List of all approved posts
- error: Any error that occurred during retrieval

The function retrieves only posts that have been reviewed and approved by admins.
*/
func (s *PostsService) GetAll() ([]db.Post, error) {
	return s.repo.GetAllPosts()
}

/*
GetAllPostsForAdmin retrieves all posts for admin review.
Admins can see both approved and pending posts.

Returns:
- []db.Post: List of all posts (approved and pending)
- error: Any error that occurred during retrieval

The function retrieves all posts regardless of their review status.
*/
func (s *PostsService) GetAllPostsForAdmin() ([]db.Post, error) {
	return s.repo.GetAllPostsForAdmin()
}

/*
GetByUser retrieves posts by a specific user.
Only returns approved posts for the specified user.

Parameters:
- userId: The ID of the user whose posts to retrieve

Returns:
- []db.Post: List of the user's approved posts
- error: Any error that occurred during retrieval

The function retrieves only posts that belong to the specified user and have been approved.
*/
func (s *PostsService) GetByUser(userId string) ([]db.Post, error) {
	return s.repo.GetPostsByUserId(userId)
}

/*
ReviewPost reviews a post (admin operation).
Admins can approve or reject posts.

Parameters:
- postId: The ID of the post to review
- action: The review action ("approve" or "reject")

Returns:
- error: Any error that occurred during review

The function:
1. Validates that post ID and action are provided
2. Updates the post's reviewed status in the database
3. Returns any error that occurred

Note: When a post is approved, it becomes visible to the public.
When a post is rejected, it remains hidden from public view.
*/
func (s *PostsService) ReviewPost(postId string, action string) error {
	if postId == "" || action == "" {
		return fmt.Errorf("post ID and action are required")
	}
	return s.repo.ReviewPost(postId, action)
}
