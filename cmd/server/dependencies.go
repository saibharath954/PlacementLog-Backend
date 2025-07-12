package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/varnit-ta/PlacementLog/internal/db"
	"github.com/varnit-ta/PlacementLog/internal/posts"
	userauth "github.com/varnit-ta/PlacementLog/internal/userAuth"
)

type App struct {
	userAuthHandler *userauth.UserAuthHandler
	postHandler     *posts.PostsHandler
}

func InitApp() (*App, error) {
	db, err := db.InitDatabse()

	if err != nil {
		return nil, err
	}

	userAuthRepo := userauth.NewUserAuthRepo(db)
	userAuthService := userauth.NewUserAuthService(userAuthRepo)
	userAuthHandler := userauth.NewUserAuthHandler(userAuthService)

	postRepo := posts.NewPostsRepo(db)
	postService := posts.NewPostsService(postRepo)
	postHandler := posts.NewPostsHandler(postService)

	return &App{
		userAuthHandler: userAuthHandler,
		postHandler:     postHandler,
	}, nil
}

func (a App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Post("/auth/login", a.userAuthHandler.Login)
		r.Post("/auth/register", a.userAuthHandler.Register)
	})

	r.Group(func(r chi.Router) {
		r.Post("/posts", a.postHandler.AddPost)
		r.Get("/posts", a.postHandler.GetAll)
		r.Get("/posts/user", a.postHandler.GetByUser)
		r.Delete("/posts", a.postHandler.DeletePost)
	})

	return r
}
