package s4storage

import "fmt"

type Key struct {
	Address string
	SlotID  uint8
}

// Trivial in-memory store.
type S4APIService struct {
	store map[Key][]byte
}

func NewS4Service() *S4APIService {
	return &S4APIService{
		store: make(map[Key][]byte),
	}
}

func (s *S4APIService) Put(address string, slotId uint8, payload []byte, nonce uint32) error {
	s.store[Key{address, slotId}] = payload
	return nil
}

func (s *S4APIService) Get(address string, slotId uint8) ([]byte, error) {
	val, ok := s.store[Key{address, slotId}]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return val, nil
}

func (s *S4APIService) GetNextNonce(address string, slotId uint8) (uint32, error) {
	return 0, nil
}
