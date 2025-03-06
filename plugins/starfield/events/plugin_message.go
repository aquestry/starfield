package events

import (
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"starfield/plugins/starfield/logger"
)

func PluginMessage(e *proxy.PluginMessageEvent) {
	logger.L.Info("event", "type", "PluginMessageEvent", "identifier", e.Identifier())
}
