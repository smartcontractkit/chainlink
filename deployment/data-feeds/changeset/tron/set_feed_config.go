package tron

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// SetFeedConfigChangeset is a changeset that sets a feed configuration on DataFeedsCache contract.
var SetFeedConfigChangeset = cldf.CreateChangeSet(setFeedConfigLogic, setFeedConfigPrecondition)

func setFeedConfigLogic(env cldf.Environment, c types.SetFeedDecimalTronConfig) (cldf.ChangesetOutput, error) {
	/*chain := env.BlockChains.TronChains()[c.ChainSelector]

	parsedABI, err := abi.JSON(strings.NewReader(cache.DataFeedsCacheABI))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to parse ABI: %w", err)
	}

	dataIDs, err := changeset.FeedIDsToBytes16(c.DataIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data ids: %s, %w", c.DataIDs, err)
	}

	calldata, err := parsedABI.Pack(
		"setDecimalFeedConfigs",
		dataIDs,
		c.Descriptions,
		c.WorkflowMetadata,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to pack calldata: %w", err)
	}

	triggerResp, err := chain.Client.TriggerSmartContractWithData(chain.Address, c.CacheAddress, fmt.Sprintf("%x", calldata), c.TriggerOptions.FeeLimit, c.TriggerOptions.TAmount)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	txInfo, err := chain.SendAndConfirm(context.Background(), triggerResp.Transaction, c.TriggerOptions.ConfirmRetryOptions)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", txInfo.ID, err)
	}*/

	return cldf.ChangesetOutput{}, nil
}

func setFeedConfigPrecondition(env cldf.Environment, c types.SetFeedDecimalTronConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
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
