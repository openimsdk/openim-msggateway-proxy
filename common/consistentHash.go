package common

import (
	"context"
	"fmt"
	"github.com/stathat/consistent"
	log "github.com/xuexihuang/new_log15"
	"os"
	"strconv"
	"strings"
)

var GatewayConsistent *consistent.Consistent

func NewGatewayConsistent() *consistent.Consistent {
	gatewayConsistent := consistent.New()
	gatewayHosts := getMsgGatewayHost(context.Background())
	for _, v := range gatewayHosts {
		gatewayConsistent.Add(v)
	}
	return gatewayConsistent
}

// like openimserver-openim-msggateway-0.openimserver-openim-msggateway-headless.openim-lin.svc.cluster.local:80
func getMsgGatewayHost(ctx context.Context) []string {
	instance := "openimserver"
	selfPodName := os.Getenv("MY_POD_NAME")
	replicas := os.Getenv("MY_MSGGATEWAY_REPLICACOUNT")
	ns := os.Getenv("MY_POD_NAMESPACE")
	nReplicas, _ := strconv.Atoi(replicas)
	podInfo := strings.Split(selfPodName, "-")
	instance = podInfo[0]
	var ret []string
	for i := 0; i < nReplicas; i++ {
		host := fmt.Sprintf("%s-openim-msggateway-%d.%s-openim-msggateway-headless.%s.svc.cluster.local", instance, i, instance, ns)
		ret = append(ret, host)
	}
	log.Info("msggateway host info", "ret", ret)
	return ret
}
