package metatx_test

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/bank_erc20"
	forwarder_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/metatx"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
)

func TestMetaERC20SameChain(t *testing.T) {
	chainID := uint64(1337)

	// deploys and owns contract
	_, contractOwner := generateKeyAndTransactor(t, chainID)
	// holder1Key sends tokens to holder2
	holder1Key, holder1 := generateKeyAndTransactor(t, chainID)
	// holder2Key receives tokens
	_, holder2 := generateKeyAndTransactor(t, chainID)
	// relayKey is the relayer that submits signed meta-transaction to the forwarder contract on-chain
	_, relay := generateKeyAndTransactor(t, chainID)
	// ccipFeeProvider can withdraw native tokens from BankERC20 contract
	_, ccipFeeProvider := generateKeyAndTransactor(t, chainID)

	chain := backends.NewSimulatedBackend(core.GenesisAlloc{
		contractOwner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(1e18)),
		},
	}, ethconfig.Defaults.Miner.GasCeil)

	// deploys forwarder that verifies meta transaction signature and forwards requests to token
	forwarderAddress, forwarder := setUpForwarder(t, contractOwner, chain)

	totalTokens := big.NewInt(1e9)
	// CCIP router address is not needed because the test is for same-chain transfer
	tokenAddress, token := setUpBankERC20(t, contractOwner, chain, forwarderAddress, common.HexToAddress("0x1"), ccipFeeProvider.From, totalTokens, chainID)
	// amount to transfer
	amount := assets.Ether(1).ToInt()

	// fund BankERC20 contract with native ETH. Test setup may use a low Eth price, important to send enough.
	ccipFeeBudget := assets.Ether(3).ToInt()
	transferNative(t, contractOwner, tokenAddress, 50_000, ccipFeeBudget, chain)

	sourceTokenEthBal, err := chain.BalanceAt(testutils.Context(t), tokenAddress, nil)
	require.NoError(t, err)
	require.Equal(t, ccipFeeBudget, sourceTokenEthBal)

	t.Run("single same-chain meta transfer", func(t *testing.T) {
		// transfer BankERC20 from contract owner to holder1
		transferToken(t, token, contractOwner, holder1, amount, chain)

		deadline := big.NewInt(int64(chain.Blockchain().CurrentHeader().Time + uint64(time.Hour)))

		calldata, calldataHash, err := metatx.GenerateMetaTransferCalldata(holder2.From, amount, chainID)
		require.NoError(t, err)

		signature, domainSeparatorHash, typeHash, forwarderNonce, err := metatx.SignMetaTransfer(*forwarder,
			holder1Key.ToEcdsaPrivKey(),
			holder1.From,
			tokenAddress,
			calldataHash,
			deadline,
			metatx.BankERC20TokenName,
			metatx.BankERC20TokenVersion,
		)
		require.NoError(t, err)

		forwardRequest := forwarder_wrapper.IForwarderForwardRequest{
			From:           holder1.From,
			Target:         tokenAddress,
			Nonce:          forwarderNonce,
			Data:           calldata,
			ValidUntilTime: deadline,
		}

		transferNative(t, contractOwner, relay.From, 21_000, ccipFeeBudget, chain)

		holder1BalanceBefore, err := token.BalanceOf(nil, holder1.From)
		require.NoError(t, err)

		// send meta transaction to forwarder
		_, err = forwarder.Execute(relay, forwardRequest, domainSeparatorHash, typeHash, nil, signature)
		require.NoError(t, err)
		chain.Commit()

		holder2Balance, err := token.BalanceOf(nil, holder2.From)
		require.NoError(t, err)
		require.Equal(t, holder2Balance, amount)

		holder1Balance, err := token.BalanceOf(nil, holder1.From)
		require.NoError(t, err)
		require.Equal(t, holder1Balance, holder1BalanceBefore.Sub(holder1BalanceBefore, amount))

		totalSupplyAfter, err := token.TotalSupply(nil)
		require.NoError(t, err)
		require.Equal(t, totalSupplyAfter, big.NewInt(0).Mul(totalTokens, big.NewInt(1e18)))
	})
}

func TestMetaERC20CrossChain(t *testing.T) {
	ccipContracts := integrationtesthelpers.SetupCCIPIntegrationTH(t, testhelpers.SourceChainID, testhelpers.SourceChainSelector, testhelpers.DestChainID, testhelpers.DestChainSelector)

	// holder1Key sends tokens to holder2
	holder1Key, holder1 := generateKeyAndTransactor(t, ccipContracts.Source.Chain.Blockchain().Config().ChainID.Uint64())
	// holder2Key receives tokens
	_, holder2 := generateKeyAndTransactor(t, ccipContracts.Dest.Chain.Blockchain().Config().ChainID.Uint64())
	// relayKey is the relayer that submits signed meta-transaction to the forwarder contract on-chain
	_, relay := generateKeyAndTransactor(t, ccipContracts.Source.Chain.Blockchain().Config().ChainID.Uint64())
	// ccipFeeProvider can withdraw native tokens from BankERC20 contract
	_, ccipFeeProvider := generateKeyAndTransactor(t, ccipContracts.Source.Chain.Blockchain().Config().ChainID.Uint64())

	forwarderAddress, forwarder := setUpForwarder(t, ccipContracts.Source.User, ccipContracts.Source.Chain)

	totalTokens := big.NewInt(1e9)
	sourceTokenAddress, sourceToken := setUpBankERC20(t, ccipContracts.Source.User, ccipContracts.Source.Chain, forwarderAddress, ccipContracts.Source.Router.Address(), ccipFeeProvider.From, totalTokens, testhelpers.SourceChainID)

	sourcePoolAddress, destToken, err := ccipContracts.SetupLockAndMintTokenPool(sourceTokenAddress, "WrappedBankToken", "WBANK")
	require.NoError(t, err)

	amount := assets.Ether(1).ToInt()
	ccipFeeBudget := assets.Ether(3).ToInt()
	transferNative(t, ccipContracts.Source.User, sourceTokenAddress, 50_000, ccipFeeBudget, ccipContracts.Source.Chain)

	linkUSD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"UsdPerLink": "8000000000000000000"}`))
		require.NoError(t, err)
	}))
	ethUSD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"UsdPerETH": "2000000000000000000000"}`))
		require.NoError(t, err)
	}))
	wrappedDestTokenUSD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"UsdPerWrappedDestToken": "500000000000000000"}`))
		require.NoError(t, err)
	}))
	bankERC20USD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"UsdPerBankERC20": "5000000000000000000"}`))
		require.NoError(t, err)
	}))
	wrapped, err := ccipContracts.Source.Router.GetWrappedNative(nil)
	require.NoError(t, err)
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
		linkUSD.URL, ethUSD.URL, wrappedDestTokenUSD.URL, bankERC20USD.URL, ccipContracts.Dest.LinkToken.Address(), wrapped, destToken.Address(), sourceTokenAddress)
	defer linkUSD.Close()
	defer ethUSD.Close()
	defer wrappedDestTokenUSD.Close()
	defer bankERC20USD.Close()

	ccipContracts.SetUpNodesAndJobs(t, tokenPricesUSDPipeline, 29599)

	geCurrentSeqNum := 1

	t.Run("single cross-chain meta transfer", func(t *testing.T) {
		// transfer BankERC20 from owner to holder1
		transferToken(t, sourceToken, ccipContracts.Source.User, holder1, amount, ccipContracts.Source.Chain)

		deadline := big.NewInt(int64(ccipContracts.Source.Chain.Blockchain().CurrentHeader().Time + uint64(time.Hour)))

		calldata, calldataHash, err := metatx.GenerateMetaTransferCalldata(holder2.From, amount, ccipContracts.Dest.ChainSelector)
		require.NoError(t, err)

		signature, domainSeparatorHash, typeHash, forwarderNonce, err := metatx.SignMetaTransfer(
			*forwarder,
			holder1Key.ToEcdsaPrivKey(),
			holder1.From,
			sourceTokenAddress,
			calldataHash,
			deadline,
			metatx.BankERC20TokenName,
			metatx.BankERC20TokenVersion,
		)
		require.NoError(t, err)

		forwardRequest := forwarder_wrapper.IForwarderForwardRequest{
			From:           holder1.From,
			Target:         sourceTokenAddress,
			Nonce:          forwarderNonce,
			Data:           calldata,
			ValidUntilTime: deadline,
		}

		transferNative(t, ccipContracts.Source.User, relay.From, 21_000, ccipFeeBudget, ccipContracts.Source.Chain)

		// send meta transaction to forwarder
		_, err = forwarder.Execute(relay, forwardRequest, domainSeparatorHash, typeHash, []byte{}, signature)
		require.NoError(t, err)
		ccipContracts.Source.Chain.Commit()

		gomega.NewWithT(t).Eventually(func() bool {
			ccipContracts.Source.Chain.Commit()
			ccipContracts.Dest.Chain.Commit()
			holder2Balance, err2 := destToken.BalanceOf(nil, holder2.From)
			require.NoError(t, err2)
			return holder2Balance.Cmp(amount) == 0
		}, testutils.WaitTimeout(t), 5*time.Second).Should(gomega.BeTrue())

		ccipContracts.AllNodesHaveReqSeqNum(t, geCurrentSeqNum)
		ccipContracts.EventuallyReportCommitted(t, geCurrentSeqNum)

		executionLogs := ccipContracts.AllNodesHaveExecutedSeqNums(t, geCurrentSeqNum, geCurrentSeqNum)
		assert.Len(t, executionLogs, 1)
		ccipContracts.AssertExecState(t, executionLogs[0], testhelpers.ExecutionStateSuccess)

		// source token is locked in the token pool
		lockedTokenBal, err := sourceToken.BalanceOf(nil, sourcePoolAddress)
		require.NoError(t, err)
		require.Equal(t, lockedTokenBal, amount)

		// source total supply should stay the same
		sourceTotalSupply, err := sourceToken.TotalSupply(nil)
		require.NoError(t, err)
		require.Equal(t, sourceTotalSupply, big.NewInt(0).Mul(totalTokens, big.NewInt(1e18)))

		// new wrapped tokens minted on dest token
		destTotalSupply, err := destToken.TotalSupply(nil)
		require.NoError(t, err)
		require.Equal(t, destTotalSupply, amount)
	})
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

func generateKeyAndTransactor(t *testing.T, chainID uint64) (key ethkey.KeyV2, transactor *bind.TransactOpts) {
	key = cltest.MustGenerateRandomKey(t)
	transactor, err := bind.NewKeyedTransactorWithChainID(key.ToEcdsaPrivKey(), big.NewInt(0).SetUint64(chainID))
	require.NoError(t, err)
	return
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
	tx := types.NewTransaction(
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
