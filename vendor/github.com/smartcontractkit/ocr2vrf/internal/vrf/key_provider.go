package vrf

import (
	"github.com/smartcontractkit/ocr2vrf/internal/dkg"
	dkg_contract "github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
)

type KeyProvider interface {
	KeyLookup(p dkg_contract.KeyID) (kd dkg.KeyData)
}
