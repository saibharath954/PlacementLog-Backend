package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/varnit-ta/PlacementLog/internal/db"
	userauth "github.com/varnit-ta/PlacementLog/internal/userAuth"
)

type App struct {
	userAuthHandler *userauth.UserAuthHandler
}

func InitApp() (*App, error) {
	db, err := db.InitDatabse()

	if err != nil {
		return nil, err
	}

	userAuthRepo := userauth.NewUserAuthRepo(db)
	userAuthService := userauth.NewUserAuthService(userAuthRepo)
	userAuthHandler := userauth.NewUserAuthHandler(userAuthService)

	return &App{
		userAuthHandler: userAuthHandler,
	}, nil
}

func (a App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/auth/login", a.userAuthHandler.Login)
		r.Post("/auth/register", a.userAuthHandler.Register)
	})

	return r
}
