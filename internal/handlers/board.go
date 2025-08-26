// Package handlers
package handlers

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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

	c := templates.Board(b)
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
	// reload store so next /board reflects saved data
	h.store = store.LoadBoard(filepath.Join("./data", "db.json"))

	http.Redirect(w, r, "/board", http.StatusSeeOther)
}

// GetSettings shows a simple settings screen to edit the board or reset it.
func (h *ScoreBoardHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	b := h.store.GetBoard()
	c := templates.Settings(b)
	if err := templates.Layout(c, "Settings").Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PostSettings updates the board name and teams in one go.
// It rebuilds the board from the posted fields and saves it.
func (h *ScoreBoardHandler) PostSettings(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	boardName := strings.TrimSpace(r.FormValue("board_name"))
	if boardName == "" {
		boardName = "Untitled Board"
	}

	teams := make([]*store.Team, 0, 4)
	for i := 1; i <= 4; i++ {
		idx := strconv.Itoa(i)
		name := strings.TrimSpace(r.FormValue("team_name_" + idx))
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

		teams = append(teams, &store.Team{
			TeamName:  name,
			TeamColor: colorForIndex(len(teams)),
			Members:   members,
		})
	}

	// Build new board and persist
	nb := store.NewBoard(boardName)
	for _, t := range teams {
		nb.AddTeam(t)
	}

	_ = os.MkdirAll("./data", 0755)
	_ = nb.SaveToJSON(filepath.Join("./data", "db.json"))
	h.store = store.LoadBoard(filepath.Join("./data", "db.json"))

	http.Redirect(w, r, "/board", http.StatusSeeOther)
}

// PostResetBoard wipes the saved board and sends you to create a new one.
func (h *ScoreBoardHandler) PostResetBoard(w http.ResponseWriter, r *http.Request) {
	// Best-effort delete; if it's not there that's fine.
	_ = os.Remove(filepath.Join("./data", "db.json"))

	// Reset in-memory board too so navigation doesn't show stale data.
	h.store = store.NewBoard("")

	http.Redirect(w, r, "/board/new", http.StatusSeeOther)
}

// GetGames shows the games page to list and add games.
func (h *ScoreBoardHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	b := h.store.GetBoard()
	c := templates.Games(b)
	if err := templates.Layout(c, "Games").Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PostGames creates a game name across all teams if not present.
func (h *ScoreBoardHandler) PostGames(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	gameName := strings.TrimSpace(r.FormValue("game_name"))
	if gameName == "" {
		http.Error(w, "game name required", http.StatusBadRequest)
		return
	}
	b := h.store.GetBoard()
	for _, t := range b.Teams {
		if t == nil {
			continue
		}
		exists := false
		for _, g := range t.Games {
			if g.GameName == gameName {
				exists = true
				break
			}
		}
		if !exists {
			t.Games = append(t.Games, store.Game{GameName: gameName, Rounds: make(map[string]int)})
		}
	}
	_ = os.MkdirAll("./data", 0755)
	_ = b.SaveToJSON(filepath.Join("./data", "db.json"))
	h.store = store.LoadBoard(filepath.Join("./data", "db.json"))
	http.Redirect(w, r, "/games", http.StatusSeeOther)
}

// PostRenameGame renames a game across all teams.
func (h *ScoreBoardHandler) PostRenameGame(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	oldName := strings.TrimSpace(r.FormValue("old_name"))
	newName := strings.TrimSpace(r.FormValue("new_name"))
	if oldName == "" || newName == "" {
		http.Error(w, "old and new names required", http.StatusBadRequest)
		return
	}
	b := h.store.GetBoard()
	for _, t := range b.Teams {
		for i := range t.Games {
			if t.Games[i].GameName == oldName {
				t.Games[i].GameName = newName
			}
		}
	}
	_ = os.MkdirAll("./data", 0755)
	_ = b.SaveToJSON(filepath.Join("./data", "db.json"))
	h.store = store.LoadBoard(filepath.Join("./data", "db.json"))
	http.Redirect(w, r, "/games", http.StatusSeeOther)
}

// PostDeleteGame deletes a game across all teams.
func (h *ScoreBoardHandler) PostDeleteGame(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	b := h.store.GetBoard()
	for _, t := range b.Teams {
		filtered := make([]store.Game, 0, len(t.Games))
		for _, g := range t.Games {
			if g.GameName != name {
				filtered = append(filtered, g)
			}
		}
		t.Games = filtered
	}
	_ = os.MkdirAll("./data", 0755)
	_ = b.SaveToJSON(filepath.Join("./data", "db.json"))
	h.store = store.LoadBoard(filepath.Join("./data", "db.json"))
	http.Redirect(w, r, "/games", http.StatusSeeOther)
}

// GetTeamScores shows a page to edit a team's scores by game/round.
func (h *ScoreBoardHandler) GetTeamScores(w http.ResponseWriter, r *http.Request) {
	teamParam, _ := url.PathUnescape(chi.URLParam(r, "team"))
	b := h.store.GetBoard()
	var team *store.Team
	for _, t := range b.Teams {
		if t != nil && t.TeamName == teamParam {
			team = t
			break
		}
	}
	if team == nil {
		http.NotFound(w, r)
		return
	}
	c := templates.TeamScores(team)
	if err := templates.Layout(c, "Team Scores").Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PostTeamScores upserts a round score for a specific team and game.
func (h *ScoreBoardHandler) PostTeamScores(w http.ResponseWriter, r *http.Request) {
	teamParam, _ := url.PathUnescape(chi.URLParam(r, "team"))
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	gameName := strings.TrimSpace(r.FormValue("game_name"))
	roundName := strings.TrimSpace(r.FormValue("round_name"))
	scoreStr := strings.TrimSpace(r.FormValue("score"))
	if gameName == "" || scoreStr == "" {
		http.Error(w, "game and score required", http.StatusBadRequest)
		return
	}
	b := h.store.GetBoard()
	var team *store.Team
	for _, t := range b.Teams {
		if t != nil && t.TeamName == teamParam {
			team = t
			break
		}
	}
	if team == nil {
		http.NotFound(w, r)
		return
	}
	var game *store.Game
	for i := range team.Games {
		if team.Games[i].GameName == gameName {
			game = &team.Games[i]
			break
		}
	}
	if game == nil {
		http.Error(w, "game does not exist; add it in Games", http.StatusBadRequest)
		return
	}
	if game.Rounds == nil {
		game.Rounds = make(map[string]int)
	}
	// Determine next round using max(existing)+1 to avoid gaps after deletions
	if roundName == "" {
		next := templates.NextRoundForGame(*game)
		roundName = strconv.Itoa(next)
	}
	scoreVal, err := strconv.Atoi(scoreStr)
	if err != nil {
		http.Error(w, "score must be a number", http.StatusBadRequest)
		return
	}
	game.Rounds[roundName] = scoreVal
	_ = os.MkdirAll("./data", 0755)
	_ = b.SaveToJSON(filepath.Join("./data", "db.json"))
	h.store = store.LoadBoard(filepath.Join("./data", "db.json"))
	http.Redirect(w, r, "/board/team/"+url.PathEscape(team.TeamName), http.StatusSeeOther)
}

// PostTeamScoresBulk updates multiple rounds for a specific team/game.
func (h *ScoreBoardHandler) PostTeamScoresBulk(w http.ResponseWriter, r *http.Request) {
	teamParam, _ := url.PathUnescape(chi.URLParam(r, "team"))
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	gameName := strings.TrimSpace(r.FormValue("game_name"))
	if gameName == "" {
		http.Error(w, "game required", http.StatusBadRequest)
		return
	}
	b := h.store.GetBoard()
	var team *store.Team
	for _, t := range b.Teams {
		if t != nil && t.TeamName == teamParam {
			team = t
			break
		}
	}
	if team == nil {
		http.NotFound(w, r)
		return
	}
	var game *store.Game
	for i := range team.Games {
		if team.Games[i].GameName == gameName {
			game = &team.Games[i]
			break
		}
	}
	if game == nil {
		http.Error(w, "game does not exist; add it in Games", http.StatusBadRequest)
		return
	}
	if game.Rounds == nil {
		game.Rounds = make(map[string]int)
	}
	// Expect multiple round_name and score fields (parallel slices by order)
	roundNames := r.Form["round_name"]
	scores := r.Form["score"]
	for i := 0; i < len(roundNames) && i < len(scores); i++ {
		rn := strings.TrimSpace(roundNames[i])
		sc := strings.TrimSpace(scores[i])
		if rn == "" || sc == "" {
			continue
		}
		val, err := strconv.Atoi(sc)
		if err != nil {
			continue
		}
		game.Rounds[rn] = val
	}
	_ = os.MkdirAll("./data", 0755)
	_ = b.SaveToJSON(filepath.Join("./data", "db.json"))
	h.store = store.LoadBoard(filepath.Join("./data", "db.json"))
	http.Redirect(w, r, "/board/team/"+url.PathEscape(team.TeamName), http.StatusSeeOther)
}

// PostDeleteRound deletes a specific round for a team/game.
func (h *ScoreBoardHandler) PostDeleteRound(w http.ResponseWriter, r *http.Request) {
	teamParam, _ := url.PathUnescape(chi.URLParam(r, "team"))
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	gameName := strings.TrimSpace(r.FormValue("game_name"))
	roundName := strings.TrimSpace(r.FormValue("round_name"))
	if gameName == "" || roundName == "" {
		http.Error(w, "game and round required", http.StatusBadRequest)
		return
	}
	b := h.store.GetBoard()
	var team *store.Team
	for _, t := range b.Teams {
		if t != nil && t.TeamName == teamParam {
			team = t
			break
		}
	}
	if team == nil {
		http.NotFound(w, r)
		return
	}
	for i := range team.Games {
		if team.Games[i].GameName == gameName {
			if team.Games[i].Rounds != nil {
				delete(team.Games[i].Rounds, roundName)
			}
			break
		}
	}
	_ = os.MkdirAll("./data", 0755)
	_ = b.SaveToJSON(filepath.Join("./data", "db.json"))
	h.store = store.LoadBoard(filepath.Join("./data", "db.json"))
	http.Redirect(w, r, "/board/team/"+url.PathEscape(team.TeamName), http.StatusSeeOther)
}
