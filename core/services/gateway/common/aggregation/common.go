package aggregation

// StringSet is a simple set implementation for strings.
type StringSet map[string]struct{}

func (s StringSet) Add(val string) {
	s[val] = struct{}{}
}

func (s StringSet) Contains(val string) bool {
	_, exists := s[val]
	return exists
}

func (s StringSet) Remove(val string) {
	delete(s, val)
}

func (s StringSet) Values() []string {
	values := make([]string, 0, len(s))
	for k := range s {
		values = append(values, k)
	}
	return values
}
