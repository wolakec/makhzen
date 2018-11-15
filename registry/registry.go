package registry

import "github.com/wolakec/makhzen/broadcaster"

type Registry struct {
	Nodes       []Node
	Broadcaster MessageBroadcaster
}

type MessageBroadcaster interface {
	SendMessage(key string, value string, addr string) error
}

type Node struct {
	Address string
}

func (r *Registry) AddNode(node Node) Node {
	r.Nodes = append(r.Nodes, node)
	return node
}

func (r *Registry) GetNodes() []Node {
	return r.Nodes
}

func (r *Registry) Broadcast(key string, value string) {
	for _, node := range r.Nodes {
		r.Broadcaster.SendMessage(key, value, node.Address)
	}
}

func New(addresses []string) *Registry {
	var r Registry

	r.Nodes = []Node{}

	r.Broadcaster = &broadcaster.Broadcaster{}

	for _, addr := range addresses {
		r.AddNode(
			Node{
				Address: addr,
			},
		)
	}

	return &r
}
