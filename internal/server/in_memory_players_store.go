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
