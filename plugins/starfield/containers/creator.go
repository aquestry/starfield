package containers

import (
	"fmt"
	"net"

	"github.com/go-logr/logr"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Log logr.Logger
var P proxy.Proxy
var Lobby proxy.RegisteredServer

func CreateContainer(name, template string, port int) {
	node, ok := GlobalManager.GetServerNode("local")
	if ok {
		node.Run(fmt.Sprintf("docker run -d -e PAPER_VELOCITY_SECRET=%s -p %d:25565 %s", P.Config().Forwarding.VelocitySecret, port, template))
		addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		serverInfo := proxy.NewServerInfo(name, addr)
		regServer, _ := P.Register(serverInfo)
		Lobby = regServer
		GlobalManager.AddServer(name, "local", Lobby)
	}
}
