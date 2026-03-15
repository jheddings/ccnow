package preset

import "github.com/jheddings/ccglow/internal/types"

func f1Preset() []types.SegmentNode {
	return []types.SegmentNode{
		// ── Line 1: Location & Git ──────────────────────────

		// Red section: working directory
		{Type: "pwd.smart", Provider: "pwd", Style: &types.StyleAttrs{
			Color: "white", Background: "#DC0000", Prefix: " ", Suffix: " ",
		}},
		{Type: "pwd.name", Provider: "pwd", Style: &types.StyleAttrs{
			Color: "white", Background: "#DC0000", Bold: true, Suffix: " ",
		}},

		// Chevron: red → carbon
		{Type: "literal", Props: map[string]any{"text": "\ue0b0"}, Style: &types.StyleAttrs{
			Color: "#DC0000", Background: "#3A3A3A",
		}},

		// Carbon section: git info
		{Type: "git.branch", Provider: "git", Style: &types.StyleAttrs{
			Color: "white", Background: "#3A3A3A", Bold: true, Prefix: " \ue0a0 ",
		}},
		{Type: "git.insertions", Provider: "git", Style: &types.StyleAttrs{
			Color: "#50FA7B", Background: "#3A3A3A", Prefix: " +",
		}},
		{Type: "git.deletions", Provider: "git", Style: &types.StyleAttrs{
			Color: "#FF6B6B", Background: "#3A3A3A", Prefix: " -", Suffix: " ",
		}},

		// Chevron: carbon → terminal
		{Type: "literal", Props: map[string]any{"text": "\ue0b0"}, Style: &types.StyleAttrs{
			Color: "#3A3A3A",
		}},

		// Newline
		{Type: "literal", Props: map[string]any{"text": "\n"}},

		// ── Line 2: Session Info ────────────────────────────

		// Blue section: model
		{Type: "model.name", Provider: "model", Style: &types.StyleAttrs{
			Color: "white", Background: "#003DA5", Bold: true, Prefix: " ", Suffix: " ",
		}},

		// Chevron: blue → graphite
		{Type: "literal", Props: map[string]any{"text": "\ue0b0"}, Style: &types.StyleAttrs{
			Color: "#003DA5", Background: "#2D2D2D",
		}},

		// Graphite section: context window
		{Type: "context.tokens", Provider: "context", Style: &types.StyleAttrs{
			Color: "white", Background: "#2D2D2D", Bold: true, Prefix: " ",
		}},
		{Type: "context.size", Provider: "context", Style: &types.StyleAttrs{
			Color: "white", Background: "#2D2D2D", Prefix: "/",
		}},
		{Type: "context.percent", Provider: "context", Style: &types.StyleAttrs{
			Color: "white", Background: "#2D2D2D", Prefix: " (", Suffix: ") ",
		}},

		// Chevron: graphite → papaya
		{Type: "literal", Props: map[string]any{"text": "\ue0b0"}, Style: &types.StyleAttrs{
			Color: "#2D2D2D", Background: "#FF8000",
		}},

		// Papaya section: cost
		{Type: "cost.usd", Provider: "cost", Style: &types.StyleAttrs{
			Color: "white", Background: "#FF8000", Bold: true, Prefix: " ", Suffix: " ",
		}},

		// Chevron: papaya → racing green
		{Type: "literal", Props: map[string]any{"text": "\ue0b0"}, Style: &types.StyleAttrs{
			Color: "#FF8000", Background: "#006633",
		}},

		// Racing green section: duration
		{Type: "session.duration", Provider: "session", Style: &types.StyleAttrs{
			Color: "white", Background: "#006633", Prefix: " ", Suffix: " ",
		}},

		// Chevron: racing green → silver
		{Type: "literal", Props: map[string]any{"text": "\ue0b0"}, Style: &types.StyleAttrs{
			Color: "#006633", Background: "#3A3A3A",
		}},

		// Silver section: line diffs
		{Type: "session.lines-added", Provider: "session", Style: &types.StyleAttrs{
			Color: "#50FA7B", Background: "#3A3A3A", Prefix: " +",
		}},
		{Type: "session.lines-removed", Provider: "session", Style: &types.StyleAttrs{
			Color: "#FF6B6B", Background: "#3A3A3A", Prefix: " -", Suffix: " ",
		}},

		// Chevron: silver → terminal
		{Type: "literal", Props: map[string]any{"text": "\ue0b0"}, Style: &types.StyleAttrs{
			Color: "#3A3A3A",
		}},
	}
}
