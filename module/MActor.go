package module

import (
	"errors"
	"github.com/openimsdk/openim-msggateway-proxy/common"
	"github.com/openimsdk/openim-msggateway-proxy/gate"
	"github.com/openimsdk/openim-msggateway-proxy/network"
	log "github.com/xuexihuang/new_log15"
	"net/url"
	"sync"
)

const (
	WsUserID    = "sendID"
	OperationID = "operationID"
	PlatformID  = "platformID"
)

type ParamStru struct {
	UrlPath   string
	Token     string
	SessionId string
	UserId    string
	GroupId   string
	OrgId     string
	OrgName   string
}

func (p *ParamStru) GetUserID() string {
	u, err := url.Parse(p.UrlPath)
	if err != nil {
		return ""
	}
	return u.Query().Get(WsUserID)
}
func (p *ParamStru) GetOperationID() string {
	u, err := url.Parse(p.UrlPath)
	if err != nil {
		return ""
	}
	return u.Query().Get(OperationID)
}
func (p *ParamStru) GetPlatformID() string {
	u, err := url.Parse(p.UrlPath)
	if err != nil {
		return ""
	}
	return u.Query().Get(PlatformID)
}

type MActorIm struct {
	mWsClient     *network.WsClient
	param         *ParamStru
	nChanLen      int //接收数据网络缓存
	wg            sync.WaitGroup
	a             gate.Agent
	SessionId     string
	closeChan     chan bool        //主动关闭协程的通道
	ReceivMsgChan chan interface{} //接收网络层数据通道
	isclosing     bool
}

func NewMActor(a gate.Agent, sessionId string, appParam *ParamStru) (MActor, error) {
	ret := &MActorIm{param: appParam, a: a, SessionId: sessionId, closeChan: make(chan bool, 1), nChanLen: 10, ReceivMsgChan: make(chan interface{}, 10), isclosing: false}
	///////////////////////////////////////
	host, err := common.GatewayConsistent.Get(appParam.UserId)
	if err != nil {
		log.Error("get gatewayConsistent error", "appParam", appParam)
		return nil, errors.New("get gatewayConsistent error")
	}
	mWsClient, err := network.NewWSClient(sessionId, host, appParam.UrlPath)
	if err != nil {
		log.Error("network.NewWSClient error", "err", err)
		return nil, err
	}
	ret.mWsClient = mWsClient
	///////////////////////////////////////
	go ret.run()
	return ret, nil
}

func (actor *MActorIm) run() {
	actor.wg.Add(1)
	defer common.TryRecoverAndDebugPrint()
	defer actor.wg.Done()
	for {
		select {
		case <-actor.closeChan:
			log.Info("收到退出信号", "sessionId", actor.SessionId)
			actor.mWsClient.Destroy()
			return
		case recvData := <-actor.ReceivMsgChan:
			if actor.isclosing == true {
				continue
			}
			data := recvData.(*common.TWSData)
			actor.mWsClient.WriteMsg(data)
		case resp := <-actor.mWsClient.RecvMsgChan:
			if resp.MsgType == common.CloseMessage {
				actor.a.Destroy()
				actor.isclosing = true
			} else {
				actor.a.WriteMsg(resp)
			}
		}
	}
}

func (actor *MActorIm) Destroy() {
	actor.closeChan <- true
	actor.wg.Wait()
	actor.a = nil
	log.Info("退出MQPushActorIm", "sessionId", actor.SessionId)
}
func (actor *MActorIm) ProcessRecvMsg(msg interface{}) error {
	if len(actor.ReceivMsgChan) == actor.nChanLen {
		log.Error("send channel is full", "sessionId", actor.SessionId)
		return errors.New("send channel is full")
	}
	actor.ReceivMsgChan <- msg
	return nil
}
