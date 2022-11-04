package evm_test

import (
	"math/big"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type legacyEthNodeConfig struct {
	defaultChainID        *big.Int
	ethereumURL           string
	ethereumHTTPURL       *url.URL
	ethereumSecondaryURLs []url.URL
	evmNodes              string
	pg.QConfig
}

func (c legacyEthNodeConfig) DefaultChainID() *big.Int {
	return c.defaultChainID
}

func (c legacyEthNodeConfig) EthereumURL() string {
	return c.ethereumURL
}

func (c legacyEthNodeConfig) EthereumHTTPURL() *url.URL {
	return c.ethereumHTTPURL
}

func (c legacyEthNodeConfig) EthereumSecondaryURLs() []url.URL {
	return c.ethereumSecondaryURLs
}

func (c legacyEthNodeConfig) EthereumNodes() string {
	return c.evmNodes
}

func Test_ClobberDBFromEnv(t *testing.T) {
	var fixtureChains int64 = 2
	var fixtureNodes int64 = 1

	t.Run("inserts nodes from config", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)

		cfg := legacyEthNodeConfig{
			defaultChainID:        big.NewInt(42),
			ethereumURL:           "ws://example.com/foo/ws",
			ethereumHTTPURL:       cltest.MustParseURL(t, "http://example.com/foo"),
			ethereumSecondaryURLs: []url.URL{*cltest.MustParseURL(t, "http://secondary1.example/foo"), *cltest.MustParseURL(t, "https://secondary2.example/bar")},
			QConfig:               pgtest.NewQConfig(false),
		}

		err := evm.ClobberDBFromEnv(db, cfg, logger.TestLogger(t))
		require.NoError(t, err)

		cltest.AssertCount(t, db, "evm_chains", fixtureChains+1)
		cltest.AssertCount(t, db, "evm_nodes", fixtureNodes+3)

		var primaryNode evmtypes.Node
		err = db.Get(&primaryNode, `SELECT * FROM evm_nodes WHERE evm_chain_id = 42 AND NOT send_only`)
		require.NoError(t, err)

		assert.Equal(t, "primary-0-42", primaryNode.Name)
		assert.Equal(t, cfg.defaultChainID.String(), primaryNode.EVMChainID.String())
		assert.True(t, primaryNode.WSURL.Valid)
		assert.Equal(t, cfg.ethereumURL, primaryNode.WSURL.String)
		assert.True(t, primaryNode.HTTPURL.Valid)
		assert.Equal(t, cfg.ethereumHTTPURL.String(), primaryNode.HTTPURL.String)
		assert.False(t, primaryNode.SendOnly)

		var sendonlyNodes []evmtypes.Node
		err = db.Select(&sendonlyNodes, `SELECT * FROM evm_nodes WHERE evm_chain_id = 42 AND send_only ORDER BY http_url`)
		require.NoError(t, err)
		require.Len(t, sendonlyNodes, 2)

		assert.True(t, sendonlyNodes[0].SendOnly)
		assert.Equal(t, "sendonly-0-42", sendonlyNodes[0].Name)
		assert.False(t, sendonlyNodes[0].WSURL.Valid)
		assert.True(t, sendonlyNodes[0].HTTPURL.Valid)
		assert.Equal(t, "http://secondary1.example/foo", sendonlyNodes[0].HTTPURL.String)

		assert.True(t, sendonlyNodes[1].SendOnly)
		assert.Equal(t, "sendonly-1-42", sendonlyNodes[1].Name)
		assert.False(t, sendonlyNodes[1].WSURL.Valid)
		assert.True(t, sendonlyNodes[1].HTTPURL.Valid)
		assert.Equal(t, "https://secondary2.example/bar", sendonlyNodes[1].HTTPURL.String)
	})

	t.Run("ignores sendonlys if they duplicate the primary URL", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := legacyEthNodeConfig{
			defaultChainID:        big.NewInt(42),
			ethereumURL:           "ws://example.com/foo/ws",
			ethereumHTTPURL:       cltest.MustParseURL(t, "http://example.com/foo"),
			ethereumSecondaryURLs: []url.URL{*cltest.MustParseURL(t, "http://example.com/foo"), *cltest.MustParseURL(t, "https://secondary2.example/bar")},
			QConfig:               pgtest.NewQConfig(false),
		}

		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		err := evm.ClobberDBFromEnv(db, cfg, lggr)
		require.NoError(t, err)

		testutils.RequireLogMessage(t, observedLogs, "Got secondary URL http://example.com/foo which is already specified as a primary node HTTP URL. It does not make any sense to have both primary and sendonly nodes sharing the same URL and does not grant any extra redundancy. This sendonly RPC url will be ignored.")
	})
}

func TestSetupNodes(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	// Insert existing node which will be erased
	pgtest.MustExec(t, db, `INSERT INTO evm_nodes (name, evm_chain_id, ws_url, http_url, send_only, created_at, updated_at) VALUES ('foo','0','ws://example.com',null,false,NOW(),NOW())`)

	s := `
[
	{
		"name": "primary_0_1",
		"evmChainId": "0",
		"wsUrl": "ws://test1.invalid",
		"sendOnly": false
	},
	{
		"name": "primary_0_2",
		"evmChainId": "0",
		"wsUrl": "ws://test2.invalid",
		"httpUrl": "https://test2.invalid",
		"sendOnly": false
	},
	{
		"name": "primary_1337_1",
		"evmChainId": "1337",
		"wsUrl": "ws://test3.invalid",
		"httpUrl": "http://test3.invalid",
		"sendOnly": false
	},
	{
		"name": "sendonly_1337_1",
		"evmChainId": "1337",
		"httpUrl": "http://test4.invalid",
		"sendOnly": true
	},
	{
		"name": "sendonly_0_1",
		"evmChainId": "0",
		"httpUrl": "http://test5.invalid",
		"sendOnly": true
	},
	{
		"name": "primary_42_1",
		"evmChainId": "42",
		"wsUrl": "ws://test6.invalid",
		"sendOnly": false
	},
	{
		"name": "sendonly_43_1",
		"evmChainId": "43",
		"httpUrl": "http://test7.invalid",
		"sendOnly": true
	},
	{
		"name": "zzzz this will be ignored due to duplicate ws url",
		"evmChainId": "0",
		"wsUrl": "ws://test1.invalid",
		"sendOnly": false
	},
	{
		"name": "zzzz this will be ignored due to duplicate http url",
		"evmChainId": "0",
		"wsUrl": "ws://test8.invalid",
		"httpUrl": "https://test2.invalid",
		"sendOnly": false
	}
]
	`

	cfg := legacyEthNodeConfig{
		evmNodes: s,
		QConfig:  pgtest.NewQConfig(false),
	}

	err := evm.SetupNodes(db, cfg, logger.TestLogger(t))
	require.NoError(t, err)

	cltest.AssertCount(t, db, "evm_nodes", 7)

	var nodes []evmtypes.Node
	err = db.Select(&nodes, `SELECT * FROM evm_nodes ORDER BY name ASC`)
	require.NoError(t, err)

	require.Len(t, nodes, 7)

	assert.Equal(t, "primary_0_1", nodes[0].Name)
	assert.Equal(t, "primary_0_2", nodes[1].Name)
	assert.Equal(t, "primary_1337_1", nodes[2].Name)
	assert.Equal(t, "primary_42_1", nodes[3].Name)
	assert.Equal(t, "sendonly_0_1", nodes[4].Name)
	assert.Equal(t, "sendonly_1337_1", nodes[5].Name)
	assert.Equal(t, "sendonly_43_1", nodes[6].Name)
}
