package lib

import (
	"fmt"
	"strings"
	"time"
	"os/exec"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
)

func GetServerStat() (ServerStat, []error) {
	var d ServerStat
	var err []error

	errHost := d.GetHostStat()
	if errHost != nil {
		err = append(err, errHost)
	}
	errMemory := d.GetMemoryStat()
	if errMemory != nil {
		err = append(err, errMemory)
	}
	errDisk := d.GetDiskIOStat()
	if errDisk != nil {
		err = append(err, errDisk)
	}
	errApache := d.GetApacheStat()
	if errApache != nil {
		err = append(err, errApache)
	}
	errCpu := d.GetCpuStat()
	if errCpu != nil {
		err = append(err, errCpu)
	}
	d.GetTime()

	if err != nil {
		return d, err
	}

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

func (s *ServerStat) GetMemoryStat() (error) {
	m, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	s.Total = m.Total
	s.Available = m.Available
	s.UsedPercent = m.UsedPercent
	return nil
}

func (s *ServerStat) GetDiskIOStat() (error) {
	var ds []DiskStat
	i, err := disk.IOCounters()
	if err != nil {
		return err
	}
	for k, v := range i {
		var d DiskStat
		d.Name       = k
		d.IoTime     = v.IoTime
		d.WeightedIO = v.WeightedIO
		ds = append(ds, d)
	}
	s.DiskIO = ds
	return nil
}

func (s *ServerStat)GetApacheStat() (error) {
	var dataLine int
	out, err := exec.Command("apachectl", "status").Output()
	if err != nil {
		return err
	}

	d := string(out)
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
	return nil
}

func (s *ServerStat)GetTime() {
	now := time.Now()
	s.Time =fmt.Sprint(now)
}

func (s *ServerStat)GetCpuStat() (error) {
	c, err := cpu.Times(true)
	if err != nil {
		return err
	}
  s.Cpu = c
	return nil
}
