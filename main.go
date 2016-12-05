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
	
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/mem"
	// "github.com/shirou/gopsutil/cpu"
	// "github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
)

type ServerStat struct {
	// Host
	HostName             string `json:"hostname"`
	HostID               string `json:"hostid"`
	VirtualizationSystem string `json:"virtualizationSystem"`
  // mem.VirtualMemoryStat
	Total        uint64 `json:"total"`
	Available    uint64 `json:"available"`
	Used         uint64 `json:"used"`
	UsedPercent  float64 `json:"usedPercent"`
	Free         uint64 `json:"free"`
	Active       uint64 `json:"active"`
	Inactive     uint64 `json:"inactive"`
	Wired        uint64 `json:"wired"`
	Buffers      uint64 `json:"buffers"`
	Cached       uint64 `json:"cached"`
	Writeback    uint64 `json:"writeback"`
	Dirty        uint64 `json:"dirty"`
	WritebackTmp uint64 `json:"writebacktmp"`
	// Cpu
	// DiskIO
	ReadCount        uint64 `json:"readCount"`
	MergedReadCount  uint64 `json:"mergedReadCount"`
	WriteCount       uint64 `json:"writeCount"`
	MergedWriteCount uint64 `json:"mergedWriteCount"`
	ReadBytes        uint64 `json:"readBytes"`
	WriteBytes       uint64 `json:"writeBytes"`
	ReadTime         uint64 `json:"readTime"`
	WriteTime        uint64 `json:"writeTime"`
	IopsInProgress   uint64 `json:"iopsInProgress"`
	IoTime           uint64 `json:"ioTime"`
	WeightedIO       uint64 `json:"weightedIO"`
	Name             string `json:"name"`
	SerialNumber     string `json:"serialNumber"`
	
	// Memory *mem.VirtualMemoryStat  `json:"-"`
	// Cpu    []cpu.TimesStat         `json:"-"`
	// DiskIO map[string]disk.IOCountersStat    `json:"-"`
	// Host   host.InfoStat  `json:"-"`
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
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
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
	// d.GetDiskIOStat()
	
	return d
}

func (s ServerStat) GetHostStat() (ServerStat) {
	h, _ := host.Info()
	s.HostName = h.Hostname
	return s
}

func (s ServerStat) GetMemoryStat() (ServerStat){
	m, _ := mem.VirtualMemory()
	s.Total = m.Total
	return s
}

// func (s ServerStat) GetDiskIOStat() {
// 	i, _ := disk.IOCounters()
// 	s.ReadCount = i.ReadCount
// }
