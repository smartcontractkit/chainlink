package changeset

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	creforwarder "github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSet[DeployForwarderRequest] = DeployForwarder

type DeployForwarderRequest struct {
	Qualifier      string
	ChainSelectors []uint64 // filter to only deploy to these chains; if empty, deploy to all chains
}

func DeployForwarder(env cldf.Environment, cfg DeployForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	out.DataStore = datastore.NewMemoryDataStore()
	out.AddressBook = cldf.NewMemoryAddressBook() //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore

	selectors := cfg.ChainSelectors
	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(env.BlockChains.EVMChains()))
	}

	for _, sel := range selectors {
		report, err := operations.ExecuteOperation(
			env.OperationsBundle,
			creforwarder.DeployOp,
			creforwarder.DeployOpDeps{Env: &env}, creforwarder.DeployOpInput{
				ChainSelector: sel,
				Qualifier:     cfg.Qualifier,
			})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", sel, err)
		}
		out.Reports = append(out.Reports, report.ToGenericReport())
		// merge the datastore outputs
		if err := out.DataStore.Addresses().Add(report.Output.AddressRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge datastore for chain selector %d: %w", sel, err)
		}
		// merge the address book outputs
		if err := out.AddressBook.Merge(report.Output.AddressBook); err != nil { //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for chain selector %d: %w", sel, err)
		}
	}

	return out, nil
}

// DeployForwarderV2 deploys the KeystoneForwarder contract to the specified chain
func DeployForwarderV2(env cldf.Environment, req *DeployRequestV2) (cldf.ChangesetOutput, error) {
	d := func(ctx context.Context, chain cldf_evm.Chain, ab cldf.AddressBook) (*internal.DeployResponse, error) {
		report, err := operations.ExecuteOperation(
			env.OperationsBundle,
			creforwarder.DeployOp,
			creforwarder.DeployOpDeps{Env: &env}, creforwarder.DeployOpInput{
				ChainSelector: req.ChainSel,
				Qualifier:     req.Qualifier,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", req.ChainSel, err)
		}
		tv := cldf.TypeAndVersion{
			Type:    cldf.ContractType(report.Output.AddressRef.Type),
			Version: *report.Output.AddressRef.Version,
			Labels:  cldf.NewLabelSet(report.Output.AddressRef.Labels.List()...),
		}
		err = ab.Save(chain.Selector, report.Output.AddressRef.Address, tv)
		if err != nil {
			return nil, fmt.Errorf("failed to save KeystoneForwarder: %w", err)
		}

		return &internal.DeployResponse{
			Address: common.HexToAddress(report.Output.AddressRef.Address),
			Tv:      tv,
		}, nil
	}
	req.deployFn = d
	return deploy(env, req)
}

var _ cldf.ChangeSet[ConfigureForwardContractsRequest] = ConfigureForwardContracts

type ConfigureForwardContractsRequest struct {
	WFDonName string
	// workflow don node ids in the offchain client. Used to fetch and derive the signer keys
	WFNodeIDs        []string
	RegistryChainSel uint64

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains map[uint64]struct{}
}

func (r ConfigureForwardContractsRequest) Validate() error {
	if len(r.WFNodeIDs) == 0 {
		return errors.New("WFNodeIDs must not be empty")
	}
	return nil
}

func (r ConfigureForwardContractsRequest) UseMCMS() bool {
	return r.MCMSConfig != nil
}

// TODO: use crefowarder.ConfigureOP instead of internal.ConfigureForwardContracts
func ConfigureForwardContracts(env cldf.Environment, req ConfigureForwardContractsRequest) (cldf.ChangesetOutput, error) {
	wfDon, err := internal.NewRegisteredDon(env, internal.RegisteredDonConfig{
		NodeIDs:          req.WFNodeIDs,
		Name:             req.WFDonName,
		RegistryChainSel: req.RegistryChainSel,
	})
	cfg := creforwarder.DonConfiguration{
		Name: req.WFDonName,
		ID:   wfDon.Info.Id,
		F:    wfDon.Info.F,
		// use the next config version since we are going to update the config
		Version: wfDon.Info.ConfigCount,
		NodeIDs: req.WFNodeIDs,
	}

	var mcmsConfig *proposalutils.TimelockConfig
	if req.MCMSConfig != nil {
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: req.MCMSConfig.MinDuration,
		}
	}
	seqReport, err := operations.ExecuteSequence(
		env.OperationsBundle,
		creforwarder.ConfigureSeq,
		creforwarder.ConfigureSeqDeps{Env: &env},
		creforwarder.ConfigureSeqInput{
			DON:        cfg,
			MCMSConfig: mcmsConfig,
			Chains:     req.Chains,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute configure forwarder sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports:               seqReport.ExecutionReports,
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
	}, nil
}
