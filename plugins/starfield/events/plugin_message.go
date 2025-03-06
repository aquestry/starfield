package events

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func PluginMessage(e *proxy.PluginMessageEvent) {
	logger.L.Info("event", "type", "PluginMessageEvent", "identifier", e.Identifier())
}
