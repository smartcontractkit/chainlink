package retirement

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

type neverShouldRetireCache struct{}

func NewNeverShouldRetireCache() llotypes.ShouldRetireCache {
	return &neverShouldRetireCache{}
}

func (n *neverShouldRetireCache) ShouldRetire(digest ocrtypes.ConfigDigest) (bool, error) {
	return false, nil
}
