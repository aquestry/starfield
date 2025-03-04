package config

import (
	"fmt"
	"os"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/containers/node"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
)

var Logger logr.Logger

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

func LoadConfig() {
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
		Logger.Info("Default config was created! (pls configure & restart)")
		os.Exit(0)
		return
	}
	data, err := os.ReadFile(f)
	if err != nil {
		Logger.Info("config", "failed to create", "error", err)
		os.Exit(0)
		return
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		Logger.Info("config", "failed to parse file", "error", err)
		os.Exit(0)
		return
	}

	for _, n := range c.Nodes {
		switch n.Type {
		case "local":
			containers.GlobalContainers.AddNode("local", node.NewLocalNode())
		case "externKey":
			Logger.Info("config", "name", n.Name, "user", n.User, "IP", n.IP, "port", n.Port, "keypath", n.KeyPath, "passPhrase", n.PassPhrase)
			rn, err := node.NewRemoteNodeWithKey(n.Name, n.User, fmt.Sprintf("%s:%d", n.IP, n.Port), n.KeyPath, n.PassPhrase)
			if err != nil {
				Logger.Error(err, "Failed to create remote node with key", "node", n.Name)
				continue
			}
			containers.GlobalContainers.AddNode(n.Name, rn)
		case "externPass":
			rn, err := node.NewRemoteNodeWithPassword(n.Name, n.User, fmt.Sprintf("%s:%d", n.IP, n.Port), n.Password)
			if err != nil {
				Logger.Error(err, "Failed to create remote node with password", "node", n.Name)
				continue
			}
			containers.GlobalContainers.AddNode(n.Name, rn)
		}
	}
}
