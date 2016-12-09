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

func GetServerStat() (ServerStat, error) {
	var d ServerStat

  err := d.GetHostStat()
	if err != nil {
		return d, err
	}
  d.GetMemoryStat()
	d.GetDiskIOStat()
	d.GetTime()
	d.GetApacheStat()
	return d, nil
}

func (s *ServerStat) GetHostStat() (error) {
  h, err := host.Info()
	if err != nil {
		return err
	}
	s.HostName             = h.Hostname
	s.HostID               = h.HostID
	s.VirtualizationSystem = h.VirtualizationSystem
	return nil
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
