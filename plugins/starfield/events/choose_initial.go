package events

import (
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"starfield/plugins/starfield/logger"
	"starfield/plugins/starfield/orch"
)

func ChooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	s := orch.GetTargetLobby()
	e.SetInitialServer(s)
	logger.L.Info("event", "type", "PlayerChooseInitialServerEvent", "player", player.Username(), "result", s.ServerInfo().Name())
}
