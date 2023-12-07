package pkg

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type WindowManager struct {
	W           *gtk.Window
	WMChan      chan int
	PairedLabel *gtk.Label
}

func (l *LSDM) NewWindowManger() *WindowManager {
	wm := &WindowManager{}
	wm.WMChan = make(chan int)
	wm.CreateWindow()
	go wm.HandleChan()
	return wm
}

func (WM *WindowManager) HandleChan() {
	for {
		select {
		case i := <-WM.WMChan:
			switch i {
			case 0:
				WM.CloseWindow()
			case 1:
				WM.OpenWindow()
			}
		}
	}
}

func (WM *WindowManager) CreateWindow() {
	log.Printf("creating window")
	var err error
	gtk.Init(nil)
	WM.W, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("unable to create window:", err)
	}
	WM.W.Connect("destroy", func() {
		gtk.MainQuit()
	})

	WM.SetTitle("LSD")
	WM.W.SetDefaultSize(600, 800)
	WM.W.ShowAll()
}

func (WM *WindowManager) OpenWindow() {
	log.Printf("opening window")
	gtk.Main()
}

func (WM *WindowManager) CloseWindow() {
	log.Printf("closing window")
	gtk.MainQuit()
}

func (WM *WindowManager) SetTitle(t string) {
	log.Printf("changing window title to: %s", t)
	WM.W.SetTitle(t)
}

func (WM *WindowManager) SetPairedLabel(t string) {
	log.Printf("changing paired label to: %s", t)
	WM.W.SetTitle(t)
}
