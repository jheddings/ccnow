package render

import "testing"

func TestFormatValue(t *testing.T) {
	tests := []struct {
		value    any
		format   string
		expected string
	}{
		{42, "", "42"},
		{"hello", "", "hello"},
		{3.14, "", "3.14"},
		{42, "%d%%", "42%"},
		{42, "+%d", "+42"},
		{"text", "%s!", "text!"},
		{3.14, "%.1f", "3.1"},
		{nil, "", ""},
		{nil, "%v", ""},
	}

	for _, tt := range tests {
		result := FormatValue(tt.value, tt.format)
		if result != tt.expected {
			t.Errorf("FormatValue(%v, %q) = %q, want %q", tt.value, tt.format, result, tt.expected)
		}
	}
}
