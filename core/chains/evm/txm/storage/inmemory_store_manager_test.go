package storage

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
)

func TestAddAndRemove(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStoreManager(logger.Test(t), testutils.FixtureChainID)
	// Adds a new address
	m.Add(fromAddress)
	assert.Len(t, m.InMemoryStoreMap, 1)

	// Noops if address exists
	m.Add(fromAddress)
	assert.Len(t, m.InMemoryStoreMap, 1)

	// Adds multiple addresses
	fromAddress1 := testutils.NewAddress()
	fromAddress2 := testutils.NewAddress()
	addresses := []common.Address{fromAddress1, fromAddress2}
	m.Add(addresses...)
	assert.Len(t, m.InMemoryStoreMap, 3)

	// Remove an address
	m.Remove(fromAddress1)
	assert.Len(t, m.InMemoryStoreMap, 2)

	// Noops if address doesn't exist
	m.Remove(testutils.NewAddress())
	assert.Len(t, m.InMemoryStoreMap, 2)
}
