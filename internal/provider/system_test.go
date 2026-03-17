package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0B"},
		{512, "512B"},
		{1023, "1023B"},
		{1024, "1K"},
		{1536, "1.5K"},
		{10240, "10K"},
		{1048576, "1M"},
		{1572864, "1.5M"},
		{1073741824, "1G"},
		{1610612736, "1.5G"},
		{34359738368, "32G"},
		{1099511627776, "1T"},
		{1649267441664, "1.5T"},
	}

	for _, tt := range tests {
		result := FormatBytes(tt.input)
		if result != tt.expected {
			t.Errorf("FormatBytes(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		seconds  uint64
		expected string
	}{
		{0, "0m"},
		{30, "0m"},
		{60, "1m"},
		{90, "1m"},
		{3600, "1h 0m"},
		{3660, "1h 1m"},
		{7200, "2h 0m"},
		{86400, "1d 0h"},
		{90000, "1d 1h"},
		{259200, "3d 0h"},
		{307800, "3d 13h"},
	}

	for _, tt := range tests {
		result := FormatUptime(tt.seconds)
		if result != tt.expected {
			t.Errorf("FormatUptime(%d) = %q, want %q", tt.seconds, result, tt.expected)
		}
	}
}

func systemValues(result *types.ProviderResult) map[string]any {
	return result.Values["system"].(map[string]any)
}

func TestSystemProvider_ReturnsAllKeys(t *testing.T) {
	p := &systemProvider{}
	sess := &types.SessionData{CWD: "/"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	sys := systemValues(result)

	// load
	load, ok := sys["load"].(map[string]any)
	if !ok {
		t.Fatal("expected system.load map")
	}
	for _, key := range []string{"avg1", "avg5", "avg15"} {
		if _, ok := load[key].(float64); !ok {
			t.Errorf("expected system.load.%s to be float64, got %T", key, load[key])
		}
	}

	// mem
	mem, ok := sys["mem"].(map[string]any)
	if !ok {
		t.Fatal("expected system.mem map")
	}
	if _, ok := mem["used"].(string); !ok {
		t.Errorf("expected system.mem.used to be string, got %T", mem["used"])
	}
	if _, ok := mem["total"].(string); !ok {
		t.Errorf("expected system.mem.total to be string, got %T", mem["total"])
	}
	if _, ok := mem["percent"].(int); !ok {
		t.Errorf("expected system.mem.percent to be int, got %T", mem["percent"])
	}

	// disk
	disk, ok := sys["disk"].(map[string]any)
	if !ok {
		t.Fatal("expected system.disk map")
	}
	if _, ok := disk["used"].(string); !ok {
		t.Errorf("expected system.disk.used to be string, got %T", disk["used"])
	}
	if _, ok := disk["total"].(string); !ok {
		t.Errorf("expected system.disk.total to be string, got %T", disk["total"])
	}
	if _, ok := disk["percent"].(int); !ok {
		t.Errorf("expected system.disk.percent to be int, got %T", disk["percent"])
	}

	// disk should have real values for /
	if disk["total"] == "" {
		t.Error("expected non-empty system.disk.total for /")
	}

	// battery (may be zero on desktops, just check types)
	battery, ok := sys["battery"].(map[string]any)
	if !ok {
		t.Fatal("expected system.battery map")
	}
	if _, ok := battery["percent"].(int); !ok {
		t.Errorf("expected system.battery.percent to be int, got %T", battery["percent"])
	}
	if _, ok := battery["state"].(string); !ok {
		t.Errorf("expected system.battery.state to be string, got %T", battery["state"])
	}

	// uptime
	if _, ok := sys["uptime"].(string); !ok {
		t.Errorf("expected system.uptime to be string, got %T", sys["uptime"])
	}
	if sys["uptime"] == "" {
		t.Error("expected non-empty system.uptime")
	}
}

func TestSystemProvider_EmptyCWD(t *testing.T) {
	p := &systemProvider{}
	sess := &types.SessionData{CWD: ""}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	sys := systemValues(result)
	disk := sys["disk"].(map[string]any)

	// should fall back to / and still return values
	if disk["total"] == "" {
		t.Error("expected non-empty disk.total with empty CWD")
	}
}
