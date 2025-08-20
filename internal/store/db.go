// Package store
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// CreateNewBoardFile creates a new board and saves it as a JSON file in ./data/ folder
func CreateNewBoardFile(boardName string) (Database, string, error) {
	// Ensure data directory exists
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, "", fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create a safe filename from board name
	safeFilename := strings.ReplaceAll(strings.TrimSpace(boardName), " ", "_")
	safeFilename = strings.ToLower(safeFilename)
	// Remove any characters that aren't alphanumeric, underscores, or hyphens
	cleanFilename := ""
	for _, char := range safeFilename {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' || char == '-' {
			cleanFilename += string(char)
		}
	}

	// Ensure we have a valid filename
	if cleanFilename == "" {
		cleanFilename = "untitled_board"
	}

	// Create full filepath
	filename := fmt.Sprintf("%s.json", cleanFilename)
	filepath := filepath.Join(dataDir, filename)

	// Create new board
	board := NewBoard(boardName)

	// Save to JSON file
	if err := board.SaveToJSON(filepath); err != nil {
		return nil, "", fmt.Errorf("failed to save board to file: %w", err)
	}

	return board, filepath, nil
}

// ListBoardFiles returns a list of all JSON board files in the ./data/ directory
func ListBoardFiles() ([]string, error) {
	dataDir := "./data"

	// Check if data directory exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read directory contents
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	// Filter for JSON files
	var boardFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			boardFiles = append(boardFiles, filepath.Join(dataDir, file.Name()))
		}
	}

	return boardFiles, nil
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
