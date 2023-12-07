package pkg

import "log"

type LSDM struct {
	WM      *WindowManager
	LSDChan chan int
}

func NewLSDManager() *LSDM {
	l := &LSDM{}
	l.LSDChan = make(chan int)
	l.WM = l.NewWindowManger()
	go l.HandleChan()
	return l
}

func (l *LSDM) HandleChan() {
	for {
		log.Printf("listening for next command")
		select {
		case i := <-l.LSDChan:
			switch i {
			case 0:
				l.WM.WMChan <- 0
				return
			case 1:
				l.WM.WMChan <- 1
			}
		}
	}
}
