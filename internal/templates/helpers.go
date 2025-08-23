package templates

import (
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
