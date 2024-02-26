package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	RecordWin(name string)
	GetPlayerStore(name string) int
	GetPlayers() []string
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
func (ps *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// register a new multiplexer
	router := http.NewServeMux()

	// register a new endpoint `/players/{name}`
	router.Handle("/players/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ps.processPlayers(w, r)
	}))

	// register a new endpoint `/league`
	router.Handle("/league", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ps.processLeague(w, r)
	}))

	// server serves all registered endpoints
	router.ServeHTTP(w, r)
}

// Postcondition: returns a number of total player's wins & record a win for a given name (1 request - 1 win increment)
func (ps *PlayerServer) processPlayers(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		ps.recordWin(w, player)
	case http.MethodGet:
		ps.showScore(w, player)
	}
}

// Postcondition: returns a list of all players stored (format JSON)
func (ps *PlayerServer) processLeague(w http.ResponseWriter, r *http.Request) {
	var result struct {
		Players []string `json:"players"`
	}

	result.Players = ps.store.GetPlayers()

	w.Header().Set("Content-Type", "application/json")

	bytes, _ := json.Marshal(result)
	w.Write(bytes)
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
