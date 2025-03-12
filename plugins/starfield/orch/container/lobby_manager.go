package container

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/util"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var lobby proxy.RegisteredServer

func CreateLobby() {
	name := "lobby-" + util.RandomString()
	container, err := CreateContainer(name, "lobby", "anton691/simple-lobby")
	if err != nil {
		logger.L.Error(err, "error creating lobby")
	} else {
		lobby = container.Info
	}
}

func GetLobby() proxy.RegisteredServer {
	return lobby
}
