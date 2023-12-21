package network

import (
	"github.com/openimsdk/openim-msggateway-proxy/common"
	log "github.com/xuexihuang/new_log15"
	"testing"
)

// ///////////////////////
// /////////////////////////
func TestRobot3d(t *testing.T) {

	client, err := NewWSClient("", "", "/ws?type=robot3d")
	if err != nil {
		log.Info("NewWSClient error")
		return
	}
	defer client.Destroy()
	body := `{
    "type":3,
    "tts_text":"你好"
    }`

	client.WriteMsg(&common.TWSData{MsgType: common.MessageText, Msg: []byte(body)})
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
