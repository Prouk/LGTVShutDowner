package main

import (
	"flag"
	"github.com/Prouk/LGTVShutDowner/pkg"
	"log"
)

var cmd string
var cfg string

func main() {
	flag.StringVar(&cmd, "cmd", "", "Command to launch at service start time")
	flag.StringVar(&cfg, "cfg", "", "Path to config file")
	flag.Parse()
	lsd := pkg.CreateLsd(cmd, cfg)
	<-lsd.ExitChann
	log.Printf("stopping lsd service\n")
}
