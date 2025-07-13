package adminauth

import (
	"github.com/varnit-ta/PlacementLog/internal/db"
	"github.com/varnit-ta/PlacementLog/pkg/jwt"
)

type AdminService struct {
	repo *AdminRepo
}

func NewAdminService(repo *AdminRepo) *AdminService {
	return &AdminService{repo: repo}
}

func (s AdminService) Login(username, password string) (string, *db.Admin, error) {
	admin, err := s.repo.Login(username, password)
	if err != nil {
		return "", nil, err
	}

	token, err := jwt.GenerateJwtToken(admin.ID)
	if err != nil {
		return "", nil, err
	}

	return token, admin, nil
}

func (s AdminService) Register(username, password string) (string, *db.Admin, error) {
	admin, err := s.repo.Register(username, password)
	if err != nil {
		return "", nil, err
	}

	token, err := jwt.GenerateJwtToken(admin.ID)
	if err != nil {
		return "", nil, err
	}

	return token, admin, nil
}
