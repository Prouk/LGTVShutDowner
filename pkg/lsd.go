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
		"wss://"+lsd.Ip.String())
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
