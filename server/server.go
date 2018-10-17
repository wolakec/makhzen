package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ItemServer struct {
	Store ItemStore
}

type ItemStore interface {
	GetValue(key string) (string, bool)
	Set(key string, value string) string
}

type ItemBody struct {
	Value string `json:"value"`
}

func (i *ItemServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		i.updateItem(w, r)
	case http.MethodGet:
		i.getItem(w, r)
	}
}

func (i *ItemServer) updateItem(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/items/"):]
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

	v := i.Store.Set(key, item.Value)
	log.Printf("PUT - key %s, value %s", key, v)

	fmt.Fprint(w, v)
}

func (i *ItemServer) getItem(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/items/"):]
	val, ok := i.Store.GetValue(key)

	if ok == false {
		w.WriteHeader(http.StatusNotFound)
	}

	log.Printf("GET - key %s, value %s", key, val)
	fmt.Fprint(w, val)
}
