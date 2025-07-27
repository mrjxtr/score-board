// Package store
package store

import (
	"encoding/json"
	"os"
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
}

func LoadDB() Database {
	return nil
}

// NewBoard creates a new score board
func NewBoard(name string) *ScoreBoard {
	return &ScoreBoard{
		BoardName: name,
		Teams:     []*Team{},
	}
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
