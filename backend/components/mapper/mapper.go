package mapper

type Mapper[T any, R any] struct {
	fn func(T) R
}

func New[T any, R any](fn func(T) R) *Mapper[T, R] {
	return &Mapper[T, R]{fn: fn}
}

func (m *Mapper[T, R]) Map(src T) R {
	return m.fn(src)
}

func (m *Mapper[T, R]) MapSlice(src []T) []R {
	if len(src) == 0 {
		return []R{}
	}

	result := make([]R, len(src))
	for i, v := range src {
		result[i] = m.fn(v)
	}
	return result
}
