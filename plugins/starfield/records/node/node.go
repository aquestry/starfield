package node

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
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
	Addr() string
}

type LocalNode struct{}

type RemoteNode struct {
	FullAddr string
	Config   *ssh.ClientConfig
	client   *ssh.Client
}

func NewLocalNode() *LocalNode {
	return &LocalNode{}
}

func (ln *LocalNode) Addr() string {
	return "127.0.0.1"
}

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

func NewRemoteNodeWithPassword(user, fullAddr, password string) (*RemoteNode, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return &RemoteNode{FullAddr: fullAddr, Config: config}, nil
}

func NewRemoteNodeWithKey(user, fullAddr, keyPath, passphrase string) (*RemoteNode, error) {
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
	return &RemoteNode{FullAddr: fullAddr, Config: config}, nil
}

func (rn *RemoteNode) Addr() string {
	host, _, err := net.SplitHostPort(rn.FullAddr)
	if err != nil {
		return rn.FullAddr
	}
	return host
}

func (rn *RemoteNode) Port() string {
	_, port, err := net.SplitHostPort(rn.FullAddr)
	if err != nil {
		return ""
	}
	return port
}

func (rn *RemoteNode) Run(cmd string) (string, error) {
	Logger.Info("running remote", "addr", rn.FullAddr, "command", cmd)
	if err := rn.ensureConnected(); err != nil {
		Logger.Error(err, "remote connect error", "addr", rn.FullAddr)
		return "", err
	}
	session, err := rn.client.NewSession()
	if err != nil {
		Logger.Error(err, "remote session error", "addr", rn.FullAddr)
		return "", err
	}
	defer session.Close()
	var outBuf, errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf
	if err := session.Run(cmd); err != nil {
		out := outBuf.String() + errBuf.String()
		Logger.Error(err, "remote error", "addr", rn.FullAddr, "command", cmd, "output", out)
		return out, fmt.Errorf("remote: %w (output: %s)", err, out)
	}
	out := outBuf.String() + errBuf.String()
	result := strings.TrimSpace(out)
	Logger.Info("remote success", "command", cmd, "result", result)
	return result, nil
}

func (rn *RemoteNode) ensureConnected() error {
	if rn.client != nil {
		return nil
	}
	c, err := ssh.Dial("tcp", rn.FullAddr, rn.Config)
	if err != nil {
		return err
	}
	rn.client = c
	return nil
}

func (rn *RemoteNode) Close() error {
	if rn.client != nil {
		err := rn.client.Close()
		rn.client = nil
		return err
	}
	return nil
}
