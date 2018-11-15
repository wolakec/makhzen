package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/wolakec/makhzen/registry"
)

type StubItemStore struct {
	items map[string]string
}

func (s *StubItemStore) GetValue(key string) (string, bool) {
	v, ok := s.items[key]
	return v, ok
}

func (s *StubItemStore) Set(key string, v string) string {
	s.items[key] = v
	return v
}

type StubRegistry struct {
	Nodes            []registry.Node
	broadcasterCalls int
}

func (r *StubRegistry) AddNode(node registry.Node) registry.Node {
	r.Nodes = append(r.Nodes, node)
	return node
}

func (r *StubRegistry) GetNodes() []registry.Node {
	return r.Nodes
}

func (r *StubRegistry) Broadcast(key string, value string) {
	r.broadcasterCalls = r.broadcasterCalls + 1
}

func TestGETItems(t *testing.T) {
	store := StubItemStore{
		map[string]string{
			"Region":   "europe",
			"Platform": "mobile",
		},
	}
	reg := StubRegistry{
		Nodes: []registry.Node{},
	}
	server := NewMakhzenServer(&store, &reg)

	t.Run("returns value for 'Region' key", func(t *testing.T) {
		request := newGetValueRequest("Region")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := http.StatusOK
		got := response.Code

		assertStatus(t, got, want)
		assertResponseBody(t, response.Body.String(), "europe")
	})

	t.Run("returns value for 'Platform' key", func(t *testing.T) {
		request := newGetValueRequest("Platform")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := http.StatusOK
		got := response.Code

		assertStatus(t, got, want)
		assertResponseBody(t, response.Body.String(), "mobile")
	})

	t.Run("returns 404 on missing key", func(t *testing.T) {
		request := newGetValueRequest("NonExistantKey")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := http.StatusNotFound
		got := response.Code

		assertStatus(t, got, want)
	})
}

func TestPUTItems(t *testing.T) {
	store := StubItemStore{
		map[string]string{},
	}
	reg := StubRegistry{
		Nodes: []registry.Node{},
	}
	server := NewMakhzenServer(&store, &reg)

	t.Run("returns accepted on PUT", func(t *testing.T) {
		request := newPutValueRequest("LoadBalancer", "10.0.0.1:3030")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := http.StatusAccepted
		got := response.Code

		assertStatus(t, got, want)
	})

	t.Run("returns sent value 10.0.0.1:3030 on PUT", func(t *testing.T) {
		request := newPutValueRequest("LoadBalancer", "10.0.0.1:3030")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := "10.0.0.1:3030"
		got := response.Body.String()

		assertResponseBody(t, got, want)
	})

	t.Run("returns sent value 192.0.0.0:3030 on PUT", func(t *testing.T) {
		request := newPutValueRequest("LoadBalancer", "192.0.0.0:3030")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := "192.0.0.0:3030"
		got := response.Body.String()

		assertResponseBody(t, got, want)
	})
}

func TestBroadcastOnPut(t *testing.T) {
	store := StubItemStore{
		map[string]string{},
	}
	reg := StubRegistry{
		Nodes: []registry.Node{},
	}
	server := NewMakhzenServer(&store, &reg)

	t.Run("calls broadcast on PUT", func(t *testing.T) {
		request := newPutValueRequest("Region", "europe")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		request = newPutValueRequest("VPC", "iz5689")
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)

		expectedCalls := 2
		if reg.broadcasterCalls != expectedCalls {
			t.Errorf("broadcast was called %d times, expected %d", reg.broadcasterCalls, expectedCalls)
		}
	})
}

func TestPOSTMessage(t *testing.T) {
	store := StubItemStore{
		map[string]string{},
	}
	reg := StubRegistry{
		Nodes: []registry.Node{},
	}
	server := NewMakhzenServer(&store, &reg)

	t.Run("POST message returns 200", func(t *testing.T) {

		request := newPostMessageRequest("user", "test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusOK

		assertStatus(t, got, want)
	})

	t.Run("POST message stores value", func(t *testing.T) {

		wantKey := "new-key"
		wantValue := "server-5890"
		request := newPostMessageRequest(wantKey, wantValue)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got, ok := store.GetValue(wantKey)

		if ok == false {
			t.Errorf("key: %s does not exist", wantKey)
		}

		if got != wantValue {
			t.Errorf("incorrect value - got: %s, wanted: %s", got, wantValue)
		}
	})

	t.Run("POST message updates value", func(t *testing.T) {

		wantKey := "new-key"
		wantValue := "new_value"
		request := newPostMessageRequest(wantKey, wantValue)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got, ok := store.GetValue(wantKey)

		if ok == false {
			t.Errorf("key: %s does not exist", wantKey)
		}

		if got != wantValue {
			t.Errorf("incorrect value - got: %s, wanted: %s", got, wantValue)
		}
	})
}

func newPostMessageRequest(key string, value string) *http.Request {
	payload := fmt.Sprintf(`
			{
				"key": "%s",
				"value": "%s"
			}
		`, key, value)

	request, err := http.NewRequest(http.MethodPost, "/message", strings.NewReader(payload))

	if err != nil {
		log.Fatalln(err)
	}

	return request
}

func TestPUTItemsPersistsData(t *testing.T) {
	store := StubItemStore{
		map[string]string{},
	}
	reg := StubRegistry{
		Nodes: []registry.Node{},
	}
	server := NewMakhzenServer(&store, &reg)

	t.Run("returns persisted value europe on GET", func(t *testing.T) {
		request := newPutValueRequest("Region", "europe")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		request = newGetValueRequest("Region")
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "europe"

		assertResponseBody(t, got, want)
	})
}

func TestGETNodes(t *testing.T) {

	t.Run("it returns 200 on /nodes", func(t *testing.T) {
		store := StubItemStore{
			map[string]string{},
		}
		reg := StubRegistry{
			Nodes: []registry.Node{},
		}
		server := NewMakhzenServer(&store, &reg)

		request, _ := http.NewRequest(http.MethodGet, "/nodes", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("it returns a list of nodes on /nodes", func(t *testing.T) {
		wanted := []registry.Node{
			registry.Node{
				Address: "1234",
			},
			registry.Node{
				Address: "5678",
			},
		}
		store := StubItemStore{
			map[string]string{},
		}
		reg := StubRegistry{
			Nodes: wanted,
		}
		server := NewMakhzenServer(&store, &reg)

		request, _ := http.NewRequest(http.MethodGet, "/nodes", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []registry.Node

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Node, '%v'", response.Body, err)
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got %v, want %v", got, wanted)
		}

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestPOSTNode(t *testing.T) {
	t.Run("it returns a 201", func(t *testing.T) {
		store := StubItemStore{
			map[string]string{},
		}
		reg := StubRegistry{
			Nodes: []registry.Node{},
		}
		server := NewMakhzenServer(&store, &reg)

		request, _ := http.NewRequest(http.MethodPost, "/nodes", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
	})

	t.Run("it returns the created node", func(t *testing.T) {
		store := StubItemStore{
			map[string]string{},
		}
		reg := StubRegistry{
			Nodes: []registry.Node{},
		}

		server := NewMakhzenServer(&store, &reg)

		sentNode := registry.Node{
			Address: "localhost:3000",
		}

		b, err := json.Marshal(sentNode)

		if err != nil {
			t.Errorf("Unable to unmarshal %v %v", sentNode, err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/nodes", bytes.NewBuffer(b))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
	})
}

func newPutValueRequest(k string, v string) *http.Request {
	body := map[string]interface{}{
		"value": v,
	}
	bytesRepresentation, err := json.Marshal(body)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/items/%s", k), bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	return req
}

func newGetValueRequest(k string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/items/%s", k), nil)
	return req
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response status incorrect - got '%d', wanted '%d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body incorrect - got '%s', wanted '%s'", got, want)
	}
}
