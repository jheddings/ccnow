package render

import (
	"testing"

	"github.com/jheddings/ccglow/internal/condition"
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
	result := Tree(nil, seg, sess, map[string]any{}, map[string]string{}, map[string]any{})
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
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{Type: "pwd.name"},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
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
		"git.branch":     "",
		"git.insertions": "",
	}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Type: "git.branch"},
				{Type: "git.insertions"},
			},
		},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
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
	condEnv := condition.BuildNestedEnv(segmentValues)

	disabled := false
	tree := []types.SegmentNode{
		{Type: "pwd.name", Enabled: &disabled},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
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

	result := Tree(tree, seg, sess, map[string]any{}, map[string]string{}, map[string]any{})
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

	result := Tree(tree, seg, sess, map[string]any{}, map[string]string{}, map[string]any{})
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
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{{Type: "test.name"}}
	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_DataSegmentEmptyCollapses(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": ""}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{{Type: "test.name"}}
	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
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
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{{Type: "test.count", Format: "+%d"}}
	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
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
	defaultFormats := map[string]string{"fmt.pct": "%d%%"}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{{Type: "fmt.pct"}}
	result := Tree(tree, seg, sess, segmentValues, defaultFormats, condEnv)
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
	defaultFormats := map[string]string{"fmt.pct": "%d%%"}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{{Type: "fmt.pct", Format: "(%d)"}}
	result := Tree(tree, seg, sess, segmentValues, defaultFormats, condEnv)
	if result != "(85)" {
		t.Errorf("expected '(85)', got %q", result)
	}
}

func TestTree_EmptyStringCollapses(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.plain": ""}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{{Type: "test.plain"}}
	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "" {
		t.Errorf("expected empty (collapsed), got %q", result)
	}
}

func TestTree_DataSegmentWhenPasses(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.count": 75}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{Type: "test.count", When: "value >= 50"},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "75" {
		t.Errorf("expected '75', got %q", result)
	}
}

func TestTree_DataSegmentWhenFails(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.count": 25}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{Type: "test.count", When: "value >= 50"},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "" {
		t.Errorf("expected empty (when failed), got %q", result)
	}
}

func TestTree_DataSegmentWhenCrossProvider(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"test.name":  "feature",
		"test.count": 5,
	}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{Type: "test.name", When: "test.count > 0"},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "feature" {
		t.Errorf("expected 'feature', got %q", result)
	}
}

func TestTree_DataSegmentWhenText(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": "hello"}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{Type: "test.name", When: "text != ''"},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_CompositeWhen(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"git.branch": "main",
	}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{
			When: "git.branch != ''",
			Children: []types.SegmentNode{
				{Type: "git.branch"},
			},
		},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "main" {
		t.Errorf("expected 'main', got %q", result)
	}
}

func TestTree_CompositeWhenFails(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{
		"git.branch": "",
	}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{
			When: "git.branch != ''",
			Children: []types.SegmentNode{
				{Type: "literal", Props: map[string]any{"text": "should not appear"}},
			},
		},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "" {
		t.Errorf("expected empty (composite when failed), got %q", result)
	}
}

func TestTree_WhenNoExpression(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": "hello"}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{Type: "test.name"},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_WhenInvalidExpression(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	segmentValues := map[string]any{"test.name": "hello"}
	condEnv := condition.BuildNestedEnv(segmentValues)

	tree := []types.SegmentNode{
		{Type: "test.name", When: ">>>bad<<<"},
	}

	result := Tree(tree, seg, sess, segmentValues, map[string]string{}, condEnv)
	if result != "" {
		t.Errorf("expected empty (invalid when), got %q", result)
	}
}

func TestResolveAll(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}
	sess := &types.SessionData{CWD: "/tmp"}

	values, formats := ResolveAll(providers, sess)
	if values["test.name"] != "hello" {
		t.Errorf("expected test.name='hello', got %v", values["test.name"])
	}
	if formats["test.pct"] != "%d%%" {
		t.Errorf("expected test.pct format, got %q", formats["test.pct"])
	}
}

// testProvider implements DataProvider for tests.
type testProvider struct{}

func (p *testProvider) Name() string { return "test" }
func (p *testProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	return &types.ProviderResult{
		Values: map[string]any{
			"test.name":  "hello",
			"test.count": 42,
			"test.pct":   85,
		},
		Formats: map[string]string{
			"test.pct": "%d%%",
		},
	}, nil
}
