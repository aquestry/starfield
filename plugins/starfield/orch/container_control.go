package orch

import (
	"fmt"
	"github.com/aquestry/starfield/plugins/starfield/orch/node"
	"github.com/go-logr/logr"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"time"
)

var ProxyInstance *proxy.Proxy
var Logger logr.Logger

func CreateContainer(name, tag, template string) {
	start := time.Now()
	n := getNodeWithLowestInstances()
	cmd := fmt.Sprintf("docker run --name %s -d -e PAPER_VELOCITY_SECRET=%s -p 0:25565 %s",
		name, ProxyInstance.Config().Forwarding.VelocitySecret, template)
	_, err := n.Run(cmd)
	if err != nil {
		Logger.Error(err, "CreateContainer failed", "command", cmd)
		return
	}
	RegisterContainer(name, tag, n.Addr(), n, start)
}

func DeleteContainer(name string) {
	srv, err := GetContainer(name)
	if err != nil {
		Logger.Error(err, "container not found")
		return
	}
	cmd := fmt.Sprintf("docker rm %s --force", name)
	_, e := srv.Node.Run(cmd)
	if e != nil {
		Logger.Error(e, "DeleteContainer failed", "command", cmd)
	}
	Remove(name)
}

func getNodeWithLowestInstances() (selectedNode node.Node) {
	minCount := int(^uint(0) >> 1)
	counts := make(map[string]int)
	for _, srv := range GetContainers() {
		counts[srv.Node.Addr()]++
	}
	for _, n := range GetNodes() {
		c := counts[n.Addr()]
		if c < minCount {
			minCount = c
			selectedNode = n
		}
	}
	return
}
