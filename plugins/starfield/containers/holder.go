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
		return fmt.Errorf("node %q not found", nodeName)
	}
	nm.Servers[serverName] = ServerRecord{Name: serverName, Node: n, Info: info}
	return nil
}

func (nm *NodeManager) RemoveServer(serverName string) error {
	if _, exists := nm.Servers[serverName]; !exists {
		return fmt.Errorf("server %s not found", serverName)
	}
	delete(nm.Servers, serverName)
	return nil
}

func (nm *NodeManager) GetServer(serverName string) (ServerRecord, error) {
	server, exists := nm.Servers[serverName]
	if !exists {
		return ServerRecord{}, fmt.Errorf("server %s not found", serverName)
	}
	return server, nil
}

func (nm *NodeManager) GetServers() []ServerRecord {
	servers := make([]ServerRecord, 0, len(nm.Servers))
	for _, rec := range nm.Servers {
		servers = append(servers, rec)
	}
	return servers
}
