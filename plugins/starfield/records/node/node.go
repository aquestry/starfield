package node

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/crypto/ssh"
)

var Logger logr.Logger

type Node interface {
	Run(cmd string) (string, error)
}

type LocalNode struct{}

func NewLocalNode() *LocalNode { return &LocalNode{} }

func (ln *LocalNode) Run(cmd string) (string, error) {
	Logger.Info("running local", "command", cmd)
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", cmd)
	} else {
		c = exec.Command("sh", "-c", cmd)
	}
	out, err := c.CombinedOutput()
	if err != nil {
		Logger.Error(err, "local error", "command", cmd, "output", string(out))
		return string(out), fmt.Errorf("local: %w (output: %s)", err, out)
	}
	result := strings.TrimSpace(string(out))
	Logger.Info("local success", "command", cmd, "result", result)
	return result, nil
}

type RemoteNode struct {
	Addr   string
	Config *ssh.ClientConfig
	client *ssh.Client
}

func NewRemoteNodeWithPassword(user, addr, password string) (*RemoteNode, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return &RemoteNode{Addr: addr, Config: config}, nil
}

func NewRemoteNodeWithKey(user, addr, keyPath, passphrase string) (*RemoteNode, error) {
	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	var signer ssh.Signer
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(data, []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey(data)
	}
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return &RemoteNode{Addr: addr, Config: config}, nil
}

func (rn *RemoteNode) ensureConnected() error {
	if rn.client != nil {
		return nil
	}
	c, err := ssh.Dial("tcp", rn.Addr, rn.Config)
	if err != nil {
		return err
	}
	rn.client = c
	return nil
}

func (rn *RemoteNode) Run(cmd string) (string, error) {
	Logger.Info("running remote", "addr", rn.Addr, "command", cmd)
	if err := rn.ensureConnected(); err != nil {
		Logger.Error(err, "remote connect error", "addr", rn.Addr)
		return "", err
	}
	session, err := rn.client.NewSession()
	if err != nil {
		Logger.Error(err, "remote session error", "addr", rn.Addr)
		return "", err
	}
	defer session.Close()
	var outBuf, errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf
	if err := session.Run(cmd); err != nil {
		out := outBuf.String() + errBuf.String()
		Logger.Error(err, "remote error", "addr", rn.Addr, "command", cmd, "output", out)
		return out, fmt.Errorf("remote: %w (output: %s)", err, out)
	}
	out := outBuf.String() + errBuf.String()
	result := strings.TrimSpace(string(out))
	Logger.Info("local success", "command", cmd, "result", result)
	return result, nil
}

func (rn *RemoteNode) Close() error {
	if rn.client != nil {
		err := rn.client.Close()
		rn.client = nil
		return err
	}
	return nil
}
