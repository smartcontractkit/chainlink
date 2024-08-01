package llo

var _ ShouldRetireCache = &shouldRetireCache{}

type shouldRetireCache struct{}

// TODO: https://smartcontract-it.atlassian.net/browse/MERC-3386
func NewShouldRetireCache() ShouldRetireCache {
	return newShouldRetireCache()
}

func newShouldRetireCache() *shouldRetireCache {
	return &shouldRetireCache{}
}

func (c *shouldRetireCache) ShouldRetire() (bool, error) {
	// TODO
	return false, nil
}
