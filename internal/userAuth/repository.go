package userauth

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/lib/pq"
	"github.com/varnit-ta/PlacementLog/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthRepo struct {
	db *sql.DB
}

func NewUserAuthRepo(db *sql.DB) *UserAuthRepo {
	return &UserAuthRepo{
		db: db,
	}
}

func (repo UserAuthRepo) Login(username, pass string) (*db.User, error) {
	if username == "" || pass == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	re := regexp.MustCompile(`^\d{2}[a-z]{3}\d{4}$`)

	if !re.MatchString(username) {
		return nil, fmt.Errorf("not a valid registration number")
	}

	queryString := `
		SELECT id, username, password
		FROM placement_log_users
		WHERE username=$1;
	`

	var user db.User
	var hashedPass string

	err := repo.db.QueryRow(queryString, username).Scan(
		&user.ID,
		&user.Username,
		&hashedPass,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no such user exists")
	}

	if err != nil {
		return nil, fmt.Errorf("db error: %v", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass)); err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	return &user, nil
}

func (repo UserAuthRepo) Register(username, pass string) (*db.User, error) {
	if username == "" || pass == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	re := regexp.MustCompile(`^\d{2}[a-z]{3}\d{4}$`)

	if !re.MatchString(username) {
		return nil, fmt.Errorf("not a valid registration number")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("error hashing pass: %v", err)
	}

	queryString := `
		INSERT INTO placement_log_users (username, password)
		VALUES ($1, $2)
		RETURNING id;
	`

	var userId string

	err = repo.db.QueryRow(queryString, username, hashedPass).Scan(&userId)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("user already exists")
		}
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	return &db.User{
		ID:       userId,
		Username: username,
	}, nil
}
