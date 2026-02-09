package display

import "github.com/fatih/color"

var (
	Green  = color.New(color.FgGreen).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Gray   = color.New(color.FgHiBlack).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
)

// StatusColor returns the status string with appropriate color.
func StatusColor(status string) string {
	switch status {
	case "active":
		return Green(status)
	case "wip":
		return Blue(status)
	case "recent":
		return Yellow(status)
	case "paused":
		return Gray(status)
	default:
		return status
	}
}
