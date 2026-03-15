package preset

import "github.com/jheddings/ccnow/internal/types"

func defaultPreset() []types.SegmentNode {
	return []types.SegmentNode{
		{Type: "pwd.smart", Provider: "pwd", Style: &types.StyleAttrs{Color: "31"}},
		{Type: "pwd.name", Provider: "pwd", Style: &types.StyleAttrs{Color: "39", Bold: true}},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " | ", Color: "240"},
			Children: []types.SegmentNode{
				{Type: "git.branch", Provider: "git", Style: &types.StyleAttrs{Color: "whiteBright", Bold: true, Prefix: "\ue0a0 "}},
				{Type: "git.insertions", Provider: "git", Style: &types.StyleAttrs{Color: "green", Prefix: " \u00b7 +"}},
				{Type: "git.deletions", Provider: "git", Style: &types.StyleAttrs{Color: "red", Prefix: " -"}},
			},
		},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Type: "context.tokens", Provider: "context", Style: &types.StyleAttrs{Color: "white", Bold: true}},
				{Type: "context.percent", Provider: "context", Style: &types.StyleAttrs{Color: "white", Prefix: " (", Suffix: ")"}},
			},
		},
		{Type: "session.duration", Provider: "session", Style: &types.StyleAttrs{Color: "magenta", Prefix: " \u00b7 "}},
	}
}
