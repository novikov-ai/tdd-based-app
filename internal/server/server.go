package server

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	RecordWin(name string)
	GetPlayerStore(name string) int
}

type PlayerServer struct {
	store PlayerStore
}

func New(store PlayerStore) *PlayerServer {
	return &PlayerServer{
		store: store,
	}
}

// Postcondition: server registered endpoints: `/players/{name}` & `/league`
// `players/{name}` => returns a number of total player's wins & record a win for a given name (1 request - 1 win increment)
// `league` => returns a list of all players stored (format JSON)
func (ps *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// register a new multiplexer
	router := http.NewServeMux()

	// register a new endpoint `/players/{name}`
	router.Handle("/players/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := strings.TrimPrefix(r.URL.Path, "/players/")

		switch r.Method {
		case http.MethodPost:
			ps.recordWin(w, player)
		case http.MethodGet:
			ps.showScore(w, player)
		}
	}))

	// register a new endpoint `/league`
	router.Handle("/league", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	// server serves all registered endpoints
	router.ServeHTTP(w, r)
}

func (ps *PlayerServer) recordWin(w http.ResponseWriter, player string) {
	ps.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (ps *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := ps.store.GetPlayerStore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}

func GetPlayerStore(name string) int {
	switch name {
	case "james":
		return 20
	case "nick":
		return 13
	}

	return 0
}
