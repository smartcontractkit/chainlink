package v2

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/versioning"
)

const (
	// ContractWorkflowSourceName is the name used for logging and identification.
	ContractWorkflowSourceName = "ContractWorkflowSource"
)

// ContractWorkflowSource implements WorkflowMetadataSource by reading from the on-chain
// workflow registry contract.
type ContractWorkflowSource struct {
	lggr                    logger.Logger
	workflowRegistryAddress string
	chainSelector           string
	contractReaderFn        versioning.ContractReaderFactory
	contractReader          commontypes.ContractReader
	mu                      sync.RWMutex
}

// NewContractWorkflowSource creates a new contract-based workflow source.
// chainSelector is the chain selector where the workflow registry contract is deployed.
func NewContractWorkflowSource(
	lggr logger.Logger,
	contractReaderFn versioning.ContractReaderFactory,
	workflowRegistryAddress string,
	chainSelector string,
) *ContractWorkflowSource {
	return &ContractWorkflowSource{
		lggr:                    logger.Named(lggr, ContractWorkflowSourceName),
		contractReaderFn:        contractReaderFn,
		workflowRegistryAddress: workflowRegistryAddress,
		chainSelector:           chainSelector,
	}
}

// ListWorkflowMetadata fetches workflow metadata from the on-chain contract.
// It lazily initializes the contract reader on first call.
func (c *ContractWorkflowSource) ListWorkflowMetadata(ctx context.Context, don capabilities.DON) ([]WorkflowMetadataView, *commontypes.Head, error) {
	c.tryInitialize(ctx)

	c.mu.RLock()
	reader := c.contractReader
	c.mu.RUnlock()

	if reader == nil {
		return nil, nil, errors.New("contract reader not initialized")
	}

	contractBinding := commontypes.BoundContract{
		Address: c.workflowRegistryAddress,
		Name:    WorkflowRegistryContractName,
	}

	readIdentifier := contractBinding.ReadIdentifier(GetWorkflowsByDONMethodName)
	var headAtLastRead *commontypes.Head
	var allWorkflows []WorkflowMetadataView

	for _, family := range don.Families {
		params := GetWorkflowListByDONParams{
			DonFamily: family,
			Start:     big.NewInt(0),
			Limit:     big.NewInt(MaxResultsPerQuery),
		}

		for {
			var err error
			var workflows struct {
				List []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView
			}

			headAtLastRead, err = reader.GetLatestValueWithHeadData(ctx, readIdentifier, primitives.Finalized, params, &workflows)
			if err != nil {
				return []WorkflowMetadataView{}, &commontypes.Head{Height: "0"}, fmt.Errorf("failed to get latest value with head data: %w", err)
			}

			for _, wfMeta := range workflows.List {
				// Skip workflows with incomplete/invalid metadata - this can indicate stale metadata
				// from deleted workflows in the contract (known contract bug where deleted workflows
				// aren't fully removed from contract state)
				if !isValidWorkflowMetadata(wfMeta) {
					c.lggr.Warnw("Workflow has incomplete metadata from contract, skipping",
						"source", ContractWorkflowSourceName,
						"workflowID", hex.EncodeToString(wfMeta.WorkflowId[:]),
						"workflowName", wfMeta.WorkflowName,
						"owner", hex.EncodeToString(wfMeta.Owner.Bytes()),
						"binaryURL", wfMeta.BinaryUrl,
						"configURL", wfMeta.ConfigUrl,
						"status", wfMeta.Status)
					continue
				}

				// TODO: https://smartcontract-it.atlassian.net/browse/CAPPL-1021 load balance across workflow nodes in DON Family
				allWorkflows = append(allWorkflows, WorkflowMetadataView{
					WorkflowID:   wfMeta.WorkflowId,
					Owner:        wfMeta.Owner.Bytes(),
					CreatedAt:    wfMeta.CreatedAt,
					Status:       ContractStatusToInternal(wfMeta.Status),
					WorkflowName: wfMeta.WorkflowName,
					BinaryURL:    wfMeta.BinaryUrl,
					ConfigURL:    wfMeta.ConfigUrl,
					Tag:          wfMeta.Tag,
					Attributes:   wfMeta.Attributes,
					DonFamily:    wfMeta.DonFamily,
					Source:       c.SourceIdentifier(),
				})
			}

			// if less workflows than limit, then we have reached the end of the list
			if int64(len(workflows.List)) < MaxResultsPerQuery {
				break
			}

			// otherwise, increment the start parameter and continue to fetch more workflows
			params.Start.Add(params.Start, big.NewInt(int64(len(workflows.List))))
		}
	}

	c.lggr.Debugw("Loaded workflows from contract",
		"address", c.workflowRegistryAddress,
		"count", len(allWorkflows),
		"donFamilies", don.Families)

	if headAtLastRead == nil {
		return allWorkflows, &commontypes.Head{Height: "0"}, nil
	}

	return allWorkflows, headAtLastRead, nil
}

func (c *ContractWorkflowSource) Name() string {
	return ContractWorkflowSourceName
}

// SourceIdentifier returns the source identifier used in WorkflowMetadataView.Source.
// Format: contract:{chain_selector}:{contract_address}
func (c *ContractWorkflowSource) SourceIdentifier() string {
	return fmt.Sprintf("contract:%s:%s", c.chainSelector, c.workflowRegistryAddress)
}

// Ready returns nil if the contract reader is initialized.
func (c *ContractWorkflowSource) Ready() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.contractReader == nil {
		return errors.New("contract reader not initialized")
	}
	return nil
}

// tryInitialize attempts to initialize the contract reader. Returns true if ready.
func (c *ContractWorkflowSource) tryInitialize(ctx context.Context) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.contractReader != nil {
		return true
	}

	reader, err := c.newWorkflowRegistryContractReader(ctx)
	if err != nil {
		c.lggr.Debugw("Contract reader not yet available", "error", err)
		return false
	}

	c.contractReader = reader
	c.lggr.Debugw("Contract reader initialized successfully")
	return true
}

// newWorkflowRegistryContractReader creates a new contract reader configured for the workflow registry.
func (c *ContractWorkflowSource) newWorkflowRegistryContractReader(ctx context.Context) (commontypes.ContractReader, error) {
	contractReaderCfg := config.ChainReaderConfig{
		Contracts: map[string]config.ChainContractReader{
			WorkflowRegistryContractName: {
				ContractABI: workflow_registry_wrapper_v2.WorkflowRegistryABI,
				Configs: map[string]*config.ChainReaderDefinition{
					GetWorkflowsByDONMethodName: {
						ChainSpecificName: GetWorkflowsByDONMethodName,
						ReadType:          config.Method,
					},
				},
			},
		},
	}

	marshalledCfg, err := json.Marshal(contractReaderCfg)
	if err != nil {
		return nil, err
	}

	reader, err := c.contractReaderFn(ctx, marshalledCfg)
	if err != nil {
		return nil, err
	}

	bc := commontypes.BoundContract{
		Name:    WorkflowRegistryContractName,
		Address: c.workflowRegistryAddress,
	}

	// bind contract to contract reader
	if err := reader.Bind(ctx, []commontypes.BoundContract{bc}); err != nil {
		return nil, err
	}

	if err := reader.Start(ctx); err != nil {
		return nil, err
	}

	return reader, nil
}

// isValidWorkflowMetadata checks if workflowID and workflowOwner are valid
// in the metadata pulled from the contract. In the case of contract deletion bugs
// (where deleted workflows retain stale metadata with zero addresses), this func
// filters out noisy deploys/workflow events.
func isValidWorkflowMetadata(wfMeta workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView) bool {
	return !isEmptyWorkflowID(wfMeta.WorkflowId) && !isZeroOwner(wfMeta.Owner.Bytes())
}
