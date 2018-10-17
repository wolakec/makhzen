package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestGETItems(t *testing.T) {
	store := StubItemStore{
		map[string]string{
			"Region":   "europe",
			"Platform": "mobile",
		},
	}
	server := &ItemServer{&store}

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
	server := &ItemServer{&store}

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

func TestPUTItemsPersistsData(t *testing.T) {
	store := StubItemStore{
		map[string]string{},
	}
	server := &ItemServer{&store}

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
