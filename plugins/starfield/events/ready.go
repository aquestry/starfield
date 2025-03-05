package events

import (
	"starfield/plugins/starfield/config"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/logger"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Ready(e *proxy.ReadyEvent) {
	config.LoadConfig()
	containers.CreateLobby()
	logger.L.Info("event", "type", "ReadyEvent")
}
