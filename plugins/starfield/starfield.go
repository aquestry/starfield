package starfield

import (
	"context"
	"errors"
	"github.com/aquestry/starfield/plugins/starfield/commands"
	"github.com/aquestry/starfield/plugins/starfield/events"
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch/container"
	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/brigodier"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/edition/java/proxy/message"
	"time"
)

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		logger.L = logr.FromContextOrDiscard(ctx)
		container.ProxyInstance = p

		channel, _ := message.ChannelIdentifierFrom("nebula:main")
		p.ChannelRegistrar().Register(channel)

		event.Subscribe(p.Event(), 0, events.PlayerChooseInitialServer)
		event.Subscribe(p.Event(), 0, events.ServerConnected)
		event.Subscribe(p.Event(), 0, events.PluginMessage)
		event.Subscribe(p.Event(), 0, events.Shutdown)
		event.Subscribe(p.Event(), 0, events.Ready)

		// Update online state of the containers
		go func() {
			for {
				start := time.Now()
				container.Update()
				time.Sleep(time.Until(start.Add(1 * time.Second)))
			}
		}()

		p.Command().Register(commands.ContainerCommand())
		p.Command().Register(commands.LobbyCommand())
		brigodier.ErrDispatcherUnknownArgument = errors.New("Incorrect argument for that command.")

		return nil
	},
}
