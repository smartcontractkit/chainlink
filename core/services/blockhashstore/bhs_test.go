package blockhashstore_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
)

func TestStoreRotatesFromAddresses(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	cfg := configtest.NewTestGeneralConfig(t)
	kst := cltest.NewKeyStore(t, db)
	require.NoError(t, kst.Unlock(ctx, cltest.Password))
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		KeyStore:       kst.Eth(),
		DB:             db,
		Client:         ethClient,
	})
	chain, err := legacyChains.Get(cltest.FixtureChainID.String())
	require.NoError(t, err)
	coreKS := keystest.NewMemoryChainStore()
	ks := keys.NewStore(coreKS)
	addr1 := coreKS.MustCreate(t)
	addr2 := coreKS.MustCreate(t)
	fromAddresses := []types.EIP55Address{types.EIP55AddressFromAddress(addr1), types.EIP55AddressFromAddress(addr2)}
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
		ks,
	)
	require.NoError(t, err)

	txm.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx txmgr.TxRequest) bool {
		return tx.FromAddress.String() == addr1.String()
	})).Once().Return(txmgr.Tx{}, nil)

	txm.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx txmgr.TxRequest) bool {
		return tx.FromAddress.String() == addr2.String()
	})).Once().Return(txmgr.Tx{}, nil)

	// store 2 blocks
	err = bhs.Store(ctx, 1)
	require.NoError(t, err)
	err = bhs.Store(ctx, 2)
	require.NoError(t, err)
}
