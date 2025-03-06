package orch

import (
	"starfield/plugins/starfield/orch/node"
)

var nlist []node.Node

func RegisterNode(n node.Node) {
	nlist = append(nlist, n)
}

func GetNodes() []node.Node {
	return nlist
}
