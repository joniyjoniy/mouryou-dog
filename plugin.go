package main

import (
	"encoding/json"
	"bytes"
	"flag"
	"log"
	"time"
	"os"
	"os/signal"
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
	HostName             string  `json:"hostname"`
	HostID               string  `json:"hostid"`
	VirtualizationSystem string  `json:"virtualizationSystem"`
  // mem.VirtualMemoryStat
	Total                uint64  `json:"total"`
	Available            uint64  `json:"available"`
	UsedPercent          float64 `json:"usedPercent"`
	// DiskIO map[string]disk.IOCountersStat
	DiskIO               string  `json:"diskIO"`
	// IoTime            uint64  `json:"ioTime"`
	// WeightedIO        uint64  `json:"weightedIO"`
	// Time
	Time                 string  `json:"time"`

	// Cpu	
	// Cpu    []cpu.TimesStat         `json:"-"`

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

	d = d.GetHostStat()
  d = d.GetMemoryStat()
	d = d.GetDiskIOStat()
	d = d.GetTime()
	return d
}

func (s ServerStat) GetHostStat() (ServerStat) {
	h, _ := host.Info()
	s.HostName             = h.Hostname
	s.HostID               = h.HostID
	s.VirtualizationSystem = h.VirtualizationSystem
	return s
}

func (s ServerStat) GetMemoryStat() (ServerStat) {
	m, _ := mem.VirtualMemory()
	s.Total = m.Total
	s.Available = m.Available
	s.UsedPercent = m.UsedPercent
	return s
}

func (s ServerStat) GetDiskIOStat() (ServerStat) {
	i, _ := disk.IOCounters()
	s.DiskIO = ConvertMapToString(i)
	return s
}

func ConvertMapToString(m map[string]disk.IOCountersStat) (string) {
	var str string

	str = "{"
	for k, v := range m {
		str  = str + string(k) + ":{"
		str = str + "ioTime:" + fmt.Sprint(v.IoTime) + ","
		str = str + "weightedIO:" + fmt.Sprint(v.WeightedIO) + "},"
	}
  str = strings.TrimRight(str, ",")
	str = str + "}"
	return str
}

func (s ServerStat)GetTime() (ServerStat) {
	now := time.Now()
	s.Time =fmt.Sprint(now)
	return s
}
