package events

import (
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/logger"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ChooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	s := containers.GetTargetLobby()
	e.SetInitialServer(s)
	logger.L.Info("event", "type", "PlayerChooseInitialServerEvent", "player", player.Username(), "result", s.ServerInfo().Name())
}
