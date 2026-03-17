package provider

import (
	"fmt"

	"github.com/jheddings/ccglow/internal/types"
)

type systemProvider struct{}

func (p *systemProvider) Name() string { return "system" }

func (p *systemProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	sys := map[string]any{
		"load": map[string]any{
			"avg1":  0.0,
			"avg5":  0.0,
			"avg15": 0.0,
		},
		"mem": map[string]any{
			"used":    "",
			"total":   "",
			"percent": 0,
		},
		"disk": map[string]any{
			"used":    "",
			"total":   "",
			"percent": 0,
		},
		"battery": map[string]any{
			"percent": 0,
			"state":   "",
		},
		"uptime": "",
	}

	result := &types.ProviderResult{
		Values: map[string]any{"system": sys},
		Formats: map[string]string{
			"system.load.avg1":        "%.2f",
			"system.load.avg5":        "%.2f",
			"system.load.avg15":       "%.2f",
			"system.mem.percent":      "%d%%",
			"system.disk.percent":     "%d%%",
			"system.battery.percent":  "%d%%",
		},
	}

	cwd := session.CWD
	if cwd == "" {
		cwd = "/"
	}

	fillLoadAvg(sys["load"].(map[string]any))
	fillMemory(sys["mem"].(map[string]any))
	fillDisk(sys["disk"].(map[string]any), cwd)
	fillBattery(sys["battery"].(map[string]any))
	fillUptime(sys)

	return result, nil
}

// FormatBytes formats a byte count into a human-readable string (e.g. "1.5G").
func FormatBytes(bytes uint64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
		tb = 1024 * gb
	)

	switch {
	case bytes >= tb:
		v := float64(bytes) / float64(tb)
		if v == float64(uint64(v)) {
			return fmt.Sprintf("%dT", uint64(v))
		}
		return fmt.Sprintf("%.1fT", v)
	case bytes >= gb:
		v := float64(bytes) / float64(gb)
		if v == float64(uint64(v)) {
			return fmt.Sprintf("%dG", uint64(v))
		}
		return fmt.Sprintf("%.1fG", v)
	case bytes >= mb:
		v := float64(bytes) / float64(mb)
		if v == float64(uint64(v)) {
			return fmt.Sprintf("%dM", uint64(v))
		}
		return fmt.Sprintf("%.1fM", v)
	case bytes >= kb:
		v := float64(bytes) / float64(kb)
		if v == float64(uint64(v)) {
			return fmt.Sprintf("%dK", uint64(v))
		}
		return fmt.Sprintf("%.1fK", v)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

// FormatUptime formats a duration in seconds into a human-readable string (e.g. "3d 14h").
func FormatUptime(seconds uint64) string {
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours%24)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes%60)
	}
	return fmt.Sprintf("%dm", minutes)
}
