package main

import (
	"github.com/Prouk/LGTVShutDowner/pkg"
	"log"
	"time"
)

func main() {
	log.Printf("lsd starting")
	l := pkg.NewLSDManager()
	l.LSDChan <- 1
	time.Sleep(5 * time.Second)
	l.LSDChan <- 0
	log.Printf("lsd exiting")
}
