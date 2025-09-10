package vault

import (
	"sync"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
)

type LazyPublicKey struct {
	publicKey *tdh2easy.PublicKey
	mu        sync.Mutex
}

func (p *LazyPublicKey) Get() *tdh2easy.PublicKey {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.publicKey
}

func (p *LazyPublicKey) Set(pk *tdh2easy.PublicKey) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.publicKey = pk
}

func NewLazyPublicKey() *LazyPublicKey {
	return &LazyPublicKey{}
}
