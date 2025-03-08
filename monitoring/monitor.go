package monitoring

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type SystemStats struct {
	// CPU
	CPUUsage float64

	// Memory
	MemoryTotal     uint64
	MemoryUsed      uint64
	MemoryFree      uint64
	MemoryUsagePerc float64

	// Network
	BytesSent    uint64
	BytesRecv    uint64
	PacketsSent  uint64
	PacketsRecv  uint64
	NetworkSpeed float64 // bytes/sec

	// Go Runtime
	NumGoroutines int
	NumCPU        int

	// Process specific
	ProcessCPU    float64
	ProcessMemory uint64

	Timestamp time.Time
}

type SystemMonitor struct {
	stats      *SystemStats
	process    *process.Process
	lastUpdate time.Time
	lastNet    net.IOCountersStat
}

func NewSystemMonitor() (*SystemMonitor, error) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}

	return &SystemMonitor{
		stats:   &SystemStats{},
		process: proc,
	}, nil
}

func (sm *SystemMonitor) Collect() (*SystemStats, error) {
	now := time.Now()

	// CPU Usage (total system)
	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		sm.stats.CPUUsage = cpuPercent[0]
	}

	// Memory stats
	if vmstat, err := mem.VirtualMemory(); err == nil {
		sm.stats.MemoryTotal = vmstat.Total
		sm.stats.MemoryUsed = vmstat.Used
		sm.stats.MemoryFree = vmstat.Free
		sm.stats.MemoryUsagePerc = vmstat.UsedPercent
	}

	// Network stats
	if netStats, err := net.IOCounters(false); err == nil && len(netStats) > 0 {
		current := netStats[0]

		// Calculate network speed if we have previous measurements
		if !sm.lastUpdate.IsZero() {
			duration := now.Sub(sm.lastUpdate).Seconds()
			if duration > 0 {
				bytesDelta := current.BytesSent + current.BytesRecv -
					(sm.lastNet.BytesSent + sm.lastNet.BytesRecv)
				sm.stats.NetworkSpeed = float64(bytesDelta) / duration
			}
		}

		sm.stats.BytesSent = current.BytesSent
		sm.stats.BytesRecv = current.BytesRecv
		sm.stats.PacketsSent = current.PacketsSent
		sm.stats.PacketsRecv = current.PacketsRecv
		sm.lastNet = current
	}

	// Process specific stats
	if cpuPercent, err := sm.process.CPUPercent(); err == nil {
		sm.stats.ProcessCPU = cpuPercent
	}
	if memInfo, err := sm.process.MemoryInfo(); err == nil {
		sm.stats.ProcessMemory = memInfo.RSS
	}

	// Go runtime stats
	sm.stats.NumGoroutines = runtime.NumGoroutine()
	sm.stats.NumCPU = runtime.NumCPU()

	sm.stats.Timestamp = now
	sm.lastUpdate = now

	return sm.stats, nil
}

func (s *SystemStats) String() string {
	return fmt.Sprintf(
		"System Stats:\n"+
			"CPU Usage: %.2f%%\n"+
			"Memory Usage: %.2f%% (Used: %d MB, Free: %d MB)\n"+
			"Network: ↑%s ↓%s (Speed: %s/s)\n"+
			"Goroutines: %d\n"+
			"Process CPU: %.2f%%\n"+
			"Process Memory: %d MB",
		s.CPUUsage,
		s.MemoryUsagePerc,
		s.MemoryUsed/(1024*1024),
		s.MemoryFree/(1024*1024),
		formatBytes(s.BytesSent),
		formatBytes(s.BytesRecv),
		formatBytes(uint64(s.NetworkSpeed)),
		s.NumGoroutines,
		s.ProcessCPU,
		s.ProcessMemory/(1024*1024),
	)
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
