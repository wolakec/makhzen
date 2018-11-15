package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/wolakec/makhzen/registry"
	"github.com/wolakec/makhzen/server"
	"github.com/wolakec/makhzen/store"
)

func main() {

	port := flag.String("port", "5000", "a port number")
	cluster := flag.String("cluster", "", "http://127.0.0.1:3001,http://127.0.0.1:3002")
	flag.Parse()

	formattedPort := ":" + *port
	instances := strings.Split(*cluster, ",")

	itemStore := store.New()
	r := registry.New(instances)

	s := server.NewMakhzenServer(itemStore, r)

	handler := http.HandlerFunc(s.ServeHTTP)
	fmt.Printf("listening on port %s \n", *port)
	if err := http.ListenAndServe(formattedPort, handler); err != nil {
		log.Fatalf("could not listen on port %v %v", *port, err)
	}
}
