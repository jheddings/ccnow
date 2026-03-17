package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func fillLoadAvg(m map[string]any) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return
	}

	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return
	}

	if v, err := strconv.ParseFloat(fields[0], 64); err == nil {
		m["avg1"] = v
	}
	if v, err := strconv.ParseFloat(fields[1], 64); err == nil {
		m["avg5"] = v
	}
	if v, err := strconv.ParseFloat(fields[2], 64); err == nil {
		m["avg15"] = v
	}
}

func fillMemory(m map[string]any) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	info := parseMeminfo(string(data))

	total := info["MemTotal"]
	available := info["MemAvailable"]
	if available == 0 {
		// Fallback for kernels < 3.14
		available = info["MemFree"] + info["Buffers"] + info["Cached"]
	}

	if total == 0 {
		return
	}

	// /proc/meminfo reports in kB
	totalBytes := total * 1024
	used := total - available
	usedBytes := used * 1024

	m["total"] = FormatBytes(totalBytes)
	m["used"] = FormatBytes(usedBytes)
	m["percent"] = int(used * 100 / total)
}

func parseMeminfo(data string) map[string]uint64 {
	result := make(map[string]uint64)
	for _, line := range strings.Split(data, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		valStr := strings.TrimSpace(parts[1])
		valStr = strings.TrimSuffix(valStr, " kB")
		if v, err := strconv.ParseUint(strings.TrimSpace(valStr), 10, 64); err == nil {
			result[key] = v
		}
	}
	return result
}

func fillDisk(m map[string]any, path string) {
	var stat unix.Statfs_t
	if err := unix.Statfs(path, &stat); err != nil {
		return
	}

	bsize := uint64(stat.Bsize)
	total := stat.Blocks * bsize
	free := stat.Bavail * bsize
	used := total - free

	m["total"] = FormatBytes(total)
	m["used"] = FormatBytes(used)
	if total > 0 {
		m["percent"] = int(used * 100 / total)
	}
}

func fillBattery(m map[string]any) {
	matches, err := filepath.Glob("/sys/class/power_supply/BAT*")
	if err != nil || len(matches) == 0 {
		return
	}

	bat := matches[0]

	capData, err := os.ReadFile(filepath.Join(bat, "capacity"))
	if err == nil {
		if v, err := strconv.Atoi(strings.TrimSpace(string(capData))); err == nil {
			m["percent"] = v
		}
	}

	statusData, err := os.ReadFile(filepath.Join(bat, "status"))
	if err == nil {
		status := strings.TrimSpace(string(statusData))
		switch strings.ToLower(status) {
		case "charging":
			m["state"] = "charging"
		case "discharging":
			m["state"] = "discharging"
		case "full":
			m["state"] = "full"
		case "not charging":
			m["state"] = "full"
		}
	}
}

func fillUptime(sys map[string]any) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return
	}

	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return
	}

	var seconds float64
	if _, err := fmt.Sscanf(fields[0], "%f", &seconds); err != nil {
		return
	}

	sys["uptime"] = FormatUptime(uint64(seconds))
}
