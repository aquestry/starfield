package container

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/util"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var fallback proxy.RegisteredServer

func CreateLobby() {
	name := "lobby-" + util.RandomString()
	container, err := CreateContainer(name, "lobby", "anton691/simple-lobby")
	if err != nil {
		logger.L.Error(err, "error creating lobby", "name", name)
	} else {
		fallback = container.Info
	}
}

func GetLobby() proxy.RegisteredServer {
	m := int(^uint(0) >> 1)
	result := fallback
	for _, c := range GetContainers() {
		if c.Tag == "lobby" && c.Online {
			if cnt := c.Info.Players().Len(); cnt < m {
				m = cnt
				result = c.Info
			}
		}
	}
	return result
}
