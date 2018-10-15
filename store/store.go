package store

type Store struct {
	items map[string]string
}

func (s *Store) Set(k string, v string) string {
	s.items[k] = v
	return v
}

func (s *Store) Get(k string) (string, bool) {
	item, exists := s.items[k]

	if exists {
		return item, true
	}

	return "", false
}

func New() *Store {
	var s Store
	s.items = make(map[string]string)

	return &s
}
