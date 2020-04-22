package vrf

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestMeasureRandomnessRequestGasCost(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, _, fee := registerProvingKey(t, coordinator)

	estimate := estimateGas(t, coordinator.backend, common.Address{},
		coordinator.consumerContractAddress, coordinator.consumerABI,
		"requestRandomness", common.BytesToHash(keyHash[:]), fee, one)

	assert.Greater(t, estimate, uint64(174000),
		"requestRandomness tx gas cost lower than expected")
	assert.Less(t, estimate, uint64(180000),
		"requestRandomness tx gas cost higher than expected")
}
