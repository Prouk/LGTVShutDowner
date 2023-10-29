package main

import (
	"flag"
	"github.com/Prouk/LGTVShutDowner/pkg"
	"log"
)

var cmd string

func main() {
	flag.StringVar(&cmd, "c", "", "Command to launch at service start time")
	flag.Parse()
	lsd := pkg.CreateLsd(cmd)
	<-lsd.ExitChann
	log.Printf("stopping lsd service\n")
}
