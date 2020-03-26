package vrf

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestMeasureRandomnessRequestGasCost(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash_, _, fee := registerProvingKey(t, coordinator)

	estimate := estimateGas(t, coordinator.backend, common.Address{},
		coordinator.consumerContractAddress, coordinator.consumerABI,
		"requestRandomness", common.BytesToHash(keyHash_[:]), fee, one)

	assert.Greater(t, estimate, uint64(175000),
		"requestRandomness tx gas cost lower than expected")
	assert.Less(t, estimate, uint64(176000),
		"requestRandomness tx gas cost higher than expected")
}
