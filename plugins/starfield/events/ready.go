package events

import (
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"starfield/plugins/starfield/config"
	"starfield/plugins/starfield/logger"
	"starfield/plugins/starfield/orch"
)

func Ready(e *proxy.ReadyEvent) {
	config.LoadConfig()
	orch.CreateLobby()
	logger.L.Info("event", "type", "ReadyEvent", "address", e.Addr())
}
