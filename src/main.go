package main

import (
	"flag"
	"github.com/Prouk/LGTVShutDowner/pkg"
	"log"
)

var cmd string

func main() {
	flag.StringVar(&cmd, "c", "PowerOn", "The command to send to send to the TV (default : 'PowerOn')")
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
}
