package server

import (
	"github.com/varnit-ta/PlacementLog/internal/db"
)

type App struct {
}

func InitApp() (*App, error) {
	_, err := db.InitDatabse()

	if err != nil {
		return nil, err
	}

	return &App{}, nil
}
