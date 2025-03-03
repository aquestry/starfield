package containers

import (
	"fmt"
	"net"

	"github.com/go-logr/logr"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var P proxy.Proxy
var Lobby proxy.RegisteredServer
var Log logr.Logger

func CreateContainer(name, template string, port int) {
	n, ok := GlobalManager.Nodes["local"]
	if !ok {
		Log.Error(fmt.Errorf("node 'local' not found"), "CreateContainer: local node not found")
		return
	}
	cmd := fmt.Sprintf("docker run --name %s -d -e PAPER_VELOCITY_SECRET=%s -p %d:25565 %s", name, P.Config().Forwarding.VelocitySecret, port, template)
	_, err := n.Run(cmd)
	if err != nil {
		Log.Error(err, "CreateContainer: docker command failed", "command", cmd)
		return
	}
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		Log.Error(err, "CreateContainer: failed to resolve TCP address", "port", port)
		return
	}
	serverInfo := proxy.NewServerInfo(name, addr)
	regServer, err := P.Register(serverInfo)
	if err != nil {
		Log.Error(err, "CreateContainer: failed to register server", "serverInfo", serverInfo)
		return
	}
	Lobby = regServer
	err = GlobalManager.AddServer(name, "local", Lobby)
	if err != nil {
		Log.Error(err, "CreateContainer: failed to add server to GlobalManager", "serverName", name)
	}
}

func DeleteContainer(name string) {
	n, ok := GlobalManager.Nodes["local"]
	if !ok {
		Log.Error(fmt.Errorf("node 'local' not found"), "CreateContainer: local node not found")
		return
	}
	cmd := fmt.Sprintf("docker rm %s --force", name)
	_, err := n.Run(cmd)
	if err != nil {
		Log.Error(err, "CreateContainer: docker command failed", "command", cmd)
		return
	}
}
