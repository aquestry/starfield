package orch

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/util"
)

func Check() {
	for _, c := range GetContainers() {
		motd, online := util.GetState(c.Info.ServerInfo().Addr().String())
		logger.L.Info("checker", "container", c.Name, "motd", motd, "online", online)
	}
}
