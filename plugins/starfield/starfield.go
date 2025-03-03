package starfield

import (
	"context"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/events"
	"starfield/plugins/starfield/records/node"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var log logr.Logger

type plugin struct {
	proxy *proxy.Proxy
}

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log = logr.FromContextOrDiscard(ctx)
		containers.Log = log
		events.Log = log
		node.Logger = log
		pl := &plugin{proxy: p}
		containers.P = *pl.proxy
		event.Subscribe(p.Event(), 0, events.ChooseInitial)
		event.Subscribe(p.Event(), 0, events.PreShutdownEvent)
		event.Subscribe(p.Event(), 0, pl.init)
		return nil
	},
}

func (p *plugin) init(e *proxy.ReadyEvent) {
	containers.GlobalManager.AddNode("local", node.NewLocalNode())
	containers.CreateContainer("lobby", "anton691/simple-lobby:latest", 25566)
}
