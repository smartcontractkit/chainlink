package vrf

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/pkg/errors"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"
	kshare "go.dedis.ch/kyber/v3/share"

	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg"
	vrf_types "github.com/smartcontractkit/ocr2vrf/types"
)

func m(x int64) *mod.Int { return mod.NewInt64(x, bn256.P) }

var three = m(3)

func affineCoordinates(p kyber.Point) (*big.Int, *big.Int) {
	b := altbn_128.LongMarshal(p)
	if len(b) != 64 {
		panic("wrong length for marshaled point")
	}
	return big.NewInt(0).SetBytes(b[:32]), big.NewInt(0).SetBytes(b[32:])
}

func (s *sigRequest) computePartialSig(
	block vrf_types.Block, kd dkg.KeyData,
) (kyber.Point, error) {

	seed := blsSeed(s.configDigest, block, kd.PublicKey)

	output := kd.SecretShare.Mul(seed)

	pk := s.i.Index(kd.Shares).(kshare.PubShare).V

	if !validateSignature(s.pairing, seed, pk, output) {
		return nil, errors.Errorf(failedVerifyOwnContributionMsg)
	}
	return output, nil
}

func blsSeed(
	domainSeparator common.Hash,
	block vrf_types.Block,
	pk kyber.Point,
) kyber.Point {

	h := block.VRFHash(domainSeparator, pk)

	hpoint := altbn_128.NewHashProof(h).HashPoint
	return hpoint
}

const (
	failedVerifyOwnContributionMsg = "could not verify own contribution to signature"
)
