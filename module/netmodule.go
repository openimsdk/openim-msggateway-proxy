package module

import (
	"errors"
	"github.com/openimsdk/openim-msggateway-proxy/common"
	"github.com/openimsdk/openim-msggateway-proxy/gate"
	log "github.com/xuexihuang/new_log15"
	"net/url"
)

type MActor interface {
	//recv消息
	ProcessRecvMsg(interface{}) error
	//关闭循环，并释放资源
	Destroy()
	run()
}

func NewAgent(a gate.Agent) {
	aUerData := a.UserData().(*common.TAgentUserData)
	log.Info("one ws connect", "sessionId", aUerData.SessionID)
	param, err := checkPath(aUerData)
	if err != nil {
		log.Error("ws path校验失败", "userData", aUerData, "sessionId", aUerData.SessionID)
		a.Destroy()
		return
	}
	log.Info("checkToken info", "param", param, "err", err)
	actor, err := NewMActor(a, param.SessionId, param)
	if err != nil {
		log.Error("NewMQActor error", "err", err, "sessionId", aUerData.SessionID)
		a.Destroy()
		return
	}
	aUerData.ProxyBody = actor
	aUerData.UserId = param.GetUserID()
	a.SetUserData(aUerData)
	log.Info("one linked", "param", param, "sessionId", aUerData.SessionID)
}

func CloseAgent(a gate.Agent) {
	aUerData := a.UserData().(*common.TAgentUserData)
	if aUerData.ProxyBody != nil {
		aUerData.ProxyBody.(MActor).Destroy()
		aUerData.ProxyBody = nil
	}
	log.Info("one dislinkder", "sessionId", a.UserData().(*common.TAgentUserData).SessionID)
}
func DataRecv(data interface{}, a gate.Agent) {
	aUerData := a.UserData().(*common.TAgentUserData)
	if aUerData.ProxyBody != nil {
		err := aUerData.ProxyBody.(MActor).ProcessRecvMsg(data)
		if err != nil {
			log.Error("溢出错误", "sessionId", aUerData.SessionID)
			a.Destroy()
		}
	}
}
func checkPath(data *common.TAgentUserData) (*ParamStru, error) {
	ret := new(ParamStru)
	ret.SessionId = data.SessionID
	u, err := url.Parse(data.AppString)
	if err != nil {
		log.Error("ws url path not correct", "sessionId", data.SessionID)
		return nil, errors.New("ws url path not correct")
	}
	q := u.Query()
	ret.UserId = q.Get("sendID")
	ret.UrlPath = data.AppString
	return ret, nil
}
