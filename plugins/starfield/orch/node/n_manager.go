package node

var nlist []Node

func RegisterNode(n Node) {
	nlist = append(nlist, n)
}

func GetNodes() []Node {
	return nlist
}
