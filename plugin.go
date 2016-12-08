package main

import (
	"encoding/json"
	"bytes"
	"flag"
	"log"
	"time"
	"os"
	"os/signal"
	"os/exec"
	"net/url"
	"strings"
	"fmt"

	"github.com/gorilla/websocket"
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

var addr = flag.String("addr", "localhost:8080", "monitoring address")

func main() {
	var buf bytes.Buffer
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d := GetServerStat()
			j, _ :=json.Marshal(d)
			buf.Write(j)
			err := c.WriteMessage(websocket.TextMessage, j)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
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
