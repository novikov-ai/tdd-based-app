package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nsf/jsondiff"
)

// TESTS FOLLOWS DESIGN
func TestServeHTTP(t *testing.T) {
	server := New(NewInMemoryPlayersStore())
	server.store.RecordWin("james", "") // warm up

	t.Run("registered endpoint: `/players/{name}`", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/james", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, response.Body.String(), "1")
	})

	t.Run("registered endpoint: `/league`", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, response.Body.String(), `{"players":["james"]}`)
	})
}

func Test_processPlayers(t *testing.T) {
	t.Run("returns a number of total player's wins", func(t *testing.T) {
		server := New(NewInMemoryPlayersStore())
		server.store.RecordWin("james", "") // warm up

		request, _ := http.NewRequest(http.MethodGet, "/players/james", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, response.Body.String(), "1")
	})

	t.Run("record a win for a given name (1 request - 1 win increment)", func(t *testing.T) {
		server := New(NewInMemoryPlayersStore())

		var response *httptest.ResponseRecorder
		for i := 0; i < 5; i++ {
			request, _ := http.NewRequest(http.MethodPost, "/players/james", nil)
			response = httptest.NewRecorder()
			server.ServeHTTP(response, request)
		}

		assertStatus(t, response.Code, http.StatusAccepted)
		assertResponse(t, response.Body.String(), "")

		request, _ := http.NewRequest(http.MethodGet, "/players/james", nil)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponse(t, response.Body.String(), "5")
	})
}

func Test_processLeague(t *testing.T) {
	t.Run("returns a list of all players stored (format JSON)", func(t *testing.T) {
		server := New(NewInMemoryPlayersStore())

		server.store.RecordWin("james", "") // warm up
		server.store.RecordWin("alex", "")

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		diff, _ := jsondiff.Compare([]byte(`{"players":["james","alex"]}`), []byte(response.Body.String()), &jsondiff.Options{})
		if diff != jsondiff.FullMatch {
			t.Errorf("got: %q, want: %q", response.Body.String(), `{"players":["james","alex"]}`)
		}
	})
}

func Test_processLeaguePlayers(t *testing.T) {
	playersLeagues := map[string]string{
		"james": "high",
		"bike":  "low",
		"alex":  "high",
	}
	t.Run("returns a list of league players stored (format JSON)", func(t *testing.T) {
		server := New(NewInMemoryPlayersStore())

		var request *http.Request
		for player, league := range playersLeagues {
			server.store.RecordWin(player, league) // warm up
			request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s?league=%s", player, league), nil)
		}

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		assertGivenLeague(t, server, "high", `{"players":["james","alex"], "league": "high"}`)
		assertGivenLeague(t, server, "low", `{"players":["bike"], "league": "low"}`)
	})
}

// TESTS BASED ON TDD
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

func assertGivenLeague(t *testing.T, server *PlayerServer, league, expected string) {
	request, _ := http.NewRequest(http.MethodGet, "/players/league?value="+league, nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	diff, _ := jsondiff.Compare([]byte(expected), []byte(response.Body.String()), &jsondiff.Options{})
	if diff != jsondiff.FullMatch {
		t.Errorf("got: %q, want: %q", response.Body.String(), expected)
	}
}

type StubPlayerStore struct {
	scores map[string]int
}

func (st *StubPlayerStore) GetPlayerScore(name string) int {
	return st.scores[name]
}

func (st *StubPlayerStore) RecordWin(name string, league string) {
	st.scores[name]++
}

func (st *StubPlayerStore) GetPlayers() []string {
	return nil
}

func (st *StubPlayerStore) GetPlayersOfLeague(league string) []string {
	return []string{}
}
