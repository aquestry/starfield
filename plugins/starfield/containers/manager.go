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
	Servers map[string]ServerRecord
}

var GlobalManager = &NodeManager{
	Nodes:   make(map[string]node.Node),
	Servers: make(map[string]ServerRecord),
}

func (nm *NodeManager) AddNode(name string, n node.Node) {
	nm.Nodes[name] = n
}

func (nm *NodeManager) AddServer(serverName, nodeName string, info proxy.RegisteredServer) error {
	n, ok := nm.Nodes[nodeName]
	if !ok {
		return fmt.Errorf("node %s not found", nodeName)
	}
	nm.Servers[serverName] = ServerRecord{Name: serverName, Node: n, Info: info}
	return nil
}

func (nm *NodeManager) GetServerNode(serverName string) (node.Node, bool) {
	rec, ok := nm.Servers[serverName]
	return rec.Node, ok
}
