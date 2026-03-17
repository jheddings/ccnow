package provider

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

const systemTimeout = 5 * time.Second

func sysExec(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), systemTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func fillLoadAvg(m map[string]any) {
	// sysctl -n vm.loadavg returns "{ 1.23 4.56 7.89 }"
	out, err := sysExec("sysctl", "-n", "vm.loadavg")
	if err != nil {
		return
	}

	out = strings.Trim(out, "{ }")
	fields := strings.Fields(out)
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
	// Total memory
	totalStr, err := sysExec("sysctl", "-n", "hw.memsize")
	if err != nil {
		return
	}
	total, err := strconv.ParseUint(strings.TrimSpace(totalStr), 10, 64)
	if err != nil {
		return
	}

	// Used memory from vm_stat
	out, err := sysExec("vm_stat")
	if err != nil {
		return
	}

	pageSize := uint64(unix.Getpagesize())
	pages := parseVMStat(out)

	active := pages["Pages active"]
	wired := pages["Pages wired down"]
	speculative := pages["Pages speculative"]
	compressed := pages["Pages occupied by compressor"]

	used := (active + wired + speculative + compressed) * pageSize

	m["total"] = FormatBytes(total)
	m["used"] = FormatBytes(used)
	if total > 0 {
		m["percent"] = int(used * 100 / total)
	}
}

func parseVMStat(output string) map[string]uint64 {
	result := make(map[string]uint64)
	for _, line := range strings.Split(output, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		valStr := strings.TrimSpace(parts[1])
		valStr = strings.TrimSuffix(valStr, ".")
		if v, err := strconv.ParseUint(valStr, 10, 64); err == nil {
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
	out, err := sysExec("ioreg", "-rc", "AppleSmartBattery")
	if err != nil || out == "" {
		return
	}

	props := parseIORegProperties(out)

	maxCap, hasMax := props["MaxCapacity"]
	curCap, hasCur := props["CurrentCapacity"]
	if hasMax && hasCur {
		max, errMax := strconv.Atoi(maxCap)
		cur, errCur := strconv.Atoi(curCap)
		if errMax == nil && errCur == nil && max > 0 {
			m["percent"] = int(cur * 100 / max)
		}
	}

	if props["IsCharging"] == "Yes" {
		m["state"] = "charging"
	} else if props["FullyCharged"] == "Yes" {
		m["state"] = "full"
	} else if props["ExternalConnected"] == "No" {
		m["state"] = "discharging"
	}
}

func parseIORegProperties(output string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "\"") {
			continue
		}
		// Lines look like: "Key" = Value
		line = strings.TrimPrefix(line, "\"")
		parts := strings.SplitN(line, "\" = ", 2)
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = parts[1]
	}
	return result
}

func fillUptime(sys map[string]any) {
	// sysctl -n kern.boottime returns "{ sec = 1234567890, usec = 0 }"
	out, err := sysExec("sysctl", "-n", "kern.boottime")
	if err != nil {
		return
	}

	var sec int64
	_, err = fmt.Sscanf(out, "{ sec = %d,", &sec)
	if err != nil || sec <= 0 {
		return
	}

	uptime := uint64(time.Now().Unix() - sec)
	sys["uptime"] = FormatUptime(uptime)
}
