package logger

import "maps"

type Fields map[string]any

func (f Fields) With(xs ...any) Fields {
	if len(xs)%2 != 0 {
		panic("expected even number of arguments")
	}
	f2 := make(Fields, len(f)+(len(xs)/2))
	maps.Copy(f2, f)
	for i := 0; i < len(xs)/2; i++ {
		key, is := xs[i*2].(string)
		if !is {
			continue
		}
		val := xs[i*2+1]
		f2[key] = val
	}
	return f2
}

func (f Fields) Merge(f2 Fields) Fields {
	f3 := make(Fields, len(f)+len(f2))
	maps.Copy(f3, f)
	maps.Copy(f3, f2)
	return f3
}

func (f Fields) Slice() []any {
	s := make([]any, len(f)*2)
	var i int
	for k, v := range f {
		s[i*2] = k
		s[i*2+1] = v
		i++
	}
	return s
}
