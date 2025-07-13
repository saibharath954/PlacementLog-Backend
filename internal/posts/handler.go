package posts

import (
	"errors"
	"net/http"

	"github.com/varnit-ta/PlacementLog/pkg/utils"
)

/*
PostsHandler handles post-related HTTP requests.
Provides endpoints for creating, reading, updating, and deleting posts.
Includes both user and admin-specific operations.
*/
type PostsHandler struct {
	srv *PostsService
}

/*
NewPostsHandler creates a new PostsHandler instance with the provided service.

Parameters:
- srv: The posts service

Returns:
- *PostsHandler: A new handler instance
*/
func NewPostsHandler(srv *PostsService) *PostsHandler {
	return &PostsHandler{
		srv: srv,
	}
}

/*
createPostRequest represents the JSON payload for creating a new post.
*/
type createPostRequest struct {
	PostBody map[string]any `json:"post_body"`
}

/*
updatePostRequest represents the JSON payload for updating an existing post.
*/
type updatePostRequest struct {
	PostBody map[string]any `json:"post_body"`
}

/*
AddPost handles post creation requests.
Creates a new post with reviewed=false, requiring admin approval.

HTTP Method: POST
Endpoint: /posts

Headers Required:
- Authorization: Bearer <user_jwt_token>

Request Body:

	{
	  "post_body": {
	    "company": "Google",
	    "role": "Software Engineer",
	    "rounds": [...]
	  }
	}

Response (201 Created):

	{
	  "id": "post_id",
	  "user_id": "user_id",
	  "post_body": {...}
	}

Returns:
- 201 Created: Post created successfully
- 400 Bad Request: Invalid request format
- 401 Unauthorized: Missing or invalid token
*/
func (h *PostsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	var req createPostRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, err)
		return
	}

	userId := r.Header.Get("X-User-ID")

	if userId == "" {
		utils.WriteError(w, http.ErrNoCookie)
		return
	}

	post, err := h.srv.AddPost(userId, req.PostBody)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, post, http.StatusCreated)
}

/*
UpdatePost handles post update requests.
Users can only update their own posts.

HTTP Method: PUT
Endpoint: /posts?id=<post_id>

Headers Required:
- Authorization: Bearer <user_jwt_token>

Query Parameters:
- id: The ID of the post to update

Request Body:

	{
	  "post_body": {
	    "company": "Updated Company",
	    "role": "Updated Role",
	    "rounds": [...]
	  }
	}

Response (200 OK):

	{
	  "id": "post_id",
	  "user_id": "user_id",
	  "post_body": {...}
	}

Returns:
- 200 OK: Post updated successfully
- 400 Bad Request: Invalid request format or missing post ID
- 401 Unauthorized: Missing or invalid token
- 403 Forbidden: User not authorized to update this post
*/
func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var req updatePostRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, err)
		return
	}

	userId := r.Header.Get("X-User-ID")
	query := r.URL.Query()
	postId := query.Get("id")

	if userId == "" {
		utils.WriteError(w, http.ErrNoCookie)
		return
	}

	if postId == "" {
		utils.WriteError(w, http.ErrMissingFile)
		return
	}

	post, err := h.srv.UpdatePost(postId, userId, req.PostBody)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, post, http.StatusOK)
}

/*
DeletePost handles post deletion requests.
Users can only delete their own posts.

HTTP Method: DELETE
Endpoint: /posts?id=<post_id>

Headers Required:
- Authorization: Bearer <user_jwt_token>

Query Parameters:
- id: The ID of the post to delete

Response (200 OK):

	{
	  "message": "post deleted"
	}

Returns:
- 200 OK: Post deleted successfully
- 400 Bad Request: Missing post ID
- 401 Unauthorized: Missing or invalid token
- 403 Forbidden: User not authorized to delete this post
*/
func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	postId := query.Get("id")
	userId := r.Header.Get("X-User-ID")

	if postId == "" {
		utils.WriteError(w, http.ErrMissingFile)
		return
	}

	if userId == "" {
		utils.WriteError(w, http.ErrNoCookie)
		return
	}

	err := h.srv.DeletePost(postId, userId)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "post deleted"}, http.StatusOK)
}

/*
GetAll handles requests to retrieve all approved posts.
This endpoint is public and doesn't require authentication.

HTTP Method: GET
Endpoint: /posts

Response (200 OK):
[

	{
	  "id": "post_id",
	  "user_id": "user_id",
	  "post_body": {...}
	}

]

Returns:
- 200 OK: List of all approved posts
- 500 Internal Server Error: Database error
*/
func (h *PostsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.srv.GetAll()
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, posts, http.StatusOK)
}

/*
GetByUser handles requests to retrieve posts by a specific user.
Users can only view their own posts.

HTTP Method: GET
Endpoint: /posts/user?user_id=<user_id>

Headers Required:
- Authorization: Bearer <user_jwt_token>

Query Parameters:
- user_id: The ID of the user whose posts to retrieve

Response (200 OK):
[

	{
	  "id": "post_id",
	  "user_id": "user_id",
	  "post_body": {...}
	}

]

Returns:
- 200 OK: List of user's approved posts
- 400 Bad Request: Missing user ID
- 401 Unauthorized: Missing or invalid token
- 403 Forbidden: User can only view their own posts
*/
func (h *PostsHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	requestedUserId := query.Get("user_id")
	authenticatedUserId := r.Header.Get("X-User-ID")

	if requestedUserId == "" {
		utils.WriteError(w, http.ErrMissingFile)
		return
	}

	if requestedUserId != authenticatedUserId {
		utils.WriteError(w, errors.New("forbidden: can only access own posts"))
		return
	}

	posts, err := h.srv.GetByUser(requestedUserId)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, posts, http.StatusOK)
}

/*
GetAllPostsForAdmin handles requests to retrieve all posts for admin review.
Admins can see both approved and pending posts.

HTTP Method: GET
Endpoint: /admin/posts

Headers Required:
- Authorization: Bearer <admin_jwt_token>

Response (200 OK):
[

	{
	  "id": "post_id",
	  "user_id": "user_id",
	  "post_body": {...},
	  "reviewed": true/false
	}

]

Returns:
- 200 OK: List of all posts (approved and pending)
- 401 Unauthorized: Missing or invalid admin token
- 500 Internal Server Error: Database error
*/
func (h *PostsHandler) GetAllPostsForAdmin(w http.ResponseWriter, r *http.Request) {
	posts, err := h.srv.GetAllPostsForAdmin()

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, posts, http.StatusOK)
}

/*
ReviewPost handles post review requests by admins.
Admins can approve or reject posts.

HTTP Method: PUT
Endpoint: /admin/posts/review?id=<post_id>&action=<action>

Headers Required:
- Authorization: Bearer <admin_jwt_token>

Query Parameters:
- id: The ID of the post to review
- action: Either "approve" or "reject"

Response (200 OK):

	{
	  "message": "post approved" or "post rejected"
	}

Returns:
- 200 OK: Post reviewed successfully
- 400 Bad Request: Missing parameters or invalid action
- 401 Unauthorized: Missing or invalid admin token
- 500 Internal Server Error: Database error
*/
func (h *PostsHandler) ReviewPost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	postId := query.Get("id")
	action := query.Get("action") // "approve" or "reject"

	if postId == "" || action == "" {
		utils.WriteError(w, http.ErrMissingFile)
		return
	}

	if action != "approve" && action != "reject" {
		utils.WriteError(w, errors.New("invalid action: must be 'approve' or 'reject'"))
		return
	}

	err := h.srv.ReviewPost(postId, action)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "post " + action + "d"}, http.StatusOK)
}

/*
DeletePostAsAdmin handles post deletion requests by admins.
Admins can delete any post, regardless of ownership.

HTTP Method: DELETE
Endpoint: /admin/posts?id=<post_id>

Headers Required:
- Authorization: Bearer <admin_jwt_token>

Query Parameters:
- id: The ID of the post to delete

Response (200 OK):

	{
	  "message": "post deleted by admin"
	}

Returns:
- 200 OK: Post deleted successfully
- 400 Bad Request: Missing post ID
- 401 Unauthorized: Missing or invalid admin token
- 500 Internal Server Error: Database error
*/
func (h *PostsHandler) DeletePostAsAdmin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	postId := query.Get("id")

	if postId == "" {
		utils.WriteError(w, http.ErrMissingFile)
		return
	}

	err := h.srv.DeletePostAsAdmin(postId)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "post deleted by admin"}, http.StatusOK)
}
