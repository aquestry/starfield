package container

import (
	"fmt"
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch/node"
	"github.com/aquestry/starfield/plugins/starfield/util"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"time"
)

var ProxyInstance *proxy.Proxy

func CreateContainer(name, tag, template string) (*Container, error) {
	start := time.Now()
	n := getNodeWithLowestInstances()
	precmd := fmt.Sprintf("docker run --name %s -d -p 0:25565 anton691/simple-lobby", name)
	_, err := n.Run(precmd)
	if err != nil {
		logger.L.Error(err, "pre create failed", "command", precmd)
		return &Container{}, err
	}
	port, err := util.GetPort(n, name)
	if err != nil {
		logger.L.Error(err, "get port failed")
	}
	c := fmt.Sprintf("docker rm %s --force", name)
	_, e := n.Run(c)
	if e != nil {
		logger.L.Error(err, "delete failed", "command", c)
	}
	cmd := fmt.Sprintf("docker run --name %s -d -e PAPER_VELOCITY_SECRET=%s --restart unless-stopped -p %d:25565 %s", name, ProxyInstance.Config().Forwarding.VelocitySecret, port, template)
	_, err = n.Run(cmd)
	if err != nil {
		logger.L.Error(err, "pre create failed", "command", cmd)
		return &Container{}, err
	}
	return RegisterContainer(name, tag, n.Addr(), port, n, start)
}

func DeleteContainer(c *Container) {
	cmd := fmt.Sprintf("docker rm %s --force", c.Name)
	_, e := c.Node.Run(cmd)
	if e != nil {
		logger.L.Error(e, "delete failed", "command", cmd)
	}
	UnregisterContainer(c.Name)
}

func getNodeWithLowestInstances() node.Node {
	minCount := int(^uint(0) >> 1)
	counts := make(map[string]int)
	for _, srv := range GetContainers() {
		counts[srv.Node.Addr()]++
	}
	var selectedNode node.Node
	for _, n := range node.GetNodes() {
		c := counts[n.Addr()]
		if c < minCount {
			minCount = c
			selectedNode = n
		}
	}
	return selectedNode
}
