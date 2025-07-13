package adminauth

import (
	"database/sql"
	"fmt"

	"github.com/varnit-ta/PlacementLog/internal/db"
	"golang.org/x/crypto/bcrypt"
)

/*
AdminRepo handles admin authentication data access operations.
Provides methods for admin login and registration with database interactions.
*/
type AdminRepo struct {
	db *sql.DB
}

/*
NewAdminRepo creates a new AdminRepo instance with the provided database connection.

Parameters:
- db: The database connection

Returns:
- *AdminRepo: A new repository instance
*/
func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

/*
Login validates admin credentials against the database.

Parameters:
- username: The admin's username
- password: The admin's password

Returns:
- *db.Admin: The authenticated admin information
- error: Any error that occurred during login

The function:
1. Validates that username and password are provided
2. Queries the database for the admin
3. Compares the provided password with the hashed password using bcrypt
4. Returns admin information upon successful authentication

Possible errors:
- "all fields are required": Missing username or password
- "admin not found": Admin not found in database
- "incorrect password": Password doesn't match
*/
func (repo AdminRepo) Login(username, password string) (*db.Admin, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	var admin db.Admin
	var hashedPass string

	query := `
		SELECT id, username, password 
		FROM placement_log_admins 
		WHERE username = $1;
	`

	err := repo.db.QueryRow(query, username).Scan(&admin.ID, &admin.Username, &hashedPass)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("admin not found")
	} else if err != nil {
		return nil, fmt.Errorf("db error: %v", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password)) != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	return &admin, nil
}

/*
Register creates a new admin account in the database.

Parameters:
- username: The admin's username
- password: The admin's password

Returns:
- *db.Admin: The newly created admin information
- error: Any error that occurred during registration

The function:
1. Validates that username and password are provided
2. Hashes the password using bcrypt with default cost
3. Inserts the new admin into the database
4. Returns the created admin information

Possible errors:
- "all fields are required": Missing username or password
- "error hashing password": Password hashing failed
- "failed to register admin": Database insertion failed
*/
func (repo AdminRepo) Register(username, password string) (*db.Admin, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	query := `
		INSERT INTO placement_log_admins (username, password)
		VALUES ($1, $2)
		RETURNING id;
	`

	var adminID string
	err = repo.db.QueryRow(query, username, hashedPass).Scan(&adminID)

	if err != nil {
		return nil, fmt.Errorf("failed to register admin: %v", err)
	}

	return &db.Admin{
		ID:       adminID,
		Username: username,
	}, nil
}
