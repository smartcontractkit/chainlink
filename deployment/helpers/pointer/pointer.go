package pointer

func Coalesce[T any](p *T, fallback T) T {
	if p != nil {
		return *p
	}
	return fallback
}

//go:fix inline
func To[T any](v T) *T {
	return new(v)
}
