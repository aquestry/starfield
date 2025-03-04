package node

import "github.com/go-logr/logr"

var Logger logr.Logger

type Node interface {
	Run(args ...string) (string, error)
	Addr() string
	Port() string
	Name() string
	GetFreePort() (int, error)
	UpdateFreePort() error
	Close() error
}
