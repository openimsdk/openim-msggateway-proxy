package main

import (
	"github.com/openimsdk/openim-msggateway-proxy/common"
	"github.com/openimsdk/openim-msggateway-proxy/network"
	log "github.com/xuexihuang/new_log15"
	"time"
)

func main() {
	client, err := network.NewWSClient("12345678", "127.0.0.1:80", "/ws?sendID=huanglin")
	if err != nil {
		log.Info("NewWSClient error")
		return
	}
	defer client.Destroy()
	body := `{
    "type":3,
    "tts_text":"你好"
    }`
	go func() {
		for {
			time.Sleep(2 * time.Second)
			client.WriteMsg(&common.TWSData{MsgType: common.MessageText, Msg: []byte(body)})
			time.Sleep(3 * time.Second)
			client.WriteMsg(&common.TWSData{MsgType: common.PingMessage, Msg: []byte("pingmsg")})
			time.Sleep(3 * time.Second)
			client.WriteMsg(&common.TWSData{MsgType: common.PongMessage, Msg: []byte("pongmsg")})
		}

	}()

	for {
		select {
		case clientMsg, _ := <-client.RecvMsgChan:
			if clientMsg.MsgType == common.CloseMessage {
				log.Info("ws net error")
				return
			} else {
				log.Info("get one recv data", "msgType", clientMsg.MsgType)
			}
		}
	}
}
