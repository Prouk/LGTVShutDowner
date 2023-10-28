package main

import (
	"github.com/Prouk/LGTVShutDowner/pkg"
	"log"
)

func main() {
	lsd := pkg.CreateLsd()
	<-lsd.ExitChann
	log.Printf("stopping lsd service\n")
}
