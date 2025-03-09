package events

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch/container"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ShutdownEvent(_ *proxy.ShutdownEvent) {
	for _, c := range container.GetContainers() {
		container.DeleteContainer(c)
	}
	logger.L.Info("event", "type", "ShutdownEvent")
}
