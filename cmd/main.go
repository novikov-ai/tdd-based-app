package main

import (
	"log"
	"net/http"
	"tdd-based-app/internal/server"
)

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerStore(name string) int {
	return 123
}

func main() {
	s := server.New(&InMemoryPlayerStore{})
	log.Fatal(http.ListenAndServe(":4000", s))
}
