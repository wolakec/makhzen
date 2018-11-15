package registry

import (
	"reflect"
	"testing"
)

type BroadcasterSpy struct {
	noCalls int
}

func (b *BroadcasterSpy) SendMessage(key string, value string, addr string) error {
	b.noCalls = b.noCalls + 1

	return nil
}

func TestNewReturnsRegistry(t *testing.T) {
	t.Run("Test New returns Registry containing added nodes", func(t *testing.T) {
		a := []string{"127.0.0.1:3001", "127.0.0.1:3002"}

		r := New(a)

		got := r.GetNodes()

		want := []Node{
			{
				Address: "127.0.0.1:3001",
			},
			{
				Address: "127.0.0.1:3002",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestAddNodeReturnsValue(t *testing.T) {
	t.Run("Test AddNode returns node with address 1234", func(t *testing.T) {
		r := &Registry{}

		newNode := r.AddNode(Node{
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
			Address: "192168",
		})
		got := newNode.Address
		want := "192168"

		if got != want {
			t.Errorf("got %v, want %s", got, want)
		}
	})
}

func TestBroadcast(t *testing.T) {
	t.Run("Test broadcast sends 1 message", func(t *testing.T) {
		broadcaster := BroadcasterSpy{}
		r := &Registry{
			Nodes: []Node{
				{
					Address: "127.0.0.1:4000",
				},
			},
			Broadcaster: &broadcaster,
		}

		r.Broadcast("key", "val")

		expectedCalls := 1
		got := broadcaster.noCalls

		if got != expectedCalls {
			t.Errorf("expected %d calls, got %d", expectedCalls, got)
		}
	})

	t.Run("Test broadcast sends 2 messages", func(t *testing.T) {
		broadcaster := BroadcasterSpy{}
		r := &Registry{
			Nodes: []Node{
				{
					Address: "127.0.0.1:4000",
				},
				{
					Address: "127.0.0.1:4002",
				},
			},
			Broadcaster: &broadcaster,
		}

		r.Broadcast("key", "val")

		expectedCalls := 2
		got := broadcaster.noCalls

		if got != expectedCalls {
			t.Errorf("expected %d calls, got %d", expectedCalls, got)
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
					Address: "192168",
				},
				Node{
					Address: "1234",
				},
			},
		}

		got := r.GetNodes()
		want := []Node{
			Node{
				Address: "192168",
			},
			Node{
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
			Address: "192168",
		})

		got := r.GetNodes()
		want := []Node{Node{
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
			Address: "1234",
		})
		r.AddNode(Node{
			Address: "5678",
		})

		got := r.GetNodes()
		want := []Node{
			Node{
				Address: "1234",
			},
			Node{
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
					Address: "1234",
				},
			},
		}

		r.AddNode(Node{
			Address: "5678",
		})
		r.AddNode(Node{
			Address: "910",
		})

		got := r.GetNodes()
		want := []Node{
			Node{
				Address: "1234",
			},
			Node{
				Address: "5678",
			},
			Node{
				Address: "910",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
