package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wolakec/makhzen/broadcaster"
	"github.com/wolakec/makhzen/registry"
)

type MakhzenServer struct {
	Store    ItemStore
	Registry NodeRegistry
	http.Handler
}

type ItemStore interface {
	GetValue(key string) (string, bool)
	Set(key string, value string) string
}

type ItemBody struct {
	Value string `json:"value"`
}

type NodeRegistry interface {
	AddNode(node registry.Node) registry.Node
	GetNodes() []registry.Node
	Broadcast(key string, value string)
}

func NewMakhzenServer(store ItemStore, registry NodeRegistry) *MakhzenServer {
	s := new(MakhzenServer)

	s.Store = store
	s.Registry = registry

	router := http.NewServeMux()

	router.Handle("/nodes", http.HandlerFunc(s.nodesHandler))
	router.Handle("/items/", http.HandlerFunc(s.itemsHandler))
	router.Handle("/message", http.HandlerFunc(s.messageHandler))

	s.Handler = router

	return s
}

func (s *MakhzenServer) nodesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		s.getNodes(w, r)
	case http.MethodPost:
		s.createNode(w, r)
	}
}

func (s *MakhzenServer) getNodes(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.Registry.GetNodes())

	w.WriteHeader(http.StatusOK)
}

func (s *MakhzenServer) createNode(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (s *MakhzenServer) messageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	var msg broadcaster.Message
	err = json.Unmarshal(b, &msg)

	if err != nil {
		log.Fatal(err)
	}

	s.Store.Set(msg.Key, msg.Value)

	log.Printf("recieved message from node: %s, key: %s, value: %s", r.RemoteAddr, msg.Key, msg.Value)
}

func (s *MakhzenServer) itemsHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/items/"):]

	switch r.Method {
	case http.MethodPut:
		s.updateItem(w, r, key)
	case http.MethodGet:
		s.getItem(w, key)
	}
}

func (s *MakhzenServer) updateItem(w http.ResponseWriter, r *http.Request, key string) {

	w.WriteHeader(http.StatusAccepted)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var item ItemBody
	err = json.Unmarshal(b, &item)
	if err != nil {
		log.Fatal(err)
	}

	v := s.Store.Set(key, item.Value)
	log.Printf("PUT - key %s, value %s", key, v)

	s.Registry.Broadcast(key, item.Value)

	fmt.Fprint(w, v)
}

func (s *MakhzenServer) getItem(w http.ResponseWriter, key string) {

	val, ok := s.Store.GetValue(key)

	if ok == false {
		w.WriteHeader(http.StatusNotFound)
	}

	log.Printf("GET - key %s, value %s", key, val)
	fmt.Fprint(w, val)
}
