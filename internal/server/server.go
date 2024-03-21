package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	RecordWin(name string, league string)
	GetPlayerScore(name string) int
	GetPlayers() []string
	GetPlayersOfLeague(league string) []string
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

	// register a new endpoint `/players/league?value=`
	router.Handle("/players/league", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ps.processLeaguePlayers(w, r)
	}))

	// server serves all registered endpoints
	router.ServeHTTP(w, r)
}

// Postcondition: returns a number of total player's wins & record a win for a given name (1 request - 1 win increment)
func (ps *PlayerServer) processPlayers(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	query := r.URL.Query()
	league := query.Get("league")

	switch r.Method {
	case http.MethodPost:
		ps.recordWin(w, player, league)
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

func (ps *PlayerServer) recordWin(w http.ResponseWriter, player, league string) {
	ps.store.RecordWin(player, league)
	w.WriteHeader(http.StatusAccepted)
}

func (ps *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := ps.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}

// Postcondition: returns a list of all players related to a given league stored (format JSON)
func (ps *PlayerServer) processLeaguePlayers(w http.ResponseWriter, r *http.Request) {
	var result struct {
		Players []string `json:"players"`
		League  string   `json:"league"`
	}

	query := r.URL.Query()
	league := query.Get("value")

	result.Players = ps.store.GetPlayersOfLeague(league)
	result.League = league

	w.Header().Set("Content-Type", "application/json")

	bytes, _ := json.Marshal(result)
	w.Write(bytes)
}
