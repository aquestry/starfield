package starfield

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"starfield/plugins/starfield/commands"
	"starfield/plugins/starfield/events"
	"starfield/plugins/starfield/logger"
	"starfield/plugins/starfield/orch"
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
