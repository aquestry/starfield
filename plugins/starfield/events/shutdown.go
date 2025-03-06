package events

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ShutdownEvent(_ *proxy.ShutdownEvent) {
	for _, c := range orch.GetContainers() {
		orch.DeleteContainer(c.Name)
	}
	logger.L.Info("event", "type", "ShutdownEvent")
}
