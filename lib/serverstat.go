package lib

import (
	"encoding/json"
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

func (d ServerStat) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func (d DiskStat) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}
