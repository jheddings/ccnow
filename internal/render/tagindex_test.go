package render

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

type testData struct {
	Name  *string `segment:"test.name"`
	Count *int    `segment:"test.count"`
	Plain string  `segment:"test.plain"`
	Skip  string
}

type testProvider struct{}

func (p *testProvider) Name() string                                    { return "test" }
func (p *testProvider) Resolve(session *types.SessionData) (any, error) { return &testData{}, nil }
func (p *testProvider) Fields() any                                     { return &testData{} }

type fmtData struct {
	Pct *int `segment:"fmt.pct,format:%d%%"`
}

type fmtProvider struct{}

func (p *fmtProvider) Name() string                                    { return "fmt" }
func (p *fmtProvider) Resolve(session *types.SessionData) (any, error) { return &fmtData{}, nil }
func (p *fmtProvider) Fields() any                                     { return &fmtData{} }

type plainProvider struct{}

func (p *plainProvider) Name() string                                    { return "plain" }
func (p *plainProvider) Resolve(session *types.SessionData) (any, error) { return nil, nil }

func TestBuildTagIndex(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}
	idx, err := BuildTagIndex(providers)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := idx["test.name"]; !ok {
		t.Error("expected test.name in index")
	}
	if _, ok := idx["test.count"]; !ok {
		t.Error("expected test.count in index")
	}
	if _, ok := idx["test.plain"]; !ok {
		t.Error("expected test.plain in index")
	}
	if _, ok := idx["test.skip"]; ok {
		t.Error("expected test.skip NOT in index")
	}
	if idx["test.name"].Provider != "test" {
		t.Errorf("expected provider 'test', got %q", idx["test.name"].Provider)
	}
}

func TestBuildTagIndex_DefaultFormat(t *testing.T) {
	providers := map[string]types.DataProvider{
		"fmt": &fmtProvider{},
	}
	idx, err := BuildTagIndex(providers)
	if err != nil {
		t.Fatal(err)
	}
	if idx["fmt.pct"].DefaultFormat != "%d%%" {
		t.Errorf("expected default format, got %q", idx["fmt.pct"].DefaultFormat)
	}
}

func TestBuildTagIndex_SkipsNonFieldProvider(t *testing.T) {
	providers := map[string]types.DataProvider{
		"plain": &plainProvider{},
	}
	idx, err := BuildTagIndex(providers)
	if err != nil {
		t.Fatal(err)
	}
	if len(idx) != 0 {
		t.Errorf("expected empty index, got %d entries", len(idx))
	}
}

func TestBuildTagIndex_DuplicateErrors(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test":  &testProvider{},
		"test2": &testProvider{},
	}
	_, err := BuildTagIndex(providers)
	if err == nil {
		t.Error("expected error for duplicate segment names")
	}
}

func TestResolveSegmentValues(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}
	idx, _ := BuildTagIndex(providers)

	name := "hello"
	count := 42
	providerData := map[string]any{
		"test": &testData{Name: &name, Count: &count, Plain: "raw"},
	}
	values := ResolveSegmentValues(idx, providerData)

	if values["test.name"] != "hello" {
		t.Errorf("expected 'hello', got %v", values["test.name"])
	}
	if values["test.count"] != 42 {
		t.Errorf("expected 42, got %v", values["test.count"])
	}
	if values["test.plain"] != "raw" {
		t.Errorf("expected 'raw', got %v", values["test.plain"])
	}
}

func TestResolveSegmentValues_NilPointer(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}
	idx, _ := BuildTagIndex(providers)
	providerData := map[string]any{
		"test": &testData{},
	}
	values := ResolveSegmentValues(idx, providerData)

	if values["test.name"] != nil {
		t.Errorf("expected nil for nil *string, got %v", values["test.name"])
	}
	if values["test.count"] != nil {
		t.Errorf("expected nil for nil *int, got %v", values["test.count"])
	}
}

func TestResolveSegmentValues_MissingProvider(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}
	idx, _ := BuildTagIndex(providers)
	values := ResolveSegmentValues(idx, map[string]any{})

	if _, ok := values["test.name"]; ok {
		t.Error("expected test.name NOT in values when provider data missing")
	}
}
