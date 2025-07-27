package handlers

import "github.com/mrjxtr-dev/score-board/internal/store"

type Handlers struct {
	Board *ScoreBoardHandler
}

func NewHandlers(db store.Database) *Handlers {
	return &Handlers{
		Board: NewScoreBoardHandler(db),
	}
}
