package devslog

import "log/slog"

const (
	resetColour  = "\033[0m"
	colourRed    = "\033[31m"
	colourYellow = "\033[33m"
	colourWhite  = "\033[37m"
	colorGray    = "\033[90m"
)

func levelColour(l slog.Level) string {
	switch l {
	case slog.LevelError:
		return colourRed
	case slog.LevelWarn:
		return colourYellow
	case slog.LevelInfo:
		return colourWhite
	case slog.LevelDebug:
		return ""
	default:
		return ""
	}
}

func text(color string, text string) string {
	return color + text + resetColour
}

func gray(text string) string {
	return colorGray + text + resetColour
}
