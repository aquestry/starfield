package events

import (
	"fmt"
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch/container"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"strings"
)

func PluginMessage(e *proxy.PluginMessageEvent) {
	message := string(e.Data())
	logger.L.Info("event", "type", "PluginMessageEvent", "identifier", e.Identifier().ID(), "message", message)
	if strings.HasPrefix(message, "queue") {
		parts := strings.Split(message, ":")
		player := container.ProxyInstance.PlayerByName(parts[1])
		if player != nil {
			c, _ := container.CreateContainer(parts[1]+"-parkour", "parkour", "anton691/simple-parkour")
			c.Pending = append(c.Pending, player)
			fmt.Println("added player to pending:", player.Username())
		}
	}
	if strings.HasPrefix(message, "lobby") {
		parts := strings.Split(message, ":")
		player := container.ProxyInstance.PlayerByName(parts[1])
		if player != nil {
			Lobby(player)
		}
	}
}

func Lobby(player proxy.Player) bool {
	c, _ := container.GetContainer(player.CurrentServer().Server().ServerInfo().Name())
	if c.Tag != "lobby" {
		player.CreateConnectionRequest(container.GetTargetLobby()).Connect(player.Context())
		return true
	} else {

		player.SendMessage(&component.Text{
			Content: "You are already in a lobby."})
	}
	return false
}
