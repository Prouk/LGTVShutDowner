package main

import (
	"flag"
	"fmt"
	"github.com/Prouk/LGTVShutDowner/pkg"
	"io"
	"log"
	"time"
)

var cmd string

func main() {
	flag.StringVar(&cmd, "c", "none", "The command to send to send to the TV.")
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
	if cmd != "none" {
		err = lsd.SendCommand(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
	for {
		if !lsd.WsConn.IsClientConn() {
			break
		}
		buf := make([]byte, 0, 4096) // big buffer
		tmp := make([]byte, 256)
		for {
			n, err := lsd.WsConn.Read(tmp)
			if err != nil {
				if err != io.EOF {
					fmt.Println("read error:", err)
				}
				break
			}
			//fmt.Println("got", n, "bytes.")
			buf = append(buf, tmp[:n]...)

		}
		if len(buf) > 0 {
			fmt.Printf("received message: %s\n", string(buf))
		}
		time.Sleep(time.Millisecond * 1000)
	}
	fmt.Printf("service Shutting down")
}
