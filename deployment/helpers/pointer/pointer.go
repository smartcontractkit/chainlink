package pointer

func Coalesce[T any](p *T, fallback T) T {
	if p != nil {
		return *p
	}
	return fallback
}

func To[T any](v T) *T {
	return &v
}
