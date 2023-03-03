package types

type Event struct {
	FromAddress Address
	Keys        []Felt
	Data        []Felt
}
