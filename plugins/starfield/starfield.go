package starfield

import (
	"context"
	"fmt"
	"os"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/containers/node"
	"starfield/plugins/starfield/events"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Nodes []NodeConfig `yaml:"nodes"`
}

type NodeConfig struct {
	Type       string `yaml:"type"`
	Name       string `yaml:"name,omitempty"`
	IP         string `yaml:"ip,omitempty"`
	User       string `yaml:"user,omitempty"`
	Port       int    `yaml:"port,omitempty"`
	KeyPath    string `yaml:"keyPath,omitempty"`
	PassPhrase string `yaml:"passPhrase,omitempty"`
	Password   string `yaml:"password,omitempty"`
}

var log logr.Logger

type plugin struct {
	p *proxy.Proxy
}

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log = logr.FromContextOrDiscard(ctx)
		pl := &plugin{p: p}

		containers.ProxyInstance = *pl.p
		events.Log = log
		node.Logger = log
		containers.Logger = log

		event.Subscribe(p.Event(), 0, events.ChooseInitial)
		event.Subscribe(p.Event(), 0, events.PreShutdownEvent)
		event.Subscribe(p.Event(), 0, events.Init)

		loadConfig()
		return nil
	},
}

func loadConfig() {
	f := "starfield.yml"
	if _, err := os.Stat(f); os.IsNotExist(err) {
		d := Config{
			Nodes: []NodeConfig{
				{Type: "local"},
				{Type: "externKey", Name: "Hetzner-1", IP: "129.92.351.42", User: "root", Port: 22, KeyPath: "path/to/key"},
				{Type: "externPass", Name: "Hetzner-2", IP: "129.92.351.43", User: "root", Port: 22, Password: "1234"},
			},
		}
		data, _ := yaml.Marshal(d)
		_ = os.WriteFile(f, data, 0644)
		log.Info("Default config was created! (pls configure & restart)")
		os.Exit(0)
	}
	data, err := os.ReadFile(f)
	if err != nil {
		log.Info("config", "failed to create", err)
		os.Exit(0)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Info("config", "failed to parse file", err)
		os.Exit(0)
	}
	for _, n := range c.Nodes {
		if n.Type == "local" {
			containers.GlobalContainers.AddNode("local", node.NewLocalNode())
		} else if n.Type == "externKey" {
			rn, err := node.NewRemoteNodeWithKey(n.Name, n.User, fmt.Sprintf("%s:%d", n.IP, n.Port), n.KeyPath, n.PassPhrase)
			if err == nil {
				containers.GlobalContainers.AddNode(n.Name, rn)
			}
		} else if n.Type == "externPass" {
			rn, err := node.NewRemoteNodeWithPassword(n.Name, n.User, fmt.Sprintf("%s:%d", n.IP, n.Port), n.Password)
			if err == nil {
				containers.GlobalContainers.AddNode(n.Name, rn)
			}
		}
	}
}
