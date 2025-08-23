// Package handlers
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mrjxtr-dev/score-board/internal/config"
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

// GetScoreBoard renders the board page or redirects to creation if empty.
func (h *ScoreBoardHandler) GetScoreBoard(w http.ResponseWriter, r *http.Request) {
	b := h.store.GetBoard()
	if b == nil || len(b.Teams) == 0 || b.BoardName == "" {
		http.Redirect(w, r, "/board/new", http.StatusSeeOther)
		return
	}

	c := templates.Board()
	err := templates.Layout(c, "Score Board").Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetNewBoard shows the form for creating a scoreboard.
func (h *ScoreBoardHandler) GetNewBoard(w http.ResponseWriter, r *http.Request) {
	c := templates.CreateBoard()
	err := templates.Layout(c, "Create Board").Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// colorForIndex picks a default color based on index 0..3.
func colorForIndex(i int) map[string]string {
	order := []string{"pink", "red", "blue", "yellow"}
	if i < 0 || i >= len(order) {
		return map[string]string{"color": "#FFFFFF"}
	}
	name := order[i]
	hex := config.DefaulColors[name]
	return map[string]string{"color": hex}
}

// PostNewBoard handles form submission to create a scoreboard.
func (h *ScoreBoardHandler) PostNewBoard(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	boardName := r.FormValue("board_name")
	if boardName == "" {
		http.Error(w, "board name required", http.StatusBadRequest)
		return
	}

	teams := make([]*store.Team, 0, 4)
	for i := 1; i <= 4; i++ {
		idx := strconv.Itoa(i)
		name := r.FormValue("team_name_" + idx)
		membersRaw := r.FormValue("team_members_" + idx)
		var members []string
		if strings.TrimSpace(membersRaw) != "" {
			parts := strings.Split(membersRaw, ",")
			for _, p := range parts {
				trim := strings.TrimSpace(p)
				if trim != "" {
					members = append(members, trim)
				}
			}
		}
		if name == "" && len(members) == 0 {
			continue
		}
		if name == "" {
			continue
		}
		teams = append(teams, &store.Team{
			TeamName:  name,
			TeamColor: colorForIndex(len(teams)),
			Members:   members,
		})
	}

	b := store.NewBoard(boardName)
	for _, t := range teams {
		b.AddTeam(t)
	}

	_ = os.MkdirAll("./data", 0755)
	_ = b.SaveToJSON(filepath.Join("./data", "db.json"))

	http.Redirect(w, r, "/board", http.StatusSeeOther)
}
