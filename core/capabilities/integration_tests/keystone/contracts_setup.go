package keystone

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/require"

	data_feeds_cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
)

func SetupForwarderContract(t *testing.T, reportCreator *framework.DON,
	backend *framework.EthBlockchain) (common.Address, *forwarder.KeystoneForwarder) {
	addr, _, fwd, err := forwarder.DeployKeystoneForwarder(backend.TransactionOpts(), backend.Client())
	require.NoError(t, err)
	backend.Commit()

	signers := make([]common.Address, 0, len(reportCreator.GetPeerIDsAndOCRSigners()))
	for _, p := range reportCreator.GetPeerIDsAndOCRSigners() {
		signers = append(signers, common.HexToAddress(p.Signer))
	}

	_, err = fwd.SetConfig(backend.TransactionOpts(), reportCreator.GetID(), reportCreator.GetConfigVersion(), reportCreator.GetF(), signers)
	require.NoError(t, err)
	backend.Commit()

	return addr, fwd
}

func SetupConsumerContract(t *testing.T, backend *framework.EthBlockchain,
	forwarderAddress common.Address, workflowOwner string, workflowName string) (common.Address, *feeds_consumer.KeystoneFeedsConsumer) {
	addr, _, consumer, err := feeds_consumer.DeployKeystoneFeedsConsumer(backend.TransactionOpts(), backend.Client())
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := common.HexToAddress(workflowOwner)

	_, err = consumer.SetConfig(backend.TransactionOpts(), []common.Address{forwarderAddress}, []common.Address{ownerAddr}, [][10]byte{nameBytes})
	require.NoError(t, err)

	backend.Commit()

	return addr, consumer
}

func SetupDataFeedsCacheContract(t *testing.T, backend *framework.EthBlockchain,
	forwarderAddress common.Address, workflowOwner string, workflowName string) (common.Address, *data_feeds_cache.DataFeedsCache) {
	addr, _, dataFeedsCache, err := data_feeds_cache.DeployDataFeedsCache(backend.TransactionOpts(), backend.Client())
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := common.HexToAddress(workflowOwner)

	_, err = dataFeedsCache.SetFeedAdmin(backend.TransactionOpts(), backend.TransactionOpts().From, true)
	require.NoError(t, err)
	backend.Commit()

	feedIDBytes := [16]byte{}
	copy(feedIDBytes[:], common.FromHex("0x04de41ba4fc9d91ad900000000000000")) // Data ID for secure mint report for chain selector 16015286601757825753 (ethereum-testnet-sepolia)

	tx, err := dataFeedsCache.SetDecimalFeedConfigs(backend.TransactionOpts(), [][16]byte{feedIDBytes}, []string{"securemint"},
		[]data_feeds_cache.DataFeedsCacheWorkflowMetadata{
			{
				AllowedSender:        forwarderAddress,
				AllowedWorkflowOwner: ownerAddr,
				AllowedWorkflowName:  nameBytes,
			},
		})
	if err != nil {
		t.Logf("Failed to set decimal feed configs: %v", err)
		rpcError := decodeDFCacheError(t, backend, tx)
		t.Fatalf("failed to set decimal feed configs: %v, rpc error: %s", err, rpcError)
	}
	backend.Commit()

	return addr, dataFeedsCache
}

func decodeDFCacheError(t *testing.T, backend *framework.EthBlockchain, tx *types.Transaction) string {
	t.Helper()

	contractAbi, txError := data_feeds_cache.DataFeedsCacheMetaData.GetAbi()
	require.NoError(t, txError)

	errString := "0xd86ad9cf0000000000000000000000000acba7ba442154d7b1d11ffa5ecb40e9457fca05"

	errString = strings.TrimPrefix(errString, "0x")
	errBytes, err := hex.DecodeString(errString)
	require.NoError(t, err)

	selector := errBytes[:4]
	for name, errType := range contractAbi.Errors {
		if bytes.Equal(selector, errType.ID.Bytes()) {
			args, err := errType.Inputs.Unpack(errBytes[4:])
			require.NoError(t, err)
			return fmt.Sprintf("Matched error: %s, args: %v", name, args)
		}

	}

	return fmt.Sprintf("No custom DF Cache error found for %v.", txError)
}

func getTxError(t *testing.T, client simulated.Client, from common.Address, tx *types.Transaction) (string, error) {
	t.Helper()

	re, err := client.TransactionReceipt(t.Context(), tx.Hash())
	require.NoError(t, err, "error getting transaction receipt")
	if err != nil {
		return "", errors.Wrap(err, "error getting transaction receipt")
	}

	call := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Data:     tx.Data(),
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
	_, err = client.CallContract(context.Background(), call, re.BlockNumber)
	if err == nil {
		panic("no error calling contract")
	}

	return parseError(err)
}

func parseError(txError error) (string, error) {
	b, err := json.Marshal(txError)
	if err != nil {
		return "", err
	}
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	if json.Unmarshal(b, &callErr) != nil {
		return "", err
	}

	return callErr.Data, nil
}
