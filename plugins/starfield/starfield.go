package starfield

import (
	"context"
	"github.com/aquestry/starfield/plugins/starfield/events"
	"github.com/aquestry/starfield/plugins/starfield/logger"
	"github.com/aquestry/starfield/plugins/starfield/orch"
	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/edition/java/proxy/message"
	"time"
)

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		logger.L = logr.FromContextOrDiscard(ctx)
		orch.ProxyInstance = p
		channel, _ := message.ChannelIdentifierFrom("nebula:main")
		p.ChannelRegistrar().Register(channel)
		event.Subscribe(p.Event(), 0, events.ChooseInitial)
		event.Subscribe(p.Event(), 0, events.ShutdownEvent)
		event.Subscribe(p.Event(), 0, events.PluginMessage)
		event.Subscribe(p.Event(), 0, events.Ready)

		go func() {
			for {
				start := time.Now()
				orch.Check()
				time.Sleep(time.Until(start.Add(1 * time.Second)))
			}
		}()

		return nil
	},
}
