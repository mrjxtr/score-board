package templates

import (
	"strconv"
	"strings"

	"github.com/mrjxtr-dev/score-board/internal/store"
)

// DefaultColorHex returns the default team color (hex) for a 1-based index.
// Keeps it simple: pink, red, blue, yellow in that order.
func DefaultColorHex(i int) string {
	switch i {
	case 1:
		return "#D50059" // pink
	case 2:
		return "#C50000" // red
	case 3:
		return "#1D03AF" // blue
	case 4:
		return "#FFBB02" // yellow
	default:
		return "#FFFFFF"
	}
}

// TeamNameAt returns the team name at zero-based index, or empty if out of range.
func TeamNameAt(b *store.ScoreBoard, idx int) string {
	if b == nil || idx < 0 || idx >= len(b.Teams) {
		return ""
	}
	return b.Teams[idx].TeamName
}

// TeamMembersCSVAt returns members as a comma-separated string for a team index.
func TeamMembersCSVAt(b *store.ScoreBoard, idx int) string {
	if b == nil || idx < 0 || idx >= len(b.Teams) {
		return ""
	}
	if len(b.Teams[idx].Members) == 0 {
		return ""
	}
	return strings.Join(b.Teams[idx].Members, ", ")
}

// UniqueGameNames scans all teams and returns unique game names in first-seen order.
func UniqueGameNames(b *store.ScoreBoard) []string {
	if b == nil || len(b.Teams) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	names := make([]string, 0)
	for _, t := range b.Teams {
		if t == nil || len(t.Games) == 0 {
			continue
		}
		for _, g := range t.Games {
			if g.GameName == "" {
				continue
			}
			if _, ok := seen[g.GameName]; ok {
				continue
			}
			seen[g.GameName] = struct{}{}
			names = append(names, g.GameName)
		}
	}
	return names
}

// NextRoundForGame returns the next round number as max(existing round as int) + 1.
// If no rounds exist or none are numeric, it returns 1.
func NextRoundForGame(g store.Game) int {
	maxRound := 0
	for rn := range g.Rounds {
		if v, err := strconv.Atoi(strings.TrimSpace(rn)); err == nil {
			if v > maxRound {
				maxRound = v
			}
		}
	}
	return maxRound + 1
}
