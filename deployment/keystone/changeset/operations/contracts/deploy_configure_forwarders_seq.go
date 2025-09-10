package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsevm "github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	mcmsOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms/ops"
	mcmsSeqs "github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms/seqs"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// Use this to deploy keystone forwarders and configure them with the given DONs on NEW CHAINS
type DeployConfigureForwardersSeqDeps struct {
	Env         *cldf.Environment
	Registry    *capabilities_registry.CapabilitiesRegistry
	RegistryRef datastore.AddressRefKey
	// this is for writer don
	WriteCapabilityConfigs []internal.CapabilityConfig
	P2pToWriteCapabilities map[p2pkey.PeerID][]capabilities_registry.CapabilitiesRegistryCapability
}

type ForwarderDeploymentOps struct {
	// Qualifier is the qualifier of the forwarder. If not empty, it will be used to check if the forwarder is already deployed.
	Qualifier string
	// Dons is the list of DONs to configure on the forwarder.
	WfDons []ConfigureKeystoneDON
}

type DeployConfigureForwardersSeqInput struct {
	// chains to deploy forwarders to
	ForwarderDeploymentChains map[uint64]ForwarderDeploymentOps
	// capabilities registry chain selector
	RegistryChainSel uint64
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
		lggr := deps.Env.Logger
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()
		batches := []mcmstypes.BatchOperation{}
		timelockAddressByChain := make(map[uint64]string)
		inspectorPerChain := map[uint64]sdk.Inspector{}
		proposerAddressByChain := make(map[uint64]string)
		proposals := []mcms.TimelockProposal{}
		evmChain := deps.Env.BlockChains.EVMChains()
		donMapCache := make(map[string]internal.RegisteredDon) // cache for registered dons
		// 1. forwarder deployment, forwarder configuration and ownership transfer to timelock if MCMS is used
		for target, ops := range input.ForwarderDeploymentChains {
			allChainAddresses := deps.Env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(target))
			var timelockAddr common.Address
			var proposerAddr common.Address
			forwarderTV := cldf.NewTypeAndVersion(internal.KeystoneForwarder, deployment.Version1_1_0)
			// check if the forwarder is already deployed
			// extract timelock and proposer addresses
			for _, addr := range allChainAddresses {
				if addr.Type == datastore.ContractType(forwarderTV.Type) && addr.Qualifier == ops.Qualifier {
					b.Logger.Infof("Skipping forwarder deployment and configuration for chain selector %d as it already exists", target)
					continue
				}
				if addr.Type == datastore.ContractType(types.RBACTimelock) {
					timelockAddr = common.HexToAddress(addr.Address)
				}
				if addr.Type == datastore.ContractType(types.ProposerManyChainMultisig) {
					proposerAddr = common.HexToAddress(addr.Address)
				}
			}

			// 1.1 deploy forwarder
			lggr.Infof("Deploying Keystone Forwarder for chain selector %d", target)
			forwarderAddress, err := deployForwarderOp(b, deps, target, ops.Qualifier, ab, as)
			if err != nil {
				return DeployConfigureForwardersSeqOutput{
					AddressBook: ab,
					Addresses:   as.Addresses(),
				}, fmt.Errorf("failed to deploy Keystone Forwarder for target %d: %w", target, err)
			}

			// 1.2 configure forwarder
			chain := evmChain[target]
			if len(ops.WfDons) != 0 {
				lggr.Infof("Configuring Keystone Forwarder for chain selector %d with %d DONs", target, len(ops.WfDons))
				err := configureForwarderOp(b, deps, input, target, forwarderAddress, chain, donMapCache, ops.WfDons, as)
				if err != nil {
					return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to configure Keystone Forwarder for target %d: %w", target, err)
				}
			}

			// check if user wants to use MCMS and if yes, if timelock and proposer are deployed
			if !input.UseMCMS() || (timelockAddr == (common.Address{}) || proposerAddr == (common.Address{})) {
				lggr.Infof("Skipping ownership transfer of forwarder as no timelock found for chain selector %d", target)
				continue
			}
			// 1.3 transfer ownership to timelock if MCMS is used
			timelockAddressByChain[target] = timelockAddr.String()
			proposerAddressByChain[target] = proposerAddr.String()
			inspectorPerChain[target] = mcmsevm.NewInspector(chain.Client)
			lggr.Infof("Transferring ownership of Keystone Forwarder to timelock for chain selector %d", target)
			err = tranferOwnershipOp(b, target, forwarderAddress, chain, &batches, timelockAddr)
			if err != nil {
				return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to transfer ownership of forwarder to timelock for chain selector %d: %w", target, err)
			}
		}
		// build proposal for transfer ownership to timelock if MCMS is used
		if input.UseMCMS() {
			b.Logger.Infof("Building MCMS proposal for ownership transfer to timelock")
			proposal, err := proposalutils.BuildProposalFromBatchesV2(
				*deps.Env,
				timelockAddressByChain, proposerAddressByChain, inspectorPerChain,
				batches, "Transfer ownership to timelock", proposalutils.TimelockConfig{
					MinDelay: input.MCMSConfig.MinDuration,
				})
			if err != nil {
				return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to build proposal for transfer ownership to timelock: %w", err)
			}
			proposals = append(proposals, *proposal)
		}

		// 2. append new write capabilities to the writer dons
		lggr.Infof("Appending new write capabilities to the writer DONs on registry chain selector %d", input.RegistryChainSel)
		err := appendCapabilitiesOp(b, deps, input, &proposals)
		if err != nil {
			return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to append capabilities: %w", err)
		}
		// 3. update writer don on registry with new capabilities
		lggr.Infof("Updating writer DON on registry with new capabilities on registry chain selector %d", input.RegistryChainSel)
		err = updateDonOp(b, deps, input, &proposals)
		if err != nil {
			return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to update writer don: %w", err)
		}
		// TODO: aggregate proposals
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

func deployForwarderOp(
	b operations.Bundle,
	deps DeployConfigureForwardersSeqDeps,
	target uint64,
	qualifier string,
	ab cldf.AddressBook,
	as *datastore.MemoryDataStore,
) (common.Address, error) {
	deployForwarderDep := DeployForwarderOpDeps{Env: deps.Env}
	deployForwarderInput := DeployForwarderOpInput{ChainSelector: target, Qualifier: qualifier}
	deployForwarderReport, err := operations.ExecuteOperation(b, DeployKeystoneForwarderOp, deployForwarderDep, deployForwarderInput)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy Keystone Forwarder for target %d: %w", target, err)
	}

	// merge address book
	err = ab.Merge(deployForwarderReport.Output.AddressBook)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to save Keystone Forwarder address on address book for target %d: %w", target, err)
	}
	// merge address store
	addrs, err := deployForwarderReport.Output.Addresses.Fetch()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to fetch Keystone Forwarder addresses for target %d: %w", target, err)
	}
	forwarderAddrRef := addrs[0]
	if addrRefErr := as.Addresses().Add(forwarderAddrRef); addrRefErr != nil {
		return common.Address{}, fmt.Errorf("failed to save Keystone Forwarder address on datastore for target %d: %w", target, addrRefErr)
	}
	return common.HexToAddress(forwarderAddrRef.Address), nil
}

func configureForwarderOp(
	b operations.Bundle,
	deps DeployConfigureForwardersSeqDeps,
	input DeployConfigureForwardersSeqInput,
	target uint64,
	forwarderAddress common.Address,
	chain evm.Chain,
	donMapCache map[string]internal.RegisteredDon,
	wfDons []ConfigureKeystoneDON,
	as *datastore.MemoryDataStore,
) error {
	chainDons, err := resolveChainDons(*deps.Env, input.RegistryChainSel, deps.Registry, wfDons, donMapCache)
	if err != nil {
		return fmt.Errorf("configure-forwarders-seq failed: failed to resolve DONs for chain %d: %w", target, err)
	}
	forwarderContract, err := crecontracts.GetOwnedContractV2[*forwarder.KeystoneForwarder](as.Addresses(), chain, forwarderAddress.String())
	if err != nil {
		return fmt.Errorf("configure-forwarders-seq failed: failed to get KeystoneForwarder contract for chain selector %d: %w", target, err)
	}
	_, err = operations.ExecuteOperation(b, ConfigureForwarderOp, ConfigureForwarderOpDeps{
		Env:      deps.Env,
		Chain:    &chain,
		Contract: forwarderContract.Contract,
		Dons:     chainDons,
	}, ConfigureForwarderOpInput{UseMCMS: input.UseMCMS(), ChainSelector: target})
	if err != nil {
		return fmt.Errorf("configure-forwarders-seq failed for chain selector %d: %w", target, err)
	}
	return nil
}

func tranferOwnershipOp(
	b operations.Bundle,
	target uint64,
	forwarderAddress common.Address,
	chain evm.Chain,
	batches *[]mcmstypes.BatchOperation,
	timelockAddr common.Address,
) error {
	_, forwarderOwnableContract, err := mcmsSeqs.LoadOwnableContract(forwarderAddress, chain.Client)
	if err != nil {
		return fmt.Errorf("failed to load ownable contract for chain selector %d: %w", target, err)
	}
	// transfer ownership to timelock (we send this on chain directly, does not need to go through MCMS)
	_, err = operations.ExecuteOperation(b, mcmsOps.OpEVMTransferOwnership,
		mcmsOps.OpEVMOwnershipDeps{
			Chain:    chain,
			OwnableC: forwarderOwnableContract,
		},
		mcmsOps.OpEVMTransferOwnershipInput{
			ChainSelector:   target,
			TimelockAddress: timelockAddr,
			Address:         forwarderAddress,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to transfer ownership of forwarder to timelock for chain selector %d: %w", target, err)
	}
	// accept ownership as timelock (timelock needs to sign this, so we send it through MCMS)
	acceptOwnershipReport, err := operations.ExecuteOperation(b, mcmsOps.OpEVMAcceptOwnership,
		mcmsOps.OpEVMOwnershipDeps{
			Chain:    chain,
			OwnableC: forwarderOwnableContract,
		},
		mcmsOps.OpEVMTransferOwnershipInput{
			ChainSelector:   target,
			TimelockAddress: timelockAddr,
			Address:         forwarderAddress,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to transfer ownership of forwarder to timelock for chain selector %d: %w", target, err)
	}
	mcmsTx := mcmstypes.Transaction{
		To:               forwarderAddress.String(),
		Data:             acceptOwnershipReport.Output.Tx.Data(),
		AdditionalFields: json.RawMessage(`{"value": 0}`), // JSON-encoded `{"value": 0}`
	}
	*batches = append(*batches, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(target),
		Transactions:  []mcmstypes.Transaction{mcmsTx},
	})
	return nil
}

func appendCapabilitiesOp(
	b operations.Bundle,
	deps DeployConfigureForwardersSeqDeps,
	input DeployConfigureForwardersSeqInput,
	proposals *[]mcms.TimelockProposal,
) error {
	appendCapabilitiesReport, err := operations.ExecuteOperation(b, AppendCapabilitiesOp, AppendCapabilitiesOpDeps{
		Env:               deps.Env,
		RegistryRef:       deps.RegistryRef,
		P2pToCapabilities: deps.P2pToWriteCapabilities,
	}, AppendCapabilitiesOpInput{
		RegistryChainSel: input.RegistryChainSel,
		MCMSConfig:       input.MCMSConfig,
	})
	if err != nil {
		return fmt.Errorf("append-capabilities-op failed: %w", err)
	}
	// if MCMS is used, append the proposal to the list of proposals
	if input.UseMCMS() {
		*proposals = append(*proposals, appendCapabilitiesReport.Output.MCMSTimelockProposals...)
	}
	return nil
}

func updateDonOp(
	b operations.Bundle,
	deps DeployConfigureForwardersSeqDeps,
	input DeployConfigureForwardersSeqInput,
	proposals *[]mcms.TimelockProposal,
) error {
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
		return fmt.Errorf("update-don-op failed: %w", err)
	}
	// if MCMS is used, append the proposal to the list of proposals
	if input.UseMCMS() {
		*proposals = append(*proposals, updateDonReport.Output.MCMSTimelockProposals...)
	}
	return nil
}
