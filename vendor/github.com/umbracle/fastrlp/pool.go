package fastrlp

import (
	"sync"
)

// DefaultParserPool is a default ParserPool
var DefaultParserPool ParserPool

// ParserPool may be used for pooling Parsers for similarly typed RLPs.
type ParserPool struct {
	pool sync.Pool
}

// Get acquires a Parser from the pool.
func (pp *ParserPool) Get() *Parser {
	v := pp.pool.Get()
	if v == nil {
		return &Parser{}
	}
	return v.(*Parser)
}

// Put releases the parser to the pool.
func (pp *ParserPool) Put(p *Parser) {
	pp.pool.Put(p)
}

// DefaultArenaPool is a default ArenaPool
var DefaultArenaPool ArenaPool

// ArenaPool may be used for pooling Arenas for similarly typed RLPs.
type ArenaPool struct {
	pool sync.Pool
}

// Get acquires an Arena from the pool.
func (ap *ArenaPool) Get() *Arena {
	v := ap.pool.Get()
	if v == nil {
		return &Arena{}
	}
	return v.(*Arena)
}

// Put releases an Arena to the pool.
func (ap *ArenaPool) Put(a *Arena) {
	a.Reset()
	ap.pool.Put(a)
}
