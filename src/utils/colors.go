package utils

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

var upstreamPalette = []string{
	colorCyan,
	colorMagenta,
	colorYellow,
	colorGreen,
	colorBlue,
	colorRed,
}

// ColorForIndex returns a deterministic color code for a given upstream index.
func ColorForIndex(i int) string {
	if len(upstreamPalette) == 0 {
		return ""
	}
	return upstreamPalette[i%len(upstreamPalette)]
}

// Colorize wraps text with a color and reset code.
func Colorize(text, color string) string {
	if color == "" {
		return text
	}
	return color + text + colorReset
}

// ColorForMethod returns a color code for a given HTTP method, following common OpenAPI-style colors.
func ColorForMethod(method string) string {
	switch method {
	case "GET":
		return colorBlue
	case "POST":
		return colorGreen
	case "PUT":
		return colorYellow
	case "DELETE":
		return colorRed
	case "PATCH":
		return colorMagenta
	default:
		return ""
	}
}


