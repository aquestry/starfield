package events

import (
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/logger"
)

func ShutdownEvent(e *proxy.ShutdownEvent) {
	for _, c := range containers.GetContainers() {
		containers.DeleteContainer(c.Name)
	}
	logger.L.Info("event", "type", "ShutdownEvent")
}
