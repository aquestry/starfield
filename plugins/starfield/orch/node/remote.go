package node

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type RemoteNode struct {
	name     string
	FullAddr string
	Config   *ssh.ClientConfig
	client   *ssh.Client
	freePort int
}

func NewRemoteNodeWithPassword(name, user, fullAddr, password string) (*RemoteNode, error) {
	c := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return &RemoteNode{name: name, FullAddr: fullAddr, Config: c}, nil
}

func NewRemoteNodeWithKey(name, user, fullAddr, keyPath, passphrase string) (*RemoteNode, error) {
	data, err := os.ReadFile(keyPath)
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
	c := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return &RemoteNode{name: name, FullAddr: fullAddr, Config: c}, nil
}

func (rn *RemoteNode) Run(args ...string) (string, error) {
	if err := rn.ensureConnected(); err != nil {
		return "", err
	}
	if len(args) == 0 {
		return "", fmt.Errorf("no command provided")
	}
	var cmd string
	if len(args) == 1 {
		cmd = args[0]
	} else {
		cmd = strings.Join(args, " ")
	}
	s, err := rn.client.NewSession()
	if err != nil {
		return "", err
	}
	defer s.Close()
	var outBuf, errBuf bytes.Buffer
	s.Stdout = &outBuf
	s.Stderr = &errBuf
	if err := s.Run(cmd); err != nil {
		o := outBuf.String() + errBuf.String()
		return o, fmt.Errorf("remote: %w (output: %s)", err, o)
	}
	return strings.TrimSpace(outBuf.String() + errBuf.String()), nil
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

func (rn *RemoteNode) Addr() string {
	host, _, err := net.SplitHostPort(rn.FullAddr)
	if err != nil {
		return rn.FullAddr
	}
	return host
}

func (rn *RemoteNode) Port() string {
	_, p, err := net.SplitHostPort(rn.FullAddr)
	if err != nil {
		return ""
	}
	return p
}

func (rn *RemoteNode) Name() string {
	return rn.name
}

func (rn *RemoteNode) GetFreePort() int {
	return rn.freePort
}

func (rn *RemoteNode) UpdateFreePort() {
	out, err := rn.Run(`python3 -c "import socket; s=socket.socket(); s.bind(('0.0.0.0', 0)); print(s.getsockname()[1]); s.close()"`)
	if err != nil {
		rn.freePort = 0
		return
	}
	p, e := strconv.Atoi(strings.TrimSpace(out))
	if e != nil {
		rn.freePort = 0
		return
	}
	rn.freePort = p
}

func (rn *RemoteNode) Close() error {
	if rn.client != nil {
		e := rn.client.Close()
		rn.client = nil
		return e
	}
	return nil
}
