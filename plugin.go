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
	"./lib"
)

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
			d := lib.GetServerStat()
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
