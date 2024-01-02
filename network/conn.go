package network

import (
	"github.com/openimsdk/openim-msggateway-proxy/common"
	"net"
)

type Conn interface {
	ReadMsg() (int, []byte, error)
	WriteMsg(args *common.TWSData) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	SetPongHandler(h func(appData string) error)
	SetPingHandler(h func(appData string) error)
}
