package main

import (
	"flag"
	"fmt"
	"github.com/Prouk/LGTVShutDowner/pkg"
	"golang.org/x/net/websocket"
	"log"
)

var cmd string

func main() {
	flag.StringVar(&cmd, "c", "PowerOn", "The command to send to send to the TV.")
	flag.Parse()
	lsd, err := pkg.NewLsd()
	if err != nil {
		log.Fatal(err)
	}
	err = lsd.ConnectTv()
	if err != nil {
		log.Fatal(err)
	}
	defer lsd.WsConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = lsd.SendCommand(cmd)
	if err != nil {
		log.Fatal(err)
	}
	for {
		if !lsd.WsConn.IsClientConn() {
			break
		}
		var msg map[string]interface{}
		err := websocket.JSON.Receive(lsd.WsConn, msg)
		if err != nil {
			fmt.Printf("Error receiving ws message: %s\n", err)
		} else {
			fmt.Printf("Message received: %s\n", msg)
		}
	}
	fmt.Printf("Service Shutting down")
}
