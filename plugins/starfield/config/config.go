package config

import (
	"fmt"
	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
	"os"
	"starfield/plugins/starfield/orch"
	"starfield/plugins/starfield/orch/node"
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
	count := 0
	local := false
	for _, n := range c.Nodes {
		switch n.Type {
		case "local":
			count++
			if !local {
				local = true
				orch.RegisterNode(node.NewLocalNode())
			}
		case "externKey":
			if n.Name == "local" {
				Logger.Info("config", "error", "Name of node can't be 'local'!")
				os.Exit(0)
			}
			rn, err := node.NewRemoteNodeWithKey(n.Name, n.User, fmt.Sprintf("%s:%d", n.IP, n.Port), n.KeyPath, n.PassPhrase)
			if err != nil {
				Logger.Error(err, "Failed to create remote node with key", "node", n.Name)
				continue
			}
			count++
			orch.RegisterNode(rn)
		case "externPass":
			if n.Name == "local" {
				Logger.Info("config", "error", "Name of node can't be 'local'!")
				os.Exit(0)
			}
			rn, err := node.NewRemoteNodeWithPassword(n.Name, n.User, fmt.Sprintf("%s:%d", n.IP, n.Port), n.Password)
			if err != nil {
				Logger.Error(err, "Failed to create remote node with password", "node", n.Name)
				continue
			}
			count++
			orch.RegisterNode(rn)
		}
	}
	Logger.Info("config", "nodes", count)
	Logger.Info("config", "local", local)
	if count == 0 {
		Logger.Info("config", "error", "You need atleast one node!")
		os.Exit(0)
	}
}
