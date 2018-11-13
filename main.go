package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/wolakec/makhzen/registry"
	"github.com/wolakec/makhzen/server"
	"github.com/wolakec/makhzen/store"
)

func main() {

	port := flag.String("port", "5000", "a port number")
	flag.Parse()

	itemStore := store.New()
	r := registry.New()

	formattedPort := ":" + *port

	s := server.NewMakhzenServer(itemStore, r)

	handler := http.HandlerFunc(s.ServeHTTP)
	fmt.Printf("listening on port %s", *port)
	if err := http.ListenAndServe(formattedPort, handler); err != nil {
		log.Fatalf("could not listen on port %v %v", *port, err)
	}
}
