package starfield

import (
	"context"
	"github.com/aquestry/starfield/plugins/starfield/commands"
	"github.com/aquestry/starfield/plugins/starfield/events"
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch"
	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		logger.L = logr.FromContextOrDiscard(ctx)
		orch.ProxyInstance = p

		event.Subscribe(p.Event(), 0, events.ChooseInitial)
		event.Subscribe(p.Event(), 0, events.ShutdownEvent)
		event.Subscribe(p.Event(), 0, events.PluginMessage)
		event.Subscribe(p.Event(), 0, events.Ready)

		p.Command().Register(commands.ContainerCommand())

		return nil

	},
}
