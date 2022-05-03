package forwarders_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/authorized_receiver"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/test-go/testify/mock"
	"github.com/test-go/testify/require"
)

var GetAuthorisedSendersABI = evmtypes.MustGetABI(authorized_receiver.AuthorizedReceiverABI).Methods["getAuthorizedSenders"]

var SimpleOracleCallABI = evmtypes.MustGetABI(operator_wrapper.OperatorABI).Methods["getChainlinkToken"]

func TestFwdMgr(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

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

	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewPGCfg(true)),
		client.NewSimulatedBackendClient(t, ec, testutils.FixtureChainID), lggr, 100*time.Millisecond, 2, 3)
	fwdMgr := forwarders.NewFwdMgr(db, ethClient, lp, lggr, pgtest.NewPGCfg(true))
	fwdMgr.ORM = forwarders.NewORM(db, logger.TestLogger(t), cfg)

	_, err = fwdMgr.ORM.CreateForwarder(forwarderAddr, utils.Big(*testutils.FixtureChainID))
	require.NoError(t, err)

	lst, err := fwdMgr.ORM.FindForwardersByChain(utils.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	require.Equal(t, len(lst), 1)
	require.Equal(t, lst[0].Address, forwarderAddr)

	// Mocking getAuthorisedSenders on forwarder to return EOA
	ethClient.On("CallContract", mock.Anything,
		ethereum.CallMsg{From: common.HexToAddress("0x0"), To: &forwarderAddr, Data: []uint8{0x24, 0x8, 0xaf, 0xaa}},
		mock.Anything).Return(genAuthorisedSenders(t, []common.Address{owner.From}), nil)

	// Mocking getAuthorisedSenders on operator to return forwarder
	ethClient.On("CallContract", mock.Anything,
		ethereum.CallMsg{From: common.HexToAddress("0x0"), To: &operatorAddr, Data: []uint8{0x24, 0x8, 0xaf, 0xaa}},
		mock.Anything).Return(genAuthorisedSenders(t, []common.Address{forwarderAddr}), nil)

	err = fwdMgr.Start()
	require.NoError(t, err)

	f, _, err := fwdMgr.MaybeForwardTransaction(owner.From, operatorAddr, getSimpleOperatorCall(t))
	require.NoError(t, err)
	require.Equal(t, f, forwarderAddr)

	err = fwdMgr.Stop()
	require.NoError(t, err)
}

func genAuthorisedSenders(t *testing.T, addrs []common.Address) []byte {
	args, err := GetAuthorisedSendersABI.Outputs.Pack(addrs)
	require.NoError(t, err)

	dataBytes := append(GetAuthorisedSendersABI.ID, args...)
	require.NotEmpty(t, dataBytes)
	return args
}

func getSimpleOperatorCall(t *testing.T) []byte {
	args, err := SimpleOracleCallABI.Inputs.Pack()
	require.NoError(t, err)

	dataBytes := append(SimpleOracleCallABI.ID, args...)
	require.NotEmpty(t, dataBytes)
	return args
}
