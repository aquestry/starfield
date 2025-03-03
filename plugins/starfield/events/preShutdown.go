package events

import (
	"starfield/plugins/starfield/containers"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func PreShutdownEvent(e *proxy.PreShutdownEvent) {
	Log.Info("Shutting Down test event logh  sdghiusdhguiopawsdhg9puasdgj9!")
	for _, server := range containers.GlobalManager.Servers {
		containers.DeleteContainer(server.Name)
	}
}
