package events

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch/container"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func PlayerChooseInitialServer(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	s := container.GetLobby()
	e.SetInitialServer(s)
	logger.L.Info("event", "type", "PlayerChooseInitialServerEvent", "player", player.Username(), "result", s.ServerInfo().Name())
}
