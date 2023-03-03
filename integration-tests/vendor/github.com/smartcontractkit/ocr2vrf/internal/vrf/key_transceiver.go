package vrf

import (
	"sync"

	"github.com/smartcontractkit/ocr2vrf/internal/dkg"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
)

type KeyTransceiver struct {
	keyID contract.KeyID
	kd    *dkg.KeyData
	mu    sync.RWMutex
}

var _ dkg.KeyConsumer = (*KeyTransceiver)(nil)
var _ KeyProvider = (*KeyTransceiver)(nil)

func NewKeyTransceiver(keyID contract.KeyID) *KeyTransceiver {
	return &KeyTransceiver{keyID, nil, sync.RWMutex{}}
}

func (kt *KeyTransceiver) KeyInvalidated(kID contract.KeyID) {

	kt.mu.Lock()
	defer kt.mu.Unlock()

	if kt.keyID == kID {
		kt.kd = nil
	}
}

func (kt *KeyTransceiver) NewKey(kID contract.KeyID, kd *dkg.KeyData) {

	kt.mu.Lock()
	defer kt.mu.Unlock()

	if kt.keyID == kID {
		kt.kd = kd.Clone()
	}
}

func (kt *KeyTransceiver) KeyLookup(p contract.KeyID) dkg.KeyData {

	kt.mu.RLock()
	defer kt.mu.RUnlock()

	if p == kt.keyID {
		if kt.kd != nil {
			return *kt.kd.Clone()
		}
		return dkg.KeyData{nil, nil, nil, 0, false}
	}

	panic("key consumer is asking for unknown key ID")
}

func (kt *KeyTransceiver) KeyGenerated() bool {
	return (kt.kd != nil) && kt.kd.Present
}
