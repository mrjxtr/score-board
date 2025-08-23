// Package routes
package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrjxtr-dev/score-board/internal/config"
	"github.com/mrjxtr-dev/score-board/internal/handlers"
	"github.com/mrjxtr-dev/score-board/internal/store"
)

func SetupRoutes(cfg *config.Config, db store.Database) *chi.Mux {
	r := chi.NewRouter()
	setupGlobalMiddleware(r)

	h := handlers.NewHandlers(db)

	fileserver := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileserver))

	r.Get("/", h.Home.GetHome)
	r.Get("/about", h.Home.GetAbout)

	// Settings: edit/update board and reset
	r.Get("/settings", h.Board.GetSettings)
	r.Post("/settings", h.Board.PostSettings)
	r.Post("/settings/reset", h.Board.PostResetBoard)

	r.Route("/board", func(r chi.Router) {
		r.Get("/", h.Board.GetScoreBoard)
		r.Get("/new", h.Board.GetNewBoard)
		r.Post("/new", h.Board.PostNewBoard)
	})
	return r
}

func setupGlobalMiddleware(r *chi.Mux) {
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/ping"),
	)
}
