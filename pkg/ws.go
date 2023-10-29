package pkg

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Ws struct {
	client *websocket.Conn
	i      int
}

func (lsd *Lsd) CreateWs() {
	var err error
	res := new(http.Response)
	ws := new(Ws)
	ws.i = 0
	uri := url.URL{Scheme: "wss", Host: lsd.Config.LGTVShutDowner.TVInfos.Ip, Path: ""}
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	ws.client, res, err = dialer.Dial(uri.String(), nil)
	if err != nil || res == nil {
		log.Printf("error connecting to screen websocket: %s\n", err)
		return
	}
	log.Printf("websocket response: %s\n", res.Status)
	lsd.Ws = ws
	go lsd.VerifyConnect()
	go lsd.ListenWs()
}

func (lsd *Lsd) ListenWs() {
	for {
		_, msg, err := lsd.Ws.client.ReadMessage()
		if err != nil {
			log.Printf("error reading websocket message: %s\n", err)
			return
		}
		lsd.HandleWsMsg(msg)
	}
}

func (lsd *Lsd) SendWs(m *Message) {
	log.Printf("sending message: %v\n", m)
	err := lsd.Ws.client.WriteJSON(m)
	if err != nil {
		return
	}
	lsd.Ws.i++
}

func (lsd *Lsd) HandleWsMsg(b []byte) {
	p := new(fastjson.Parser)
	v, err := p.ParseBytes(b)
	if err != nil {
		log.Printf("error parsing message: %s", err)
	}
	switch string(v.GetStringBytes("type")) {
	case "registered":
		log.Printf("message received: %s\n", b)
		s := v.GetObject("payload").Get("client-key").String()
		lsd.Config.LGTVShutDowner.TVInfos.ClientKey = s[1 : len(s)-1]
		lsd.SaveConfig()
	case "error":
		log.Printf("permission error, trying to reconnect: %s\n", b)
		errPld := string(v.GetStringBytes("error"))
		if strings.Contains(errPld, "401") {
			lsd.VerifyConnect()
		}
	default:
		log.Printf("no handler for message received: %s\n", b)
	}
}
