package set

var empty = struct{}{}

type set[T comparable] struct {
	data map[T]struct{}
}

func NewSet[T comparable]() *set[T] {

	return &set[T]{
		data: make(map[T]struct{}),
	}
}

func NewSetWith[T comparable](elems ...T) *set[T] {
	s := NewSet[T]()
	for _, elem := range elems {
		s.data[elem] = empty
	}
	return s
}

func (s *set[T]) Add(x T) {
	s.data[x] = empty
}

func (s *set[T]) Contains(x T) bool {
	_, exists := s.data[x]
	return exists
}

func (s *set[T]) ContainsAll(elems ...T) bool {
	for _, elem := range elems {
		if !s.Contains(elem) {
			return false
		}
	}
	return true
}

func (s *set[T]) Remove(x T) {
	delete(s.data, x)
}

func (s *set[T]) Len() int {
	return len(s.data)
}

func (s *set[T]) Clear() {
	s.data = make(map[T]struct{})
}

func (s *set[T]) Keys() []T {
	result := make([]T, 0, len(s.data))
	for k := range s.data {
		result = append(result, k)
	}

	return result
}
