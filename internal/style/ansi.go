package style

import (
	"fmt"
	"strconv"
	"strings"
)

var namedColors = map[string]string{
	"black":         "30",
	"red":           "31",
	"green":         "32",
	"yellow":        "33",
	"blue":          "34",
	"magenta":       "35",
	"cyan":          "36",
	"white":         "37",
	"blackBright":   "90",
	"redBright":     "91",
	"greenBright":   "92",
	"yellowBright":  "93",
	"blueBright":    "94",
	"magentaBright": "95",
	"cyanBright":    "96",
	"whiteBright":   "97",
}

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiItalic = "\x1b[3m"
)

func resolveColor(color string) string {
	if color == "" {
		return ""
	}

	if code, ok := namedColors[color]; ok {
		return fmt.Sprintf("\x1b[%sm", code)
	}

	if strings.HasPrefix(color, "#") && len(color) == 7 {
		r, _ := strconv.ParseInt(color[1:3], 16, 64)
		g, _ := strconv.ParseInt(color[3:5], 16, 64)
		b, _ := strconv.ParseInt(color[5:7], 16, 64)
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	}

	if n, err := strconv.Atoi(color); err == nil && n >= 0 && n <= 255 {
		return fmt.Sprintf("\x1b[38;5;%dm", n)
	}

	return ""
}
