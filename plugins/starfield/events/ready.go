package events

import (
	"starfield/plugins/starfield/containers"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Init(e *proxy.ReadyEvent) {
	containers.CreateLobby()
}
