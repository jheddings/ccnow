package types

// SessionData represents the JSON session data piped from Claude Code.
type SessionData struct {
	CWD           string           `json:"cwd"`
	SessionID     string           `json:"session_id,omitempty"`
	Model         *ModelInfo       `json:"model,omitempty"`
	Cost          *CostInfo        `json:"cost,omitempty"`
	ContextWindow *ContextWindow   `json:"context_window,omitempty"`
	Version       string           `json:"version,omitempty"`
	OutputStyle   *OutputStyleInfo `json:"output_style,omitempty"`
}

// ModelInfo contains model identification from the session.
type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// OutputStyleInfo contains the output style configuration.
type OutputStyleInfo struct {
	Name string `json:"name"`
}

// CostInfo contains session cost and timing data.
type CostInfo struct {
	TotalCostUSD       float64 `json:"total_cost_usd"`
	TotalDurationMS    float64 `json:"total_duration_ms"`
	TotalAPIDurationMS float64 `json:"total_api_duration_ms"`
	TotalLinesAdded    int     `json:"total_lines_added"`
	TotalLinesRemoved  int     `json:"total_lines_removed"`
}

// ContextWindow contains token usage data.
type ContextWindow struct {
	UsedPercentage      int           `json:"used_percentage"`
	RemainingPercentage int           `json:"remaining_percentage"`
	ContextWindowSize   int           `json:"context_window_size,omitempty"`
	TotalInputTokens    *int          `json:"total_input_tokens,omitempty"`
	TotalOutputTokens   *int          `json:"total_output_tokens,omitempty"`
	CurrentUsage        *CurrentUsage `json:"current_usage,omitempty"`
}

// CurrentUsage breaks down token counts by category.
type CurrentUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// StyleAttrs controls ANSI styling for a segment.
type StyleAttrs struct {
	Color      string `json:"color,omitempty"`
	Background string `json:"bgcolor,omitempty"`
	Bold       bool   `json:"bold,omitempty"`
	Italic     bool   `json:"italic,omitempty"`
	Prefix     string `json:"prefix,omitempty"`
	Suffix     string `json:"suffix,omitempty"`
}

// SegmentNode is a node in the render tree. Composite nodes have Children;
// atomic nodes have an Expr (evaluated against env) or a static Value.
type SegmentNode struct {
	Expr     string        `json:"expr,omitempty"`
	Value    any           `json:"value,omitempty"`
	Format   string        `json:"format,omitempty"`
	When     string        `json:"when,omitempty"`
	Enabled  *bool         `json:"enabled,omitempty"`
	Style    *StyleAttrs   `json:"style,omitempty"`
	Children []SegmentNode `json:"children,omitempty"`

	// EnabledFn is set programmatically (presets) and takes precedence over Enabled.
	EnabledFn func(*SessionData) bool `json:"-"`
}

// ProviderResult holds the values and optional default formats returned by a provider.
type ProviderResult struct {
	Values  map[string]any    // nested maps (e.g. {"git": {"branch": "main"}})
	Formats map[string]string // flat dotted keys (e.g. "context.percent.used" -> "%d%%")
}

// DataProvider fetches external data and returns named segment values.
type DataProvider interface {
	Name() string
	Resolve(session *SessionData) (*ProviderResult, error)
}
