package gate

import (
	"fmt"
	"github.com/openimsdk/openim-msggateway-proxy/network/tjson"
	"sync"
	"time"
)

var Processor = tjson.NewProcessor()

type GateNet struct {
	*Gate
	CloseSig chan bool
	Wg       sync.WaitGroup
}

func InitWSsever(wsPort int) *GateNet {
	gatenet := new(GateNet)
	gatenet.Gate = &Gate{
		MaxConnNum:      10000,
		PendingWriteNum: 1000,
		MaxMsgLen:       5000000,
		WSAddr:          ":" + fmt.Sprintf("%d", wsPort),
		HTTPTimeout:     10 * time.Second,
		CertFile:        "",
		KeyFile:         "",
		LenMsgLen:       2,
		Processor:       Processor,
	}
	gatenet.CloseSig = make(chan bool, 1)
	return gatenet
}

func (gt *GateNet) SetMsgFun(Fun1 func(Agent), Fun2 func(Agent), Fun3 func(interface{}, Agent)) {
	gt.Gate.SetFun(Fun1, Fun2, Fun3)
}
func (gt *GateNet) Runloop() {
	gt.Wg.Add(1)
	gt.Run(gt.CloseSig)
	gt.Wg.Done()
}
func (gt *GateNet) CloseGate() {
	gt.CloseSig <- true
	gt.Wg.Wait()
	gt.Gate.OnDestroy()
}
