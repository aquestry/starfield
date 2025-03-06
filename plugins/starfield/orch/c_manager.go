package orch

import (
	"fmt"
	"net"
	"starfield/plugins/starfield/logger"
	"starfield/plugins/starfield/orch/node"
	"strconv"
	"strings"
	"time"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type Container struct {
	Name string
	Tag  string
	IP   string
	Node node.Node
	Info proxy.RegisteredServer
	Port int
}

var clist []Container

func RegisterContainer(name, tag, ip string, n node.Node, start time.Time) (Container, error) {
	p, e := n.Run("docker", "port", name, "25565")
	if e != nil {
		return Container{}, e
	}
	parts := strings.SplitN(p, ":", 2)
	if len(parts) < 2 {
		return Container{}, e
	}
	port := parts[1]
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", n.Addr(), port))
	info := proxy.NewServerInfo(name, addr)
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return Container{}, e
	}
	server, err := ProxyInstance.Register(info)
	if err != nil {
		return Container{}, e
	}
	c := Container{name, tag, ip, n, server, portNumber}
	clist = append(clist, c)
	duration := time.Since(start)
	logger.L.Info("create", "type", "container", "time", duration)
	return c, nil
}

func Remove(name string) {
	for i, c := range clist {
		if c.Name == name {
			clist = append(clist[:i], clist[i+1:]...)
			return
		}
	}
}

func GetContainers() []Container {
	return clist
}

func GetContainer(name string) (Container, error) {
	for _, a := range clist {
		if a.Name == name {
			return a, nil
		}
	}
	return Container{}, fmt.Errorf("not found")
}
