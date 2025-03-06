package events

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ChooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	s := orch.GetTargetLobby()
	e.SetInitialServer(s)
	logger.L.Info("event", "type", "PlayerChooseInitialServerEvent", "player", player.Username(), "result", s.ServerInfo().Name())
}
