// Package store
package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ScoreBoard represents a score board
type ScoreBoard struct {
	BoardName string  `json:"board"`
	Teams     []*Team `json:"teams"`
}

// Team represents a team
type Team struct {
	TeamName  string            `json:"team"`
	TeamColor map[string]string `json:"color"`
	Members   []string          `json:"members"`
	Games     []Game            `json:"games"`
}

// Game represents a game
type Game struct {
	GameName string         `json:"game"`
	Rounds   map[string]int `json:"rounds"`
}

type Database interface {
	SaveToJSON(filename string) error
	GetBoard() *ScoreBoard
}

// LoadDB boots the single-app scoreboard from ./data/db.json.
// If it's missing or unreadable, it spins up a default board and saves it.
func LoadDB() Database {
	const dataDir = "./data"
	const dbFilename = "db.json"

	// Ensure data directory exists; if it fails, still return an in-memory default
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return NewBoard("Default Board")
	}

	fullpath := filepath.Join(dataDir, dbFilename)

	// If no db file yet, create one with a default board
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		board := NewBoard("Default Board")
		_ = board.SaveToJSON(fullpath)
		return board
	}

	// Load existing board
	return LoadBoard(fullpath)
}

// LoadDB loads a score board from a JSON file
func LoadBoard(filename string) Database {
	board := &ScoreBoard{}

	// Try to read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		// If file doesn't exist, return a new empty board
		return NewBoard("Default Board")
	}

	// Try to parse the JSON
	if err := json.Unmarshal(data, board); err != nil {
		// If parsing fails, return a new empty board
		return NewBoard("Default Board")
	}

	return board
}

// NewBoard creates a new score board
func NewBoard(name string) *ScoreBoard {
	return &ScoreBoard{
		BoardName: name,
		Teams:     []*Team{},
	}
}

// GetBoard returns the underlying scoreboard instance.
func (b *ScoreBoard) GetBoard() *ScoreBoard {
	return b
}

// AddTeam adds a team to the board
func (b *ScoreBoard) AddTeam(team *Team) {
	for _, t := range b.Teams {
		if t.TeamName == team.TeamName {
			return
		}
	}
	b.Teams = append(b.Teams, team)
}

// RemoveTeam removes a given team from the board
func (b *ScoreBoard) RemoveTeam(team *Team) {
	newTeams := make([]*Team, 0, len(b.Teams))

	for _, t := range b.Teams {
		if t.TeamName != team.TeamName {
			newTeams = append(newTeams, t)
		}
	}
	b.Teams = newTeams
}

// EditTeam edits the old team with the new team values
func (b *ScoreBoard) EditTeam(oldTeam, updates *Team) {
	for _, t := range b.Teams {
		if oldTeam.TeamName == t.TeamName {
			if updates.TeamName != "" {
				t.TeamName = updates.TeamName
			}
			if updates.TeamColor != nil {
				t.TeamColor = updates.TeamColor
			}
			if updates.Members != nil {
				t.Members = updates.Members
			}
			if updates.Games != nil {
				t.Members = updates.Members
			}

			return
		}
	}
}

// SaveToJSON saves the current board to a JSON file
func (b *ScoreBoard) SaveToJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(b); err != nil {
		return err
	}

	return nil
}
