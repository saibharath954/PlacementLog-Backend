package adminauth

import (
	"database/sql"
	"fmt"

	"github.com/varnit-ta/PlacementLog/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

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
