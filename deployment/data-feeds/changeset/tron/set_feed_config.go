package tron

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	fullnode "github.com/fbsobreira/gotron-sdk/pkg/http/fullnode"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// SetFeedConfigChangeset is a changeset that sets a feed configuration on DataFeedsCache contract.
var SetFeedConfigChangeset = cldf.CreateChangeSet(setFeedConfigLogic, setFeedConfigPrecondition)

func setFeedConfigLogic(env cldf.Environment, c types.SetFeedDecimalTronConfig) (cldf.ChangesetOutput, error) {
	chain := env.BlockChains.TronChains()[c.ChainSelector]

	parsedABI, err := abi.JSON(strings.NewReader(cache.DataFeedsCacheABI))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to parse ABI: %w", err)
	}

	dataIDs, err := changeset.FeedIDsToBytes16(c.DataIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data ids: %s, %w", c.DataIDs, err)
	}

	workflowMetadata := parseWorkflowMetadata(c.WorkflowMetadata)

	calldata, err := parsedABI.Pack(
		"setDecimalFeedConfigs",
		dataIDs,
		c.Descriptions,
		workflowMetadata,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to pack calldata: %w", err)
	}

	tcRequest := fullnode.TriggerSmartContractRequest{
		OwnerAddress:    chain.Address.String(),
		ContractAddress: c.CacheAddress.String(),
		Data:            hex.EncodeToString(calldata),
		FeeLimit:        c.TriggerOptions.FeeLimit,
		CallValue:       c.TriggerOptions.TAmount,
		Visible:         true,
	}
	contractResponse := fullnode.TriggerSmartContractResponse{}
	err = chain.Client.FullNodeClient().Post("/triggersmartcontract", tcRequest, &contractResponse)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	txInfo, err := chain.SendAndConfirm(context.Background(), contractResponse.Transaction, c.TriggerOptions.ConfirmRetryOptions)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %+v, %w", txInfo, err)
	}

	return cldf.ChangesetOutput{}, nil
}

func setFeedConfigPrecondition(env cldf.Environment, c types.SetFeedDecimalTronConfig) error {
	_, ok := env.BlockChains.TronChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if (len(c.DataIDs) == 0) || (len(c.Descriptions) == 0) || (len(c.WorkflowMetadata) == 0) {
		return errors.New("dataIDs, descriptions and workflowMetadata must not be empty")
	}
	if len(c.DataIDs) != len(c.Descriptions) {
		return errors.New("dataIDs and descriptions must have the same length")
	}
	_, err := changeset.FeedIDsToBytes16(c.DataIDs)
	if err != nil {
		return fmt.Errorf("failed to convert feed ids to bytes16: %w", err)
	}

	return changeset.ValidateCacheForTronChain(env, c.ChainSelector, c.CacheAddress)
}

func parseWorkflowMetadata(tronWorkflowMetadata []types.DataFeedsCacheTronWorkflowMetadata) []cache.DataFeedsCacheWorkflowMetadata {
	workflowMetadata := make([]cache.DataFeedsCacheWorkflowMetadata, len(tronWorkflowMetadata))

	for i, tronMeta := range tronWorkflowMetadata {
		workflowMetadata[i] = cache.DataFeedsCacheWorkflowMetadata{
			AllowedSender:        tronMeta.AllowedSender.EthAddress(),
			AllowedWorkflowOwner: tronMeta.AllowedWorkflowOwner.EthAddress(),
			AllowedWorkflowName:  tronMeta.AllowedWorkflowName,
		}
	}

	return workflowMetadata
}
