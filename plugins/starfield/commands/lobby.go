package commands

import (
	"github.com/aquestry/starfield/plugins/starfield/events"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func LobbyCommand() brigodier.LiteralNodeBuilder {
	lobby := command.Command(func(ctx *command.Context) error {
		player, ok := ctx.Source.(proxy.Player)
		if !ok {
			return ctx.SendMessage(&component.Text{
				Content: "Only players can use this command.",
			})
		}
		msg := "Connecting you to a lobby..."
		fine := events.Lobby(player)
		if !fine {
			msg = "You are already on a lobby."
		}
		return ctx.SendMessage(&component.Text{Content: msg})
	})
	return brigodier.Literal("lobby").Executes(lobby)
}
