package legacygasstation_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/sqlx"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/legacygasstation/generated/bank_erc20"
	forwarder_wrapper "github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/legacygasstation/generated/legacy_gas_station_forwarder"
	lgslib "github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

const (
	BankERC20TokenName    = "BankToken"
	BankERC20TokenVersion = "1"
)

func TestIntegration_LegacyGasStation_SameChainTransfer(t *testing.T) {
	// owner deploys forwarder and token contracts
	_, owner := generateKeyAndTransactor(t, testutils.SimulatedChainID)
	// relay is a CL-owned address that posts txs
	relayKey, relay := generateKeyAndTransactor(t, testutils.SimulatedChainID)
	// sender transfers token to receiver
	senderKey, sender := generateKeyAndTransactor(t, testutils.SimulatedChainID)
	// receiver receives token from sender
	_, receiver := generateKeyAndTransactor(t, testutils.SimulatedChainID)

	genesisData := core.GenesisAlloc{
		owner.From: {Balance: assets.Ether(1000).ToInt()},
		relay.From: {Balance: assets.Ether(1000).ToInt()},
	}
	gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)

	ccipChainSelector := uint64(12345)
	// no offramp or router address needed for same-chain transfer
	dummyOffRampAddress := common.HexToAddress("0x1")
	dummyCCIPRouter := common.HexToAddress("0x2")
	forwarder, bankERC20, forwarderAddress, bankERC20Address := setupTokenAndForwarderContracts(t, owner, backend, dummyCCIPRouter, ccipChainSelector)

	amount := big.NewInt(1e18)
	transferToken(t, bankERC20, owner, sender, amount, backend)

	config, db := setUpDB(t)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend, relayKey)
	require.NoError(t, app.Start(testutils.Context(t)))
	orm := legacygasstation.NewORM(db, app.Logger, app.Config.Database())

	createLegacyGasStationServerJob(t, app, forwarderAddress, testutils.SimulatedChainID.Uint64(), ccipChainSelector, relayKey)
	statusUpdateServer := lgslib.NewUnstartedStatusUpdateServer(t)
	go statusUpdateServer.Start()
	defer statusUpdateServer.Stop()
	createLegacyGasStationSidecarJob(t, app, forwarderAddress, dummyOffRampAddress, ccipChainSelector, testutils.SimulatedChainID.Uint64(), statusUpdateServer)

	t.Run("single same-chain meta transfer", func(t *testing.T) {
		req := generateRequest(t, backend, forwarder, bankERC20Address, senderKey, receiver.From, amount, ccipChainSelector, ccipChainSelector)
		requestID := sendTransaction(t, req, app.Server.URL)
		verifySameChainTransfer(t, orm, backend, bankERC20, requestID, receiver, amount, ccipChainSelector)
	})
}

func verifyTxStatus(t *testing.T, orm lgslib.ORM, backend *backends.SimulatedBackend, requestID string, sourceChainCCIPSelector uint64, status types.Status) {
	gomega.NewWithT(t).Eventually(func() bool {
		backend.Commit()
		txs, err := orm.SelectBySourceChainIDAndStatus(testutils.Context(t), sourceChainCCIPSelector, status)
		require.NoError(t, err)
		for _, tx := range txs {
			if tx.Status == status && tx.ID == requestID {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func verifySameChainTransfer(t *testing.T, orm lgslib.ORM, backend *backends.SimulatedBackend, bankERC20 *bank_erc20.BankERC20, requestID string, receiver *bind.TransactOpts, amount *big.Int, ccipChainSelector uint64) {
	// verify same-chain meta transaction
	gomega.NewWithT(t).Eventually(func() bool {
		backend.Commit()
		amountReceived, err := bankERC20.BalanceOf(nil, receiver.From)
		require.NoError(t, err)
		return amountReceived.Cmp(amount) == 0
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// verify legacy_gasless_txs has correct status
	verifyTxStatus(t, orm, backend, requestID, ccipChainSelector, types.Finalized)
}

func ptr[T any](t T) *T { return &t }

func createLegacyGasStationServerJob(
	t *testing.T,
	app chainlink.Application,
	forwarderAddress common.Address,
	evmChainID uint64,
	ccipChainSelector uint64,
	fromKey ethkey.KeyV2,
) job.Job {
	jid := uuid.New()
	jobName := fmt.Sprintf("legacygasstationserver-%s", jid.String())
	s := testspecs.GenerateLegacyGasStationServerSpec(testspecs.LegacyGasStationServerSpecParams{
		JobID:             jid.String(),
		Name:              jobName,
		ForwarderAddress:  forwarderAddress.Hex(),
		EVMChainID:        evmChainID,
		CCIPChainSelector: ccipChainSelector,
		FromAddresses:     []string{fromKey.Address.String()},
	}).Toml()
	jb, err := legacygasstation.ValidatedServerSpec(s)
	require.NoError(t, err)
	err = app.AddJobV2(testutils.Context(t), &jb)
	require.NoError(t, err)
	// Wait until all jobs are active and listening for logs
	gomega.NewWithT(t).Eventually(func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, jb := range jbs {
			return jb.Name.String == jobName
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())
	return jb
}

func createLegacyGasStationSidecarJob(
	t *testing.T,
	app chainlink.Application,
	forwarderAddress,
	offRampAddress common.Address,
	ccipChainSelector uint64,
	evmChainID uint64,
	server lgslib.TestStatusUpdateServer,
) job.Job {
	jid := uuid.New()
	jobName := fmt.Sprintf("legacygasstationsidecar-%s", jid.String())
	s := testspecs.GenerateLegacyGasStationSidecarSpec(testspecs.LegacyGasStationSidecarSpecParams{
		JobID:             jid.String(),
		Name:              jobName,
		ForwarderAddress:  forwarderAddress.Hex(),
		OffRampAddress:    offRampAddress.Hex(),
		EVMChainID:        evmChainID,
		CCIPChainSelector: ccipChainSelector,
		StatusUpdateURL:   fmt.Sprintf("http://localhost:%d/return_success", server.Port),
	}).Toml()
	jb, err := legacygasstation.ValidatedSidecarSpec(s)
	require.NoError(t, err)
	err = app.AddJobV2(testutils.Context(t), &jb)
	require.NoError(t, err)
	// Wait until all jobs are active and listening for logs
	gomega.NewWithT(t).Eventually(func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, jb := range jbs {
			return jb.Name.String == jobName
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())
	return jb
}

func setUpForwarder(t *testing.T, owner *bind.TransactOpts, chain *backends.SimulatedBackend) (common.Address, *forwarder_wrapper.LegacyGasStationForwarder) {
	// deploys EIP 2771 forwarder contract that verifies signatures from meta transaction and forwards the call to recipient contract (i.e BankERC20 token)
	forwarderAddress, _, forwarder, err := forwarder_wrapper.DeployLegacyGasStationForwarder(owner, chain)
	require.NoError(t, err)
	chain.Commit()
	// registers EIP712-compliant domain separator for BankERC20 token
	_, err = forwarder.RegisterDomainSeparator(owner, BankERC20TokenName, BankERC20TokenVersion)
	require.NoError(t, err)
	chain.Commit()

	return forwarderAddress, forwarder
}

func setUpBankERC20(t *testing.T, owner *bind.TransactOpts, chain *backends.SimulatedBackend, forwarderAddress, routerAddress, ccipFeeProvider common.Address, totalSupply *big.Int, chainID uint64) (common.Address, *bank_erc20.BankERC20) {
	// deploys BankERC20 token that enables meta transactions for same-chain and cross-chain token transfers
	tokenAddress, _, token, err := bank_erc20.DeployBankERC20(
		owner, chain, "BankToken", "BANK", big.NewInt(0).Mul(totalSupply, big.NewInt(1e18)), forwarderAddress, routerAddress, ccipFeeProvider, chainID)
	require.NoError(t, err)
	chain.Commit()
	return tokenAddress, token
}

func transferToken(t *testing.T, token *bank_erc20.BankERC20, sender, receiver *bind.TransactOpts, amount *big.Int, chain *backends.SimulatedBackend) {
	senderBalanceBefore, err := token.BalanceOf(nil, sender.From)
	require.NoError(t, err)
	chain.Commit()

	_, err = token.Transfer(sender, receiver.From, amount)
	require.NoError(t, err)
	chain.Commit()

	receiverBal, err := token.BalanceOf(nil, receiver.From)
	require.NoError(t, err)
	require.Equal(t, amount, receiverBal)

	senderBal, err := token.BalanceOf(nil, sender.From)
	require.NoError(t, err)
	require.Equal(t, senderBalanceBefore.Sub(senderBalanceBefore, amount), senderBal)
}

func setUpDB(t *testing.T) (cfg chainlink.GeneralConfig, db *sqlx.DB) {
	cfg, db = heavyweight.FullTestDBV2(t, "legacy_gas_station_integration_test", func(c *chainlink.Config, s *chainlink.Secrets) {
		require.Zero(t, testutils.SimulatedChainID.Cmp(c.EVM[0].ChainID.ToInt()))
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](2_000_000)

		c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100)
		c.EVM[0].HeadTracker.SamplingInterval = models.MustNewDuration(0) // Head sampling disabled

		c.EVM[0].Transactions.ResendAfterThreshold = models.MustNewDuration(0)
		c.EVM[0].Transactions.ReaperThreshold = models.MustNewDuration(100 * time.Millisecond)

		c.EVM[0].FinalityDepth = ptr[uint32](1)
		c.Feature.LogPoller = ptr(true)
		c.Feature.EAL = ptr(true)
	})
	return
}

func sendTransaction(t *testing.T, req types.SendTransactionRequest, url string) string {
	body, err := json.Marshal(req)
	require.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/%s", url, "gasstation/send_transaction"),
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	var jsonResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	require.NoError(t, err)
	require.Equal(t, "2.0", jsonResp["jsonrpc"])
	require.NotNil(t, jsonResp["result"])
	decodedJsonResp, ok := jsonResp["result"].(map[string]interface{})
	require.True(t, ok)
	require.NotNil(t, decodedJsonResp["request_id"])
	requestID, ok := decodedJsonResp["request_id"].(string)
	require.True(t, ok)
	return uuid.MustParse(requestID).String()
}

func setupTokenAndForwarderContracts(
	t *testing.T,
	deployer *bind.TransactOpts,
	backend *backends.SimulatedBackend,
	ccipRouterAddress common.Address,
	ccipChainSelector uint64,
) (*forwarder_wrapper.LegacyGasStationForwarder, *bank_erc20.BankERC20, common.Address, common.Address) {
	forwarderAddress, forwarder := setUpForwarder(t, deployer, backend)
	bankERC20Address, bankERC20 := setUpBankERC20(t, deployer, backend, forwarderAddress, ccipRouterAddress, deployer.From, big.NewInt(1e9), ccipChainSelector)
	return forwarder, bankERC20, forwarderAddress, bankERC20Address
}

func generateKeyAndTransactor(t *testing.T, chainID *big.Int) (key ethkey.KeyV2, transactor *bind.TransactOpts) {
	key = cltest.MustGenerateRandomKey(t)
	transactor, err := bind.NewKeyedTransactorWithChainID(key.ToEcdsaPrivKey(), chainID)
	require.NoError(t, err)
	return
}

func generateRequest(
	t *testing.T,
	backend *backends.SimulatedBackend,
	forwarder *forwarder_wrapper.LegacyGasStationForwarder,
	bankERC20Address common.Address,
	senderKey ethkey.KeyV2,
	receiver common.Address,
	amount *big.Int,
	sourceChainSelector uint64,
	destChainSelector uint64,
) types.SendTransactionRequest {
	_, calldataHash, err := lgslib.GenerateMetaTransferCalldata(receiver, amount, destChainSelector)
	require.NoError(t, err)
	deadline := big.NewInt(int64(backend.Blockchain().CurrentHeader().Time + uint64(time.Hour)))
	signature, _, _, forwarderNonce, err := lgslib.SignMetaTransfer(*forwarder, senderKey.ToEcdsaPrivKey(), senderKey.Address, bankERC20Address, calldataHash, deadline, BankERC20TokenName, BankERC20TokenVersion)
	require.NoError(t, err)

	return types.SendTransactionRequest{
		From:               senderKey.Address,
		Target:             bankERC20Address,
		TargetName:         BankERC20TokenName,
		Version:            BankERC20TokenVersion,
		Nonce:              forwarderNonce,
		Receiver:           receiver,
		Amount:             amount,
		SourceChainID:      sourceChainSelector,
		DestinationChainID: destChainSelector,
		ValidUntilTime:     deadline,
		Signature:          signature,
	}
}
