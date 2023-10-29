package pkg

import (
	"github.com/getlantern/systray"
)

type TrayMenu struct {
	ConnectMenuItem *systray.MenuItem
	PingMenuItem    *systray.MenuItem
	QuitMenuItem    *systray.MenuItem
	OnMenuItem      *systray.MenuItem
	APIMenuItem     *systray.MenuItem
	OffMenuItem     *systray.MenuItem
}

func (lsd *Lsd) InitTray(cmd string) func() {
	return func() {
		systray.SetIcon(lsd.GetIcon())
		systray.SetTitle("LSD")
		systray.SetTooltip("LGTV ShutDown")
		lsd.TrayMenu = lsd.GetTrayMenuItems()
		go lsd.ListenClick()
	}
}

func (lsd *Lsd) GetTrayMenuItems() *TrayMenu {
	tm := new(TrayMenu)
	tm.ConnectMenuItem = systray.AddMenuItem("Connect", "Connect to the screen")
	tm.OnMenuItem = systray.AddMenuItem("TurnOn", "Turn on the screen")
	tm.OffMenuItem = systray.AddMenuItem("ShutDown", "Shutdown the damn screen")
	tm.PingMenuItem = systray.AddMenuItem("Ping", "Send Ping message to screen")
	tm.APIMenuItem = systray.AddMenuItem("ListAPI", "Request a list of TV API")
	tm.QuitMenuItem = systray.AddMenuItem("Quit", "Stop lsd service")
	return tm
}

func (lsd *Lsd) ListenClick() {
	for {
		select {
		case <-lsd.TrayMenu.ConnectMenuItem.ClickedCh:
			lsd.ConnectScreen()
		case <-lsd.TrayMenu.OnMenuItem.ClickedCh:
			lsd.TurnOnScreen()
		case <-lsd.TrayMenu.OffMenuItem.ClickedCh:
			lsd.TurnOffScreen()
		case <-lsd.TrayMenu.PingMenuItem.ClickedCh:
			lsd.PingScreen()
		case <-lsd.TrayMenu.APIMenuItem.ClickedCh:
			lsd.GetAPIList()
		case <-lsd.TrayMenu.QuitMenuItem.ClickedCh:
			systray.Quit()
			lsd.ExitChann <- true
			return
		}
	}
}
