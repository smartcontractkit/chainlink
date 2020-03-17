package vrf

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeasureRandomnessRequestGasCost(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash_, _, fee := registerProvingKey(t, coordinator)
	rawData, err := coordinator.consumerABI.Pack("requestRandomness",
		common.BytesToHash(keyHash_[:]), fee, one)
	require.NoError(t, err)
	callMsg := ethereum.CallMsg{To: &coordinator.consumerContractAddress, Data: rawData}
	estimate, err := coordinator.backend.EstimateGas(context.TODO(), callMsg)
	require.NoError(t, err)
	assert.Less(t, estimate, uint64(180000))
}
