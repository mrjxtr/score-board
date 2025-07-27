// Package handlers
package handlers

import (
	"net/http"

	"github.com/mrjxtr-dev/score-board/internal/store"
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
}
