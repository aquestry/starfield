package main

import (
	"github.com/aquestry/starfield/plugins/starfield"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(proxy.Plugins,
		starfield.Plugin,
	)
	gate.Execute()
}
