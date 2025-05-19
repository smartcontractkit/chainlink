package solidity_cross_tests_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
)

func TestMeasureRandomnessRequestGasCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	coordinator := vrftesthelpers.NewVRFCoordinatorUniverse(t, key)
	keyHash_, _, fee := registerProvingKey(t, coordinator)

	estimate := estimateGas(t, coordinator.Backend.Client(), common.Address{},
		coordinator.ConsumerContractAddress, coordinator.ConsumerABI,
		"testRequestRandomness", common.BytesToHash(keyHash_[:]), fee)

	assert.Greater(t, estimate, uint64(134000),
		"requestRandomness tx gas cost lower than expected")
	// Note: changed from 160000 to 164079 in the Berlin hard fork (Geth 1.10)
	assert.Less(t, estimate, uint64(167000),
		"requestRandomness tx gas cost higher than expected")
}
