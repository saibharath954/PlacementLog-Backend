package userauth

import (
	"fmt"

	"github.com/varnit-ta/PlacementLog/internal/db"
	"github.com/varnit-ta/PlacementLog/pkg/jwt"
)

// Define UserAuthRepository interface for testability
//go:generate mockgen -destination=mock_userauth_repo.go -package=userauth . UserAuthRepository

type UserAuthRepository interface {
	Login(regno, pass string) (*db.User, error)
	Register(regno, username, pass string) (*db.User, error)
}

/*
UserAuthService handles user authentication business logic.
Provides methods for user login and registration with JWT token generation.
*/
type UserAuthService struct {
	repo UserAuthRepository
}

/*
NewUserAuthService creates a new UserAuthService instance with the provided repository.

Parameters:
- repo: The user authentication repository

Returns:
- *UserAuthService: A new service instance
*/
func NewUserAuthService(repo UserAuthRepository) *UserAuthService {
	return &UserAuthService{repo: repo}
}

/*
Login authenticates a user with the provided credentials and generates a JWT token.

Parameters:
- username: The user's username (registration number format: 22bcs1234)
- password: The user's password

Returns:
- string: JWT token for the authenticated user
- string: User ID of the authenticated user
- error: Any error that occurred during authentication

The function:
1. Validates the user credentials against the database
2. Generates a JWT token with "user" role
3. Returns the token and user ID upon successful authentication
*/
func (s *UserAuthService) Login(regno, password string) (string, *db.User, error) {
	user, err := s.repo.Login(regno, password)

	if err != nil {
		return "", nil, fmt.Errorf("login failed: %w", err)
	}

	token, err := jwt.GenerateJwtToken(user.ID, "user")

	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

/*
Register creates a new user account and generates a JWT token.

Parameters:
- username: The user's username (registration number format: 22bcs1234)
- password: The user's password

Returns:
- string: JWT token for the newly registered user
- string: User ID of the newly registered user
- error: Any error that occurred during registration

The function:
1. Validates the username format (registration number)
2. Hashes the password using bcrypt
3. Creates a new user in the database
4. Generates a JWT token with "user" role
5. Returns the token and user ID upon successful registration
*/
func (s *UserAuthService) Register(regno, name, password string) (string, string, error) {
	user, err := s.repo.Register(regno, name, password)

	if err != nil {
		return "", "", fmt.Errorf("registration failed: %w", err)
	}

	token, err := jwt.GenerateJwtToken(user.ID, "user")

	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user.ID, nil
}
