package style

import (
	"strings"

	"github.com/jheddings/ccnow/internal/types"
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

	styled := value

	if colorLevel > 0 {
		var mods strings.Builder
		if attrs.Bold {
			mods.WriteString(ansiBold)
		}
		if attrs.Italic {
			mods.WriteString(ansiItalic)
		}

		colorCode := resolveColor(attrs.Color)
		if colorCode != "" || mods.Len() > 0 {
			styled = ansiReset + mods.String() + colorCode + value + ansiReset
		}
	}

	if attrs.Prefix != "" {
		styled = attrs.Prefix + styled
	}
	if attrs.Suffix != "" {
		styled = styled + attrs.Suffix
	}

	return styled
}
