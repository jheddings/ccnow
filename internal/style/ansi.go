package style

import (
	"fmt"
	"strconv"
	"strings"
)

var namedFgColors = map[string]string{
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

var namedBgColors = map[string]string{
	"black":         "40",
	"red":           "41",
	"green":         "42",
	"yellow":        "43",
	"blue":          "44",
	"magenta":       "45",
	"cyan":          "46",
	"white":         "47",
	"blackBright":   "100",
	"redBright":     "101",
	"greenBright":   "102",
	"yellowBright":  "103",
	"blueBright":    "104",
	"magentaBright": "105",
	"cyanBright":    "106",
	"whiteBright":   "107",
}

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiItalic = "\x1b[3m"
)

func resolveColor(color string) string {
	return resolveAnsiColor(color, namedFgColors, 38)
}

func resolveBgColor(color string) string {
	return resolveAnsiColor(color, namedBgColors, 48)
}

func resolveAnsiColor(color string, named map[string]string, extBase int) string {
	if color == "" {
		return ""
	}

	if code, ok := named[color]; ok {
		return fmt.Sprintf("\x1b[%sm", code)
	}

	if strings.HasPrefix(color, "#") && len(color) == 7 {
		r, _ := strconv.ParseInt(color[1:3], 16, 64)
		g, _ := strconv.ParseInt(color[3:5], 16, 64)
		b, _ := strconv.ParseInt(color[5:7], 16, 64)
		return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", extBase, r, g, b)
	}

	if n, err := strconv.Atoi(color); err == nil && n >= 0 && n <= 255 {
		return fmt.Sprintf("\x1b[%d;5;%dm", extBase, n)
	}

	return ""
}
