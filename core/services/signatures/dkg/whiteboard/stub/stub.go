package stub

import (
	w "chainlink/core/services/signatures/dkg/whiteboard"

	"go.dedis.ch/kyber/v3"
)

// BoardKey implements w.BoardKey
type BoardKey struct {
	k string
	p kyber.Point
}

var _ w.BoardKey = &BoardKey{} // *BoardKey implements w.BoardKey

// BoardQuery implements w.BoardQuery
type BoardQuery struct {
	q string
	p kyber.Point
}

var _ w.BoardQuery = &BoardQuery{} // *BoardQuery implements w.BoardQuery

// Match implements w.BoardQuery.Match
func (q *BoardQuery) Match(k w.BoardKey) bool {
	return k.(BoardKey).k == q.q
}

// PublicKey implements w.BoardQuery.PublicKey
func (q *BoardQuery) PublicKey() w.PublicKey { return q.p }
