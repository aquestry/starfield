package starfield

import (
	"context"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/records/node"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var log logr.Logger
var VelocitySecret string

type plugin struct {
	proxy *proxy.Proxy
}

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log = logr.FromContextOrDiscard(ctx)

		containers.Log = log
		node.Logger = log

		pl := &plugin{proxy: p}
		containers.P = *pl.proxy

		event.Subscribe(p.Event(), 0, pl.chooseInitial)
		event.Subscribe(p.Event(), 0, pl.init)

		return nil
	},
}

func (p *plugin) chooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	log.Info("Choose initial server", "player", player.Username())
	e.SetInitialServer(containers.Lobby)
}

func (p *plugin) init(e *proxy.ReadyEvent) {
	containers.GlobalManager.AddNode("local", node.NewLocalNode())
	containers.CreateContainer("lobby", "anton691/simple-lobby:latest", 25566)
}
