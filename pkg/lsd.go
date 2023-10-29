package pkg

import (
	"github.com/getlantern/systray"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Lsd struct {
	Config      *Config
	TrayMenu    *TrayMenu
	Ws          *Ws
	TVConnected bool
	SigChann    chan os.Signal
	ExitChann   chan bool
}

func CreateLsd() *Lsd {
	lsd := new(Lsd)
	lsd.SigChann = make(chan os.Signal)
	lsd.ExitChann = make(chan bool)
	lsd.Config = new(Config)
	lsd.TVConnected = false
	signal.Notify(lsd.SigChann)
	lsd.LoadConfig()
	go lsd.ListenSig()
	lsd.CreateWs()
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
		systray.Quit()
		lsd.ExitChann <- true
	case syscall.SIGINT:
		systray.Quit()
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
	log.Printf("stopping tray\n")
}

func (lsd *Lsd) ConnectScreen() {
	log.Printf("connecting to screen\n")
	lsd.SendWs(&Message{
		Type: "register",
		ID:   "register_0",
		Uri:  "",
		Payload: Payload{
			ForcePairing: true,
			Volume:       0,
			PairingType:  "PROMPT",
			Manifest: Manifest{
				ManifestVersion: 1,
				AppVersion:      "1.1",
				Signed: Signed{
					Created:  "20140509",
					AppId:    "com.lge.test",
					VendorId: "com.lge",
					LocalizedAppNames: map[string]string{
						"":       "LG Remote App",
						"ko-KR":  "리모컨 앱",
						"zxx-XX": "ЛГ Rэмotэ AПП",
					},
					LocalizedVendorNames: map[string]string{
						"": "LG Electronics",
					},
					Permissions: []string{
						"",
					},
					Serial: "2f930e2d2cfe083771f68e4fe7bb07",
				},
				Permissions: []string{
					"WRITE_NOTIFICATION_TOAST",
					"CONTROL_POWER",
				},
				Signatures: []Signature{
					{
						SignatureVersion: 1,
						Signature:        "eyJhbGdvcml0aG0iOiJSU0EtU0hBMjU2Iiwia2V5SWQiOiJ0ZXN0LXNpZ25pbmctY2VydCIsInNpZ25hdHVyZVZlcnNpb24iOjF9.hrVRgjCwXVvE2OOSpDZ58hR+59aFNwYDyjQgKk3auukd7pcegmE2CzPCa0bJ0ZsRAcKkCTJrWo5iDzNhMBWRyaMOv5zWSrthlf7G128qvIlpMT0YNY+n/FaOHE73uLrS/g7swl3/qH/BGFG2Hu4RlL48eb3lLKqTt2xKHdCs6Cd4RMfJPYnzgvI4BNrFUKsjkcu+WD4OO2A27Pq1n50cMchmcaXadJhGrOqH5YmHdOCj5NSHzJYrsW0HPlpuAx/ECMeIZYDh6RMqaFM2DXzdKX9NmmyqzJ3o/0lkk/N97gfVRLW5hA29yeAwaCViZNCP8iC9aO0q9fQojoa7NQnAtw==",
					},
				},
			},
			ClientKey: "",
			Message:   "",
		},
	})
}

func (lsd *Lsd) VerifyConnect() {
	log.Printf("connecting to screen\n")
	lsd.SendWs(&Message{
		Type: "register",
		ID:   "register_0",
		Uri:  "",
		Payload: Payload{
			ForcePairing: false,
			Volume:       0,
			PairingType:  "PROMPT",
			Manifest: Manifest{
				ManifestVersion: 1,
				AppVersion:      "1.1",
				Signed: Signed{
					Created:  "20140509",
					AppId:    "com.lge.test",
					VendorId: "com.lge",
					LocalizedAppNames: map[string]string{
						"":       "LG Remote App",
						"ko-KR":  "리모컨 앱",
						"zxx-XX": "ЛГ Rэмotэ AПП",
					},
					LocalizedVendorNames: map[string]string{
						"": "LG Electronics",
					},
					Permissions: []string{
						"",
					},
					Serial: "2f930e2d2cfe083771f68e4fe7bb07",
				},
				Permissions: []string{
					"WRITE_NOTIFICATION_TOAST",
					"CONTROL_POWER",
				},
				Signatures: []Signature{
					{
						SignatureVersion: 1,
						Signature:        "eyJhbGdvcml0aG0iOiJSU0EtU0hBMjU2Iiwia2V5SWQiOiJ0ZXN0LXNpZ25pbmctY2VydCIsInNpZ25hdHVyZVZlcnNpb24iOjF9.hrVRgjCwXVvE2OOSpDZ58hR+59aFNwYDyjQgKk3auukd7pcegmE2CzPCa0bJ0ZsRAcKkCTJrWo5iDzNhMBWRyaMOv5zWSrthlf7G128qvIlpMT0YNY+n/FaOHE73uLrS/g7swl3/qH/BGFG2Hu4RlL48eb3lLKqTt2xKHdCs6Cd4RMfJPYnzgvI4BNrFUKsjkcu+WD4OO2A27Pq1n50cMchmcaXadJhGrOqH5YmHdOCj5NSHzJYrsW0HPlpuAx/ECMeIZYDh6RMqaFM2DXzdKX9NmmyqzJ3o/0lkk/N97gfVRLW5hA29yeAwaCViZNCP8iC9aO0q9fQojoa7NQnAtw==",
					},
				},
			},
			ClientKey: lsd.Config.LGTVShutDowner.TVInfos.ClientKey,
			Message:   "",
		},
	})
}

func (lsd *Lsd) PingScreen() {
	log.Printf("sending ping to screen\n")
	lsd.SendWs(&Message{
		Type: "request",
		ID:   strconv.Itoa(lsd.Ws.i),
		Uri:  "ssap://system.notifications/createToast",
		Payload: Payload{
			ForcePairing: false,
			Volume:       0,
			PairingType:  "",
			Manifest:     Manifest{},
			ClientKey:    "",
			Message:      "Ping",
		},
	})
}

func (lsd *Lsd) TurnOnScreen() {
	log.Printf("Truning on the screen\n")
	lsd.SendWs(&Message{
		Type: "request",
		ID:   strconv.Itoa(lsd.Ws.i),
		Uri:  "ssap://system.notifications/createToast",
		Payload: Payload{
			ForcePairing: false,
			Volume:       0,
			PairingType:  "",
			Manifest:     Manifest{},
			ClientKey:    "",
			Message:      "Ping",
		},
	})
}

func (lsd *Lsd) TurnOffScreen() {
	log.Printf("Truning off the screen\n")
	lsd.SendWs(&Message{
		Type: "request",
		ID:   strconv.Itoa(lsd.Ws.i),
		Uri:  "ssap://system.notifications/createToast",
		Payload: Payload{
			ForcePairing: false,
			Volume:       0,
			PairingType:  "",
			Manifest:     Manifest{},
			ClientKey:    "",
			Message:      "Ping",
		},
	})
}
