package adminauth

import (
	"github.com/varnit-ta/PlacementLog/internal/db"
	"github.com/varnit-ta/PlacementLog/pkg/jwt"
)

/*
AdminService handles admin authentication business logic.
Provides methods for admin login and registration with JWT token generation.
*/
type AdminService struct {
	repo *AdminRepo
}

/*
NewAdminService creates a new AdminService instance with the provided repository.

Parameters:
- repo: The admin authentication repository

Returns:
- *AdminService: A new service instance
*/
func NewAdminService(repo *AdminRepo) *AdminService {
	return &AdminService{repo: repo}
}

/*
Login authenticates an admin with the provided credentials and generates a JWT token.

Parameters:
- username: The admin's username
- password: The admin's password

Returns:
- string: JWT token for the authenticated admin
- *db.Admin: The authenticated admin information
- error: Any error that occurred during authentication

The function:
1. Validates the admin credentials against the database
2. Generates a JWT token with "admin" role
3. Returns the token and admin information upon successful authentication
*/
func (s AdminService) Login(username, password string) (string, *db.Admin, error) {
	admin, err := s.repo.Login(username, password)
	if err != nil {
		return "", nil, err
	}

	token, err := jwt.GenerateJwtToken(admin.ID, "admin")
	if err != nil {
		return "", nil, err
	}

	return token, admin, nil
}

/*
Register creates a new admin account and generates a JWT token.

Parameters:
- username: The admin's username
- password: The admin's password

Returns:
- string: JWT token for the newly registered admin
- *db.Admin: The newly registered admin information
- error: Any error that occurred during registration

The function:
1. Hashes the password using bcrypt
2. Creates a new admin in the database
3. Generates a JWT token with "admin" role
4. Returns the token and admin information upon successful registration
*/
func (s AdminService) Register(username, password string) (string, *db.Admin, error) {
	admin, err := s.repo.Register(username, password)
	if err != nil {
		return "", nil, err
	}

	token, err := jwt.GenerateJwtToken(admin.ID, "admin")
	if err != nil {
		return "", nil, err
	}

	return token, admin, nil
}
