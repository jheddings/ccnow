package render

import "fmt"

func FormatValue(value any, format string) string {
	if value == nil {
		return ""
	}
	if format == "" {
		return fmt.Sprintf("%v", value)
	}
	return fmt.Sprintf(format, value)
}
