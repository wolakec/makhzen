package registry

import (
	"reflect"
	"testing"
)

func TestAddNodeReturnsValue(t *testing.T) {
	t.Run("Test AddNode returns node with address 1234", func(t *testing.T) {
		r := &Registry{}

		newNode := r.AddNode(Node{
			Id:      "1337",
			Address: "1234",
		})
		got := newNode.Address
		want := "1234"

		if got != want {
			t.Errorf("got %v, want %s", got, want)
		}
	})

	t.Run("Test AddNode returns 192168", func(t *testing.T) {
		r := &Registry{}

		newNode := r.AddNode(Node{
			Id:      "1337",
			Address: "192168",
		})
		got := newNode.Address
		want := "192168"

		if got != want {
			t.Errorf("got %v, want %s", got, want)
		}
	})
}

func TestGetNodes(t *testing.T) {
	t.Run("Test get nodes returns nothing", func(t *testing.T) {
		r := &Registry{
			Nodes: []Node{},
		}

		got := r.GetNodes()
		want := []Node{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Test get nodes returns ips", func(t *testing.T) {
		r := &Registry{
			Nodes: []Node{
				Node{
					Id:      "1337",
					Address: "192168",
				},
				Node{
					Id:      "1338",
					Address: "1234",
				},
			},
		}

		got := r.GetNodes()
		want := []Node{
			Node{
				Id:      "1337",
				Address: "192168",
			},
			Node{
				Id:      "1338",
				Address: "1234",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Test get nodes returns added node", func(t *testing.T) {
		r := &Registry{
			Nodes: []Node{},
		}

		r.AddNode(Node{
			Id:      "1337",
			Address: "192168",
		})

		got := r.GetNodes()
		want := []Node{Node{
			Id:      "1337",
			Address: "192168",
		}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Test get nodes returns added nodes", func(t *testing.T) {
		r := &Registry{
			Nodes: []Node{},
		}

		r.AddNode(Node{
			Id:      "1",
			Address: "1234",
		})
		r.AddNode(Node{
			Id:      "2",
			Address: "5678",
		})

		got := r.GetNodes()
		want := []Node{
			Node{
				Id:      "1",
				Address: "1234",
			},
			Node{
				Id:      "2",
				Address: "5678",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Test get nodes returns added nodes when already populated", func(t *testing.T) {
		r := &Registry{
			Nodes: []Node{
				Node{
					Id:      "1",
					Address: "1234",
				},
			},
		}

		r.AddNode(Node{
			Id:      "2",
			Address: "5678",
		})
		r.AddNode(Node{
			Id:      "3",
			Address: "910",
		})

		got := r.GetNodes()
		want := []Node{
			Node{
				Id:      "1",
				Address: "1234",
			},
			Node{
				Id:      "2",
				Address: "5678",
			},
			Node{
				Id:      "3",
				Address: "910",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
