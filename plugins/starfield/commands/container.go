package commands

import (
	"starfield/plugins/starfield/containers"
	"strconv"

	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
)

func ContainerCommand() brigodier.LiteralNodeBuilder {
	listContainers := command.Command(func(ctx *command.Context) error {
		containerList := containers.GetContainers()
		if len(containerList) == 0 {
			message := "No containers are currently registered."
			ctx.Source.SendMessage(&component.Text{Content: message})
			return nil
		}
		message := "Available containers:\n"
		for _, container := range containerList {
			message += "- " + container.Name + " (" + container.IP + ":" + strconv.Itoa(container.Port) + ")\n"
		}
		ctx.Source.SendMessage(&component.Text{Content: message})
		return nil
	})
	return brigodier.Literal("container").
		Then(brigodier.Literal("list").Executes(listContainers))
}
