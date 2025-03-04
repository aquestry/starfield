package containers

import (
	"fmt"
	"net"
	"starfield/plugins/starfield/containers/node"

	"github.com/go-logr/logr"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var ProxyInstance proxy.Proxy
var Logger logr.Logger

func CreateContainer(name, tag, template string) {
	n := getNodeWithLowestInstances()
	port, _ := n.GetFreePort()
	n.UpdateFreePort()
	cmd := fmt.Sprintf("docker run --name %s -d -e PAPER_VELOCITY_SECRET=%s -p %d:25565 %s",
		name, ProxyInstance.Config().Forwarding.VelocitySecret, port, template)
	_, err := n.Run(cmd)
	if err != nil {
		Logger.Error(err, "CreateContainer failed", "command", cmd)
		return
	}
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", n.Addr(), port))
	if err != nil {
		Logger.Error(err, "ResolveTCPAddr failed", "port", port)
		return
	}
	info := proxy.NewServerInfo(name, addr)
	reg, err := ProxyInstance.Register(info)
	if err != nil {
		Logger.Error(err, "Register failed", "serverInfo", info)
		return
	}
	err = GlobalContainers.AddContainer(name, tag, n.Name(), reg)
	if err != nil {
		Logger.Error(err, "AddContainer failed", "containerName", name)
	}
}

func DeleteContainer(name string) {
	srv, err := GlobalContainers.GetContainer(name)
	if err != nil {
		Logger.Error(err, "container not found")
		return
	}
	cmd := fmt.Sprintf("docker rm %s --force", name)
	_, e := srv.Node.Run(cmd)
	if e != nil {
		Logger.Error(e, "DeleteContainer failed", "command", cmd)
	}
}

func getNodeWithLowestInstances() (selectedNode node.Node) {
	minCount := int(^uint(0) >> 1)
	counts := make(map[string]int)
	for _, srv := range GlobalContainers.Containers {
		counts[srv.Node.Addr()]++
	}
	for _, n := range GlobalContainers.Nodes {
		c := counts[n.Addr()]
		if c < minCount {
			minCount = c
			selectedNode = n
		}
	}
	return
}
