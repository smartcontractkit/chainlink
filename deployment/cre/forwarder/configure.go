package forwarder

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type ConfigureSeqDeps struct {
	Env *cldf.Environment
}

type DonConfiguration struct {
	Name    string   // name of the DON
	ID      uint32   // the DON id as registered in the capabilities registry. Is an id corresponding to a DON that run consensus capability
	F       uint8    // the F value for the DON as registered in the capabilities registry
	Version uint32   // the config version for the DON as registered in the capabilities registry
	NodeIDs []string // node IDs (JD IDs or PeerIDs starting with `p2p_` or csa keys) of the nodes in the DON
}

func (d DonConfiguration) ForwarderConfig(chainFamily string, c offchain.Client) (Config, error) {
	signers, err := Signers(d.NodeIDs, c, chainFamily)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get signers for DON %s: %w", d.Name, err)
	}
	if len(signers) == 0 {
		return Config{}, fmt.Errorf("no signers found for DON %s", d.Name)
	}
	return Config{
		DonID:         d.ID,
		F:             d.F,
		ConfigVersion: d.Version,
		Signers:       signers,
	}, nil
}

type ConfigureSeqInput struct {
	DON DonConfiguration // the DON to configuration for the forwarder to accept

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *proposalutils.TimelockConfig
	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains map[uint64]struct{}
}

func (i ConfigureSeqInput) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type ConfigureSeqOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
	Config                Config
}

// ConfigureSeq is a sequence that configures Keystone Forwarder contracts for a given DON.
// TODO this is mostly copied from keystone/changeset/operations/contracts/configure_forwarders_seq.go
// now that this is independent of the registry, we should be able to use this in the impl there while maintaining the existing api if needed
var ConfigureSeq = operations.NewSequence[ConfigureSeqInput, ConfigureSeqOutput, ConfigureSeqDeps](
	"configure-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Configure Keystone Forwarders",
	func(b operations.Bundle, deps ConfigureSeqDeps, input ConfigureSeqInput) (ConfigureSeqOutput, error) {
		evmChain := deps.Env.BlockChains.EVMChains()
		opPerChain := make(map[uint64]*mcmstypes.BatchOperation)
		forwarderContracts := make(map[uint64]*contracts.OwnedContract[*forwarder.KeystoneForwarder])

		cfg, err := input.DON.ForwarderConfig("evm", deps.Env.Offchain)
		out := ConfigureSeqOutput{
			Config: cfg,
		}
		if err != nil {
			return ConfigureSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: failed to get forwarder config for DON %s: %w", input.DON.Name, err)
		}
		for _, chain := range evmChain {
			if _, shouldInclude := input.Chains[chain.Selector]; len(input.Chains) > 0 && !shouldInclude {
				continue
			}

			addressesRefs := deps.Env.DataStore.Addresses().Filter(
				datastore.AddressRefByChainSelector(chain.Selector),
				datastore.AddressRefByType(datastore.ContractType(contracts.KeystoneForwarder)),
			)
			if len(addressesRefs) == 0 {
				return ConfigureSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: no KeystoneForwarder contract found for chain selector %d", chain.Selector)
			}

			for _, addrRef := range addressesRefs {
				contract, err := contracts.GetOwnedContractV2[*forwarder.KeystoneForwarder](deps.Env.DataStore.Addresses(), chain, addrRef.Address)
				if err != nil {
					return ConfigureSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: failed to get KeystoneForwarder contract for chain selector %d: %w", chain.Selector, err)
				}

				fwrReport, err := operations.ExecuteOperation(b, ConfigureOp, ConfigureOpDeps{
					Env:      deps.Env,
					Chain:    &chain,
					Contract: contract.Contract,
					Config:   cfg,
				}, ConfigureOpInput{
					UseMCMS:       input.UseMCMS(),
					ChainSelector: chain.Selector, // here to skip the check for the previous report, since unless inputs are different they are treated as the same and skipped
				})
				if err != nil {
					return ConfigureSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed for chain selector %d: %w", chain.Selector, err)
				}

				opPerChain[chain.Selector] = fwrReport.Output.BatchOperation
				forwarderContracts[chain.Selector] = contract
			}
		}

		if input.UseMCMS() {
			if len(opPerChain) == 0 {
				return out, errors.New("configure-forwarders-seq failed: no operations generated for MCMS")
			}

			for chainSelector, op := range opPerChain {
				fwr, ok := forwarderContracts[chainSelector]
				if !ok {
					return out, fmt.Errorf("configure-forwarders-seq failed: expected configured forwarder address for chain selector %d", chainSelector)
				}
				if fwr.McmsContracts == nil {
					return out, fmt.Errorf("configure-forwarders-seq failed: expected forwarder contract %s to be owned by MCMS for chain selector %d", fwr.Contract.Address(), chainSelector)
				}
				timelocksPerChain := map[uint64]string{
					chainSelector: fwr.McmsContracts.Timelock.Address().Hex(),
				}
				proposerMCMSes := map[uint64]string{
					chainSelector: fwr.McmsContracts.ProposerMcm.Address().Hex(),
				}
				inspector, err := proposalutils.McmsInspectorForChain(*deps.Env, chainSelector)
				if err != nil {
					return out, err
				}
				inspectorPerChain := map[uint64]mcmssdk.Inspector{
					chainSelector: inspector,
				}

				proposal, err := proposalutils.BuildProposalFromBatchesV2(
					*deps.Env,
					timelocksPerChain,
					proposerMCMSes,
					inspectorPerChain,
					[]mcmstypes.BatchOperation{*op},
					"proposal to set forwarder config",
					*input.MCMSConfig,
				)
				if err != nil {
					return out, fmt.Errorf("configure-forwarders-seq failed: failed to build proposal: %w", err)
				}
				out.MCMSTimelockProposals = append(out.MCMSTimelockProposals, *proposal)
			}
		}

		return out, nil
	},
)

type ConfigureOpDeps struct {
	Env      *cldf.Environment
	Chain    *evm.Chain
	Contract *forwarder.KeystoneForwarder
	Config   Config
}

type ConfigureOpInput struct {
	UseMCMS       bool
	ChainSelector uint64
}

type ConfigureOpOutput struct {
	BatchOperation *mcmstypes.BatchOperation // if using MCMS, the batch operation to propose the change

	Forwarder common.Address
	Config    Config
}

// ConfigureOp is an operation that configures a Keystone Forwarder contract.
var ConfigureOp = operations.NewOperation[ConfigureOpInput, ConfigureOpOutput, ConfigureOpDeps](
	"configure-forwarder-op",
	semver.MustParse("1.0.0"),
	"Configure Keystone Forwarder",
	func(b operations.Bundle, deps ConfigureOpDeps, input ConfigureOpInput) (ConfigureOpOutput, error) {
		r, err := configureForwarder(b.Logger, *deps.Chain, deps.Contract, deps.Config, input.UseMCMS)
		if err != nil {
			return ConfigureOpOutput{}, fmt.Errorf("configure-forwarder-op failed: failed to configure forwarder for chain selector %d: %w", deps.Chain.Selector, err)
		}
		return ConfigureOpOutput{BatchOperation: r.BatchOperation}, nil
	},
)

// Signers returns the onchain public keys of the given node IDs by retrieving the node info from the offchain client
// nodeIDs can be JD IDs or PeerIDs starting with `p2p_ or csa keys`.
func Signers(nodeIDs []string, c offchain.Client, chainFamily string) ([]common.Address, error) {
	// load the nodes from the offchain client
	nodes, err := deployment.NodeInfo(nodeIDs, c)
	if err != nil {
		return nil, fmt.Errorf("failed to get node info for node IDs %v: %w", nodeIDs, err)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].PeerID.String() < nodes[j].PeerID.String()
	})
	var out []common.Address
	for _, n := range nodes {
		if n.IsBootstrap {
			continue
		}
		var found bool
		var registryChainDetails chainsel.ChainDetails
		for details := range n.SelToOCRConfig {
			if family, err := chainsel.GetSelectorFamily(details.ChainSelector); err == nil && family == chainFamily {
				found = true
				registryChainDetails = details
			}
		}
		if !found {
			return nil, fmt.Errorf("chainType not found: %v", chainFamily)
		}
		// eth address is the first 20 bytes of the Signer
		config, exists := n.SelToOCRConfig[registryChainDetails]
		if !exists {
			return nil, fmt.Errorf("chainID not found: %v", registryChainDetails)
		}
		signer := config.OnchainPublicKey
		signerAddress := common.BytesToAddress(signer)
		out = append(out, signerAddress)
	}
	return out, nil
}

type ConfigureForwarders struct{}

var _ cldf.ChangeSetV2[ConfigureSeqInput] = ConfigureForwarders{}

func (c ConfigureForwarders) VerifyPreconditions(e cldf.Environment, config ConfigureSeqInput) error {
	for chainSel := range config.Chains {
		if _, ok := e.BlockChains.EVMChains()[chainSel]; !ok {
			return fmt.Errorf("chain selector %d not found in environment", chainSel)
		}
	}

	if config.DON.Name == "" {
		return errors.New("DON name cannot be empty")
	}
	if len(config.DON.NodeIDs) == 0 {
		return errors.New("DON must have at least one node ID")
	}

	return nil
}

func (c ConfigureForwarders) Apply(e cldf.Environment, config ConfigureSeqInput) (cldf.ChangesetOutput, error) {
	// Use ConfigureSeq which handles all dependency resolution internally
	deps := ConfigureSeqDeps{
		Env: &e,
	}

	configureReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ConfigureSeq,
		deps,
		config,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               configureReport.ExecutionReports,
		MCMSTimelockProposals: configureReport.Output.MCMSTimelockProposals,
	}, nil
}
