package cache

type Cache[Key comparable, Data any] struct {
	Cache map[Key]Data
}

func NewCache[Key comparable, Data any]() *Cache[Key, Data] {
	return &Cache[Key, Data]{
		make(map[Key]Data),
	}
}

func (s *Cache[Key, Data]) Get(key Key) (Data, bool) {
	data, ok := s.Cache[key]
	return data, ok
}

func (s *Cache[Key, Data]) Set(key Key, data Data) {
	s.Cache[key] = data
}

func (s *Cache[Key, Data]) Delete(key Key) {
	delete(s.Cache, key)
}
