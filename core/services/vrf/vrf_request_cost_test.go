package vrf_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestMeasureRandomnessRequestGasCost(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash_, _, fee := registerProvingKey(t, coordinator)

	estimate := estimateGas(t, coordinator.backend, common.Address{},
		coordinator.consumerContractAddress, coordinator.consumerABI,
		"requestRandomness", common.BytesToHash(keyHash_[:]), fee, big.NewInt(1))

	assert.Greater(t, estimate, uint64(134000),
		"requestRandomness tx gas cost lower than expected")
	assert.Less(t, estimate, uint64(160000),
		"requestRandomness tx gas cost higher than expected")
}
