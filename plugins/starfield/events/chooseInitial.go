package events

import (
	"starfield/plugins/starfield/containers"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ChooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	Log.Info("event", "type", "PlayerChooseInitialServerEvent", "player", player.Username(), "result", "lobby")
	e.SetInitialServer(containers.Lobby)
}
