package forwarder

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	mock_forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/mock_forwarder"
)

var _ cldf.ChangeSetV2[DeployMockForwardersInput] = DeployMockForwarders{}

// DeployMockForwardersInput is the input for deploying MockKeystoneForwarder contracts.
type DeployMockForwardersInput struct {
	Targets   []uint64 `json:"targets" yaml:"targets"`
	Qualifier string   `json:"qualifier" yaml:"qualifier"`
}

// DeployMockForwarders is a ChangeSetV2 that deploys MockKeystoneForwarder contracts.
type DeployMockForwarders struct{}

func (d DeployMockForwarders) VerifyPreconditions(env cldf.Environment, input DeployMockForwardersInput) error {
	if input.Qualifier == "" {
		return errors.New("qualifier is required")
	}
	for _, sel := range input.Targets {
		if _, err := chain_selectors.GetChainIDFromSelector(sel); err != nil {
			return fmt.Errorf("could not resolve chain selector %d: %w", sel, err)
		}
		if _, ok := env.BlockChains.EVMChains()[sel]; !ok {
			return fmt.Errorf("chain selector %d not found in environment", sel)
		}
	}
	return nil
}

func (d DeployMockForwarders) Apply(env cldf.Environment, input DeployMockForwardersInput) (cldf.ChangesetOutput, error) {
	seqReport, err := operations.ExecuteSequence(
		env.OperationsBundle,
		DeployMockSequence,
		DeploySequenceDeps{Env: &env},
		DeploySequenceInput(input),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute deploy mock forwarders sequence: %w", err)
	}

	ds := datastore.NewMemoryDataStore()
	addrs, err := seqReport.Output.Addresses.Fetch()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch addresses from sequence output: %w", err)
	}
	for _, addr := range addrs {
		if err := ds.Addresses().Add(addr); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add address ref to mutable datastore: %w", err)
		}
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   seqReport.ExecutionReports,
	}, nil
}

type DeployMockForwarderSequenceOutput struct {
	Addresses datastore.AddressRefStore
	Datastore datastore.DataStore
}

// DeployMockSequence deploys MockKeystoneForwarder contracts to multiple chains concurrently.
var DeployMockSequence = operations.NewSequence(
	"deploy-mock-keystone-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Deploy Mock Keystone Forwarders",
	func(b operations.Bundle, deps DeploySequenceDeps, input DeploySequenceInput) (DeployMockForwarderSequenceOutput, error) {
		as := datastore.NewMemoryDataStore()
		contractErrGroup := &errgroup.Group{}
		for _, target := range input.Targets {
			contractErrGroup.Go(func() error {
				r, err := operations.ExecuteOperation(b, DeployMockOp, DeployOpDeps(deps), DeployOpInput{
					ChainSelector: target,
					Qualifier:     input.Qualifier,
				})
				if err != nil {
					return err
				}
				addrs, err := r.Output.Addresses.Fetch()
				if err != nil {
					return fmt.Errorf("failed to fetch MockKeystoneForwarder addresses for target %d: %w", target, err)
				}
				for _, addr := range addrs {
					if addrRefErr := as.AddressRefStore.Add(addr); addrRefErr != nil {
						return fmt.Errorf("failed to save MockKeystoneForwarder address on datastore for target %d: %w", target, addrRefErr)
					}
				}

				return nil
			})
		}
		if err := contractErrGroup.Wait(); err != nil {
			return DeployMockForwarderSequenceOutput{Addresses: as.Addresses()}, fmt.Errorf("failed to deploy MockKeystoneForwarder contracts: %w", err)
		}
		return DeployMockForwarderSequenceOutput{Addresses: as.Addresses(), Datastore: as.Seal()}, nil
	},
)

type DeployMockForwarderOpOutput struct {
	Addresses  datastore.AddressRefStore
	AddressRef datastore.AddressRef // The address ref of the deployed Keystone Forwarder
}

// DeployMockOp is an operation that deploys the MockKeystoneForwarder contract.
var DeployMockOp = operations.NewOperation(
	"deploy-mock-keystone-forwarder-op",
	semver.MustParse("1.0.0"),
	"Deploy MockKeystoneForwarder Contract",
	func(b operations.Bundle, deps DeployOpDeps, input DeployOpInput) (DeployMockForwarderOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployMockForwarderOpOutput{}, fmt.Errorf("deploy-mock-keystone-forwarder-op failed: chain selector %d not found in environment", input.ChainSelector)
		}
		addr, tv, err := deployMock(b.GetContext(), chain.DeployerKey, chain)
		if err != nil {
			return DeployMockForwarderOpOutput{}, fmt.Errorf("deploy-mock-keystone-forwarder-op failed: %w", err)
		}
		labels := tv.Labels.List()
		labels = append(labels, input.Labels...)
		r := datastore.AddressRef{
			ChainSelector: input.ChainSelector,
			Address:       addr.String(),
			Type:          datastore.ContractType(tv.Type),
			Version:       &tv.Version,
			Qualifier:     input.Qualifier,
			Labels:        datastore.NewLabelSet(labels...),
		}
		ds := datastore.NewMemoryDataStore()
		if err := ds.AddressRefStore.Add(r); err != nil {
			return DeployMockForwarderOpOutput{}, fmt.Errorf("deploy-mock-keystone-forwarder-op failed: failed to add address ref to datastore: %w", err)
		}

		return DeployMockForwarderOpOutput{
			Addresses:  ds.Addresses(),
			AddressRef: r,
		}, nil
	},
)

func deployMock(ctx context.Context, auth *bind.TransactOpts, chain evm.Chain) (*common.Address, *cldf.TypeAndVersion, error) {
	forwarderAddr, tx, mockForwarderContract, err := mock_forwarder.DeployMockKeystoneForwarder(
		auth,
		chain.Client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy MockKeystoneForwarder: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to confirm and save MockKeystoneForwarder: %w", err)
	}
	tvStr, err := mockForwarderContract.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get type and version: %w", err)
	}
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}
	txHash := tx.Hash()
	txReceipt, err := chain.Client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}
	hashLabel := fmt.Sprintf("%s: %s", DeploymentHashLabel, txHash.Hex())
	blockLabel := fmt.Sprintf("%s: %s", DeploymentBlockLabel, txReceipt.BlockNumber.String())
	tv.Labels.Add(blockLabel)
	tv.Labels.Add(hashLabel)

	return &forwarderAddr, &tv, nil
}
