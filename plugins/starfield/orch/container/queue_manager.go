package container

import (
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"strings"
)

func QueueRequest(message string) {
	if strings.HasPrefix(message, "queue:") {
		request := strings.Replace(message, "queue:", "", 1)
		parts := strings.Split(request, ":")
		if len(parts) == 2 {
			player := ProxyInstance.PlayerByName(parts[0])
			queue := parts[1]
			if player != nil && queue != "" {
				AddToQueue(player, queue)
			}
		}
	}
}

func AddToQueue(player proxy.Player, queue string) {
	logger.L.Info("queue", "type", "add", "player", player.Username(), "queue", queue)
	// finish
}
