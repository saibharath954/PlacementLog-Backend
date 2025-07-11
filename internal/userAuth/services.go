package userauth

import (
	"fmt"

	"github.com/varnit-ta/PlacementLog/pkg/jwt"
)

type UserAuthService struct {
	repo *UserAuthRepo
}

func NewUserAuthService(repo *UserAuthRepo) *UserAuthService {
	return &UserAuthService{
		repo: repo,
	}
}

func (s *UserAuthService) Login(username, password string) (string, string, error) {
	user, err := s.repo.Login(username, password)

	if err != nil {
		return "", "", fmt.Errorf("login failed: %w", err)
	}

	token, err := jwt.GenerateJwtToken(user.ID)

	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user.ID, nil
}

func (s *UserAuthService) Register(username, password string) (string, string, error) {
	user, err := s.repo.Register(username, password)

	if err != nil {
		return "", "", fmt.Errorf("registration failed: %w", err)
	}

	token, err := jwt.GenerateJwtToken(user.ID)

	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user.ID, nil
}
