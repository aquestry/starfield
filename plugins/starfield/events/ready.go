package events

import (
	"starfield/plugins/starfield/config"
	"starfield/plugins/starfield/containers"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Init(e *proxy.ReadyEvent) {
	config.LoadConfig()
	containers.CreateLobby()
}
