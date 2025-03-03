package containers

import (
	"starfield/plugins/starfield"
	"starfield/plugins/starfield/records/node"
)

func CreateContainer() {
	starfield.Log.Info("test")
	ln := node.NewLocalNode()
	out, err := ln.Run("echo 'Hello from local'")
	if err != nil {
		starfield.Log.Info("manager: local error: %v\n", err)
	} else {
		starfield.Log.Info("manager:Local: %s\n", out)
	}
}
