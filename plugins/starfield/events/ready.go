package events

import (
	"starfield/plugins/starfield/config"
	"starfield/plugins/starfield/containers"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Ready(e *proxy.ReadyEvent) {
	config.LoadConfig()
	containers.CreateLobby()
	Log.Info("event", "type", "ReadyEvent")
}
