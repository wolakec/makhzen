package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wolakec/makhzen/server"
	"github.com/wolakec/makhzen/store"
)

func main() {
	itemStore := store.New()

	s := &server.ItemServer{Store: itemStore}

	handler := http.HandlerFunc(s.ServeHTTP)
	fmt.Println("listening on port 5000")
	if err := http.ListenAndServe(":5000", handler); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
