package containers

import (
	"fmt"
	"starfield/plugins/starfield/records/node"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type ServerRecord struct {
	Name string
	Node node.Node
	Info proxy.RegisteredServer
}

type NodeManager struct {
	Nodes   map[string]node.Node
	Servers []ServerRecord
}

var GlobalManager = &NodeManager{
	Nodes:   make(map[string]node.Node),
	Servers: []ServerRecord{},
}

func (nm *NodeManager) AddNode(name string, n node.Node) {
	nm.Nodes[name] = n
}

func (nm *NodeManager) AddServer(serverName, nodeName string, info proxy.RegisteredServer) error {
	n, ok := nm.Nodes[nodeName]
	if !ok {
		return fmt.Errorf("node %s not found", nodeName)
	}
	nm.Servers = append(nm.Servers, ServerRecord{Name: serverName, Node: n, Info: info})
	return nil
}

func (nm *NodeManager) RemoveServer(serverName string) error {
	index := -1
	for i, rec := range nm.Servers {
		if rec.Name == serverName {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("server %s not found", serverName)
	}
	nm.Servers = append(nm.Servers[:index], nm.Servers[index+1:]...)
	return nil
}

func (nm *NodeManager) GetServers() []ServerRecord {
	return nm.Servers
}
