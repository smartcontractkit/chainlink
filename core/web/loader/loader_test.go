package loader

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmORMMocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	coremocks "github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestLoader_Chains(t *testing.T) {
	t.Parallel()

	emvORM := &evmORMMocks.ORM{}
	app := &coremocks.Application{}
	ctx := InjectDataloader(context.Background(), app)

	defer t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, app, emvORM)
	})

	id := utils.Big{}
	err := id.UnmarshalText([]byte("123"))
	require.NoError(t, err)

	chain := types.Chain{
		ID:      id,
		Enabled: true,
	}

	emvORM.On("GetChainsByIDs", []utils.Big{id}).Return([]types.Chain{
		chain,
	}, nil)
	app.On("EVMORM").Return(emvORM)

	found, err := GetChainByID(ctx, "123")
	require.NoError(t, err)

	assert.Equal(t, chain, *found)
}

func TestLoader_Nodes(t *testing.T) {
	t.Parallel()

	emvORM := &evmORMMocks.ORM{}
	app := &coremocks.Application{}
	ctx := InjectDataloader(context.Background(), app)

	defer t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t, app, emvORM)
	})

	id := int32(1)
	chainId := utils.Big{}
	err := chainId.UnmarshalText([]byte("123"))
	require.NoError(t, err)

	node1 := types.Node{
		ID:         id,
		Name:       "test-node-1",
		EVMChainID: chainId,
	}
	node2 := types.Node{
		ID:         id,
		Name:       "test-node-1",
		EVMChainID: chainId,
	}

	emvORM.On("GetNodesByChainIDs", []utils.Big{chainId}).Return([]types.Node{
		node1, node2,
	}, nil)
	app.On("EVMORM").Return(emvORM)

	found, err := GetNodesByChainID(ctx, "123")
	require.NoError(t, err)

	assert.Equal(t, node1, found[0])
	assert.Equal(t, node2, found[1])
}
