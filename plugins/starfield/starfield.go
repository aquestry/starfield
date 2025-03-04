package starfield

import (
	"context"
	"fmt"
	"os"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/events"
	"starfield/plugins/starfield/records/node"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"gopkg.in/yaml.v2"
)

var log logr.Logger

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

func loadConfig() Config {
	filePath := "starfield.yml"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		defaultConfig := Config{
			Nodes: []NodeConfig{
				{Type: "local"},
				{Type: "externKey", Name: "Hetzner-1", IP: "129.92.351.42", User: "root", Port: 22, KeyPath: "path/to/key", PassPhrase: ""},
				{Type: "externPass", Name: "Hetzner-2", IP: "129.92.351.43", User: "root", Port: 22, Password: "1234"},
			},
		}
		data, _ := yaml.Marshal(defaultConfig)
		_ = os.WriteFile(filePath, data, 0644)
		os.Exit(0)
		return defaultConfig
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		return Config{}
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Println("Failed to parse config file:", err)
		return Config{}
	}
	return config
}

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
	config := loadConfig()
	for _, n := range config.Nodes {
		if n.Type == "local" {
			containers.GlobalManager.AddNode("local", node.NewLocalNode())
		} else if n.Type == "externKey" {
			remoteNode, err := node.NewRemoteNodeWithKey(n.Name, n.User, fmt.Sprintf("%s:%d", n.IP, n.Port), n.KeyPath, n.PassPhrase)
			if err == nil {
				containers.GlobalManager.AddNode(n.Name, remoteNode)
			}
		} else if n.Type == "externPass" {
			remoteNode, err := node.NewRemoteNodeWithPassword(n.Name, n.User, fmt.Sprintf("%s:%d", n.IP, n.Port), n.Password)
			if err == nil {
				containers.GlobalManager.AddNode(n.Name, remoteNode)
			}
		}
	}
	containers.CreateContainer("lobby", "anton691/simple-lobby:latest", 25566)
}
