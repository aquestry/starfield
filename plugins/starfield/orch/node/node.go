package node

type Node interface {
	Run(args ...string) (string, error)
	Addr() string
	Port() string
	Name() string
	Close() error
}
