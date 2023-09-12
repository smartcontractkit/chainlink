package eal_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/eal/generated/asset_catalog"
	forwarder_wrapper "github.com/smartcontractkit/capital-markets-projects/core/gethwrappers/eal/generated/forwarder"
	eallib "github.com/smartcontractkit/capital-markets-projects/lib/services/eal"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/eal/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/eal"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"
)

const (
	DomainSeparatorName    = "AssetCatalog"
	DomainSeparatorVersion = "1"
)

func TestIntegration_SendTransaction(t *testing.T) {
	// owner deploys forwarder and token contracts
	_, owner := generateKeyAndTransactor(t, testutils.SimulatedChainID)
	// relay is a CL-owned address that posts txs
	relayKey, relay := generateKeyAndTransactor(t, testutils.SimulatedChainID)
	// signer initiates EAL SendTransaction API call
	signerKey, _ := generateKeyAndTransactor(t, testutils.SimulatedChainID)

	genesisData := core.GenesisAlloc{
		owner.From: {Balance: assets.Ether(1000).ToInt()},
		relay.From: {Balance: assets.Ether(1000).ToInt()},
	}
	gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)

	forwarder, catalog, forwarderAddress, catalogAddress := setupAssetCatalogAndForwarderContracts(t, owner, backend)

	config, _ := setUpDB(t)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend, relayKey)
	require.NoError(t, app.Start(testutils.Context(t)))

	ccipChainSelector := uint64(30) //random value
	assetName := "ABC"
	assetPrice := big.NewInt(100)
	createEALJob(t, app, forwarderAddress, testutils.SimulatedChainID.Uint64(), ccipChainSelector, relayKey)

	req := generateRequest(t, backend, forwarder, catalogAddress, signerKey, assetName, assetPrice, ccipChainSelector)
	requestID := sendTransaction(t, req, app.Server.URL)
	require.NotEmpty(t, requestID)
	// assertions
	price, err := catalog.GetAssetPrice(nil, assetName)
	require.NoError(t, err)
	require.Equal(t, assetPrice, price)
}

func createEALJob(
	t *testing.T,
	app chainlink.Application,
	forwarderAddress common.Address,
	evmChainID uint64,
	ccipChainSelector uint64,
	fromKey ethkey.KeyV2,
) job.Job {
	jid := uuid.New()
	jobName := fmt.Sprintf("eal-%s", jid.String())
	s := testspecs.GenerateEALSpec(testspecs.EALSpecParams{
		JobID:             jid.String(),
		Name:              jobName,
		ForwarderAddress:  forwarderAddress.Hex(),
		EVMChainID:        evmChainID,
		CCIPChainSelector: ccipChainSelector,
		FromAddresses:     []string{fromKey.Address.String()},
	}).Toml()
	jb, err := eal.ValidatedEALSpec(s)
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
	forwarderAddress, _, forwarder, err := forwarder_wrapper.DeployForwarder(owner, chain)
	require.NoError(t, err)
	chain.Commit()
	// registers EIP712-compliant domain separator for BankERC20 token
	_, err = forwarder.RegisterDomainSeparator(owner, DomainSeparatorName, DomainSeparatorVersion)
	require.NoError(t, err)
	chain.Commit()

	return forwarderAddress, forwarder
}

func setUpAssetCatalog(t *testing.T, owner *bind.TransactOpts, chain *backends.SimulatedBackend) (common.Address, *asset_catalog.AssetCatalog) {
	tokenAddress, _, token, err := asset_catalog.DeployAssetCatalog(owner, chain)
	require.NoError(t, err)
	chain.Commit()
	return tokenAddress, token
}

func setUpDB(t *testing.T) (cfg chainlink.GeneralConfig, db *sqlx.DB) {
	cfg, db = heavyweight.FullTestDBV2(t, "eal_integration_test", func(c *chainlink.Config, s *chainlink.Secrets) {
		require.Zero(t, testutils.SimulatedChainID.Cmp(c.EVM[0].ChainID.ToInt()))
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

	resp, err := http.Post(fmt.Sprintf("%s/%s", url, "eal/v1/send_transaction"),
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

func setupAssetCatalogAndForwarderContracts(
	t *testing.T,
	deployer *bind.TransactOpts,
	backend *backends.SimulatedBackend,
) (*forwarder_wrapper.Forwarder, *asset_catalog.AssetCatalog, common.Address, common.Address) {
	forwarderAddress, forwarder := setUpForwarder(t, deployer, backend)
	assetCatalogAddress, assetCatalog := setUpAssetCatalog(t, deployer, backend)
	return forwarder, assetCatalog, forwarderAddress, assetCatalogAddress
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
	catalogAddress common.Address,
	signerKey ethkey.KeyV2,
	assetName string,
	assetPrice *big.Int,
	chainSelector uint64,
) types.SendTransactionRequest {
	_, calldataHash, err := GenerateSetAssetPriceCalldata(assetName, assetPrice)
	require.NoError(t, err)
	deadline := big.NewInt(int64(backend.Blockchain().CurrentHeader().Time + uint64(time.Hour)))
	gasLimit := big.NewInt(500_000)
	signature, encodedData, _, _, forwarderNonce, err := eallib.ForwarderEIP712Signature(
		*forwarder,
		signerKey.ToEcdsaPrivKey(),
		signerKey.Address,
		catalogAddress,
		calldataHash,
		gasLimit,
		deadline,
		DomainSeparatorName,
		DomainSeparatorVersion,
	)
	require.NoError(t, err)

	return types.SendTransactionRequest{
		From:                   signerKey.Address,
		Target:                 catalogAddress,
		DomainSeparatorName:    DomainSeparatorName,
		DomainSeparatorVersion: DomainSeparatorVersion,
		Nonce:                  forwarderNonce,
		Data:                   hex.EncodeToString(encodedData),
		Gas:                    gasLimit,
		ValidUntilTime:         deadline,
		ChainID:                chainSelector,
		Signature:              hex.EncodeToString(signature),
	}
}

func GenerateSetAssetPriceCalldata(assetName string, assetPrice *big.Int) ([]byte, [32]byte, error) {
	calldataAbi, err := abi.JSON(strings.NewReader(asset_catalog.AssetCatalogABI))
	if err != nil {
		return nil, [32]byte{}, err
	}

	calldata, err := calldataAbi.Pack("setAssetPrice", assetName, assetPrice)
	if err != nil {
		return nil, [32]byte{}, err
	}

	calldataHashRaw := crypto.Keccak256(calldata)

	var calldataHash [32]byte
	copy(calldataHash[:], calldataHashRaw[:])

	return calldata, calldataHash, nil
}

func ptr[T any](t T) *T { return &t }
