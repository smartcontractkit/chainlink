package dkg

import (
	"crypto/rand"
	"io"
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
	"github.com/smartcontractkit/ocr2vrf/types/hash"

	"go.dedis.ch/kyber/v3"
	kshare "go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type shareRecords map[hash.Hash]*shareRecord

func newShareRecords() shareRecords {
	return map[hash.Hash]*shareRecord{}
}

func (rs shareRecords) set(r *shareRecord, h hash.Hash) error {
	if h == hash.Zero {
		m, err := r.marshal()
		if err != nil {
			return errors.Wrap(err, "could not marshal share record to get content address")
		}
		h = hash.GetHash(m)
	}
	rs[h] = r
	return nil
}

func (rs shareRecords) getRandom(givenHashes []hash.Hash, randomness io.Reader) (hash.Hash, error) {
	extantHashes := make([]hash.Hash, 0, len(givenHashes))
	for _, h := range givenHashes {
		if _, ok := rs[h]; ok {
			extantHashes = append(extantHashes, h)
		}
	}
	if len(extantHashes) == 0 {
		return hash.Hash{}, errors.Errorf(
			"don't know any of the share records in the chosen key",
		)
	}
	hIdx, err := rand.Int(randomness, big.NewInt(int64(len(extantHashes))))
	if err != nil {
		return hash.Hash{}, errors.Wrap(err, "could not choose random hash")
	}
	return extantHashes[hIdx.Uint64()], nil
}

func (rs shareRecords) allKeysPresent(hs []hash.Hash) bool {
	for _, h := range hs {
		if _, ok := rs[h]; !ok {
			return false
		}
	}
	return true
}

func (rs shareRecords) recoverDistributedKeyShare(
	encryptionSecretKey kyber.Scalar,
	receiver player_idx.PlayerIdx,
	keyData *contract.KeyData,
	keyGroup anon.Suite,
	domainSep types.ConfigDigest,
) (*kshare.PriShare, error) {
	acc := keyGroup.Scalar().Zero()
	for _, h := range keyData.Hashes {
		publicShare, present := rs[h]
		if !present {
			return nil, errors.Errorf("no share record for key hash 0x%x", h)
		}
		recvShare, err := publicShare.shareSet.Decrypt(
			receiver, encryptionSecretKey, keyGroup, domainSep,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decrypt share for key hash 0x%x", h)
		}
		acc = acc.Clone().Add(acc, recvShare.V)
	}
	rv := receiver.PriShare(acc)
	return &rv, nil
}

func (rs shareRecords) recoverPublicShares(
	keyData *contract.KeyData,
) ([]kyber.Point, error) {
	var partialShares [][]kyber.Point
	if len(keyData.Hashes) == 0 {
		return nil, errors.Errorf("can't reconstruct shares from 0 share records")
	}
	for _, h := range keyData.Hashes {
		sr, ok := rs[h]
		if !ok {
			return nil, errors.Errorf("no share found for hash %s", h)
		}

		partialShares = append(partialShares, sr.shareSet.PublicShares())
	}
	rv := make([]kyber.Point, len(partialShares[0]))
	for _, ps := range partialShares {
		for playerIdx, partialShare := range ps {
			if rv[playerIdx] == nil {
				rv[playerIdx] = partialShare.Clone().Null()
			}

			rv[playerIdx] = partialShare.Clone().Add(rv[playerIdx], partialShare)
		}
	}
	return rv, nil
}
