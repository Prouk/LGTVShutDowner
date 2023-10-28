package main

import (
	"github.com/Prouk/LGTVShutDowner/pkg"
	"log"
)

func main() {
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
	err = lsd.SendCommand("PowerOff")
	if err != nil {
		log.Fatal(err)
	}
}
