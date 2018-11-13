package registry

type Registry struct {
	Nodes []Node
}

type Node struct {
	Id      string
	Address string
}

func (r *Registry) AddNode(node Node) Node {
	r.Nodes = append(r.Nodes, node)
	return node
}

func (r *Registry) GetNodes() []Node {
	return r.Nodes
}

func New() *Registry {
	var r Registry

	r.Nodes = []Node{}

	return &r
}
