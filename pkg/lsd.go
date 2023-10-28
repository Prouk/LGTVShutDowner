package pkg

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/websocket"
	"net"
)

type Lsd struct {
	Config   *Config
	Commands *Commands
	Ip       net.IP
	WsConn   *websocket.Conn
}

func NewLsd() (*Lsd, error) {
	fmt.Printf("creating lsd manager\n")
	lsd := &Lsd{}
	lsd.Config = &Config{}
	lsd.Commands = &Commands{}
	fmt.Printf("reading conf file\n")
	err := lsd.Config.ReadConf()
	if err != nil {
		fmt.Printf("couldn't read config file: %s\n", err)
		return lsd, err
	}
	err = lsd.Commands.ReadCommands()
	if err != nil {
		fmt.Printf("couldn't read commands file: %s\n", err)
		return lsd, err
	}
	lsd.Ip, err = GetLocalIP()
	if err != nil {
		return lsd, err
	}
	return lsd, nil
}

func (lsd *Lsd) ConnectTv() error {
	fmt.Printf("connecting to device: %s :: %s\n",
		lsd.Config.LGTVShutDowner.TVInfos.Ip,
		lsd.Config.LGTVShutDowner.TVInfos.Mac)
	config, err := websocket.NewConfig("wss://"+lsd.Config.LGTVShutDowner.TVInfos.Ip+":3001",
		"wss://"+lsd.Ip.String()+":3001")
	config.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	if err != nil {
		return err
	}
	ws, err := websocket.DialConfig(config)
	if err != nil {
		return err
	}
	if ws.IsClientConn() {
		fmt.Printf("websocket client connected\n")
	}
	lsd.WsConn = ws
	return nil
}

func (lsd *Lsd) GetClientKey() error {
	message := "{\"type\":\"register\",\"id\":\"register_0\",\"payload\":{\"forcePairing\":false,\"pairingType\":\"PROMPT\",\"manifest\":{\"manifestVersion\":1,\"appVersion\":\"1.1\",\"signed\":{\"created\":\"20140509\",\"appId\":\"com.lge.test\",\"vendorId\":\"com.lge\",\"localizedAppNames\":{\"\":\"LG Remote App\",\"ko-KR\":\"리모컨 앱\",\"zxx-XX\":\"ЛГ Rэмotэ AПП\"},\"localizedVendorNames\":{\"\":\"LG Electronics\"},\"permissions\":[\"TEST_SECURE\",\"CONTROL_INPUT_TEXT\",\"CONTROL_MOUSE_AND_KEYBOARD\",\"READ_INSTALLED_APPS\",\"READ_LGE_SDX\",\"READ_NOTIFICATIONS\",\"SEARCH\",\"WRITE_SETTINGS\",\"WRITE_NOTIFICATION_ALERT\",\"CONTROL_POWER\",\"READ_CURRENT_CHANNEL\",\"READ_RUNNING_APPS\",\"READ_UPDATE_INFO\",\"UPDATE_FROM_REMOTE_APP\",\"READ_LGE_TV_INPUT_EVENTS\",\"READ_TV_CURRENT_TIME\"],\"serial\":\"2f930e2d2cfe083771f68e4fe7bb07\"},\"permissions\":[\"LAUNCH\",\"LAUNCH_WEBAPP\",\"APP_TO_APP\",\"CLOSE\",\"TEST_OPEN\",\"TEST_PROTECTED\",\"CONTROL_AUDIO\",\"CONTROL_DISPLAY\",\"CONTROL_INPUT_JOYSTICK\",\"CONTROL_INPUT_MEDIA_RECORDING\",\"CONTROL_INPUT_MEDIA_PLAYBACK\",\"CONTROL_INPUT_TV\",\"CONTROL_POWER\",\"READ_APP_STATUS\",\"READ_CURRENT_CHANNEL\",\"READ_INPUT_DEVICE_LIST\",\"READ_NETWORK_STATE\",\"READ_RUNNING_APPS\",\"READ_TV_CHANNEL_LIST\",\"WRITE_NOTIFICATION_TOAST\",\"READ_POWER_STATE\",\"READ_COUNTRY_INFO\"],\"signatures\":[{\"signatureVersion\":1,\"signature\":\"eyJhbGdvcml0aG0iOiJSU0EtU0hBMjU2Iiwia2V5SWQiOiJ0ZXN0LXNpZ25pbmctY2VydCIsInNpZ25hdHVyZVZlcnNpb24iOjF9.hrVRgjCwXVvE2OOSpDZ58hR+59aFNwYDyjQgKk3auukd7pcegmE2CzPCa0bJ0ZsRAcKkCTJrWo5iDzNhMBWRyaMOv5zWSrthlf7G128qvIlpMT0YNY+n/FaOHE73uLrS/g7swl3/qH/BGFG2Hu4RlL48eb3lLKqTt2xKHdCs6Cd4RMfJPYnzgvI4BNrFUKsjkcu+WD4OO2A27Pq1n50cMchmcaXadJhGrOqH5YmHdOCj5NSHzJYrsW0HPlpuAx/ECMeIZYDh6RMqaFM2DXzdKX9NmmyqzJ3o/0lkk/N97gfVRLW5hA29yeAwaCViZNCP8iC9aO0q9fQojoa7NQnAtw==\"}]}}}"
	err := websocket.Message.Send(lsd.WsConn, message)
	if err != nil {
		return err
	}
	return nil
}

func (lsd *Lsd) SendCommand(s string) error {
	cmd, err := lsd.Commands.GetCommand(s)
	if err != nil {
		return err
	}
	fmt.Printf("sending command to TV: %s\n", cmd.Name)
	message := map[string]string{"uri": cmd.URI}
	err = websocket.JSON.Send(lsd.WsConn, message)
	if err != nil {
		return err
	}
	return nil
}

func GetLocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddress := conn.LocalAddr().(*net.UDPAddr)
	return localAddress.IP, nil
}
