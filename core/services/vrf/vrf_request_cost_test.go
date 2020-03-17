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

}
