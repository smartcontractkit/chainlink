package forwarders_test

import (
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers2/testocr2aggregator"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

var GetAuthorisedSendersABI = evmtypes.MustGetABI(authorized_receiver.AuthorizedReceiverABI).Methods["getAuthorizedSenders"]

var SimpleOracleCallABI = evmtypes.MustGetABI(operator_wrapper.OperatorABI).Methods["getChainlinkToken"]

func TestFwdMgr_MaybeForwardTransaction(t *testing.T) {
	lggr := logger.Test(t)
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	owner := testutils.MustNewSimTransactor(t)
	ctx := testutils.Context(t)

	ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	t.Cleanup(func() { ec.Close() })
	linkAddr := common.HexToAddress("0x01BE23585060835E02B77ef475b0Cc51aA1e0709")
	operatorAddr, _, _, err := operator_wrapper.DeployOperator(owner, ec, linkAddr, owner.From)
	require.NoError(t, err)
	forwarderAddr, _, forwarder, err := authorized_forwarder.DeployAuthorizedForwarder(owner, ec, linkAddr, owner.From, operatorAddr, []byte{})
	require.NoError(t, err)
	ec.Commit()
	_, err = forwarder.SetAuthorizedSenders(owner, []common.Address{owner.From})
	require.NoError(t, err)
	ec.Commit()
	authorized, err := forwarder.GetAuthorizedSenders(nil)
	require.NoError(t, err)
	t.Log(authorized)

	evmClient := client.NewSimulatedBackendClient(t, ec, testutils.FixtureChainID)

	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), evmClient, lggr, lpOpts)
	fwdMgr := forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg.EVM())
	fwdMgr.ORM = forwarders.NewORM(db)

	fwd, err := fwdMgr.ORM.CreateForwarder(ctx, forwarderAddr, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	lst, err := fwdMgr.ORM.FindForwardersByChain(ctx, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	require.Equal(t, len(lst), 1)
	require.Equal(t, lst[0].Address, forwarderAddr)

	require.NoError(t, fwdMgr.Start(testutils.Context(t)))
	addr, err := fwdMgr.ForwarderFor(ctx, owner.From)
	require.NoError(t, err)
	require.Equal(t, addr.String(), forwarderAddr.String())
	err = fwdMgr.Close()
	require.NoError(t, err)

	cleanupCalled := false
	cleanup := func(tx sqlutil.DataSource, evmChainId int64, addr common.Address) error {
		require.Equal(t, testutils.FixtureChainID.Int64(), evmChainId)
		require.Equal(t, forwarderAddr, addr)
		require.NotNil(t, tx)
		cleanupCalled = true
		return nil
	}

	err = fwdMgr.ORM.DeleteForwarder(ctx, fwd.ID, cleanup)
	assert.NoError(t, err)
	assert.True(t, cleanupCalled)
}

func TestFwdMgr_AccountUnauthorizedToForward_SkipsForwarding(t *testing.T) {
	lggr := logger.Test(t)
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	cfg := configtest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	owner := testutils.MustNewSimTransactor(t)
	ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	t.Cleanup(func() { ec.Close() })
	linkAddr := common.HexToAddress("0x01BE23585060835E02B77ef475b0Cc51aA1e0709")
	operatorAddr, _, _, err := operator_wrapper.DeployOperator(owner, ec, linkAddr, owner.From)
	require.NoError(t, err)

	forwarderAddr, _, _, err := authorized_forwarder.DeployAuthorizedForwarder(owner, ec, linkAddr, owner.From, operatorAddr, []byte{})
	require.NoError(t, err)
	ec.Commit()

	evmClient := client.NewSimulatedBackendClient(t, ec, testutils.FixtureChainID)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), evmClient, lggr, lpOpts)
	fwdMgr := forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg.EVM())
	fwdMgr.ORM = forwarders.NewORM(db)

	_, err = fwdMgr.ORM.CreateForwarder(ctx, forwarderAddr, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	lst, err := fwdMgr.ORM.FindForwardersByChain(ctx, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	require.Equal(t, len(lst), 1)
	require.Equal(t, lst[0].Address, forwarderAddr)

	err = fwdMgr.Start(testutils.Context(t))
	require.NoError(t, err)
	addr, err := fwdMgr.ForwarderFor(ctx, owner.From)
	require.ErrorIs(t, err, forwarders.ErrForwarderForEOANotFound)
	require.True(t, utils.IsZero(addr))
	err = fwdMgr.Close()
	require.NoError(t, err)
}

func TestFwdMgr_InvalidForwarderForOCR2FeedsStates(t *testing.T) {
	lggr := logger.Test(t)
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	cfg := configtest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	owner := testutils.MustNewSimTransactor(t)
	ec := backends.NewSimulatedBackend(map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	t.Cleanup(func() { ec.Close() })
	linkAddr := common.HexToAddress("0x01BE23585060835E02B77ef475b0Cc51aA1e0709")
	operatorAddr, _, _, err := operator_wrapper.DeployOperator(owner, ec, linkAddr, owner.From)
	require.NoError(t, err)

	forwarderAddr, _, forwarder, err := authorized_forwarder.DeployAuthorizedForwarder(owner, ec, linkAddr, owner.From, operatorAddr, []byte{})
	require.NoError(t, err)
	ec.Commit()

	accessAddress, _, _, err := testocr2aggregator.DeploySimpleWriteAccessController(owner, ec)
	require.NoError(t, err, "failed to deploy test access controller contract")
	ocr2Address, _, ocr2, err := testocr2aggregator.DeployOCR2Aggregator(
		owner,
		ec,
		linkAddr,
		big.NewInt(0),
		big.NewInt(10),
		accessAddress,
		accessAddress,
		9,
		"TEST",
	)
	require.NoError(t, err, "failed to deploy ocr2 test aggregator")
	ec.Commit()

	evmClient := client.NewSimulatedBackendClient(t, ec, testutils.FixtureChainID)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), evmClient, lggr, lpOpts)
	fwdMgr := forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg.EVM())
	fwdMgr.ORM = forwarders.NewORM(db)

	_, err = fwdMgr.ORM.CreateForwarder(ctx, forwarderAddr, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	lst, err := fwdMgr.ORM.FindForwardersByChain(ctx, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	require.Equal(t, len(lst), 1)
	require.Equal(t, lst[0].Address, forwarderAddr)

	fwdMgr = forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg.EVM())
	require.NoError(t, fwdMgr.Start(testutils.Context(t)))
	// cannot find forwarder because it isn't authorized nor added as a transmitter
	addr, err := fwdMgr.ForwarderForOCR2Feeds(ctx, owner.From, ocr2Address)
	require.ErrorIs(t, err, forwarders.ErrForwarderForEOANotFound)
	require.True(t, utils.IsZero(addr))

	_, err = forwarder.SetAuthorizedSenders(owner, []common.Address{owner.From})
	require.NoError(t, err)
	ec.Commit()

	authorizedSenders, err := forwarder.GetAuthorizedSenders(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	require.Equal(t, owner.From, authorizedSenders[0])

	// cannot find forwarder because it isn't added as a transmitter
	addr, err = fwdMgr.ForwarderForOCR2Feeds(ctx, owner.From, ocr2Address)
	require.ErrorIs(t, err, forwarders.ErrForwarderForEOANotFound)
	require.True(t, utils.IsZero(addr))

	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(big.NewInt(0), big.NewInt(10))
	require.NoError(t, err)

	_, err = ocr2.SetConfig(owner,
		[]common.Address{testutils.NewAddress(), testutils.NewAddress(), testutils.NewAddress(), testutils.NewAddress()},
		[]common.Address{forwarderAddr, testutils.NewAddress(), testutils.NewAddress(), testutils.NewAddress()},
		1,
		onchainConfig,
		0,
		[]byte{})
	require.NoError(t, err)
	ec.Commit()

	transmitters, err := ocr2.GetTransmitters(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	require.True(t, slices.Contains(transmitters, forwarderAddr))

	// create new fwd to have an empty cache that has to fetch authorized forwarders from log poller
	fwdMgr = forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg.EVM())
	require.NoError(t, fwdMgr.Start(testutils.Context(t)))
	addr, err = fwdMgr.ForwarderForOCR2Feeds(ctx, owner.From, ocr2Address)
	require.NoError(t, err, "forwarder should be valid and found because it is both authorized and set as a transmitter")
	require.Equal(t, forwarderAddr, addr)
	require.NoError(t, fwdMgr.Close())
}
