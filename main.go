package main

import (
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"starfield/plugins/starfield"
)

func main() {
	proxy.Plugins = append(proxy.Plugins,
		starfield.Plugin,
	)
	gate.Execute()
}
