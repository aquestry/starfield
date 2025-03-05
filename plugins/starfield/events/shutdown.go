package events

import (
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/logger"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ShutdownEvent(e *proxy.ShutdownEvent) {
	for _, c := range containers.GetContainers() {
		containers.DeleteContainer(c.Name)
	}
	logger.L.Info("event", "type", "ShutdownEvent")
}
