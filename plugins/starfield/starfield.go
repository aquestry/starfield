package starfield

import (
	"context"
	"net"
	"os"
	"starfield/plugins/starfield/containers"
	"starfield/plugins/starfield/records/node"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"gopkg.in/yaml.v2"
)

var log logr.Logger
var pl *plugin
var lobby proxy.RegisteredServer

type plugin struct {
	proxy *proxy.Proxy
}

var Plugin = proxy.Plugin{
	Name: "Starfield",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log := logr.FromContextOrDiscard(ctx)
		containers.Log = log
		node.Logger = log

		log.Info("Hello from Stafield :)")

		log.Info(readSecret("config.yml"))

		containers.CreateContainer()

		pl := &plugin{proxy: p}
		event.Subscribe(p.Event(), 0, pl.chooseInitial)

		addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:25566")
		serverInfo := proxy.NewServerInfo("lobby", addr)
		regServer, _ := pl.proxy.Register(serverInfo)
		lobby = regServer

		return nil
	},
}

func readSecret(path string) string {
	// Log attempt to read config file
	log.Info("Attempting to read config file from: %s", path)

	data, err := os.ReadFile(path)
	if err != nil {
		log.Info("Error reading file: %v", err)
		return ""
	}

	log.Info("File content length: %d bytes", len(data))

	// Try to unmarshal as a full structure
	var rawConfig map[string]interface{}
	err = yaml.Unmarshal(data, &rawConfig)
	if err != nil {
		log.Info("Error parsing YAML: %v", err)
		return ""
	}

	// Log the raw structure to help debug
	log.Info("Parsed YAML top level keys: %v", getMapKeys(rawConfig))

	// Check if "config" key exists
	configMap, ok := rawConfig["config"].(map[string]interface{})
	if !ok {
		log.Info("Missing 'config' key or not a map")
		// Try without the top-level "config" key
		configMap = rawConfig
	}

	log.Info("Config map keys: %v", getMapKeys(configMap))

	// Check if "forwarding" key exists
	forwardingMap, ok := configMap["forwarding"].(map[string]interface{})
	if !ok {
		log.Info("Missing 'forwarding' key or not a map")
		return ""
	}

	log.Info("Forwarding map keys: %v", getMapKeys(forwardingMap))

	// Check if "velocitySecret" key exists
	secret, ok := forwardingMap["velocitySecret"].(string)
	if !ok {
		log.Info("Missing 'velocitySecret' key or not a string")
		return ""
	}

	log.Info("Successfully found velocitySecret")
	return secret
}

// Helper function to get map keys for logging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (p *plugin) chooseInitial(e *proxy.PlayerChooseInitialServerEvent) {
	player := e.Player()
	log.Info("Choose initial server event for: ", player.Username())
	e.SetInitialServer(lobby)
}
