package events

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch/container"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ServerConnected(e *proxy.ServerConnectedEvent) {
	c := container.GetContainer(e.Server().ServerInfo().Name())
	if c != nil {
		for i, player := range c.Pending {
			if player == e.Player() {
				c.Pending = append(c.Pending[:i], c.Pending[i+1:]...)
			}
		}
	}
	logger.L.Info("event", "type", "ServerConnectedEvent")
}
