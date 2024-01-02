package network

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/openimsdk/openim-msggateway-proxy/common"
	"github.com/openimsdk/openim-msggateway-proxy/network/tjson"
	log "github.com/xuexihuang/new_log15"
	"net"
	"sync"
	"time"
)

type WsClient struct {
	sync.Mutex
	conn        *websocket.Conn
	Addr        string
	dialer      websocket.Dialer
	closeFlag   bool
	sendMsgChan chan *common.TWSData //send msg chan
	wg          sync.WaitGroup
	RecvMsgChan chan *common.TWSData
	Processor   Processor
	sessionId   string
}

func NewWSClient(sessionID string, host string, path string) (*WsClient, error) {

	ret := &WsClient{
		Addr:        "ws://" + host + path,
		sendMsgChan: make(chan *common.TWSData, 100),
		RecvMsgChan: make(chan *common.TWSData, 100),
		closeFlag:   false,
		dialer:      websocket.Dialer{HandshakeTimeout: 3 * time.Second},
		Processor:   tjson.NewProcessor(),
		sessionId:   sessionID,
	}

	fmt.Println("ret.addr is:", ret.Addr)
	conn, err := ret.dial()
	if err != nil {
		log.Error("ret.dial error", "sessionId", sessionID)
		return nil, err
	}
	ret.conn = conn
	conn.SetReadLimit(5000000)
	ret.wg.Add(1)
	go func() {
		defer common.TryRecoverAndDebugPrint()
		for b := range ret.sendMsgChan {
			if b == nil {
				break
			}
			var err error
			if b.MsgType == common.MessageBinary {
				err = conn.WriteMessage(websocket.BinaryMessage, b.Msg)
			} else if b.MsgType == common.MessageText {
				err = conn.WriteMessage(websocket.TextMessage, b.Msg)
			} else if b.MsgType == common.PongMessage {
				err = conn.WriteMessage(websocket.PongMessage, b.Msg)
			} else if b.MsgType == common.PingMessage {
				err = conn.WriteMessage(websocket.PingMessage, b.Msg)
			}
			if err != nil {
				break
			}
		}

		conn.Close()
		ret.Lock()
		ret.closeFlag = true
		ret.Unlock()
		ret.wg.Done()
	}()
	go ret.run()
	return ret, nil
}

func (client *WsClient) dial() (*websocket.Conn, error) {

	conn, _, err := client.dialer.Dial(client.Addr, nil)
	if err == nil || client.closeFlag {
		return conn, nil
	}
	log.Info("connect error", "sessionId", client.sessionId, "client.Addr", client.Addr, "err", err)
	return nil, err
}
func (a *WsClient) run() {
	defer common.TryRecoverAndDebugPrint()
	a.conn.SetPongHandler(func(appData string) error {
		msg, err := a.Processor.UnmarshalMul(common.PongMessage, []byte(appData))
		if err != nil {
			log.Error("unmarshal message error", "sessionId", a.sessionId, "err", err)
			return nil
		}
		a.RecvMsgChan <- msg.(*common.TWSData)
		return nil
	})
	a.conn.SetPingHandler(func(appData string) error {
		msg, err := a.Processor.UnmarshalMul(common.PingMessage, []byte(appData))
		if err != nil {
			log.Error("unmarshal message error", "sessionId", a.sessionId, "err", err)
			return nil
		}
		a.RecvMsgChan <- msg.(*common.TWSData)
		return nil
	})
	for {
		nType, data, err := a.conn.ReadMessage()
		if err != nil {
			//log.Debug("read message: %v", err)
			log.Info("read message error", "sessionId", a.sessionId, "error", err)
			break
		}
		if a.Processor != nil {
			msg, err := a.Processor.UnmarshalMul(nType, data)
			if err != nil {
				log.Error("unmarshal message error", "sessionId", a.sessionId, "err", err)
				break
			}
			a.RecvMsgChan <- msg.(*common.TWSData)
		}
	}
	//////////////////////////////
	a.RecvMsgChan <- &common.TWSData{MsgType: common.CloseMessage, Msg: nil}
	/////////////////////////////
	a.Lock()
	defer a.Unlock()
	if a.closeFlag {
		return
	}
	a.doWrite(nil)
	a.closeFlag = true
}
func (client *WsClient) doDestroy() {
	client.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	client.conn.Close()

	if !client.closeFlag {
		close(client.sendMsgChan)
		client.closeFlag = true
	}
}

func (client *WsClient) doWrite(b *common.TWSData) {
	if len(client.sendMsgChan) == cap(client.sendMsgChan) {
		//log.Debug("close conn: channel full")
		log.Error("close conn: channel full", "sessionId", client.sessionId)
		client.doDestroy()
		return
	}

	client.sendMsgChan <- b
}

func (client *WsClient) Destroy() {
	client.Lock()
	client.doDestroy()
	client.Unlock()
	client.wg.Wait()
}
func (client *WsClient) WriteMsg(data *common.TWSData) {
	client.Lock()
	defer client.Unlock()
	if client.closeFlag {
		return
	}
	client.doWrite(data)
}
