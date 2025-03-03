package containers

import (
	"starfield/plugins/starfield/records/node"

	"github.com/go-logr/logr"
)

var Log logr.Logger

func CreateContainer() {
	ln := node.NewLocalNode()
	ln.Run("docker ps -a")
}
