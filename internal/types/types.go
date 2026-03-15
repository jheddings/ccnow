package types

// SessionData represents the JSON session data piped from Claude Code.
type SessionData struct {
	CWD           string         `json:"cwd"`
	Model         *ModelInfo     `json:"model,omitempty"`
	Cost          *CostInfo      `json:"cost,omitempty"`
	ContextWindow *ContextWindow `json:"context_window,omitempty"`
}

// ModelInfo contains model identification from the session.
type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
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
	UsedPercentage    int           `json:"used_percentage"`
	ContextWindowSize int           `json:"context_window_size,omitempty"`
	CurrentUsage      *CurrentUsage `json:"current_usage,omitempty"`
}

// CurrentUsage breaks down token counts by category.
type CurrentUsage struct {
	InputTokens              int `json:"input_tokens"`
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
// atomic nodes have a Type that maps to a Segment implementation.
type SegmentNode struct {
	Type     string         `json:"segment,omitempty"`
	Provider string         `json:"provider,omitempty"`
	Enabled  *bool          `json:"enabled,omitempty"`
	Style    *StyleAttrs    `json:"style,omitempty"`
	Props    map[string]any `json:"props,omitempty"`
	Children []SegmentNode  `json:"children,omitempty"`

	// EnabledFn is set programmatically (presets) and takes precedence over Enabled.
	EnabledFn func(*SessionData) bool `json:"-"`
}

// SegmentContext is passed to Segment.Render with resolved data.
type SegmentContext struct {
	Session  *SessionData
	Provider any
	Props    map[string]any
}

// Segment renders a single atomic value.
type Segment interface {
	Name() string
	Render(ctx *SegmentContext) *string
}

// DataProvider lazily fetches external data (git, pwd, etc).
type DataProvider interface {
	Name() string
	Resolve(session *SessionData) (any, error)
}
