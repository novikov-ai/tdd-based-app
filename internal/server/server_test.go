package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayersGames(t *testing.T) {
	server := PlayerServer{
		store: &StubPlayerStore{
			scores: map[string]int{
				"james": 20,
				"nick":  13,
				"alex":  34,
			},
		},
	}

	t.Run("return james score", func(t *testing.T) {
		request := newGETScoreRequest("james")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, response.Body.String(), "20")
	})

	t.Run("return nick score", func(t *testing.T) {
		request := newGETScoreRequest("nick")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, response.Body.String(), "13")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGETScoreRequest("katty")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPOSTRecordWins(t *testing.T) {
	server := PlayerServer{
		store: &StubPlayerStore{
			scores: map[string]int{
				"james": 20,
				"nick":  13,
				"alex":  34,
			},
		},
	}
	t.Run("record players wins", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/james", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}

func newGETScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func assertResponse(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("incorect status, got %d, want %d", got, want)
	}
}

type StubPlayerStore struct {
	scores map[string]int
}

func (st *StubPlayerStore) GetPlayerStore(name string) int {
	return st.scores[name]
}

func (st *StubPlayerStore) RecordWin(name string) {
	st.scores[name]++
}
