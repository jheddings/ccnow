package config

import "testing"

func TestParse_Valid(t *testing.T) {
	input := `{
		"segments": [
			{"expr": "pwd.name", "style": {"color": "red"}},
			{"expr": "git.branch"}
		]
	}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Expr != "pwd.name" {
		t.Errorf("expected pwd.name, got %s", nodes[0].Expr)
	}
}

func TestParse_WithChildren(t *testing.T) {
	input := `{
		"segments": [
			{
				"children": [
					{"expr": "git.branch"},
					{"expr": "git.insertions"}
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

func TestParse_ValueNode(t *testing.T) {
	input := `{"segments": [{"value": "|", "style": {"color": "240"}}]}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if nodes[0].Value != "|" {
		t.Errorf("expected value '|', got %v", nodes[0].Value)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParse_WithFormat(t *testing.T) {
	input := `{"segments": [{"expr": "context.percent.used", "format": "%d%%"}]}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if nodes[0].Format != "%d%%" {
		t.Errorf("expected format, got %q", nodes[0].Format)
	}
}

func TestParse_SkipsEmptyNodes(t *testing.T) {
	input := `{"segments": [{"style": {"color": "red"}}, {"expr": "pwd.name"}]}`

	nodes, err := Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node (empty skipped), got %d", len(nodes))
	}
}
