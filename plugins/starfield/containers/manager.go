package containers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-logr/logr"
)

var Log logr.Logger

func CreateContainer() {
	Log.Info("Creating new container...")
	cmd := exec.Command("docker", "ps")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("could not run command: ", err)
	}
}
