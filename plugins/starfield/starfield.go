package starfield

import (
	"context"
	"starfield/plugins/starfield/containers"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Log logr.Logger

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		Log := logr.FromContextOrDiscard(ctx)
		Log.Info("Hello from Stafield :)")
		containers.CreateContainer()
		event.Subscribe(p.Event(), 0, chooseInitial())
		return nil
	},
}

func chooseInitial() func(*proxy.PlayerChooseInitialServerEvent) {
	return func(e *proxy.PlayerChooseInitialServerEvent) {
		player := e.Player()
		Log.Info("Choose initial server event for: ", player.Username())
	}
}
