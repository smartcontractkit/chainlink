package forwarders_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/authorized_receiver"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var GetAuthorisedSendersABI = evmtypes.MustGetABI(authorized_receiver.AuthorizedReceiverABI).Methods["getAuthorizedSenders"]

var SimpleOracleCallABI = evmtypes.MustGetABI(operator_wrapper.OperatorABI).Methods["getChainlinkToken"]

func TestFwdMgr_MaybeForwardTransaction(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
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
	_, err = forwarder.SetAuthorizedSenders(owner, []common.Address{owner.From})
	require.NoError(t, err)
	ec.Commit()
	authorized, err := forwarder.GetAuthorizedSenders(nil)
	require.NoError(t, err)
	t.Log(authorized)

	evmClient := client.NewSimulatedBackendClient(t, ec, testutils.FixtureChainID)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewQConfig(true)), evmClient, lggr, 100*time.Millisecond, 2, 3, 2, 1000)
	fwdMgr := forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg)
	fwdMgr.ORM = forwarders.NewORM(db, logger.TestLogger(t), cfg)

	_, err = fwdMgr.ORM.CreateForwarder(forwarderAddr, utils.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	lst, err := fwdMgr.ORM.FindForwardersByChain(utils.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	require.Equal(t, len(lst), 1)
	require.Equal(t, lst[0].Address, forwarderAddr)

	require.NoError(t, fwdMgr.Start(testutils.Context(t)))
	addr, err := fwdMgr.GetForwarderForEOA(owner.From)
	require.NoError(t, err)
	require.Equal(t, addr, forwarderAddr)
	err = fwdMgr.Close()
	require.NoError(t, err)
}

func TestFwdMgr_AccountUnauthorizedToForward_SkipsForwarding(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
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
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewQConfig(true)), evmClient, lggr, 100*time.Millisecond, 2, 3, 2, 1000)
	fwdMgr := forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg)
	fwdMgr.ORM = forwarders.NewORM(db, logger.TestLogger(t), cfg)

	_, err = fwdMgr.ORM.CreateForwarder(forwarderAddr, utils.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	lst, err := fwdMgr.ORM.FindForwardersByChain(utils.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	require.Equal(t, len(lst), 1)
	require.Equal(t, lst[0].Address, forwarderAddr)

	err = fwdMgr.Start(testutils.Context(t))
	require.NoError(t, err)
	addr, err := fwdMgr.GetForwarderForEOA(owner.From)
	require.ErrorContains(t, err, "Cannot find forwarder for given EOA")
	require.Equal(t, addr, common.Address{})
	err = fwdMgr.Close()
	require.NoError(t, err)
}
