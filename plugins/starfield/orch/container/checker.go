package container

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/util"
)

func Update() {
	for _, c := range GetContainers() {
		motd, online := util.GetState(c.Info.ServerInfo().Addr().String())
		if c.Online != online {
			c.Online = online
			logger.L.Info("checker", "name", c.Name, "motd", motd, "online", online)
			if online {
				for _, player := range c.Pending {
					player.CreateConnectionRequest(c.Info).Connect(player.Context())
				}
			}
		}
	}
}
