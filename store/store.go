package store

type Store struct {
	items map[string]string
}

func (s *Store) Set(k string, v string) string {
	s.items[k] = v
	return v
}

func (s *Store) GetValue(k string) (string, bool) {
	v, ok := s.items[k]

	return v, ok
}

func New() *Store {
	var s Store
	s.items = make(map[string]string)

	return &s
}
