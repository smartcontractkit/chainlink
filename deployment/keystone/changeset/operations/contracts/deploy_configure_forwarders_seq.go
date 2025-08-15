package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

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
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// maybe lets make this a bit simple
// lets assume we only do this for new chains
// so if you are deploying a forwarder, you are configuring only that forwarder and then transferring ownership to mcms

type DeployConfigureForwardersSeqDeps struct {
	Env         *cldf.Environment
	Registry    *capabilities_registry.CapabilitiesRegistry
	RegistryRef datastore.AddressRefKey
	// this is for writer don
	WriteCapabilityConfigs []internal.CapabilityConfig
	P2pToWriteCapabilities map[p2pkey.PeerID][]capabilities_registry.CapabilitiesRegistryCapability
}

type ForwarderDeploymentOps struct {
	Override  bool
	Qualifier string
	Dons      []ConfigureKeystoneDON
}
type DeployConfigureForwardersSeqInput struct {
	// chains to deploy forwarders to
	ForwaderDeploymentChains map[uint64]ForwarderDeploymentOps
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
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()
		batches := []mcmstypes.BatchOperation{}
		timelockAddressByChain := make(map[uint64]string)
		inspectorPerChain := map[uint64]sdk.Inspector{}
		proposerAddressByChain := make(map[uint64]string)
		evmChain := deps.Env.BlockChains.EVMChains()
		var proposals []mcms.TimelockProposal
		donMap := make(map[string]internal.RegisteredDon)
		// forwarder deployment
		for target, ops := range input.ForwaderDeploymentChains {
			// check if the forwarder is already deployed
			forwarderTV := cldf.NewTypeAndVersion(internal.KeystoneForwarder, deployment.Version1_1_0)
			_, err := deps.Env.DataStore.Addresses().Get(datastore.NewAddressRefKey(
				target,
				datastore.ContractType(forwarderTV.Type),
				semver.MustParse(forwarderTV.Version.String()),
				"",
			))
			if err == nil {
				b.Logger.Infof("Skipping forwarder deployment for chain selector %d as it already exists", target)
			}

			// deploy forwarder
			dep := DeployForwarderOpDeps{
				Env: deps.Env,
			}
			r, err := operations.ExecuteOperation(b, DeployKeystoneForwarderOp, dep, DeployForwarderOpInput{
				ChainSelector: target,
			})
			if err != nil {
				return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to deploy Keystone Forwarder for target %d: %w", target, err)
			}

			// merge address book
			err = ab.Merge(r.Output.AddressBook)
			if err != nil {
				return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to save Keystone Forwarder address on address book for target %d: %w", target, err)
			}
			// merge address store
			addrs, err := r.Output.Addresses.Fetch()
			if err != nil {
				return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to fetch Keystone Forwarder addresses for target %d: %w", target, err)
			}
			forwarderAddrRef := addrs[0]
			if addrRefErr := as.AddressRefStore.Add(forwarderAddrRef); addrRefErr != nil {
				return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to save Keystone Forwarder address on datastore for target %d: %w", target, addrRefErr)
			}
			chain := evmChain[target]
			if len(ops.Dons) != 0 {
				chainDons, err := resolveChainDons(*deps.Env, input.RegistryChainSel, deps.Registry, ops.Dons, donMap)
				if err != nil {
					return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: failed to resolve DONs for chain %d: %w", target, err)
				}
				contract, err := changeset.GetOwnedContractV2[*forwarder.KeystoneForwarder](as.Addresses(), chain, forwarderAddrRef.Address)
				if err != nil {
					return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: failed to get KeystoneForwarder contract for chain selector %d: %w", target, err)
				}
				_, err = operations.ExecuteOperation(b, ConfigureForwarderOp, ConfigureForwarderOpDeps{
					Env:      deps.Env,
					Chain:    &chain,
					Contract: contract.Contract,
					Dons:     chainDons,
				}, ConfigureForwarderOpInput{
					UseMCMS:       input.UseMCMS(),
					ChainSelector: target, // here to skip the check for the previous report, since unless inputs are different they are treated as the same and skipped
				})
				if err != nil {
					return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed for chain selector %d: %w", target, err)
				}
			}

			fmt.Println(deps.Env.DataStore.Addresses().Fetch())

			// check for timelock
			allChainAddresses, err := deps.Env.DataStore.Addresses().Fetch()
			if err != nil {
				return DeployConfigureForwardersSeqOutput{}, fmt.Errorf("failed to fetch all chain addresses: %w", err)
			}
			var timelockAddr common.Address
			var proposerAddr common.Address
			for _, addr := range allChainAddresses {
				if addr.Type == datastore.ContractType(types.RBACTimelock) {
					timelockAddr = common.HexToAddress(addr.Address)
				}
				if addr.Type == datastore.ContractType(types.ProposerManyChainMultisig) {
					proposerAddr = common.HexToAddress(addr.Address)
				}
			}

			if timelockAddr == (common.Address{}) || proposerAddr == (common.Address{}) {
				b.Logger.Infof("Skipping ownership transfer of forwarder as no timelock found for chain selector %d", target)
				continue
			}

			timelockAddressByChain[target] = timelockAddr.String()
			proposerAddressByChain[target] = proposerAddr.String()
			inspectorPerChain[target] = evm.NewInspector(evmChain[target].Client)
			forwarderAddress := common.HexToAddress(forwarderAddrRef.Address)
			_, c, err := mcmsSeqs.LoadOwnableContract(forwarderAddress, evmChain[target].Client)
			if err != nil {
				return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to load ownable contract for chain selector %d: %w", target, err)
			}
			// transfer ownership to timelock
			_, err = operations.ExecuteOperation(b, mcmsOps.OpEVMTransferOwnership,
				mcmsOps.OpEVMOwnershipDeps{
					Chain:    chain,
					OwnableC: c,
				},
				mcmsOps.OpEVMTransferOwnershipInput{
					ChainSelector:   target,
					TimelockAddress: timelockAddr,
					Address:         forwarderAddress,
				},
			)
			// accept ownership as timelock
			opReport, err := operations.ExecuteOperation(b, mcmsOps.OpEVMAcceptOwnership,
				mcmsOps.OpEVMOwnershipDeps{
					Chain:    chain,
					OwnableC: c,
				},
				mcmsOps.OpEVMTransferOwnershipInput{
					ChainSelector:   target,
					TimelockAddress: timelockAddr,
					Address:         forwarderAddress,
				},
			)
			if err != nil {
				return DeployConfigureForwardersSeqOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to transfer ownership of forwarder to timelock for chain selector %d: %w", target, err)
			}
			mcmsTx := mcmstypes.Transaction{
				To:               forwarderAddress.Hex(),
				Data:             opReport.Output.Tx.Data(),
				AdditionalFields: json.RawMessage(`{"value": 0}`), // JSON-encoded `{"value": 0}`
			}
			batches = append(batches, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(target),
				Transactions:  []mcmstypes.Transaction{mcmsTx},
			})
		}

		if input.UseMCMS() {
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
