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

func (s *Cache[Key, Data]) Create(key Key, data Data) bool {
	if _, ok := s.Cache[key]; ok {
		return false
	}
	s.Cache[key] = data
	return true
}

func (s *Cache[Key, Data]) Update(key Key, data Data) bool {
	if _, ok := s.Cache[key]; !ok {
		return false
	}
	s.Cache[key] = data
	return true
}

func (s *Cache[Key, Data]) Delete(key Key) bool {
	if _, ok := s.Cache[key]; !ok {
		return false
	}
	delete(s.Cache, key)
	return true
}
