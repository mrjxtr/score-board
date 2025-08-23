package templates

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
