package events

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch/container"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func PluginMessage(e *proxy.PluginMessageEvent) {
	message := string(e.Data())
	container.QueueRequest(message)
	logger.L.Info("event", "type", "PluginMessageEvent", "identifier", e.Identifier().ID(), "message", message)
}
