package preset

import "github.com/jheddings/ccnow/internal/types"

func minimalPreset() []types.SegmentNode {
	return []types.SegmentNode{
		{Type: "pwd.name", Provider: "pwd", Style: &types.StyleAttrs{Color: "39"}},
		{Type: "git.branch", Provider: "git", Style: &types.StyleAttrs{Color: "whiteBright", Bold: true, Prefix: " | "}},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Type: "context.tokens", Provider: "context", Style: &types.StyleAttrs{Color: "white"}},
				{Type: "context.size", Provider: "context", Style: &types.StyleAttrs{Color: "white", Prefix: "/"}},
			},
		},
	}
}
