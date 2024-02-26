package main

import (
	"fmt"
	"log"
	"net/http"
	"tdd-based-app/internal/server"
)

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerStore(name string) int {
	return 123
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	fmt.Println("recording...")
}

func main() {
	s := server.New(&InMemoryPlayerStore{})

	log.Fatal(http.ListenAndServe(":4000", s))
}
