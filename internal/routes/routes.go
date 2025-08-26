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

// SetupRoutes wires all HTTP routes. If staticFS is non-nil, it will serve
// files from it at /static/; otherwise it falls back to the local ./static dir.
func SetupRoutes(cfg *config.Config, db store.Database, staticFS http.FileSystem) *chi.Mux {
	r := chi.NewRouter()
	setupGlobalMiddleware(r)

	h := handlers.NewHandlers(db)

	var fs http.FileSystem
	if staticFS != nil {
		fs = staticFS
	} else {
		fs = http.Dir("./static")
	}
	fileserver := http.FileServer(fs)
	r.Handle("/static/*", http.StripPrefix("/static/", fileserver))

	r.Get("/", h.Home.GetHome)
	r.Get("/about", h.Home.GetAbout)

	// Games: list/add/rename/delete
	r.Get("/games", h.Board.GetGames)
	r.Post("/games", h.Board.PostGames)
	r.Post("/games/rename", h.Board.PostRenameGame)
	r.Post("/games/delete", h.Board.PostDeleteGame)

	// Settings: edit/update board and reset
	r.Get("/settings", h.Board.GetSettings)
	r.Post("/settings", h.Board.PostSettings)
	r.Post("/settings/reset", h.Board.PostResetBoard)

	r.Route("/board", func(r chi.Router) {
		r.Get("/", h.Board.GetScoreBoard)
		r.Get("/new", h.Board.GetNewBoard)
		r.Post("/new", h.Board.PostNewBoard)
		// Team scores
		r.Get("/team/{team}", h.Board.GetTeamScores)
		r.Post("/team/{team}/scores", h.Board.PostTeamScores)
		r.Post("/team/{team}/scores/bulk", h.Board.PostTeamScoresBulk)
		r.Post("/team/{team}/scores/delete", h.Board.PostDeleteRound)
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
