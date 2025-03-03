package events

import (
	"reflect"
	"starfield/plugins/starfield/containers"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ChooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	Log.Info("event", "type", reflect.TypeOf(e).Name(), "player", player.Username(), "result", "lobby")
	e.SetInitialServer(containers.Lobby)
}
