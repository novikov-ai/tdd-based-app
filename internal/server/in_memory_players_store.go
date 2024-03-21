package server

type InMemoryPlayerStore struct {
	store          map[string]int
	playersLeagues map[string]string
}

// Postcondition: returns a new empty in-memory storage
func NewInMemoryPlayersStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		store:          make(map[string]int),
		playersLeagues: make(map[string]string),
	}
}

// Postcondition: increment total wins of given player (new player created if not exists);
// player inserted with a given league
func (s *InMemoryPlayerStore) RecordWin(name, league string) {
	s.store[name]++
	s.playersLeagues[name] = league
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

// Precondition: given league exists
// Postcondition: returns all players of given league
func (s *InMemoryPlayerStore) GetPlayersOfLeague(league string) []string {
	result := make([]string, 0, len(s.playersLeagues))
	for player, playersLeague := range s.playersLeagues {
		if playersLeague == league {
			result = append(result, player)
		}
	}
	return result
}
