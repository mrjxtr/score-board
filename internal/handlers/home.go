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

func NewHomeHandler(db store.Database) *ScoreBoardHandler {
	return &ScoreBoardHandler{
		store: db,
	}
}

func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	c := templates.Home()
	err := templates.Layout(c, "Home").Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
