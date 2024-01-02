package main

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/openimsdk/openim-msggateway-proxy/common"
	"github.com/openimsdk/openim-msggateway-proxy/gate"
	"github.com/openimsdk/openim-msggateway-proxy/module"
	log "github.com/xuexihuang/new_log15"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetOutLevel(log.LvlInfo)
	fmt.Println("客户端启动....")
	log.Info("客户端启动....")
	snode, err := snowflake.NewNode(200)
	if err != nil {
		log.Error("snowflake.NewNode error", "err", err)
		return
	}
	common.G_flakeNode = snode
	common.GatewayConsistent = common.NewGatewayConsistent()
	///////////////////////////////////
	gatenet := gate.InitWSsever(80)
	gatenet.SetMsgFun(module.NewAgent, module.CloseAgent, module.DataRecv)
	go gatenet.Runloop()
	///////////////////////////////////////////
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)
	sig := <-c
	log.Info("wsconn server closing down ", "sig", sig)
	gatenet.CloseGate()

}
