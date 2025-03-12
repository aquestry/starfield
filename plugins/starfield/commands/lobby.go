package commands

import (
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func LobbyCommand() brigodier.LiteralNodeBuilder {
	lobby := command.Command(func(ctx *command.Context) error {
		_, ok := ctx.Source.(proxy.Player)
		if !ok {
			return ctx.SendMessage(&component.Text{
				Content: "Only players can use this command.",
			})
		}
		msg := "Connecting you to a lobby..."
		return ctx.SendMessage(&component.Text{Content: msg})
	})
	return brigodier.Literal("lobby").Executes(lobby)
}
