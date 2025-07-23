package userauth

import (
	"errors"
	"strings"
	"testing"

	"github.com/varnit-ta/PlacementLog/internal/db"
)

type mockDB struct{}

// Add methods as needed to satisfy the sql.DB interface for the test, or just mock the repository methods directly.

type mockUserAuthRepo struct {
	LoginFunc    func(regno, pass string) (*db.User, error)
	RegisterFunc func(regno, username, pass string) (*db.User, error)
}

func (m *mockUserAuthRepo) Login(regno, pass string) (*db.User, error) {
	return m.LoginFunc(regno, pass)
}
func (m *mockUserAuthRepo) Register(regno, username, pass string) (*db.User, error) {
	return m.RegisterFunc(regno, username, pass)
}

func TestUserAuthRepo_Login_TableDriven(t *testing.T) {
	repo := &mockUserAuthRepo{
		LoginFunc: func(regno, pass string) (*db.User, error) {
			if regno == "" || pass == "" {
				return nil, errors.New("all fields are required")
			}
			if !strings.HasPrefix(regno, "22bcs") {
				return nil, errors.New("not a valid registration number")
			}
			if pass == "wrongpass" {
				return nil, errors.New("incorrect password")
			}
			return &db.User{ID: "11111111-1111-1111-1111-111111111111", Regno: regno, Username: "testuser"}, nil
		},
	}
	cases := []struct {
		regno, pass string
		wantErr     string
	}{
		{"", "", "all fields are required"},
		{"badregno", "pass", "not a valid registration number"},
		{"22bcs9999", "wrongpass", "incorrect password"},
		{"22bcs9999", "password", ""},
	}
	for _, c := range cases {
		_, err := repo.Login(c.regno, c.pass)
		if c.wantErr == "" && err != nil {
			t.Errorf("Login(%q, %q) unexpected error: %v", c.regno, c.pass, err)
		}
		if c.wantErr != "" && (err == nil || !strings.Contains(err.Error(), c.wantErr)) {
			t.Errorf("Login(%q, %q) = %v; want error containing %q", c.regno, c.pass, err, c.wantErr)
		}
	}
}

func TestUserAuthRepo_Register_TableDriven(t *testing.T) {
	repo := &mockUserAuthRepo{
		RegisterFunc: func(regno, username, pass string) (*db.User, error) {
			if regno == "" || username == "" || pass == "" {
				return nil, errors.New("all fields are required")
			}
			if !strings.HasPrefix(regno, "22bcs") {
				return nil, errors.New("not a valid registration number")
			}
			return &db.User{ID: "11111111-1111-1111-1111-111111111111", Regno: regno, Username: username}, nil
		},
	}
	cases := []struct {
		regno, username, pass string
		wantErr               string
	}{
		{"", "", "", "all fields are required"},
		{"badregno", "user", "pass", "not a valid registration number"},
		{"22bcs9999", "user", "password", ""},
	}
	for _, c := range cases {
		_, err := repo.Register(c.regno, c.username, c.pass)
		if c.wantErr == "" && err != nil {
			t.Errorf("Register(%q, %q, %q) unexpected error: %v", c.regno, c.username, c.pass, err)
		}
		if c.wantErr != "" && (err == nil || !strings.Contains(err.Error(), c.wantErr)) {
			t.Errorf("Register(%q, %q, %q) = %v; want error containing %q", c.regno, c.username, c.pass, err, c.wantErr)
		}
	}
}
