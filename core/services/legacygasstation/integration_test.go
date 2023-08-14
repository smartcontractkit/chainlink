package legacygasstation_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	geth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/sqlx"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/bank_erc20"
	forwarder_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/metatx"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	statusUpdateServer := legacygasstation.NewUnstartedStatusUpdateServer(t)
	go statusUpdateServer.Start()
	defer statusUpdateServer.Stop()
	createLegacyGasStationSidecarJob(t, app, forwarderAddress, dummyOffRampAddress, ccipChainSelector, testutils.SimulatedChainID.Uint64(), statusUpdateServer)

	t.Run("single same-chain meta transfer", func(t *testing.T) {
		req := generateRequest(t, backend, forwarder, bankERC20Address, senderKey, receiver.From, amount, ccipChainSelector, ccipChainSelector)
		requestID := sendTransaction(t, req, app.Server.URL)
		verifySameChainTransfer(t, orm, backend, bankERC20, requestID, receiver, amount, ccipChainSelector)
	})
}

func verifyTxStatus(t *testing.T, orm legacygasstation.ORM, backend *backends.SimulatedBackend, requestID string, sourceChainCCIPSelector uint64, status types.Status) {
	gomega.NewWithT(t).Eventually(func() bool {
		backend.Commit()
		txs, err := orm.SelectBySourceChainIDAndStatus(sourceChainCCIPSelector, status)
		require.NoError(t, err)
		for _, tx := range txs {
			if tx.Status == status && tx.ID == requestID {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func verifySameChainTransfer(t *testing.T, orm legacygasstation.ORM, backend *backends.SimulatedBackend, bankERC20 *bank_erc20.BankERC20, requestID string, receiver *bind.TransactOpts, amount *big.Int, ccipChainSelector uint64) {
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

// TestIntegration_LegacyGasStation_CrossChainTransfer_SourceChain tests cross chain transfer
// CCIP DON is not involved in this test and validations are limited to source chain
// i.e. Ensures that sidecar updates transaction status to SourceFinalized
func TestIntegration_LegacyGasStation_CrossChainTransfer_SourceChain(t *testing.T) {
	sourceChainID := uint64(1337)
	destChainID := uint64(1000)

	ccipContracts := testhelpers.SetupCCIPContracts(t, sourceChainID, destChainID)

	// relay is a CL-owned address that posts txs
	relayKey, relay := generateKeyAndTransactor(t, ccipContracts.Source.Chain.Blockchain().Config().ChainID)
	// sender transfers token to receiver
	senderKey, sender := generateKeyAndTransactor(t, ccipContracts.Source.Chain.Blockchain().Config().ChainID)
	// receiver receives token from sender
	_, receiver := generateKeyAndTransactor(t, ccipContracts.Dest.Chain.Blockchain().Config().ChainID)

	owner := ccipContracts.Source.User
	sourceBackend := ccipContracts.Source.Chain

	forwarder, bankERC20, forwarderAddress, bankERC20Address := setupTokenAndForwarderContracts(t, owner, sourceBackend, ccipContracts.Source.Router.Address(), sourceChainID)
	_, _, err := ccipContracts.SetupLockAndMintTokenPool(bankERC20Address, "WrappedBankToken", "WBANK")
	require.NoError(t, err)

	amount := big.NewInt(1e18)
	transferToken(t, bankERC20, owner, sender, amount, sourceBackend)

	config, db := setUpDB(t)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, sourceBackend, relayKey)
	require.NoError(t, app.Start(testutils.Context(t)))

	ccipFeeBudget := big.NewInt(3e18)
	transferNative(t, owner, bankERC20Address, 50_000, ccipFeeBudget, sourceBackend)
	transferNative(t, owner, relay.From, 21_000, amount, sourceBackend)

	dummySourceOffRampAddress := common.HexToAddress("0x1")
	createLegacyGasStationServerJob(t, app, forwarderAddress, sourceChainID, sourceChainID, relayKey)
	statusUpdateServer := legacygasstation.NewUnstartedStatusUpdateServer(t)
	go statusUpdateServer.Start()
	defer statusUpdateServer.Stop()
	createLegacyGasStationSidecarJob(t, app, forwarderAddress, dummySourceOffRampAddress, sourceChainID, sourceChainID, statusUpdateServer)

	req := generateRequest(t, sourceBackend, forwarder, bankERC20Address, senderKey, receiver.From, amount, sourceChainID, destChainID)
	requestID := sendTransaction(t, req, app.Server.URL)

	orm := legacygasstation.NewORM(db, app.GetLogger(), app.GetConfig().Database())
	// verify that sender balance has been decremented
	// the tokens are transferred to token pool
	gomega.NewWithT(t).Eventually(func() bool {
		ccipContracts.Source.Chain.Commit()
		senderBalance, err := bankERC20.BalanceOf(nil, sender.From)
		require.NoError(t, err)
		return senderBalance.Cmp(amount) == -1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// verify legacy_gasless_txs has correct status
	verifyTxStatus(t, orm, ccipContracts.Source.Chain, requestID, sourceChainID, types.SourceFinalized)
}

// TestIntegration_LegacyGasStation_CrossChainTransfer_DestChain tests cross chain transfer
// sets up CCIP DON for cross-chain transfer and validates destination chain sidecar has successfully picked up
// the destination chain off-ramp event and updated status to Finalized
func TestIntegration_LegacyGasStation_CrossChainTransfer_DestChain(t *testing.T) {
	destCCIPChainSelector := testhelpers.DestChainID
	ccipContracts := integrationtesthelpers.SetupCCIPIntegrationTH(t, testhelpers.SourceChainID, testhelpers.DestChainID)
	// relay is a CL-owned address that posts txs
	relayKey, relay := generateKeyAndTransactor(t, ccipContracts.Source.Chain.Blockchain().Config().ChainID)
	// sender transfers token to receiver
	senderKey, sender := generateKeyAndTransactor(t, ccipContracts.Source.Chain.Blockchain().Config().ChainID)
	// receiver receives token from sender
	_, receiver := generateKeyAndTransactor(t, ccipContracts.Dest.Chain.Blockchain().Config().ChainID)

	owner := ccipContracts.Source.User
	sourceBackend := ccipContracts.Source.Chain
	destBackend := ccipContracts.Dest.Chain

	forwarder, bankERC20, forwarderAddress, bankERC20Address := setupTokenAndForwarderContracts(t, owner, sourceBackend, ccipContracts.Source.Router.Address(), testhelpers.SourceChainID)
	_, destToken, err := ccipContracts.SetupLockAndMintTokenPool(bankERC20Address, "WrappedBankToken", "WBANK")
	require.NoError(t, err)

	amount := big.NewInt(1e18)
	transferToken(t, bankERC20, owner, sender, amount, sourceBackend)

	config, db := setUpDB(t)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, destBackend, relayKey)
	require.NoError(t, app.Start(testutils.Context(t)))

	// fund BankERC20 contract with native ETH. Test setup may use a low Eth price, important to send enough.
	ccipFeeBudget := big.NewInt(3e18)
	transferNative(t, owner, bankERC20Address, 50_000, ccipFeeBudget, sourceBackend)
	transferNative(t, owner, relay.From, 21_000, amount, sourceBackend)
	transferNative(t, ccipContracts.Dest.User, relay.From, 21_000, amount, destBackend)

	wrapped, err := ccipContracts.Source.Router.GetWrappedNative(nil)
	require.NoError(t, err)

	linkUSD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err2 := w.Write([]byte(`{"UsdPerLink": "8000000000000000000"}`))
		require.NoError(t, err2)
	}))
	defer linkUSD.Close()
	ethUSD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err2 := w.Write([]byte(`{"UsdPerETH": "2000000000000000000000"}`))
		require.NoError(t, err2)
	}))
	defer ethUSD.Close()
	wrappedDestTokenUSD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err2 := w.Write([]byte(`{"UsdPerWrappedDestToken": "500000000000000000"}`))
		require.NoError(t, err2)
	}))
	defer wrappedDestTokenUSD.Close()
	bankERC20USD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err2 := w.Write([]byte(`{"UsdPerBankERC20": "5000000000000000000"}`))
		require.NoError(t, err2)
	}))
	defer bankERC20USD.Close()
	tokenPricesUSDPipeline := fmt.Sprintf(`
	link [type=http method=GET url="%s"];
	link_parse [type=jsonparse path="UsdPerLink"];
	link->link_parse;
	eth [type=http method=GET url="%s"];
	eth_parse [type=jsonparse path="UsdPerETH"];
	eth->eth_parse;
	wrapDest [type=http method=GET url="%s"];
	wrapDest_parse [type=jsonparse path="UsdPerWrappedDestToken"];
	wrapDest->wrapDest_parse;
	bankERC20 [type=http method=GET url="%s"];
	bankERC20_parse [type=jsonparse path="UsdPerBankERC20"];
	bankERC20->bankERC20_parse
	merge [type=merge left="{}" right="{\\\"%s\\\":$(link_parse), \\\"%s\\\":$(eth_parse), \\\"%s\\\":$(wrapDest_parse), \\\"%s\\\":$(bankERC20_parse)}"];`,
		linkUSD.URL, ethUSD.URL, wrappedDestTokenUSD.URL, bankERC20USD.URL, ccipContracts.Dest.LinkToken.Address(), wrapped, destToken.Address(), bankERC20Address)
	ccipContracts.SetUpNodesAndJobs(t, tokenPricesUSDPipeline, int64(legacygasstation.GetFreePort(t)))
	dummyDestForwarderRouter := common.HexToAddress("0x2")
	statusUpdateServer := legacygasstation.NewUnstartedStatusUpdateServer(t)
	go statusUpdateServer.Start()
	defer statusUpdateServer.Stop()
	createLegacyGasStationSidecarJob(t, app, dummyDestForwarderRouter, ccipContracts.Dest.OffRamp.Address(), destCCIPChainSelector, ccipContracts.Dest.ChainID, statusUpdateServer)

	calldata, calldataHash, err := metatx.GenerateMetaTransferCalldata(receiver.From, amount, destCCIPChainSelector)
	require.NoError(t, err)

	deadline := big.NewInt(int64(sourceBackend.Blockchain().CurrentHeader().Time + uint64(time.Hour)))
	signature, domainSeparatorHash, typeHash, forwarderNonce, err := metatx.SignMetaTransfer(*forwarder, senderKey.ToEcdsaPrivKey(), sender.From, bankERC20Address, calldataHash, deadline, metatx.BankERC20TokenName, metatx.BankERC20TokenVersion)
	require.NoError(t, err)

	forwardRequest := forwarder_wrapper.IForwarderForwardRequest{
		From:           sender.From,
		Target:         bankERC20Address,
		Nonce:          forwarderNonce,
		Data:           calldata,
		ValidUntilTime: deadline,
	}

	orm := legacygasstation.NewORM(db, app.GetLogger(), app.GetConfig().Database())

	// send meta transaction to forwarder
	_, err = forwarder.Execute(relay, forwardRequest, domainSeparatorHash, typeHash, []byte{}, signature)
	require.NoError(t, err)
	ccipContracts.Source.Chain.Commit()

	// get CCIP message ID by querying forwarder event
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   (*big.Int)(ccipContracts.Source.Chain.Blockchain().CurrentBlock().Number),
		Addresses: []common.Address{
			forwarderAddress,
		},
		Topics: [][]common.Hash{
			{
				forwarder_wrapper.ForwarderForwardSucceeded{}.Topic(),
			},
		},
	}
	logs, err := ccipContracts.Source.Chain.FilterLogs(testutils.Context(t), query)
	require.NoError(t, err)
	require.True(t, len(logs) == 1)
	log, err := forwarder.ParseForwardSucceeded(logs[0])
	require.NoError(t, err)
	ccipMessageID := common.Hash(log.ReturnValue)

	// verify that token was transferred to receiver on destination chain
	gomega.NewWithT(t).Eventually(func() bool {
		ccipContracts.Source.Chain.Commit()
		destBackend.Commit()
		receiverBal, err2 := destToken.BalanceOf(nil, receiver.From)
		require.NoError(t, err2)
		return receiverBal.Cmp(amount) == 0
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, app.KeyStore.Eth(), 0)
	txStore := cltest.NewTestTxStore(t, app.GetSqlxDB(), app.Config.Database())
	ethTx := cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)

	tx := types.LegacyGaslessTx{
		ID:                 "ID",
		Forwarder:          forwarderAddress,
		From:               sender.From,
		Target:             bankERC20Address,
		Receiver:           receiver.From,
		Nonce:              utils.NewBig(big.NewInt(1)),
		Amount:             utils.NewBig(amount),
		SourceChainID:      testhelpers.SourceChainID,
		DestinationChainID: testhelpers.DestChainID,
		ValidUntilTime:     utils.NewBig(big.NewInt(2)),
		Status:             types.Submitted,
		TokenName:          metatx.BankERC20TokenName,
		TokenVersion:       metatx.BankERC20TokenVersion,
		Signature:          []byte("signature"),
		EthTxID:            ethTx.GetID(),
	}
	// create a legacy gasless tx entry in DB so that sidecar can pick up pending requests
	err = orm.InsertLegacyGaslessTx(tx)
	require.NoError(t, err)
	tx.Status = types.SourceFinalized
	tx.CCIPMessageID = &ccipMessageID
	err = orm.UpdateLegacyGaslessTx(tx)
	require.NoError(t, err)

	// verify that transaction eventually becomes finalized
	gomega.NewWithT(t).Eventually(func() bool {
		sourceBackend.Commit()
		destBackend.Commit()
		txs, err := orm.SelectByDestChainIDAndStatus(ccipContracts.Dest.ChainID, types.Finalized)
		require.NoError(t, err)
		for _, tx := range txs {
			return tx.Status == types.Finalized
		}
		return false
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
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
	server legacygasstation.TestStatusUpdateServer,
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

func setUpForwarder(t *testing.T, owner *bind.TransactOpts, chain *backends.SimulatedBackend) (common.Address, *forwarder_wrapper.Forwarder) {
	// deploys EIP 2771 forwarder contract that verifies signatures from meta transaction and forwards the call to recipient contract (i.e BankERC20 token)
	forwarderAddress, _, forwarder, err := forwarder_wrapper.DeployForwarder(owner, chain)
	require.NoError(t, err)
	chain.Commit()
	// registers EIP712-compliant domain separator for BankERC20 token
	_, err = forwarder.RegisterDomainSeparator(owner, metatx.BankERC20TokenName, metatx.BankERC20TokenVersion)
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

func transferNative(t *testing.T, sender *bind.TransactOpts, receiverAddress common.Address, gasLimit uint64, amount *big.Int, chain *backends.SimulatedBackend) {
	nonce, err := chain.NonceAt(testutils.Context(t), sender.From, nil)
	require.NoError(t, err)
	tx := geth_types.NewTransaction(
		nonce, receiverAddress,
		amount,
		gasLimit,
		assets.GWei(1).ToInt(),
		nil)
	signedTx, err := sender.Signer(sender.From, tx)
	require.NoError(t, err)
	err = chain.SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	chain.Commit()

	receiverBalance, err := chain.BalanceAt(context.Background(), receiverAddress, nil)
	require.NoError(t, err)
	require.Equal(t, amount, receiverBalance)
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
		c.Feature.LegacyGasStation = ptr(true)
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
) (*forwarder_wrapper.Forwarder, *bank_erc20.BankERC20, common.Address, common.Address) {
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
	forwarder *forwarder_wrapper.Forwarder,
	bankERC20Address common.Address,
	senderKey ethkey.KeyV2,
	receiver common.Address,
	amount *big.Int,
	sourceChainSelector uint64,
	destChainSelector uint64,
) types.SendTransactionRequest {
	_, calldataHash, err := metatx.GenerateMetaTransferCalldata(receiver, amount, destChainSelector)
	require.NoError(t, err)
	deadline := big.NewInt(int64(backend.Blockchain().CurrentHeader().Time + uint64(time.Hour)))
	signature, _, _, forwarderNonce, err := metatx.SignMetaTransfer(*forwarder, senderKey.ToEcdsaPrivKey(), senderKey.Address, bankERC20Address, calldataHash, deadline, metatx.BankERC20TokenName, metatx.BankERC20TokenVersion)
	require.NoError(t, err)

	return types.SendTransactionRequest{
		From:               senderKey.Address,
		Target:             bankERC20Address,
		TargetName:         metatx.BankERC20TokenName,
		Version:            metatx.BankERC20TokenVersion,
		Nonce:              forwarderNonce,
		Receiver:           receiver,
		Amount:             amount,
		SourceChainID:      sourceChainSelector,
		DestinationChainID: destChainSelector,
		ValidUntilTime:     deadline,
		Signature:          signature,
	}
}
