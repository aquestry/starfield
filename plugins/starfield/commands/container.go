package commands

import (
	"fmt"
	"github.com/aquestry/starfield/plugins/starfield/orch/container"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func ContainerCommand() brigodier.LiteralNodeBuilder {
	listContainers := command.Command(func(ctx *command.Context) error {
		containerList := container.GetContainers()
		msg := "Containers:\n"
		for _, c := range containerList {
			msg += fmt.Sprintf("- %s (%s:%d)\n", c.Name, c.IP, c.Port)
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
		c := container.GetContainer(containerName)
		if c != nil {
			return ctx.SendMessage(&component.Text{
				Content: fmt.Sprintf("Container '%s' not found.", containerName),
			})
		}
		here := false
		c.Info.Players().Range(func(p proxy.Player) bool {
			here = true
			return true
		})
		if here {
			return ctx.SendMessage(&component.Text{
				Content: fmt.Sprintf("You are already on '%s'.", containerName),
			})
		}
		player.CreateConnectionRequest(c.Info).Connect(ctx)
		return ctx.SendMessage(&component.Text{
			Content: fmt.Sprintf("Connecting to c '%s'...", containerName),
		})
	})

	return brigodier.Literal("container").
		Then(
			brigodier.Literal("list").Executes(listContainers),
		).Then(
		brigodier.Literal("connect").
			Then(brigodier.Argument("name", brigodier.StringWord).Executes(connectContainer).Suggests(containers())),
	)
}

func containers() brigodier.SuggestionProvider {
	return command.SuggestFunc(func(ctx *command.Context, builder *brigodier.SuggestionsBuilder) *brigodier.Suggestions {
		for _, c := range container.GetContainers() {
			here := false
			c.Info.Players().Range(func(p proxy.Player) bool {
				here = true
				return true
			})
			if !here {
				builder.Suggest(c.Name)
			}
		}
		return builder.Build()
	})
}
