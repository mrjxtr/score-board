// Package handlers
package handlers

import (
	"net/http"

	"github.com/mrjxtr-dev/score-board/internal/store"
	"github.com/mrjxtr-dev/score-board/internal/templates"
)

type HomeHandler struct {
	store store.Database
}

// NewHomeHandler creates a HomeHandler bound to the database.
func NewHomeHandler(db store.Database) *HomeHandler {
	return &HomeHandler{
		store: db,
	}
}

// GetHome renders the landing page.
func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	c := templates.Home()
	err := templates.Layout(c, "Home").Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAbout renders the about page.
func (h *HomeHandler) GetAbout(w http.ResponseWriter, r *http.Request) {
	c := templates.About()
	err := templates.Layout(c, "About").Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
