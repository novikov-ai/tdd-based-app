# Код и тесты должны следовать дизайну, логической архитектуре (рефлексия)

### Пишу код по TDD (test-driven development)

Следуя практикам TDD: сначала пишу максимально "глупые" тесты, которые падают, а потом заставляю их проходить, внося итерационные правки в код (делаю самый минимум, из-за которого тесты не проходят в моменте). Использую метод: красный-зеленый-рефактор.

В такой практике довольно сложно четко следовать итерационному подходу, при котором делаешь ровно 1 изменение в коде за раз. После первого провала теста появляется большое желание исправить код и постараться написать сразу законченную версию, так как в голове появляется сразу множество тест-кейсов, которые обязательно нужно будет учесть дальше (приходится бить себя по рукам). 

В конечном счете, появляется законченный вариант, очевидными плюсами которого являются - наличие тестов. Наш код при таком подходе все время следовал за тестами, но следуют ли концептуально наши тесты проектному дизайну?

~~~go
// server_test.go
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
~~~

~~~go
// server.go
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
~~~

### Пишу код с акцентом на то, чтобы код и тесты следовали дизайну и логической архитектуре приложения

Чтобы код и тесты следовали проектному дизайну, я из каждого требования к приложению формирую спецификацию, основанную на пред и пост условиях. 

Когда пред/пост условия написаны к основным функциям, легко оценить соответствует ли этот код ровно одному дизайну или же множеству. Стремимся к правилу: 1 код - 1 дизайн. 

Обозначив пост-условие для функции `ServeHTTP()` я сразу заметил, что логически у меня функция как будто спроектирована запутано и при этом как будто не соответствует той роли, которую ожидает от нее система. 

Подобную операцию произвел для множества других функций и в итоге после рефакторинга в соответствии со спецификацией функции как будто сами выстроились в логичный пазл, который удобно читать и развивать в будущем. 

После изменения основного файла с логикой сервера перешел к реализации тестов и сразу заметил, что тесты написанные при помощи TDD теперь не годятся для текущей архитектуры приложения. 

Новый дизайн функций потребовал концептуально новых тестов. Если раньше я проверял в `TestGETPlayersGames()` отправку GET-запроса и получения результата определенного игрока, то теперь у меня появились тесты, которые проверяют регистрацию ендпоинтов `TestServeHTTP()` и процессинг каждого из них в отдельности с соблюдением всех пост-условий: `Test_processPlayers()` , `Test_processLeague()`. 

С этого момента всегда буду делать акцент на дизайне и спецификации приложения, как основы будущих тестов и будущего кода.  

~~~go
// server.go
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	RecordWin(name string)
	GetPlayerScore(name string) int
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
	score := ps.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}
~~~

~~~go
// server_test.go
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
	server.store.RecordWin("james") // warm up

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
		server.store.RecordWin("james") // warm up

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

		server.store.RecordWin("james") // warm up
		server.store.RecordWin("alex")

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
~~~

~~~go
// in_memory_players_store.go
package server

type InMemoryPlayerStore struct {
	store map[string]int
}

// Postcondition: returns a new empty in-memory storage
func NewInMemoryPlayersStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		make(map[string]int),
	}
}

// Postcondition: increment total wins of given player (new player created if not exists)
func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.store[name]++
}

// Precondition: given player exists (if not - returns 0)
// Postcondition: returns a total score of given player
func (s *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return s.store[name]
}

// Postcondition: returns a list of all players with total score more than 0
func (s *InMemoryPlayerStore) GetPlayers() []string {
	result := []string{}

	for name := range s.store {
		result = append(result, name)
	}

	return result
}
~~~