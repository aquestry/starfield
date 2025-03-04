package events

import (
	"starfield/plugins/starfield/containers"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func PreShutdownEvent(e *proxy.PreShutdownEvent) {
	for _, c := range containers.GlobalContainers.Containers {
		containers.DeleteContainer(c.Name)
	}
	Log.Info("event", "type", "PlayerChooseInitialServerEvent")
}
