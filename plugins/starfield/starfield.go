package starfield

import (
	"context"
	"log"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log := logr.FromContextOrDiscard(ctx)
		log.Info("Hello from Stafield :)")
		event.Subscribe(p.Event(), 0, chooseInitial())
		return nil
	},
}

func chooseInitial() func(*proxy.PlayerChooseInitialServerEvent) {
	return func(e *proxy.PlayerChooseInitialServerEvent) {
		player := e.Player()
		log.Println("Choose initial server event for: ", player.Username())
	}
}
