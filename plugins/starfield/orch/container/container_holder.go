package container

import (
	"fmt"
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch/node"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"net"
	"time"
)

type Container struct {
	Name    string
	Tag     string
	IP      string
	Node    node.Node
	Info    proxy.RegisteredServer
	Port    int
	Online  bool
	Pending []proxy.Player
}

var clist []*Container

func RegisterContainer(name, tag, ip string, port int, n node.Node, start time.Time) (*Container, error) {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", n.Addr(), port))
	info := proxy.NewServerInfo(name, addr)
	server, err := ProxyInstance.Register(info)
	if err != nil {
		return nil, err
	}
	c := &Container{name, tag, ip, n, server, port, false, make([]proxy.Player, 0)}
	clist = append(clist, c)
	duration := time.Since(start)
	logger.L.Info("create", "name", name, "address", addr.String(), "time", duration)
	return c, nil
}

func Remove(name string) {
	for i, c := range clist {
		if c.Name == name {
			clist = append(clist[:i], clist[i+1:]...)
			ProxyInstance.Unregister(c.Info.ServerInfo())
			logger.L.Info("delete", "name", name)
			return
		}
	}
}

func GetContainers() []*Container {
	return clist
}

func GetContainer(name string) *Container {
	for _, c := range clist {
		if c.Name == name {
			return c
		}
	}
	return nil
}
