package pkg

import (
	"bytes"
	"encoding/binary"
	"github.com/getlantern/systray"
	"github.com/sabhiram/go-wol/wol"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Lsd struct {
	Config         *Config
	ConfigPath     string
	TrayMenu       *TrayMenu
	Ws             *Ws
	TVConnected    bool
	SigChann       chan os.Signal
	ExitChann      chan bool
	ConfigFilePath string
}

func CreateLsd(cmd string, cfg string) *Lsd {
	var err error
	lsd := new(Lsd)
	lsd.SigChann = make(chan os.Signal)
	lsd.ExitChann = make(chan bool)
	lsd.Config = new(Config)
	lsd.TVConnected = false
	signal.Notify(lsd.SigChann)
	if len(cfg) > 0 {
		lsd.ConfigPath = cfg
	} else {
		lsd.ConfigPath, err = os.UserConfigDir()
		if err != nil {
			log.Fatal(err)
		}
		lsd.ConfigPath += "/LGTVShutDowner"
		err := os.MkdirAll(lsd.ConfigPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	lsd.ConfigFilePath = lsd.ConfigPath + "/config.yaml"
	go lsd.ListenSig()
	lsd.LoadConfig()
	switch cmd {
	case "PowerOn":
		go lsd.TurnOnScreen()
		err = lsd.CreateWs()
	default:
		log.Printf("error: command %s unknow", cmd)
		err = lsd.CreateWs()
	}
	if err != nil {
		log.Printf("retrying\n")
		for i := 0; i < 10; i++ { // start of the execution block
			time.Sleep(time.Millisecond * 1000)
			err = lsd.CreateWs()
			if err == nil {
				break
			}
		}
		if err != nil {
			lsd.ExitChann <- true
			return lsd
		}
		return lsd
	}
	systray.Run(lsd.InitTray(cmd), lsd.Close)
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
		lsd.TurnOffScreen()
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
	data, err := os.ReadFile(lsd.ConfigPath + "/screen.png")
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
						"LAUNCH",
						"LAUNCH_WEBAPP",
						"APP_TO_APP",
						"CLOSE",
						"TEST_OPEN",
						"TEST_PROTECTED",
						"CONTROL_AUDIO",
						"CONTROL_DISPLAY",
						"CONTROL_INPUT_JOYSTICK",
						"CONTROL_INPUT_MEDIA_RECORDING",
						"CONTROL_INPUT_MEDIA_PLAYBACK",
						"CONTROL_INPUT_TV",
						"CONTROL_POWER",
						"READ_APP_STATUS",
						"READ_CURRENT_CHANNEL",
						"READ_INPUT_DEVICE_LIST",
						"READ_NETWORK_STATE",
						"READ_RUNNING_APPS",
						"READ_TV_CHANNEL_LIST",
						"WRITE_NOTIFICATION_TOAST",
						"READ_POWER_STATE",
						"READ_COUNTRY_INFO",
						"READ_SETTINGS",
						"CONTROL_TV_SCREEN",
						"CONTROL_TV_STANBY",
						"CONTROL_FAVORITE_GROUP",
						"CONTROL_USER_INFO",
						"CHECK_BLUETOOTH_DEVICE",
						"CONTROL_BLUETOOTH",
						"CONTROL_TIMER_INFO",
						"STB_INTERNAL_CONNECTION",
						"CONTROL_RECORDING",
						"READ_RECORDING_STATE",
						"WRITE_RECORDING_LIST",
						"READ_RECORDING_LIST",
						"READ_RECORDING_SCHEDULE",
						"WRITE_RECORDING_SCHEDULE",
						"READ_STORAGE_DEVICE_LIST",
						"READ_TV_PROGRAM_INFO",
						"CONTROL_BOX_CHANNEL",
						"READ_TV_ACR_AUTH_TOKEN",
						"READ_TV_CONTENT_STATE",
						"READ_TV_CURRENT_TIME",
						"ADD_LAUNCHER_CHANNEL",
						"SET_CHANNEL_SKIP",
						"RELEASE_CHANNEL_SKIP",
						"CONTROL_CHANNEL_BLOCK",
						"DELETE_SELECT_CHANNEL",
						"CONTROL_CHANNEL_GROUP",
						"SCAN_TV_CHANNELS",
						"CONTROL_TV_POWER",
						"CONTROL_WOL",
					},
					Serial: "2f930e2d2cfe083771f68e4fe7bb07",
				},
				Permissions: []string{
					"TEST_SECURE",
					"CONTROL_INPUT_TEXT",
					"CONTROL_MOUSE_AND_KEYBOARD",
					"READ_INSTALLED_APPS",
					"READ_LGE_SDX",
					"READ_NOTIFICATIONS",
					"SEARCH",
					"WRITE_SETTINGS",
					"WRITE_NOTIFICATION_ALERT",
					"WRITE_NOTIFICATION_TOAST",
					"CONTROL_POWER",
					"READ_CURRENT_CHANNEL",
					"READ_RUNNING_APPS",
					"READ_UPDATE_INFO",
					"UPDATE_FROM_REMOTE_APP",
					"READ_LGE_TV_INPUT_EVENTS",
					"READ_TV_CURRENT_TIME",
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
						"LAUNCH",
						"LAUNCH_WEBAPP",
						"APP_TO_APP",
						"CLOSE",
						"TEST_OPEN",
						"TEST_PROTECTED",
						"CONTROL_AUDIO",
						"CONTROL_DISPLAY",
						"CONTROL_INPUT_JOYSTICK",
						"CONTROL_INPUT_MEDIA_RECORDING",
						"CONTROL_INPUT_MEDIA_PLAYBACK",
						"CONTROL_INPUT_TV",
						"CONTROL_POWER",
						"READ_APP_STATUS",
						"READ_CURRENT_CHANNEL",
						"READ_INPUT_DEVICE_LIST",
						"READ_NETWORK_STATE",
						"READ_RUNNING_APPS",
						"READ_TV_CHANNEL_LIST",
						"WRITE_NOTIFICATION_TOAST",
						"READ_POWER_STATE",
						"READ_COUNTRY_INFO",
						"READ_SETTINGS",
						"CONTROL_TV_SCREEN",
						"CONTROL_TV_STANBY",
						"CONTROL_FAVORITE_GROUP",
						"CONTROL_USER_INFO",
						"CHECK_BLUETOOTH_DEVICE",
						"CONTROL_BLUETOOTH",
						"CONTROL_TIMER_INFO",
						"STB_INTERNAL_CONNECTION",
						"CONTROL_RECORDING",
						"READ_RECORDING_STATE",
						"WRITE_RECORDING_LIST",
						"READ_RECORDING_LIST",
						"READ_RECORDING_SCHEDULE",
						"WRITE_RECORDING_SCHEDULE",
						"READ_STORAGE_DEVICE_LIST",
						"READ_TV_PROGRAM_INFO",
						"CONTROL_BOX_CHANNEL",
						"READ_TV_ACR_AUTH_TOKEN",
						"READ_TV_CONTENT_STATE",
						"READ_TV_CURRENT_TIME",
						"ADD_LAUNCHER_CHANNEL",
						"SET_CHANNEL_SKIP",
						"RELEASE_CHANNEL_SKIP",
						"CONTROL_CHANNEL_BLOCK",
						"DELETE_SELECT_CHANNEL",
						"CONTROL_CHANNEL_GROUP",
						"SCAN_TV_CHANNELS",
						"CONTROL_TV_POWER",
						"CONTROL_WOL",
					},
					Serial: "2f930e2d2cfe083771f68e4fe7bb07",
				},
				Permissions: []string{
					"TEST_SECURE",
					"CONTROL_INPUT_TEXT",
					"CONTROL_MOUSE_AND_KEYBOARD",
					"READ_INSTALLED_APPS",
					"READ_LGE_SDX",
					"READ_NOTIFICATIONS",
					"SEARCH",
					"WRITE_SETTINGS",
					"WRITE_NOTIFICATION_ALERT",
					"WRITE_NOTIFICATION_TOAST",
					"CONTROL_POWER",
					"READ_CURRENT_CHANNEL",
					"READ_RUNNING_APPS",
					"READ_UPDATE_INFO",
					"UPDATE_FROM_REMOTE_APP",
					"READ_LGE_TV_INPUT_EVENTS",
					"READ_TV_CURRENT_TIME",
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
			Message:      "Hey, you cutie ;)",
		},
	})
}

func (lsd *Lsd) TurnOnScreen() {
	packet, err := wol.New(lsd.Config.LGTVShutDowner.TVInfos.Mac)
	if err != nil {
		log.Printf("failed to turn off screen")
	}
	// Fill our byte buffer with the bytes in our MagicPacket
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, packet)

	// Get a UDPAddr to send the broadcast to
	udpAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
	if err != nil {
		log.Printf("Unable to get a UDP address for %s\n", "255.255.255.255:9")
		return
	}

	// Open a UDP connection, and defer its cleanup
	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Printf("Unable to dial UDP address for %s\n", "255.255.255.255:9")
		return
	}
	defer connection.Close()

	// Write the bytes of the MagicPacket to the connection
	bytesWritten, err := connection.Write(buf.Bytes())
	if err != nil {
		log.Printf("Unable to write packet to connection\n")
		return
	} else if bytesWritten != 102 {
		log.Printf("warning: %d bytes written, %d expected!\n", bytesWritten, 102)
	}
}

func (lsd *Lsd) GetAPIList() {
	log.Printf("Getting API List\n")
	lsd.SendWs(&Message{
		Type: "request",
		ID:   strconv.Itoa(lsd.Ws.i),
		Uri:  "ssap://api/getServiceList",
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
		Uri:  "ssap://system/turnOff",
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
