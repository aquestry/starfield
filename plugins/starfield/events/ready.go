package events

import (
	"github.com/aquestry/starfield/plugins/starfield/config"
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Ready(e *proxy.ReadyEvent) {
	config.LoadConfig()
	orch.CreateLobby()
	logger.L.Info("event", "type", "ReadyEvent", "address", e.Addr())
}
