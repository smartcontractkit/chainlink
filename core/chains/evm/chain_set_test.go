package evm_test

import (
	"context"
	"testing"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v4"
)

func Test_ChainSet_AddNode(t *testing.T) {
	c := configtest.NewTestGeneralConfig(t)

	t.Run("adding primary node", func(t *testing.T) {
		db := pgtest.NewGormDB(t)
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: ethClient, GeneralConfig: c, ChainCfg: evmtypes.ChainCfg{}})

		ethClient.On("AddNodeToPool", mock.Anything, mock.MatchedBy(func(n eth.Node) bool {
			return n.Name() == "test primary"
		})).Return(nil)

		newNode := evmtypes.NewNode{
			Name:  "test primary",
			WSURL: null.StringFrom("ws://test.invalid/ws"),
		}
		n, err := cc.AddNode(context.Background(), newNode)
		require.NoError(t, err)

		assert.Equal(t, "test primary", n.Name)

		ethClient.AssertExpectations(t)
	})

	t.Run("adding send only node", func(t *testing.T) {
		db := pgtest.NewGormDB(t)
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: ethClient, GeneralConfig: c, ChainCfg: evmtypes.ChainCfg{}})

		ethClient.On("AddSendOnlyNodeToPool", mock.Anything, mock.MatchedBy(func(n eth.SendOnlyNode) bool {
			return n.Name() == "test sendonly"
		})).Return(nil)

		newNode := evmtypes.NewNode{
			Name:     "test sendonly",
			HTTPURL:  null.StringFrom("http://test.invalid"),
			SendOnly: true,
		}
		n, err := cc.AddNode(context.Background(), newNode)
		require.NoError(t, err)

		assert.Equal(t, "test sendonly", n.Name)

		ethClient.AssertExpectations(t)
	})

	t.Run("adding invalid node returns error", func(t *testing.T) {
		db := pgtest.NewGormDB(t)
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: ethClient, GeneralConfig: c, ChainCfg: evmtypes.ChainCfg{}})
		newNode := evmtypes.NewNode{}
		_, err := cc.AddNode(context.Background(), newNode)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nodes_name_check")

		ethClient.AssertExpectations(t)
	})
}

func Test_ChainSet_RemoveNode(t *testing.T) {
	t.Fatal("todo")
}
