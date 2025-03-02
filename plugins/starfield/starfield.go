package starfield

import (
	"context"
	"starfield/plugins/starfield/containers"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var log logr.Logger

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log := logr.FromContextOrDiscard(ctx)
		containers.Log = log
		log.Info("Hello from Stafield :)")
		event.Subscribe(p.Event(), 0, chooseInitial())
		containers.CreateContainer()
		return nil
	},
}

func chooseInitial() func(*proxy.PlayerChooseInitialServerEvent) {
	return func(e *proxy.PlayerChooseInitialServerEvent) {
		player := e.Player()
		log.Info("Choose initial server event for: ", player.Username())
	}
}
