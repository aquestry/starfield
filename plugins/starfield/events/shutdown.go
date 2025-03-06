package events

import (
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"starfield/plugins/starfield/logger"
	"starfield/plugins/starfield/orch"
)

func ShutdownEvent(e *proxy.ShutdownEvent) {
	for _, c := range orch.GetContainers() {
		orch.DeleteContainer(c.Name)
	}
	logger.L.Info("event", "type", "ShutdownEvent")
}
