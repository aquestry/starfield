package starfield

import (
	"context"
	"starfield/plugins/starfield/config"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/containers/node"
	"starfield/plugins/starfield/events"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var log logr.Logger

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log = logr.FromContextOrDiscard(ctx)

		containers.ProxyInstance = p
		events.Log = log
		node.Logger = log
		containers.Logger = log
		config.Logger = log

		event.Subscribe(p.Event(), 0, events.ChooseInitial)
		event.Subscribe(p.Event(), 0, events.PreShutdownEvent)
		event.Subscribe(p.Event(), 0, events.Ready)

		return nil
	},
}
