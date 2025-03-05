package events

import (
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/logger"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ChooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	srv := containers.GetLobbyForPlayer(player)
	e.SetInitialServer(srv)
	logger.L.Info("event", "type", "PlayerChooseInitialServerEvent", "player", player.Username(), "result", srv.ServerInfo().Name())
}
