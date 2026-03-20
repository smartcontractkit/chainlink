package vrfv2plus

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/vrfv2plus"
)

func TestVRFv2PlusMultipleSendingKeys(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrf2plus-out.toml"

	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load env-out.toml")

	cfg, err := products.LoadOutput[vrfv2plus.Configurator](outputFile)
	require.NoError(t, err, "failed to load VRFv2Plus config from env-out.toml")
	require.NotEmpty(t, cfg.Config, "vrfv2_plus config must not be empty")

	c := cfg.Config[0]
	require.NotEmpty(t, c.DeployedContracts.Coordinator, "coordinator address must not be empty")
	require.NotEmpty(t, c.VRFKeyData.KeyHash, "key hash must not be empty")

	// Verify 3 TX key addresses: nodeEVMKey + 2 extras
	require.Len(t, c.VRFKeyData.TxKeyAddresses, 3,
		"expected 3 TX key addresses (nodeEVMKey + 2 extra), got %d", len(c.VRFKeyData.TxKeyAddresses))

	keyHashBytes := common.HexToHash(c.VRFKeyData.KeyHash)
	var keyHash [32]byte
	copy(keyHash[:], keyHashBytes[:])

	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err)
	chainIDBig := new(big.Int).SetUint64(chainID)

	bcNode := in.Blockchains[0].Out.Nodes[0]
	ctx := t.Context()
	chainClient, err := products.InitSeth(bcNode.ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err, "failed to init Seth client")

	coord, err := contracts.LoadVRFCoordinatorV2_5(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err, "failed to load coordinator")

	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err, "failed to load LINK token")

	// Make one request per TX key address (deploy fresh consumer + sub for each)
	numKeys := len(c.VRFKeyData.TxKeyAddresses)
	requestIDs := make([]*big.Int, numKeys)
	for i := range numKeys {
		consumer, subID := newConsumerAndSub(ctx, t, chainClient, coord, linkToken, c)
		reqID, rErr := consumer.RequestRandomness(
			keyHash, subID, c.MinimumConfirmations,
			defaultCallbackGasLimit, false, defaultNumWords, defaultRequestCount,
		)
		require.NoError(t, rErr, "RequestRandomness failed for key index %d", i)
		require.NotNil(t, reqID)
		requestIDs[i] = reqID
	}

	// Wait for all fulfillments and collect their TX hashes
	fulfillTxHashes := make([]common.Hash, numKeys)
	for i, reqID := range requestIDs {
		idx := i
		id := reqID
		var txHash common.Hash
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			event, fErr := coord.FilterRandomWordsFulfilled(
				&bind.FilterOpts{Context: ctx},
				id,
			)
			if fErr != nil {
				return false
			}
			require.True(t, event.Success, "fulfillment %d should succeed", idx)
			txHash = event.Raw.TxHash
			return true
		}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
			"timed out waiting for RandomWordsFulfilled for requestID %s", id)
		fulfillTxHashes[idx] = txHash
	}

	// Recover sender address for each fulfillment TX
	signer := types.LatestSignerForChainID(chainIDBig)
	actualSenders := make(map[string]struct{}, numKeys)
	for _, txHash := range fulfillTxHashes {
		tx, _, tErr := chainClient.Client.TransactionByHash(ctx, txHash)
		require.NoError(t, tErr, "failed to get tx %s", txHash.Hex())
		sender, sErr := types.Sender(signer, tx)
		require.NoError(t, sErr, "failed to recover sender from tx %s", txHash.Hex())
		actualSenders[sender.Hex()] = struct{}{}
	}

	// Assert sender set equals configured TX key addresses (order-independent)
	expectedSet := make(map[string]struct{}, numKeys)
	for _, addr := range c.VRFKeyData.TxKeyAddresses {
		expectedSet[common.HexToAddress(addr).Hex()] = struct{}{}
	}
	require.Equal(t, expectedSet, actualSenders,
		"fulfillment senders should match configured TX key addresses")
}
