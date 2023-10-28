package pkg

import "github.com/getlantern/systray"

type TrayMenu struct {
	ConnectMenuItem *systray.MenuItem
	PingMenuItem    *systray.MenuItem
	QuitMenuItem    *systray.MenuItem
}

func (lsd *Lsd) InitTray() {
	systray.SetIcon(lsd.GetIcon())
	systray.SetTitle("LSD")
	systray.SetTooltip("LGTV ShutDown")
	lsd.TrayMenu = lsd.GetTrayMenuItems()
	go lsd.ListenClick()
}

func (lsd *Lsd) GetTrayMenuItems() *TrayMenu {
	tm := new(TrayMenu)
	tm.ConnectMenuItem = systray.AddMenuItem("Connect", "Connect to the screen")
	tm.PingMenuItem = systray.AddMenuItem("Ping", "Send Ping message to screen")
	tm.QuitMenuItem = systray.AddMenuItem("Quit", "Stop lsd service")
	return tm
}

func (lsd *Lsd) ListenClick() {
	for {
		select {
		case <-lsd.TrayMenu.PingMenuItem.ClickedCh:
			lsd.PingScreen()
		case <-lsd.TrayMenu.QuitMenuItem.ClickedCh:
			systray.Quit()
			lsd.ExitChann <- true
			return
		}
	}
}
