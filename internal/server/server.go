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

func (ps *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		ps.recordWin(w, player)
	case http.MethodGet:
		ps.showScore(w, player)
	}
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
