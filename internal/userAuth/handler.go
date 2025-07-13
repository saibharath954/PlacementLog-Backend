package userauth

import (
	"net/http"

	"github.com/varnit-ta/PlacementLog/pkg/utils"
)

/*
UserAuthHandler handles user authentication HTTP requests.
Provides endpoints for user login, registration, and logout.
*/
type UserAuthHandler struct {
	srv *UserAuthService
}

/*
requestPayload represents the JSON payload for login and registration requests.
*/
type requestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
responsePayload represents the JSON response for successful authentication.
*/
type responsePayload struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

/*
NewUserAuthHandler creates a new UserAuthHandler instance with the provided service.

Parameters:
- srv: The user authentication service

Returns:
- *UserAuthHandler: A new handler instance
*/
func NewUserAuthHandler(srv *UserAuthService) *UserAuthHandler {
	return &UserAuthHandler{
		srv: srv,
	}
}

/*
Login handles user login requests.
Validates user credentials and returns a JWT token upon successful authentication.

HTTP Method: POST
Endpoint: /auth/login

Request Body:

	{
	  "username": "22bcs1234",
	  "password": "password123"
	}

Response (200 OK):

	{
	  "username": "user_id",
	  "token": "jwt_token_here"
	}

Returns:
- 200 OK: Successful login with token
- 400 Bad Request: Invalid request format
- 401 Unauthorized: Invalid credentials
*/
func (h *UserAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.WriteError(w, err)
		return
	}

	token, userId, err := h.srv.Login(payload.Username, payload.Password)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	resp := responsePayload{
		UserID:   userId,
		Username: payload.Username,
		Token:    token,
	}

	utils.WriteJSON(w, resp, http.StatusOK)
}

/*
Register handles user registration requests.
Creates a new user account and returns a JWT token upon successful registration.

HTTP Method: POST
Endpoint: /auth/register

Request Body:

	{
	  "username": "22bcs1234",
	  "password": "password123"
	}

Response (201 Created):

	{
	  "username": "user_id",
	  "token": "jwt_token_here"
	}

Returns:
- 201 Created: Successful registration with token
- 400 Bad Request: Invalid request format
- 409 Conflict: Username already exists
*/
func (h *UserAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.WriteError(w, err)
		return
	}

	token, userId, err := h.srv.Register(payload.Username, payload.Password)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	resp := responsePayload{
		UserID:   userId,
		Username: payload.Username,
		Token:    token,
	}

	utils.WriteJSON(w, resp, http.StatusCreated)
}

/*
Logout handles user logout requests.
For JWT tokens, logout is typically handled client-side by removing the token.
This endpoint provides a standardized way to handle logout requests.

HTTP Method: POST
Endpoint: /auth/logout

Headers Required:
- Authorization: Bearer <jwt_token>

Response (200 OK):

	{
	  "message": "logged out successfully"
	}

Note: The client should remove the JWT token from local storage after calling this endpoint.
*/
func (h *UserAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, map[string]string{"message": "logged out successfully"}, http.StatusOK)
}
