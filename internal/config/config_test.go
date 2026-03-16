package config

import "testing"

func TestParse_Valid(t *testing.T) {
	input := `{
		"segments": [
			{"segment": "pwd.name", "style": {"color": "red"}},
			{"segment": "git.branch"}
		]
	}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Type != "pwd.name" {
		t.Errorf("expected pwd.name, got %s", nodes[0].Type)
	}
}

func TestParse_WithChildren(t *testing.T) {
	input := `{
		"segments": [
			{
				"segment": "group",
				"children": [
					{"segment": "git.branch"},
					{"segment": "git.insertions"}
				]
			}
		]
	}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if len(nodes[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(nodes[0].Children))
	}
}

func TestParse_LiteralNoProvider(t *testing.T) {
	input := `{"segments": [{"segment": "literal", "props": {"text": "hi"}}]}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParse_WithFormat(t *testing.T) {
	input := `{"segments": [{"segment": "context.percent", "format": "%d%%"}]}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if nodes[0].Format != "%d%%" {
		t.Errorf("expected format, got %q", nodes[0].Format)
	}
}
