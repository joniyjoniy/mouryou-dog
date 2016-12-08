package lib

import (
	"fmt"
	"strings"
	"time"
	"os/exec"

	"github.com/shirou/gopsutil/mem"
	// "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
)

type ServerStat struct {
	// Host
	HostName             string   `json:"hostname"`
	HostID               string   `json:"hostid"`
	VirtualizationSystem string   `json:"virtualizationSystem"`
  // mem.VirtualMemoryStat
	Total                uint64   `json:"total"`
	Available            uint64   `json:"available"`
	UsedPercent          float64  `json:"usedPercent"`
	// DiskIO map[string]disk.IOCountersStat
	DiskIO              []DiskStat `json:"diskIO"`
	// Time
	Time                 string   `json:"time"`

	// Cpu
	// Cpu    []cpu.TimesStat         `json:"-"`

	ApacheStat float64 `json:"apacheStat"`
}

type DiskStat struct {
	Name       string `json:"name"`
  IoTime     uint64 `json:"ioTime"`
	WeightedIO uint64 `json:"weightedIO"`
}

func GetServerStat() (ServerStat) {
	var d ServerStat

	d.GetHostStat()
  d.GetMemoryStat()
	d.GetDiskIOStat()
	d.GetTime()
	d.GetApacheStat()
	return d
}

func (s *ServerStat) GetHostStat() {
	h, _ := host.Info()
	s.HostName             = h.Hostname
	s.HostID               = h.HostID
	s.VirtualizationSystem = h.VirtualizationSystem
}

func (s *ServerStat) GetMemoryStat() {
	m, _ := mem.VirtualMemory()
	s.Total = m.Total
	s.Available = m.Available
	s.UsedPercent = m.UsedPercent
}

func (s *ServerStat) GetDiskIOStat() {
	var ds []DiskStat
	i, _ := disk.IOCounters()
	for k, v := range i {
		var d DiskStat
		d.Name       = k
		d.IoTime     = v.IoTime
		d.WeightedIO = v.WeightedIO
		ds = append(ds, d)
	}
	s.DiskIO = ds
}

func (s *ServerStat)GetApacheStat() {
	var dataLine int
	out, _ := exec.Command("apachectl", "status").Output()
	d :=string(out)

	lines := strings.Split(strings.TrimRight(d, "\n"), "\n")

	for k, v := range lines {
		if v == "Scoreboard Key:" {
			dataLine = k
			break
		}
	}

	board := lines[dataLine-4]
	board = board + lines[dataLine-3]
	board = board + lines[dataLine-2]
	all := len(strings.Split(board, ""))
	idles := strings.Count(board, "_") + strings.Count(board, ".")

	r := float64((all - idles)) / float64(all)

	s.ApacheStat = r
}

func (s *ServerStat)GetTime() {
	now := time.Now()
	s.Time =fmt.Sprint(now)
}
