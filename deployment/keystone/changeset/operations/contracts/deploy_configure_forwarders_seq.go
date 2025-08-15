package contracts

import (
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type DeployConfigureForwardersSeqDeps struct {
	Env         *cldf.Environment
	Registry    *capabilities_registry.CapabilitiesRegistry
	RegistryRef datastore.AddressRefKey
	// this is for writer don
	WriteCapabilityConfigs []internal.CapabilityConfig
	P2pToWriteCapabilities map[p2pkey.PeerID][]capabilities_registry.CapabilitiesRegistryCapability
}

type DeployConfigureForwardersSeqInput struct {
	// chains to deploy forwarders to
	ForwaderDeploymentChains []uint64
	// capabilities registry chain selector
	RegistryChainSel uint64
	// configure specific chain forwarders to specific workflow dons
	Chain2WfDonMap map[uint64][]ConfigureKeystoneDON
	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *changeset.MCMSConfig
}

func (i DeployConfigureForwardersSeqInput) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type DeployConfigureForwardersSeqOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
	Addresses             datastore.AddressRefStore
	AddressBook           cldf.AddressBook // The address book containing the deployed Keystone Forwarders
}

var DeployConfigureForwardersSeq = operations.NewSequence[DeployConfigureForwardersSeqInput, DeployConfigureForwardersSeqOutput, DeployConfigureForwardersSeqDeps](
	"deploy-configure-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Configure Keystone Forwarders",
	func(b operations.Bundle, deps DeployConfigureForwardersSeqDeps, input DeployConfigureForwardersSeqInput) (DeployConfigureForwardersSeqOutput, error) {
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()
		var proposals []mcms.TimelockProposal

		// forwarder deployment
		contractErrGroup := &errgroup.Group{}
		for _, target := range input.ForwaderDeploymentChains {
			contractErrGroup.Go(func() error {
				dep := DeployForwarderOpDeps{
					Env: deps.Env,
				}
				r, err := operations.ExecuteOperation(b, DeployKeystoneForwarderOp, dep, DeployForwarderOpInput{
					ChainSelector: target,
				})
				if err != nil {
					return err
				}
				// merge address book
				err = ab.Merge(r.Output.AddressBook)
				if err != nil {
					return fmt.Errorf("failed to save Keystone Forwarder address on address book for target %d: %w", target, err)
				}
				// merge address store
				addrs, err := r.Output.Addresses.Fetch()
				if err != nil {
					return fmt.Errorf("failed to fetch Keystone Forwarder addresses for target %d: %w", target, err)
				}
				for _, addr := range addrs {
					if addrRefErr := as.AddressRefStore.Add(addr); addrRefErr != nil {
						return fmt.Errorf("failed to save Keystone Forwarder address on datastore for target %d: %w", target, addrRefErr)
					}
				}

				return nil
			})
		}
		if err := contractErrGroup.Wait(); err != nil {
			return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to deploy Keystone contracts: %w", err)
		}

		// forwarder configuration
		evmChain := deps.Env.BlockChains.EVMChains()
		var out DeployConfigureForwardersSeqOutput

		donMap := make(map[string]internal.RegisteredDon)
		for chainSel, chainDonsToConfigure := range input.Chain2WfDonMap {
			chainDons, err := resolveChainDons(*deps.Env, input.RegistryChainSel, deps.Registry, chainDonsToConfigure, donMap)
			if err != nil {
				return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: failed to resolve DONs for chain %d: %w", chainSel, err)
			}
			chain := evmChain[chainSel]
			contracts, err := resolveForwarderContracts(*deps.Env, chain, as.Addresses(), input.ForwaderDeploymentChains)
			if err != nil {
				return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: failed to resolve forwarder contracts for chain %d: %w", chainSel, err)
			}
			if len(contracts) == 0 {
				return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: no KeystoneForwarder contract found for chain selector %d", chainSel)
			}
			// for each forwarder -> execute op -> create proposal if using MCMS
			for _, contract := range contracts {
				b.Logger.Info("ENTERING HERE during ExecuteOperation", chainDons)
				fwrReport, err := operations.ExecuteOperation(b, ConfigureForwarderOp, ConfigureForwarderOpDeps{
					Env:      deps.Env,
					Chain:    &chain,
					Contract: contract.Contract,
					Dons:     chainDons,
				}, ConfigureForwarderOpInput{
					UseMCMS:       input.UseMCMS(),
					ChainSelector: chainSel, // here to skip the check for the previous report, since unless inputs are different they are treated as the same and skipped
				})
				if err != nil {
					return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed for chain selector %d: %w", chainSel, err)
				}
				// no mcms required or newly deployed forwarder -> skip proposal generation
				if !input.UseMCMS() || slices.Contains(input.ForwaderDeploymentChains, chainSel) {
					continue
				}
				// we are doing this for each forwarder because i dont know if they are all owned by the same MCMS
				// so we need to build a proposal for each forwarder
				timelocksPerChain, proposerMCMSes, inspectorPerChain, err := resolveMCMSContracts(*deps.Env, chain, contract)
				if err != nil {
					return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: failed to resolve MCMS contracts for chain selector %d: %w", chainSel, err)
				}

				proposal, err := proposalutils.BuildProposalFromBatchesV2(
					*deps.Env,
					timelocksPerChain,
					proposerMCMSes,
					inspectorPerChain,
					[]mcmstypes.BatchOperation{fwrReport.Output.BatchOperation},
					"proposal to set forwarder config",
					proposalutils.TimelockConfig{
						MinDelay: input.MCMSConfig.MinDuration,
					},
				)
				if err != nil {
					return out, fmt.Errorf("configure-forwarders-seq failed: failed to build proposal: %w", err)
				}
				proposals = append(proposals, *proposal)
			}
		}

		// check for mcms on chain
		// if mcms are not found, deploy them
		// transfer ownership of forwarder to mcms

		// append capabilities to the donNames
		appendCapabilitiesReport, err := operations.ExecuteOperation(b, AppendCapabilitiesOp, AppendCapabilitiesOpDeps{
			Env:               deps.Env,
			RegistryRef:       deps.RegistryRef,
			P2pToCapabilities: deps.P2pToWriteCapabilities,
		}, AppendCapabilitiesOpInput{
			RegistryChainSel: input.RegistryChainSel,
			MCMSConfig:       input.MCMSConfig,
		})
		if err != nil {
			return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("append-capabilities-op failed: %w", err)
		}
		if input.UseMCMS() {
			proposals = append(proposals, appendCapabilitiesReport.Output.MCMSTimelockProposals...)
		}

		// update don
		p2pIDs := make([]p2pkey.PeerID, 0, len(deps.P2pToWriteCapabilities))
		for p2pID := range deps.P2pToWriteCapabilities {
			p2pIDs = append(p2pIDs, p2pID)
		}
		updateDonReport, err := operations.ExecuteOperation(b, UpdateDonOp, UpdateDonOpDeps{
			Env:               deps.Env,
			RegistryRef:       deps.RegistryRef,
			P2PIDs:            p2pIDs,
			CapabilityConfigs: deps.WriteCapabilityConfigs,
		}, UpdateDonOpInput{
			RegistryChainSel: input.RegistryChainSel,
			MCMSConfig:       input.MCMSConfig,
		})
		if err != nil {
			return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("update-don-op failed: %w", err)
		}
		if input.UseMCMS() {
			proposals = append(proposals, updateDonReport.Output.MCMSTimelockProposals...)
		}

		return DeployConfigureForwardersSeqOutput{
			MCMSTimelockProposals: proposals,
			AddressBook:           ab,
			Addresses:             as.Addresses(),
		}, nil
	},
)

// resolveChainDons resolves DONs for a given chain, using caching to avoid recreating identical DONs
func resolveChainDons(
	env cldf.Environment,
	registryChainSel uint64,
	registry *capabilities_registry.CapabilitiesRegistry,
	chainDonsToConfigure []ConfigureKeystoneDON,
	donCache map[string]internal.RegisteredDon,
) ([]internal.RegisteredDon, error) {
	chainDons := make([]internal.RegisteredDon, 0, len(chainDonsToConfigure))

	for _, don := range chainDonsToConfigure {
		// Check if DON already exists in cache
		if cachedDon, exists := donCache[don.Name]; exists {
			chainDons = append(chainDons, cachedDon)
			continue
		}

		// Create new RegisteredDon
		donConfig := internal.RegisteredDonConfig{
			NodeIDs:          don.NodeIDs,
			Name:             don.Name,
			RegistryChainSel: registryChainSel,
			Registry:         registry,
		}

		registeredDon, err := internal.NewRegisteredDon(env, donConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create registered DON %s: %w", don.Name, err)
		}

		// Cache for future use
		donCache[don.Name] = *registeredDon
		chainDons = append(chainDons, *registeredDon)
	}

	return chainDons, nil
}

func resolveForwarderContracts(
	env cldf.Environment,
	chain evm.Chain,
	newlyDeployedForwarders datastore.AddressRefStore,
	forwaderDeploymentChains []uint64,
) ([]*changeset.OwnedContract[*forwarder.KeystoneForwarder], error) {
	chainSelector := chain.Selector
	contracts := make([]*changeset.OwnedContract[*forwarder.KeystoneForwarder], 0)
	if slices.Contains(forwaderDeploymentChains, chainSelector) {
		addressesRefs := newlyDeployedForwarders.Filter(
			datastore.AddressRefByChainSelector(chainSelector),
			datastore.AddressRefByType(datastore.ContractType(changeset.KeystoneForwarder)),
		)
		contract, err := changeset.GetOwnedContractV2[*forwarder.KeystoneForwarder](newlyDeployedForwarders, chain, addressesRefs[0].Address)
		if err != nil {
			return nil, fmt.Errorf("configure-forwarders-seq failed: failed to get KeystoneForwarder contract for chain selector %d: %w", chainSelector, err)
		}
		contracts = append(contracts, contract)
		return contracts, nil
	}

	// existing forwarders
	addressesRefs := env.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(datastore.ContractType(changeset.KeystoneForwarder)),
	)
	for _, addrRef := range addressesRefs {
		contract, err := changeset.GetOwnedContractV2[*forwarder.KeystoneForwarder](env.DataStore.Addresses(), chain, addrRef.Address)
		if err != nil {
			return nil, fmt.Errorf("configure-forwarders-seq failed: failed to get KeystoneForwarder contract for chain selector %d: %w", chainSelector, err)
		}
		contracts = append(contracts, contract)
	}
	return contracts, nil
}

func resolveMCMSContracts(
	env cldf.Environment,
	chain evm.Chain,
	forwarderContract *changeset.OwnedContract[*forwarder.KeystoneForwarder],
) (map[uint64]string, map[uint64]string, map[uint64]mcmssdk.Inspector, error) {
	if forwarderContract.McmsContracts == nil {
		return nil, nil, nil, fmt.Errorf("configure-forwarders-seq failed: expected forwarder contract %s to be owned by MCMS for chain selector %d", forwarderContract.Contract.Address(), chain.Selector)
	}
	timelocksPerChain := map[uint64]string{
		chain.Selector: forwarderContract.McmsContracts.Timelock.Address().Hex(),
	}
	proposerMCMSes := map[uint64]string{
		chain.Selector: forwarderContract.McmsContracts.ProposerMcm.Address().Hex(),
	}
	inspector, err := proposalutils.McmsInspectorForChain(env, chain.Selector)
	if err != nil {
		return nil, nil, nil, err
	}
	inspectorPerChain := map[uint64]mcmssdk.Inspector{
		chain.Selector: inspector,
	}
	return timelocksPerChain, proposerMCMSes, inspectorPerChain, nil
}
