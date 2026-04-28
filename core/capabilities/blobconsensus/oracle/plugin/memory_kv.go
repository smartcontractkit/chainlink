package plugin

import "fmt"

// memoryKV is an in-memory KeyValueStateReadWriter used by the OCR3 adapter
// to bridge OutcomeContext.PreviousOutcome into OCR 3.1-style state reads.
type memoryKV struct {
	m map[string][]byte
}

func newMemoryKVFromPreviousOutcome(previousOutcome []byte) *memoryKV {
	m := map[string][]byte{}
	if len(previousOutcome) > 0 {
		m[string(prevOutcomeStateKey)] = append([]byte(nil), previousOutcome...)
	}
	return &memoryKV{m: m}
}

func (s *memoryKV) Read(key []byte) ([]byte, error) {
	if s.m == nil {
		return nil, nil
	}
	v, ok := s.m[string(key)]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (s *memoryKV) Write(key, value []byte) error {
	if s.m == nil {
		return fmt.Errorf("memoryKV: Write on nil map")
	}
	s.m[string(key)] = append([]byte(nil), value...)
	return nil
}

func (s *memoryKV) Delete(key []byte) error {
	if s.m == nil {
		return nil
	}
	delete(s.m, string(key))
	return nil
}
