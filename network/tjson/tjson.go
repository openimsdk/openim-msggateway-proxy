package tjson

import (
	"errors"
	"github.com/openimsdk/openim-msggateway-proxy/common"
)

type Login struct {
	UserName string
	PassWord string
}

type Processor struct {
}

func NewProcessor() *Processor {
	p := new(Processor)
	return p
}

func (p *Processor) UsePacketMode() bool {
	return false
}
func (p *Processor) Marshal(msg interface{}) (*common.TWSData, error) {
	////////////////////////////////////////////////////////////////////////////
	tsend := msg.(*common.TWSData)
	if tsend.MsgType != common.MessageText && tsend.MsgType != common.MessageBinary && tsend.MsgType != common.PongMessage && tsend.MsgType != common.PingMessage {
		return nil, errors.New("msg is not correct")
	}
	return tsend, nil
}
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
	return &Login{UserName: "nihao", PassWord: "huanglin"}, nil
}
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	return nil
}
func (p *Processor) UnmarshalMul(nType int, data []byte) (interface{}, error) {
	ret := &common.TWSData{}
	if nType != common.MessageText && nType != common.MessageBinary && nType != common.PingMessage && nType != common.PongMessage {
		return nil, errors.New("ws PackType is not correct")
	}
	ret.MsgType = nType
	ret.Msg = data
	return ret, nil
}
