package userauth

import (
	"testing"

	"github.com/varnit-ta/PlacementLog/internal/db"
)

type fakeUserAuthRepo struct{}

func (f *fakeUserAuthRepo) Login(regno, pass string) (*db.User, error) {
	return &db.User{ID: "11111111-1111-1111-1111-111111111111", Regno: regno, Username: "testuser"}, nil
}
func (f *fakeUserAuthRepo) Register(regno, username, pass string) (*db.User, error) {
	return nil, nil // not used in this test
}

func TestUserAuthService_Login(t *testing.T) {
	s := NewUserAuthService(&fakeUserAuthRepo{})
	token, user, err := s.Login("22bcs1234", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Regno != "22bcs1234" {
		t.Errorf("expected regno 22bcs1234, got %s", user.Regno)
	}
	if token == "" {
		t.Error("expected a token, got empty string")
	}
}
