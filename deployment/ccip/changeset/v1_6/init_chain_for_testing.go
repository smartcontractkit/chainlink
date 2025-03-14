package v1_6

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
	mcmslib "github.com/smartcontractkit/mcms"
)

var InitChainForTestingChangeset = deployment.CreateChangeSet(initChainForTestingLogic, initChainForTestingPrecondition)

var (
	PREFIX      = new(big.Int).Lsh(big.NewInt(0x000a), 240) // 0x000a << 240
	PREFIX_MASK = new(big.Int).Lsh(big.NewInt(0xFFFF), 240) // 0xFFFF << 240
)

type ChainDefinition struct {
	// RMNVerificationDisabled is true if we do not want the RMN to bless messages from this chain.
	RMNVerificationDisabled bool
	// AllowListEnabled is true if we want an allowlist to dictate who can send messages to this chain.
	AllowListEnabled bool
	// Selector is the chain selector of this chain.
	Selector uint64
	// GasPrice defines the USD price (18 decimals) per unit gas for this chain as a destination.
	GasPrice *big.Int
	// TokenPrices define the USD price (18 decimals) per 1e18 of the smallest token denomination for various tokens on this chain.
	TokenPrices map[common.Address]*big.Int
	// FeeQuoterDestChainConfig is the configuration on a fee quoter for this destination.
	FeeQuoterDestChainConfig fee_quoter.FeeQuoterDestChainConfig
}

type NewChainDefinition struct {
	ChainDefinition
	ChainContractParams
	ExistingContracts []commoncs.Contract
	ConfigOnHome      ChainConfig
	OCRParams         CCIPOCRParams
}

// InitChainForTestingConfig is a configuration struct for InitChainForTesting.
type InitChainForTestingConfig struct {
	HomeChainSelector       uint64
	HomeChainID             int64
	HomeConfigType          globals.ConfigType
	CCIPHomeDeploymentBlock uint64
	FeedChainSelector       uint64
	NewChain                NewChainDefinition
	RemoteChains            []ChainDefinition
	MCMSConfig              commontypes.MCMSWithTimelockConfigV2
	MCMSAction              timelock.TimelockOperation
}

func initChainForTestingPrecondition(e deployment.Environment, c InitChainForTestingConfig) error {
	// TODO

	return nil
}

func initChainForTestingLogic(e deployment.Environment, c InitChainForTestingConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	var allProposals []mcmslib.TimelockProposal
	mcmsConfig := &changeset.MCMSConfig{
		MinDelay:   time.Duration(c.MCMSConfig.TimelockMinDelay.Int64()),
		MCMSAction: c.MCMSAction,
	}

	// 1. Deploy the prerequisite contracts to the new chain
	err := runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return changeset.DeployPrerequisitesChangeset(e, changeset.DeployPrerequisiteConfig{
			Configs: []changeset.DeployPrerequisiteConfigPerChain{
				changeset.DeployPrerequisiteConfigPerChain{
					ChainSelector: c.NewChain.Selector,
				},
			},
		})
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployPrerequisitesChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 2. Save existing contracts
	err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return commoncs.SaveExistingContractsChangeset(e, commoncs.ExistingContractsConfig{
			ExistingContracts: c.NewChain.ExistingContracts,
		})
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SaveExistingContractsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 3. Deploy MCMS contracts
	err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return commoncs.DeployMCMSWithTimelockV2(e, map[uint64]commontypes.MCMSWithTimelockConfigV2{
			c.NewChain.Selector: c.MCMSConfig,
		})
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployMCMSWithTimelockV2 on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 4. Deploy chain contracts to the new chain
	err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return DeployChainContractsChangeset(e, DeployChainContractsConfig{
			HomeChainSelector: c.HomeChainSelector,
			ContractParamsPerChain: map[uint64]ChainContractParams{
				c.NewChain.Selector: c.NewChain.ChainContractParams,
			},
		})
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployChainContractsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 5. Update the fee quoter prices on the new chain
	gasPrices := make(map[uint64]*big.Int, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		gasPrices[remoteChain.Selector] = remoteChain.GasPrice
	}
	_, err = UpdateFeeQuoterPricesChangeset(e, UpdateFeeQuoterPricesConfig{
		PricesByChain: map[uint64]FeeQuoterPriceUpdatePerSource{
			c.NewChain.Selector: FeeQuoterPriceUpdatePerSource{
				TokenPrices: c.NewChain.TokenPrices,
				GasPrices:   gasPrices,
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterPricesChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 6. Update the fee quoter destinations on the new chain
	destChainConfigs := make(map[uint64]fee_quoter.FeeQuoterDestChainConfig, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		destChainConfigs[remoteChain.Selector] = remoteChain.FeeQuoterDestChainConfig
	}
	_, err = UpdateFeeQuoterDestsChangeset(e, UpdateFeeQuoterDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			c.NewChain.Selector: destChainConfigs,
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterDestsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 7. Add new chain config to the home chain
	out, err := UpdateChainConfigChangeset(e, UpdateChainConfigConfig{
		HomeChainSelector: c.HomeChainSelector,
		RemoteChainAdds: map[uint64]ChainConfig{
			c.NewChain.Selector: c.NewChain.ConfigOnHome,
		},
		MCMS: mcmsConfig,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateChainConfigChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	////////////////////////
	// START: HOME CHAIN PREP
	////////////////////////

	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	// Fetch the next DON ID from the capabilities registry
	donID, err := state.Chains[c.HomeChainSelector].CapabilityRegistry.GetNextDONId(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get next DON ID: %w", err)
	}
	// Precompute config digest for each plugin (TODO: compute static config)
	configCount, err := countConfigSetEvents(state.Chains[c.HomeChainSelector].CCIPHome, &bind.FilterOpts{
		Start:   c.CCIPHomeDeploymentBlock,
		End:     nil,
		Context: e.GetContext(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to count ConfigSet events emitted by CCIPHome: %w", err)
	}
	offRampAddress, err := state.GetOffRampAddressBytes(c.NewChain.Selector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get off ramp address bytes: %w", err)
	}
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get node info: %w", err)
	}
	newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
		state.Chains[c.HomeChainSelector].CCIPHome,
		e.OCRSecrets,
		offRampAddress,
		c.NewChain.Selector,
		nodes.NonBootstraps(),
		state.Chains[c.HomeChainSelector].RMNHome.Address(),
		c.NewChain.OCRParams.OCRParameters,
		c.NewChain.OCRParams.CommitOffChainConfig,
		c.NewChain.OCRParams.ExecuteOffChainConfig,
	)

	commitPluginOCR3Config, ok := newDONArgs[types.PluginTypeCCIPCommit]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing plugin %s in newDONArgs", types.PluginTypeCCIPCommit.String())
	}
	configCount++
	commitConfigDigest, err := calculateConfigDigest(donID, types.PluginTypeCCIPCommit, commitPluginOCR3Config, configCount, c.HomeChainID, state.Chains[c.HomeChainSelector].CCIPHome.Address().Hex())

	execPluginOCR3Config, ok := newDONArgs[types.PluginTypeCCIPExec]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing plugin %s in newDONArgs", types.PluginTypeCCIPExec.String())
	}
	configCount++
	execConfigDigest, err := calculateConfigDigest(donID, types.PluginTypeCCIPExec, execPluginOCR3Config, configCount, c.HomeChainID, state.Chains[c.HomeChainSelector].CCIPHome.Address().Hex())

	////////////////////////
	// END: HOME CHAIN PREP
	////////////////////////

	// 8. Add the DON to the registry and set candidate for the commit plugin
	out, err = AddDonAndSetCandidateChangeset(e, AddDonAndSetCandidateChangesetConfig{
		SetCandidateConfigBase: SetCandidateConfigBase{
			HomeChainSelector: c.HomeChainSelector,
			FeedChainSelector: c.FeedChainSelector,
			MCMS:              mcmsConfig,
		},
		PluginInfo: SetCandidatePluginInfo{
			PluginType: types.PluginTypeCCIPCommit,
			OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
				c.NewChain.Selector: c.NewChain.OCRParams,
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run AddDonAndSetCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// 9. Set the candidate for the exec plugin
	out, err = SetCandidateChangeset(e, SetCandidateChangesetConfig{
		SetCandidateConfigBase: SetCandidateConfigBase{
			HomeChainSelector: c.HomeChainSelector,
			FeedChainSelector: c.FeedChainSelector,
			MCMS:              mcmsConfig,
		},
		PluginInfo: []SetCandidatePluginInfo{
			{
				PluginType: types.PluginTypeCCIPExec,
				OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
					c.NewChain.Selector: c.NewChain.OCRParams,
				},
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SetCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// 10. Promote the candidates for the commit and exec plugins
	out, err = PromoteCandidateChangeset(e, PromoteCandidateChangesetConfig{
		HomeChainSelector: c.HomeChainSelector,
		MCMS:              mcmsConfig,
		PluginInfo: []PromoteCandidatePluginInfo{
			{
				PluginType:           types.PluginTypeCCIPCommit,
				RemoteChainSelectors: []uint64{c.NewChain.Selector},
				DigestOverrides: map[uint64]DigestAndDonID{
					c.NewChain.Selector: {
						ConfigDigest: commitConfigDigest,
						DonID:        donID,
					},
				},
			},
			{
				PluginType:           types.PluginTypeCCIPExec,
				RemoteChainSelectors: []uint64{c.NewChain.Selector},
				DigestOverrides: map[uint64]DigestAndDonID{
					c.NewChain.Selector: {
						ConfigDigest: execConfigDigest,
						DonID:        donID,
					},
				},
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run PromoteCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// 11. Set the OCR3 config on the off ramp on the new chain
	out, err = SetOCR3OffRampChangeset(e, SetOCR3OffRampConfig{
		HomeChainSel:       c.HomeChainSelector,
		RemoteChainSels:    []uint64{c.NewChain.Selector},
		CCIPHomeConfigType: c.HomeConfigType,
		MCMS:               mcmsConfig,
		ActiveDigestOverrides: map[uint64]map[types.PluginType][32]byte{
			c.NewChain.Selector: map[types.PluginType][32]byte{
				types.PluginTypeCCIPCommit: commitConfigDigest,
				types.PluginTypeCCIPExec:   execConfigDigest,
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SetOCR3OffRampChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 12. Update the fee quoter prices and destinations on the remote chains
	for _, remoteChain := range c.RemoteChains {
		out, err := UpdateFeeQuoterPricesChangeset(e, UpdateFeeQuoterPricesConfig{
			PricesByChain: map[uint64]FeeQuoterPriceUpdatePerSource{
				remoteChain.Selector: FeeQuoterPriceUpdatePerSource{
					TokenPrices: remoteChain.TokenPrices,
					GasPrices:   map[uint64]*big.Int{c.NewChain.Selector: c.NewChain.GasPrice},
				},
			},
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterPricesChangeset on chain with selector %d: %w", remoteChain.Selector, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		_, err = UpdateFeeQuoterDestsChangeset(e, UpdateFeeQuoterDestsConfig{
			UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
				c.NewChain.Selector: map[uint64]fee_quoter.FeeQuoterDestChainConfig{
					remoteChain.Selector: c.NewChain.FeeQuoterDestChainConfig,
				},
			},
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterDestsChangeset on chain with selector %d: %w", remoteChain.Selector, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)
	}

	// 13. Connect the new chain to the existing chains (use the test router)
	testRouter := true
	connections := make(map[uint64]ConnectionConfig, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		connections[remoteChain.Selector] = ConnectionConfig{
			RMNVerificationDisabled: remoteChain.RMNVerificationDisabled,
			AllowListEnabled:        remoteChain.AllowListEnabled,
		}
	}
	cfg := ConnectNewChainConfig{
		RemoteChains:     connections,
		NewChainSelector: c.NewChain.Selector,
		NewChainConnectionConfig: ConnectionConfig{
			RMNVerificationDisabled: c.NewChain.RMNVerificationDisabled,
			AllowListEnabled:        c.NewChain.AllowListEnabled,
		},
		TestRouter: &testRouter,
		MCMSConfig: mcmsConfig,
	}
	err = ConnectNewChainChangeset.VerifyPreconditions(e, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run ConnectNewChainChangeset precondition: %w", err)
	}
	out, err = ConnectNewChainChangeset.Apply(e, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run ConnectNewChainChangeset: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	/*
		TODO:
		Add RMN deployment on the new chain
		- SetRMNRemoteConfigChangeset
		- SetRMNRemoteOnRMNProxyChangeset
	*/

	return deployment.ChangesetOutput{}, nil
}

func runAndSaveAddresses(fn func() (deployment.ChangesetOutput, error), new deployment.AddressBook, existing deployment.AddressBook) error {
	output, err := fn()
	if err != nil {
		return fmt.Errorf("failed to run changeset: %w", err)
	}
	err = new.Merge(output.AddressBook)
	if err != nil {
		return fmt.Errorf("failed to update new address book: %w", err)
	}
	err = existing.Merge(output.AddressBook)
	if err != nil {
		return fmt.Errorf("failed to update existing address book: %w", err)
	}

	return nil
}

// calculateConfigDigest is a utility function that calculates the config digest for a given plugin key, static config, and version.
// This function is used in the CCIPHome contract to calculate the config digest for a plugin:
// https://github.com/smartcontractkit/chainlink/blob/82ab04f1a99e398c9c737f0a67fbd7f407c95b02/contracts/src/v0.8/ccip/capability/CCIPHome.sol#L425
func calculateConfigDigest(donID uint32, pluginType types.PluginType, staticConfig ccip_home.CCIPHomeOCR3Config, version uint32, chainID int64, contractAddress string) ([32]byte, error) {
	// "EVM" as bytes32 (padded to 32 bytes)
	var evmPadded common.Hash
	copy(evmPadded[:], []byte("EVM"))

	// ABI JSON definition for encoding the values
	abiStr := `[{"type":"bytes32"}, {"type":"uint256"}, {"type": "address"}, {"type": "uint32"}, {"type": "uint8"}, {"type": "uint32"}]`

	// Define the values you want to encode
	values := []any{
		evmPadded,
		big.NewInt(chainID),
		common.HexToAddress(contractAddress),
		donID,
		pluginType,
		version,
	}

	// Encode the values using abiEncode
	encodedData, err := abiEncode(abiStr, values...)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to abi encode values: %w", err)
	}

	// ABI JSON definition for encoding the staticConfig
	abiStr = `[
		{
			"type": "tuple",
			"components": [
				{"name": "pluginType", "type": "uint8"},
				{"name": "chainSelector", "type": "uint64"},
        		{"name": "FRoleDON", "type": "uint8"},
        		{"name": "offchainConfigVersion", "type": "uint64"},
        		{"name": "offrampAddress", "type": "bytes"},
        		{"name": "rmnHomeAddress", "type": "bytes"},
				{ 
					"name": "nodes", 
					"type": "tuple[]", 
					"components": [
						{ "name": "p2pId", "type": "bytes32" },
						{ "name": "signerKey", "type": "bytes" },
						{ "name": "transmitterKey", "type": "bytes" }
					]
				},
				{ "name": "offchainConfig", "type": "bytes" }
			]
		}
	]`
	values = []any{staticConfig}
	staticConfigBytes, err := abiEncode(abiStr, values...)

	// Append the staticConfig bytes to the encoded data
	finalData := append(encodedData, staticConfigBytes...)

	// Compute Keccak-256 hash
	hash := crypto.Keccak256Hash(finalData)

	// Convert first 32 bytes of hash to big.Int
	hashInt := new(big.Int).SetBytes(hash[:])

	// Apply PREFIX and PREFIX_MASK
	mask := new(big.Int).Not(PREFIX_MASK) // ~PREFIX_MASK
	hashInt.And(hashInt, mask)            // hashInt & ~PREFIX_MASK
	hashInt.Or(hashInt, PREFIX)           // PREFIX | (hashInt & ~PREFIX_MASK)

	// Convert final result back to [32]byte
	var finalHash [32]byte
	hashBytes := hashInt.Bytes()
	copy(finalHash[32-len(hashBytes):], hashBytes) // Ensure correct padding to 32 bytes

	return finalHash, nil
}

// abiEncode is the equivalent of abi.encode, see a full set of examples:
// https://github.com/ethereum/go-ethereum/blob/420b78659bef661a83c5c442121b13f13288c09f/accounts/abi/packing_test.go#L31
func abiEncode(abiStr string, values ...any) ([]byte, error) {
	// Create a dummy method with arguments
	inDef := fmt.Sprintf(`[{"name": "method", "type": "function", "inputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}

	res, err := inAbi.Pack("method", values...)
	if err != nil {
		return nil, err
	}

	return res[4:], nil
}

// countConfigSetEvents counts the number of ConfigSet events emitted by the CCIPHome contract
func countConfigSetEvents(ccipHome *ccip_home.CCIPHome, opts *bind.FilterOpts) (uint32, error) {
	iter, err := ccipHome.FilterConfigSet(opts, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to filter ConfigSet events: %w", err)
	}

	count := uint32(0)
	for iter.Next() {
		count++
	}

	return count, nil
}
