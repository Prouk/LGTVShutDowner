package pkg

import (
	"github.com/getlantern/systray"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Lsd struct {
	TrayMenu  *TrayMenu
	SigChann  chan os.Signal
	ExitChann chan bool
}

func CreateLsd() *Lsd {
	lsd := new(Lsd)
	lsd.SigChann = make(chan os.Signal)
	lsd.ExitChann = make(chan bool)
	signal.Notify(lsd.SigChann)
	go lsd.ListenSig()
	systray.Run(lsd.InitTray, lsd.Close)
	return lsd
}

func (lsd *Lsd) ListenSig() {
	for {
		sig := <-lsd.SigChann
		log.Printf("%s\n", sig)
		lsd.HandleSig(sig)
	}
}

func (lsd *Lsd) HandleSig(sig os.Signal) {
	switch sig {
	case syscall.SIGTERM:
		lsd.ExitChann <- true
	case syscall.SIGINT:
		lsd.ExitChann <- true
	default:
		return
	}
}

func (lsd *Lsd) GetIcon() []byte {
	data, err := os.ReadFile("./assets/screen.png")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func (lsd *Lsd) Close() {
	log.Printf("stopping tray")
}

func (lsd *Lsd) PingScreen() {
	log.Printf("sending ping to screen")
}
