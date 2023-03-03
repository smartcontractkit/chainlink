package pstoremem

import (
	"errors"
	"sync"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"

	pstore "github.com/libp2p/go-libp2p-core/peerstore"
)

type memoryKeyBook struct {
	sync.RWMutex // same lock. wont happen a ton.
	pks          map[peer.ID]ic.PubKey
	sks          map[peer.ID]ic.PrivKey
}

var _ pstore.KeyBook = (*memoryKeyBook)(nil)

// noop new, but in the future we may want to do some init work.
func NewKeyBook() *memoryKeyBook {
	return &memoryKeyBook{
		pks: map[peer.ID]ic.PubKey{},
		sks: map[peer.ID]ic.PrivKey{},
	}
}

func (mkb *memoryKeyBook) PeersWithKeys() peer.IDSlice {
	mkb.RLock()
	ps := make(peer.IDSlice, 0, len(mkb.pks)+len(mkb.sks))
	for p := range mkb.pks {
		ps = append(ps, p)
	}
	for p := range mkb.sks {
		if _, found := mkb.pks[p]; !found {
			ps = append(ps, p)
		}
	}
	mkb.RUnlock()
	return ps
}

func (mkb *memoryKeyBook) PubKey(p peer.ID) ic.PubKey {
	mkb.RLock()
	pk := mkb.pks[p]
	mkb.RUnlock()
	if pk != nil {
		return pk
	}
	pk, err := p.ExtractPublicKey()
	if err == nil {
		mkb.Lock()
		mkb.pks[p] = pk
		mkb.Unlock()
	}
	return pk
}

func (mkb *memoryKeyBook) AddPubKey(p peer.ID, pk ic.PubKey) error {
	// check it's correct first
	if !p.MatchesPublicKey(pk) {
		return errors.New("ID does not match PublicKey")
	}

	mkb.Lock()
	mkb.pks[p] = pk
	mkb.Unlock()
	return nil
}

func (mkb *memoryKeyBook) PrivKey(p peer.ID) ic.PrivKey {
	mkb.RLock()
	sk := mkb.sks[p]
	mkb.RUnlock()
	return sk
}

func (mkb *memoryKeyBook) AddPrivKey(p peer.ID, sk ic.PrivKey) error {
	if sk == nil {
		return errors.New("sk is nil (PrivKey)")
	}

	// check it's correct first
	if !p.MatchesPrivateKey(sk) {
		return errors.New("ID does not match PrivateKey")
	}

	mkb.Lock()
	mkb.sks[p] = sk
	mkb.Unlock()
	return nil
}
