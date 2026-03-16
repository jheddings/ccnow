package render

import (
	"testing"

	"github.com/jheddings/ccglow/internal/segment"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
)

func setupTestRegistries() *segment.Registry {
	reg := segment.NewRegistry()
	segment.RegisterBuiltin(reg)
	return reg
}

func TestTree_Empty(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}
	result := Tree(nil, seg, sess, map[string]any{}, map[string]any{}, TagIndex{})
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestTree_AtomicNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/home/user/project"}

	segmentValues := map[string]any{
		"pwd.name": "project",
	}
	tagIdx := TagIndex{
		"pwd.name": fieldAccessor{Provider: "pwd", FieldIndex: 0},
	}

	tree := []types.SegmentNode{
		{Type: "pwd.name"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "project" {
		t.Errorf("expected project, got %q", result)
	}
}

func TestTree_CompositeCollapse(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"git.branch":     nil,
		"git.insertions": nil,
	}
	tagIdx := TagIndex{
		"git.branch":     fieldAccessor{Provider: "git", FieldIndex: 0},
		"git.insertions": fieldAccessor{Provider: "git", FieldIndex: 1},
	}

	tree := []types.SegmentNode{
		{
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Type: "git.branch"},
				{Type: "git.insertions"},
			},
		},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (collapsed composite), got %q", result)
	}
}

func TestTree_DisabledNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"pwd.name": "tmp",
	}
	tagIdx := TagIndex{
		"pwd.name": fieldAccessor{Provider: "pwd", FieldIndex: 0},
	}

	disabled := false
	tree := []types.SegmentNode{
		{Type: "pwd.name", Enabled: &disabled},
	}

	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty for disabled node, got %q", result)
	}
}

func TestTree_Literal(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	tree := []types.SegmentNode{
		{Type: "literal", Props: map[string]any{"text": "hello"}},
	}

	result := Tree(tree, seg, sess, map[string]any{}, map[string]any{}, TagIndex{})
	if result != "hello" {
		t.Errorf("expected hello, got %q", result)
	}
}

func TestTree_MissingSegment(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	tree := []types.SegmentNode{
		{Type: "nonexistent.segment"},
	}

	result := Tree(tree, seg, sess, map[string]any{}, map[string]any{}, TagIndex{})
	if result != "" {
		t.Errorf("expected empty for missing segment, got %q", result)
	}
}

func TestTree_DataSegment(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": "hello"}
	tagIdx := TagIndex{"test.name": fieldAccessor{Provider: "test", FieldIndex: 0}}

	tree := []types.SegmentNode{{Type: "test.name"}}
	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_DataSegmentNilCollapses(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": nil}
	tagIdx := TagIndex{"test.name": fieldAccessor{Provider: "test", FieldIndex: 0}}

	tree := []types.SegmentNode{{Type: "test.name"}}
	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (collapsed), got %q", result)
	}
}

func TestTree_DataSegmentWithFormat(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.count": 42}
	tagIdx := TagIndex{"test.count": fieldAccessor{Provider: "test", FieldIndex: 1}}

	tree := []types.SegmentNode{{Type: "test.count", Format: "+%d"}}
	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "+42" {
		t.Errorf("expected '+42', got %q", result)
	}
}

func TestTree_DataSegmentDefaultFormat(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"fmt.pct": 85}
	tagIdx := TagIndex{"fmt.pct": fieldAccessor{Provider: "fmt", FieldIndex: 0, DefaultFormat: "%d%%"}}

	tree := []types.SegmentNode{{Type: "fmt.pct"}}
	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "85%" {
		t.Errorf("expected '85%%', got %q", result)
	}
}

func TestTree_DataSegmentFormatOverridesDefault(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"fmt.pct": 85}
	tagIdx := TagIndex{"fmt.pct": fieldAccessor{Provider: "fmt", FieldIndex: 0, DefaultFormat: "%d%%"}}

	tree := []types.SegmentNode{{Type: "fmt.pct", Format: "(%d)"}}
	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "(85)" {
		t.Errorf("expected '(85)', got %q", result)
	}
}

func TestTree_EmptyStringCollapses(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.plain": ""}
	tagIdx := TagIndex{"test.plain": fieldAccessor{Provider: "test", FieldIndex: 2}}

	tree := []types.SegmentNode{{Type: "test.plain"}}
	result := Tree(tree, seg, sess, map[string]any{}, segmentValues, tagIdx)
	if result != "" {
		t.Errorf("expected empty (collapsed), got %q", result)
	}
}

func TestCollectProviderNames_TagIndex(t *testing.T) {
	tagIdx := TagIndex{
		"git.branch":     fieldAccessor{Provider: "git", FieldIndex: 0},
		"context.tokens": fieldAccessor{Provider: "context", FieldIndex: 0},
	}

	tree := []types.SegmentNode{
		{Type: "git.branch"},
		{Type: "context.tokens"},
		{Type: "literal", Props: map[string]any{"text": "hi"}},
	}

	names := CollectProviderNames(tree, tagIdx)
	if !names["git"] {
		t.Error("expected git provider")
	}
	if !names["context"] {
		t.Error("expected context provider")
	}
	if len(names) != 2 {
		t.Errorf("expected 2 providers, got %d", len(names))
	}
}
