package style

import (
	"strings"

	"github.com/jheddings/ccglow/internal/types"
)

var colorLevel = 1

// SetColorLevel controls ANSI output. 0 disables colors (plain mode).
func SetColorLevel(level int) {
	colorLevel = level
}

// Apply wraps a value with ANSI escape codes and prefix/suffix.
func Apply(value string, attrs *types.StyleAttrs) string {
	if attrs == nil {
		return value
	}

	if attrs.Prefix != "" {
		value = attrs.Prefix + value
	}
	if attrs.Suffix != "" {
		value = value + attrs.Suffix
	}

	if colorLevel > 0 {
		var mods strings.Builder
		if attrs.Bold {
			mods.WriteString(ansiBold)
		}
		if attrs.Italic {
			mods.WriteString(ansiItalic)
		}

		colorCode := resolveColor(attrs.Color)
		bgCode := resolveBgColor(attrs.Background)
		if colorCode != "" || bgCode != "" || mods.Len() > 0 {
			return ansiReset + mods.String() + colorCode + bgCode + value + ansiReset
		}
	}

	return value
}
