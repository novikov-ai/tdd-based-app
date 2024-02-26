package server

type InMemoryPlayerStore struct {
	store map[string]int
}

func NewInMemoryPlayersStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		make(map[string]int),
	}
}

func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.store[name]++
}

func (s *InMemoryPlayerStore) GetPlayerStore(name string) int {
	return s.store[name]
}
