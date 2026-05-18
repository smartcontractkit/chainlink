package blockhashstore_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestStoreRotatesFromAddresses(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := configtest.NewTestGeneralConfig(t)
	kst := cltest.NewKeyStore(t, db)
	require.NoError(t, kst.Unlock(ctx, cltest.Password))
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, KeyStore: kst.Eth(), GeneralConfig: cfg, Client: ethClient})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	chain, err := legacyChains.Get(cltest.FixtureChainID.String())
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	ks := keystore.New(db, utils.FastScryptParams, lggr)
	require.NoError(t, ks.Unlock(ctx, "blah"))
	k1, err := ks.Eth().Create(ctx, &cltest.FixtureChainID)
	require.NoError(t, err)
	k2, err := ks.Eth().Create(ctx, &cltest.FixtureChainID)
	require.NoError(t, err)
	fromAddresses := []types.EIP55Address{k1.EIP55Address, k2.EIP55Address}
	txm := new(txmmocks.MockEvmTxManager)
	bhsAddress := common.HexToAddress("0x31Ca8bf590360B3198749f852D5c516c642846F6")

	store, err := blockhash_store.NewBlockhashStore(bhsAddress, chain.Client())
	require.NoError(t, err)
	bhs, err := blockhashstore.NewBulletproofBHS(
		chain.Config().EVM().GasEstimator(),
		cfg.Database(),
		fromAddresses,
		txm,
		store,
		nil,
		&cltest.FixtureChainID,
		ks.Eth(),
	)
	require.NoError(t, err)

	txm.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx txmgr.TxRequest) bool {
		return tx.FromAddress.String() == k1.Address.String()
	})).Once().Return(txmgr.Tx{}, nil)

	txm.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx txmgr.TxRequest) bool {
		return tx.FromAddress.String() == k2.Address.String()
	})).Once().Return(txmgr.Tx{}, nil)

	// store 2 blocks
	err = bhs.Store(ctx, 1)
	require.NoError(t, err)
	err = bhs.Store(ctx, 2)
	require.NoError(t, err)
}
