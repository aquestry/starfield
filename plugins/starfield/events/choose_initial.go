package events

import (
	"starfield/plugins/starfield/containers"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ChooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	srv := containers.GetLobbyForPlayer(player)
	e.SetInitialServer(srv)
	Log.Info("event", "type", "PlayerChooseInitialServerEvent", "player", player.Username(), "result", srv.ServerInfo().Name())
}
