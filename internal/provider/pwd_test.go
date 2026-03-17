package provider

import (
	"os"
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestPwdProvider(t *testing.T) {
	p := &pwdProvider{}
	sess := &types.SessionData{CWD: "/home/user/projects/myapp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	pwd := result.Values["pwd"].(map[string]any)
	if pwd["name"] != "myapp" {
		t.Errorf("expected myapp, got %s", pwd["name"])
	}
	if pwd["path"] != "/home/user/projects/" {
		t.Errorf("expected /home/user/projects/, got %s", pwd["path"])
	}
}

func TestSmartPrefix(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		cwd      string
		expected string
	}{
		// Root and top-level
		{"/", ""},
		{"/tmp", ""},
		{"/usr", ""},

		// Absolute paths (not under home)
		{"/usr/local", "/usr/"},
		{"/usr/local/bin", "/usr/local/"},
		{"/var/log/syslog", "/var/log/"},

		// Home directory itself
		{home, ""},

		// First level under home (the bug case -- was producing "~//")
		{home + "/Projects", "~/"},

		// Two levels under home
		{home + "/Projects/myapp", "~/Projects/"},

		// Three levels under home
		{home + "/Projects/myapp/src", "~/Projects/myapp/"},

		// Four levels under home (abbreviation kicks in)
		{home + "/Projects/myapp/src/pkg", "~/P/m/…/"},
	}

	for _, tt := range tests {
		result := smartPrefix(tt.cwd)
		if result != tt.expected {
			t.Errorf("smartPrefix(%q) = %q, want %q", tt.cwd, result, tt.expected)
		}
	}
}
