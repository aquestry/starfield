package events

import (
	"starfield/plugins/starfield/config"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/logger"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func PluginMessage(e *proxy.PluginMessageEvent) {
	config.LoadConfig()
	containers.CreateLobby()
	logger.L.Info("event", "type", "PluginMessageEvent")
}
