package starfield

import (
	"context"
	"starfield/plugins/starfield/commands"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/events"
	"starfield/plugins/starfield/logger"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		logger.L = logr.FromContextOrDiscard(ctx)
		containers.ProxyInstance = p

		event.Subscribe(p.Event(), 0, events.ChooseInitial)
		event.Subscribe(p.Event(), 0, events.ShutdownEvent)
		event.Subscribe(p.Event(), 0, events.Ready)

		p.Command().Register(commands.ContainerCommand())

		return nil
	},
}
