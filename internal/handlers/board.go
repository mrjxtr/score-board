// Package handlers
package handlers

import (
	"net/http"

	"github.com/mrjxtr-dev/score-board/internal/store"
	"github.com/mrjxtr-dev/score-board/internal/templates"
)

type ScoreBoardHandler struct {
	store store.Database
}

func NewScoreBoardHandler(db store.Database) *ScoreBoardHandler {
	return &ScoreBoardHandler{
		store: db,
	}
}

func (h *ScoreBoardHandler) GetScoreBoard(w http.ResponseWriter, r *http.Request) {
	c := templates.Test()
	err := templates.Layout(c, "Just A Test").Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
