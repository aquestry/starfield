package node

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type LocalNode struct {
	freePort int
}

func NewLocalNode() *LocalNode {
	return &LocalNode{}
}

func (ln *LocalNode) Run(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no command provided")
	}
	var cmd *exec.Cmd
	if len(args) == 1 && runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", args[0])
	} else if len(args) == 1 {
		cmd = exec.Command("sh", "-c", args[0])
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("local: %w (output: %s)", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func (ln *LocalNode) Addr() string {
	return "127.0.0.1"
}

func (ln *LocalNode) Port() string {
	return ""
}

func (ln *LocalNode) Name() string {
	return "local"
}

func (ln *LocalNode) GetFreePort() (int, error) {
	if ln.freePort == 0 {
		out, err := ln.Run("python3", "-c", "import socket; s=socket.socket(); s.bind(('0.0.0.0', 0)); print(s.getsockname()[1]); s.close()")
		if err != nil {
			return 0, err
		}
		p, e := strconv.Atoi(strings.TrimSpace(out))
		if e != nil {
			return 0, e
		}
		ln.freePort = p
	}
	return ln.freePort, nil
}

func (ln *LocalNode) UpdateFreePort() error {
	ln.freePort = 0
	_, err := ln.GetFreePort()
	return err
}

func (ln *LocalNode) Close() error {
	return nil
}
