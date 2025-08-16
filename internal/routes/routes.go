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

	r.Route("/", func(r chi.Router) {
		r.Get("/board", h.Board.GetScoreBoard)
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
