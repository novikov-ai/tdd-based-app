package main

import (
	"log"
	"net/http"
	"tdd-based-app/internal/server"
)

func main() {
	s := server.New(server.NewInMemoryPlayersStore())

	log.Fatal(http.ListenAndServe(":4000", s))
}
