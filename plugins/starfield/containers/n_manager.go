package containers

import (
	"fmt"
	"starfield/plugins/starfield/containers/node"
)

var nlist []node.Node

func RegisterNode(n node.Node) {
	nlist = append(nlist, n)
}

func RemoveNode(name string) {
	for i, n := range nlist {
		if n.Name() == name {
			nlist = append(nlist[:i], nlist[i+1:]...)
			return
		}
	}
}

func GetNodes() []node.Node {
	return nlist
}

func GetNode(name string) (node.Node, error) {
	for _, n := range nlist {
		if n.Name() == name {
			return n, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
