package commands

import (
	"fmt"

	"starfield/plugins/starfield/orch"

	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ContainerCommand() brigodier.LiteralNodeBuilder {
	listContainers := command.Command(func(ctx *command.Context) error {
		containerList := orch.GetContainers()
		if len(containerList) == 0 {
			return ctx.SendMessage(&component.Text{
				Content: "No orch are currently registered.",
			})
		}
		msg := "Available orch:\n"
		for _, container := range containerList {
			msg += fmt.Sprintf("- %s (%s:%d)\n", container.Name, container.IP, container.Port)
		}
		return ctx.SendMessage(&component.Text{Content: msg})
	})

	connectContainer := command.Command(func(ctx *command.Context) error {
		player, ok := ctx.Source.(proxy.Player)
		if !ok {
			return ctx.SendMessage(&component.Text{
				Content: "Only players can use this command.",
			})
		}
		containerName := ctx.String("name")
		container, err := orch.GetContainer(containerName)
		if err != nil {
			return ctx.SendMessage(&component.Text{
				Content: fmt.Sprintf("Container '%s' not found.", containerName),
			})
		}
		player.CreateConnectionRequest(container.Info).Connect(ctx)
		return ctx.SendMessage(&component.Text{
			Content: fmt.Sprintf("Connecting to container '%s'...", containerName),
		})
	})

	return brigodier.Literal("container").
		Then(
			brigodier.Literal("list").Executes(listContainers),
		).Then(
		brigodier.Literal("connect").
			Then(brigodier.Argument("name", brigodier.StringWord).Executes(connectContainer)),
	)
}
