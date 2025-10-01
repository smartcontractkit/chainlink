package v1_6

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"golang.org/x/exp/maps"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/don_id_claimer"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	_ cldf.ChangeSet[AddDonAndSetCandidateChangesetConfig] = AddDonAndSetCandidateChangeset
	_ cldf.ChangeSet[PromoteCandidateChangesetConfig]      = PromoteCandidateChangeset
	_ cldf.ChangeSet[SetCandidateChangesetConfig]          = SetCandidateChangeset
	_ cldf.ChangeSet[RevokeCandidateChangesetConfig]       = RevokeCandidateChangeset
	_ cldf.ChangeSet[UpdateChainConfigConfig]              = UpdateChainConfigChangeset

	DeployDonIDClaimerChangeset = cldf.CreateChangeSet(deployDonIDClaimerChangesetLogic, deployDonIDClaimerPrecondition)
	DonIDClaimerOffSetChangeset = cldf.CreateChangeSet(donIDClaimerOffSetChangesetLogic, donIDClaimerOffSetChangesetPrecondition)
)

func findTokenInfo(tokens []shared.TokenDetails, address common.Address) (string, uint8, error) {
	for _, token := range tokens {
		if token.Address() == address {
			tokenSymbol, err := token.Symbol(nil)
			if err != nil {
				return "", 0, fmt.Errorf("fetch token symbol for token %s: %w", address, err)
			}
			// TODO think of better solution
			// there are tokens which have diff symbols in testnet and mainnet
			if symbol, ok := shared.TokenSymbolSubstitute[tokenSymbol]; ok {
				tokenSymbol = symbol
			}
			tokenDecimals, err := token.Decimals(nil)
			if err != nil {
				return "", 0, fmt.Errorf("fetch token decimals for token %s: %w", address, err)
			}
			return tokenSymbol, tokenDecimals, nil
		}
	}
	return "", 0, fmt.Errorf("token %s not found in available tokens", address)
}

func validateExecOffchainConfig(e cldf.Environment, c *pluginconfig.ExecuteOffchainConfig, selector uint64, state stateview.CCIPOnChainState) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid execute off-chain config: %w", err)
	}
	// get offRamp
	if err := state.ValidateRamp(selector, shared.OffRamp); err != nil {
		return fmt.Errorf("validate offRamp: %w", err)
	}

	for _, observerConfig := range c.TokenDataObservers {
		switch observerConfig.Type {
		case pluginconfig.USDCCCTPHandlerType:
			if err := validateUSDCConfig(observerConfig.USDCCCTPObserverConfig, state); err != nil {
				return fmt.Errorf("invalid USDC config: %w", err)
			}
		case pluginconfig.LBTCHandlerType:
			if err := validateLBTCConfig(e, observerConfig.LBTCObserverConfig, state); err != nil {
				return fmt.Errorf("invalid LBTC config: %w", err)
			}
		default:
			return fmt.Errorf("unknown token observer config type: %s", observerConfig.Type)
		}
	}
	return nil
}

func validateCommitOffchainConfig(c *pluginconfig.CommitOffchainConfig, selector uint64, feedChainSel uint64, state stateview.CCIPOnChainState) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid commit off-chain config: %w", err)
	}

	family, err := chain_selectors.GetSelectorFamily(selector)
	if err != nil {
		return err
	}
	if family != chain_selectors.FamilyEVM {
		// TODO: implement more proper validation
		return nil
	}

	for tokenAddr, tokenConfig := range c.TokenInfo {
		tokenUnknownAddr, err := ccipocr3.NewUnknownAddressFromHex(string(tokenAddr))
		if err != nil {
			return fmt.Errorf("invalid token address %s: %w", tokenAddr, err)
		}

		aggregatorAddr := common.HexToAddress(string(tokenConfig.AggregatorAddress))
		token := common.HexToAddress(tokenUnknownAddr.String())
		tokenInfos := make([]shared.TokenDetails, 0)
		onchainState := state.Chains[selector]
		for _, tk := range onchainState.BurnMintTokens677 {
			tokenInfos = append(tokenInfos, tk)
		}
		for _, tk := range onchainState.ERC20Tokens {
			tokenInfos = append(tokenInfos, tk)
		}
		for _, tk := range onchainState.ERC677Tokens {
			tokenInfos = append(tokenInfos, tk)
		}
		var linkTokenInfo shared.TokenDetails
		linkTokenInfo = onchainState.LinkToken
		if onchainState.LinkToken == nil {
			linkTokenInfo = onchainState.StaticLinkToken
		}
		tokenInfos = append(tokenInfos, linkTokenInfo)
		tokenInfos = append(tokenInfos, onchainState.Weth9)
		symbol, decimal, err := findTokenInfo(tokenInfos, token)
		if err != nil {
			return fmt.Errorf("chain %d- %w", selector, err)
		}
		if decimal != tokenConfig.Decimals {
			return fmt.Errorf("token %s -address %s has %d decimals in provided token config, expected %d",
				symbol, token.String(), tokenConfig.Decimals, decimal)
		}
		feedChainState := state.Chains[feedChainSel]
		aggregatorInState := feedChainState.USDFeeds[shared.TokenSymbol(symbol)]
		if aggregatorAddr == (common.Address{}) {
			return fmt.Errorf("token %s -address %s has no aggregator in provided token config", symbol, token.String())
		}
		if aggregatorInState == nil {
			return fmt.Errorf("token %s -address %s has no aggregator in state,"+
				" but the aggregator %s is provided in token config", symbol, token.String(), aggregatorAddr.String())
		}
		if aggregatorAddr != aggregatorInState.Address() {
			return fmt.Errorf("token %s -address %s has aggregator %s in provided token config, expected %s",
				symbol, token.String(), aggregatorAddr.String(), aggregatorInState.Address().String())
		}
	}
	return nil
}

func validateUSDCConfig(usdcConfig *pluginconfig.USDCCCTPObserverConfig, state stateview.CCIPOnChainState) error {
	for sel, token := range usdcConfig.Tokens {
		family, err := chain_selectors.GetSelectorFamily(uint64(sel))
		if err != nil {
			return fmt.Errorf("failed to find family for selector %d: %w", sel, err)
		}
		switch family {
		case chain_selectors.FamilyEVM:
			onchainState, ok := state.Chains[uint64(sel)]
			if !ok {
				return fmt.Errorf("chain %d does not exist in EVM chain state but provided in USDCCCTPObserverConfig", sel)
			}
			if onchainState.USDCTokenPools == nil {
				return fmt.Errorf("chain %d does not have any USDC token pools deployed", sel)
			}

			var sourcePoolAddress common.Address
			if pool, ok := onchainState.USDCTokenPoolsV1_6[deployment.Version1_6_2]; ok {
				sourcePoolAddress = pool.Address()
			} else if pool, ok := onchainState.USDCTokenPools[deployment.Version1_5_1]; ok {
				sourcePoolAddress = pool.Address()
			} else {
				return fmt.Errorf("chain %d does not have USDC token pool deployed with version %s or %s", sel, deployment.Version1_5_1, deployment.Version1_6_2)
			}

			if common.HexToAddress(token.SourcePoolAddress) != sourcePoolAddress {
				return fmt.Errorf("chain %d has latest USDC token pool deployed at %s, "+
					"but SourcePoolAddress %s is provided in USDCCCTPObserverConfig",
					sel, onchainState.USDCTokenPools[deployment.Version1_5_1].Address().String(), token.SourcePoolAddress)
			}
		case chain_selectors.FamilySolana:
			onchainState, ok := state.SolChains[uint64(sel)]
			if !ok {
				return fmt.Errorf("chain %d does not exist in Solana chain state but provided in USDCCCTPObserverConfig", sel)
			}
			if onchainState.CCTPTokenPool.IsZero() {
				return fmt.Errorf("chain %d does not have a CCTP token pool deployed", sel)
			}
			if onchainState.USDCToken.IsZero() {
				return fmt.Errorf("chain %d does not have a USDC token in state", sel)
			}
			// Calculate the token pool config address for the USDCTokenPool since Solana expects the source pool to be this PDA instead of the actual pool
			tokenPoolConfig, err := tokens.TokenPoolConfigAddress(onchainState.USDCToken, onchainState.CCTPTokenPool)
			if err != nil {
				return fmt.Errorf("failed to calculate token pool config address: %w", err)
			}

			// Token source pool address is stored as a hex representation of the address bytes
			// This string is NOT a base58 address
			sourcePoolAddrByteStr := token.SourcePoolAddress
			if strings.HasPrefix(sourcePoolAddrByteStr, "0x") || strings.HasPrefix(sourcePoolAddrByteStr, "0X") {
				sourcePoolAddrByteStr = sourcePoolAddrByteStr[2:]
			}
			sourcePoolBytes, err := hex.DecodeString(sourcePoolAddrByteStr)
			if err != nil {
				return fmt.Errorf("failed to decode source pool bytes hex: %w", err)
			}

			sourcePoolPubkey := solana.PublicKeyFromBytes(sourcePoolBytes)
			if !sourcePoolPubkey.Equals(tokenPoolConfig) {
				return fmt.Errorf("solana chain %d has CCTP token pool config deployed at %s, "+
					"but SourcePoolAddress provided in USDCCCTPObserverConfig decodes to %s",
					sel, tokenPoolConfig.String(), sourcePoolPubkey.String())
			}
		default:
			return fmt.Errorf("USDC configs not supported for chain family %s", family)
		}
	}
	return nil
}

func validateLBTCConfig(e cldf.Environment, lbtcConfig *pluginconfig.LBTCObserverConfig, state stateview.CCIPOnChainState) error {
	for sel, sourcePool := range lbtcConfig.SourcePoolAddressByChain {
		_, ok := state.Chains[uint64(sel)]
		if !ok {
			return fmt.Errorf("chain %d does not exist in state but provided in LBTCObserverConfig", sel)
		}
		sourcePoolAddr := common.HexToAddress(sourcePool)
		sourcePool, err := token_pool.NewTokenPool(sourcePoolAddr, e.BlockChains.EVMChains()[uint64(sel)].Client)
		if err != nil {
			return fmt.Errorf("chain %d has an error while requesting LBTC source token pool %s: %w", sel, sourcePoolAddr, err)
		}
		_, err = sourcePool.GetToken(nil)
		if err != nil {
			return fmt.Errorf("chain %d has an error while requesting LBTC token address: %w", sel, err)
		}
	}
	return nil
}

type CCIPOCRParams struct {
	// OCRParameters contains the parameters for the OCR plugin.
	OCRParameters commontypes.OCRParameters `json:"ocrParameters"`
	// CommitOffChainConfig contains pointers to Arb feeds for prices.
	CommitOffChainConfig *pluginconfig.CommitOffchainConfig `json:"commitOffChainConfig,omitempty"`
	// ExecuteOffChainConfig contains USDC config.
	ExecuteOffChainConfig *pluginconfig.ExecuteOffchainConfig `json:"executeOffChainConfig,omitempty"`
}

func (c CCIPOCRParams) Copy() CCIPOCRParams {
	newC := CCIPOCRParams{
		OCRParameters: c.OCRParameters,
	}
	if c.CommitOffChainConfig != nil {
		commit := *c.CommitOffChainConfig
		newC.CommitOffChainConfig = &commit
	}
	if c.ExecuteOffChainConfig != nil {
		exec := *c.ExecuteOffChainConfig
		newC.ExecuteOffChainConfig = &exec
	}
	return newC
}

func (c CCIPOCRParams) Validate(e cldf.Environment, selector uint64, feedChainSel uint64, state stateview.CCIPOnChainState) error {
	if err := c.OCRParameters.Validate(); err != nil {
		return fmt.Errorf("invalid OCR parameters: %w", err)
	}
	if c.CommitOffChainConfig == nil && c.ExecuteOffChainConfig == nil {
		return errors.New("at least one of CommitOffChainConfig or ExecuteOffChainConfig must be set")
	}
	if c.CommitOffChainConfig != nil {
		if err := validateCommitOffchainConfig(c.CommitOffChainConfig, selector, feedChainSel, state); err != nil {
			return fmt.Errorf("invalid commit off-chain config: %w", err)
		}
	}
	if c.ExecuteOffChainConfig != nil {
		if err := validateExecOffchainConfig(e, c.ExecuteOffChainConfig, selector, state); err != nil {
			return fmt.Errorf("invalid execute off-chain config: %w", err)
		}
	}
	return nil
}

type PromoteCandidatePluginInfo struct {
	// RemoteChainSelectors is the chain selector of the DONs that we want to promote the candidate config of.
	// Note that each (chain, ccip capability version) pair has a unique DON ID.
	RemoteChainSelectors    []uint64         `json:"remoteChainSelectors"`
	PluginType              types.PluginType `json:"pluginType"`
	AllowEmptyConfigPromote bool             `json:"allowEmptyConfigPromote"` // safe guard to prevent promoting empty config to active
}

type PromoteCandidateChangesetConfig struct {
	HomeChainSelector uint64 `json:"homeChainSelector"`

	PluginInfo []PromoteCandidatePluginInfo `json:"pluginInfo"`
	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *proposalutils.TimelockConfig `json:"mcms,omitempty"`
}

func (p PromoteCandidateChangesetConfig) Validate(e cldf.Environment) (map[uint64]uint32, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	if err := stateview.ValidateChain(e, state, p.HomeChainSelector, p.MCMS); err != nil {
		return nil, fmt.Errorf("home chain invalid: %w", err)
	}
	homeChainState := state.Chains[p.HomeChainSelector]
	if err := commoncs.ValidateOwnership(e.GetContext(), p.MCMS != nil, e.BlockChains.EVMChains()[p.HomeChainSelector].DeployerKey.From, homeChainState.Timelock.Address(), homeChainState.CapabilityRegistry); err != nil {
		return nil, err
	}

	donIDs := make(map[uint64]uint32)
	for _, plugin := range p.PluginInfo {
		if plugin.PluginType != types.PluginTypeCCIPCommit &&
			plugin.PluginType != types.PluginTypeCCIPExec {
			return nil, errors.New("PluginType must be set to either CCIPCommit or CCIPExec")
		}
		for _, chainSelector := range plugin.RemoteChainSelectors {
			if err := cldf.IsValidChainSelector(chainSelector); err != nil {
				return nil, fmt.Errorf("don chain selector invalid: %w", err)
			}
			if err := state.ValidateRamp(chainSelector, shared.OffRamp); err != nil {
				return nil, err
			}

			donID, err := internal.DonIDForChain(
				state.Chains[p.HomeChainSelector].CapabilityRegistry,
				state.Chains[p.HomeChainSelector].CCIPHome,
				chainSelector,
			)
			if err != nil {
				return nil, fmt.Errorf("fetch don id for chain: %w", err)
			}
			if donID == 0 {
				return nil, fmt.Errorf("don doesn't exist in CR for chain %d", chainSelector)
			}
			// Check that candidate digest and active digest are not both zero - this is enforced onchain.
			pluginConfigs, err := state.Chains[p.HomeChainSelector].CCIPHome.GetAllConfigs(&bind.CallOpts{
				Context: e.GetContext(),
			}, donID, uint8(plugin.PluginType))
			if err != nil {
				return nil, fmt.Errorf("fetching %s configs from cciphome: %w", plugin.PluginType.String(), err)
			}
			// If promoteCandidate is called with AllowEmptyConfigPromote set to false and
			// the CandidateConfig config digest is zero, do not promote the candidate config to active.
			if !plugin.AllowEmptyConfigPromote && pluginConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
				return nil, fmt.Errorf("%s candidate config digest is empty", plugin.PluginType.String())
			}

			// If the active and candidate config digests are both zero, we should not promote the candidate config to active.
			if pluginConfigs.ActiveConfig.ConfigDigest == [32]byte{} &&
				pluginConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
				return nil, fmt.Errorf("%s active and candidate config digests are both zero", plugin.PluginType.String())
			}
			donIDs[chainSelector] = donID
		}
	}
	if len(e.NodeIDs) == 0 {
		return nil, errors.New("NodeIDs must be set")
	}
	if state.Chains[p.HomeChainSelector].CCIPHome == nil {
		return nil, errors.New("CCIPHome contract does not exist")
	}
	if state.Chains[p.HomeChainSelector].CapabilityRegistry == nil {
		return nil, errors.New("CapabilityRegistry contract does not exist")
	}

	return donIDs, nil
}

// PromoteCandidateChangeset generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// Note that a DON must exist prior to being able to use this changeset effectively,
// i.e AddDonAndSetCandidateChangeset must be called first.
// This can also be used to promote a 0x0 candidate config to be the active, effectively shutting down the DON.
// At that point you can call the RemoveDON changeset to remove the DON entirely from the capability registry.
// PromoteCandidateChangeset is NOT idempotent, once candidate config is promoted to active, if it's called again,
// It might promote empty candidate config to active, which is not desired.
func PromoteCandidateChangeset(
	e cldf.Environment,
	cfg PromoteCandidateChangesetConfig,
) (cldf.ChangesetOutput, error) {
	donIDs, err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("%w: %w", cldf.ErrInvalidConfig, err)
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("fetch node info: %w", err)
	}
	nonBootstraps := nodes.NonBootstraps()

	dons := make([]ccipseqs.DONUpdatePromotion, 0)
	for _, plugin := range cfg.PluginInfo {
		for _, donID := range donIDs {
			digest, err := state.Chains[cfg.HomeChainSelector].CCIPHome.GetCandidateDigest(nil, donID, uint8(plugin.PluginType))
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			if digest == [32]byte{} && !plugin.AllowEmptyConfigPromote {
				return cldf.ChangesetOutput{}, errors.New("candidate config digest is zero, promoting empty config is not allowed")
			}
			allConfigs, err := state.Chains[cfg.HomeChainSelector].CCIPHome.GetAllConfigs(nil, donID, uint8(plugin.PluginType))
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			e.Logger.Infow("Promoting candidate for plugin "+plugin.PluginType.String(), "digest", digest)

			dons = append(dons, ccipseqs.DONUpdatePromotion{
				ID:              donID,
				PluginType:      uint8(plugin.PluginType),
				ChainSelector:   allConfigs.CandidateConfig.Config.ChainSelector,
				PeerIDs:         nonBootstraps.PeerIDs(),
				F:               nonBootstraps.DefaultF(),
				IsPublic:        false,
				CandidateDigest: allConfigs.CandidateConfig.ConfigDigest,
				ActiveDigest:    allConfigs.ActiveConfig.ConfigDigest,
			})
		}
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.PromoteCandidateSequence,
		ccipseqs.DONSequenceDeps{
			HomeChain: e.BlockChains.EVMChains()[cfg.HomeChainSelector],
		},
		ccipseqs.PromoteCandidateSequenceInput{
			CapabilitiesRegistry: state.Chains[cfg.HomeChainSelector].CapabilityRegistry.Address(),
			NoSend:               cfg.MCMS != nil,
			DONs:                 dons,
		},
	)
	return opsutil.AddEVMCallSequenceToCSOutput(
		e,
		cldf.ChangesetOutput{},
		report,
		err,
		state.EVMMCMSStateByChain(),
		cfg.MCMS,
		"PromoteCandidateChangeset",
	)
}

type SetCandidatePluginInfo struct {
	// OCRConfigPerRemoteChainSelector is the chain selector of the chain where the DON will be added.
	OCRConfigPerRemoteChainSelector map[uint64]CCIPOCRParams `json:"ocrConfigPerRemoteChainSelector"`
	PluginType                      types.PluginType         `json:"pluginType"`

	// SkipChainConfigValidation skips validation of the config for chain on CCIPHome.
	// WARNING: Never enable this parameter if running this changeset in isolation.
	// This is only meant to be enabled when running this changeset as part of a larger changeset that groups multiple proposals together.
	SkipChainConfigValidation bool `json:"skipChainConfigValidation"`
}

func (p SetCandidatePluginInfo) String() string {
	allchains := maps.Keys(p.OCRConfigPerRemoteChainSelector)
	return fmt.Sprintf("PluginType: %s, Chains: %v", p.PluginType.String(), allchains)
}

func (p SetCandidatePluginInfo) Validate(e cldf.Environment, state stateview.CCIPOnChainState, homeChain uint64, feedChain uint64) error {
	if p.PluginType != types.PluginTypeCCIPCommit &&
		p.PluginType != types.PluginTypeCCIPExec {
		return errors.New("PluginType must be set to either CCIPCommit or CCIPExec")
	}
	for chainSelector, params := range p.OCRConfigPerRemoteChainSelector {
		if _, exists := state.SupportedChains()[chainSelector]; !exists {
			return fmt.Errorf("chain %d does not exist in state", chainSelector)
		}
		if err := cldf.IsValidChainSelector(chainSelector); err != nil {
			return fmt.Errorf("don chain selector invalid: %w", err)
		}
		if err := state.ValidateRamp(chainSelector, shared.OffRamp); err != nil {
			return err
		}
		if p.PluginType == types.PluginTypeCCIPCommit && params.CommitOffChainConfig == nil {
			return errors.New("commit off-chain config must be set")
		}
		if p.PluginType == types.PluginTypeCCIPExec && params.ExecuteOffChainConfig == nil {
			return errors.New("execute off-chain config must be set")
		}

		if !p.SkipChainConfigValidation {
			chainConfig, err := state.Chains[homeChain].CCIPHome.GetChainConfig(nil, chainSelector)
			if err != nil {
				return fmt.Errorf("can't get chain config for %d: %w", chainSelector, err)
			}
			// FChain should never be zero if a chain config is set in CCIPHome
			if chainConfig.FChain == 0 {
				return fmt.Errorf("chain config not set up for new chain %d", chainSelector)
			}
			if len(chainConfig.Readers) == 0 {
				return errors.New("readers must be set")
			}
			decodedChainConfig, err := chainconfig.DecodeChainConfig(chainConfig.Config)
			if err != nil {
				return fmt.Errorf("can't decode chain config: %w", err)
			}
			if err := decodedChainConfig.Validate(); err != nil {
				return fmt.Errorf("invalid chain config: %w", err)
			}
		}

		err := params.Validate(e, chainSelector, feedChain, state)
		if err != nil {
			return fmt.Errorf("invalid ccip ocr params: %w", err)
		}
	}
	return nil
}

// SetCandidateConfigBase is a common base config struct for AddDonAndSetCandidateChangesetConfig and SetCandidateChangesetConfig.
// This is extracted to deduplicate most of the validation logic.
// Remaining validation logic is done in the specific config structs that inherit from this.
type SetCandidateConfigBase struct {
	HomeChainSelector uint64 `json:"homeChainSelector"`
	FeedChainSelector uint64 `json:"feedChainSelector"`

	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *proposalutils.TimelockConfig `json:"mcms,omitempty"`
}

func (s SetCandidateConfigBase) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	if err := cldf.IsValidChainSelector(s.HomeChainSelector); err != nil {
		return fmt.Errorf("home chain selector invalid: %w", err)
	}
	if err := cldf.IsValidChainSelector(s.FeedChainSelector); err != nil {
		return fmt.Errorf("feed chain selector invalid: %w", err)
	}
	homeChainState, exists := state.Chains[s.HomeChainSelector]
	if !exists {
		return fmt.Errorf("home chain %d does not exist", s.HomeChainSelector)
	}
	if err := commoncs.ValidateOwnership(e.GetContext(), s.MCMS != nil, e.BlockChains.EVMChains()[s.HomeChainSelector].DeployerKey.From, homeChainState.Timelock.Address(), homeChainState.CapabilityRegistry); err != nil {
		return err
	}

	if len(e.NodeIDs) == 0 {
		return errors.New("nodeIDs must be set")
	}
	if state.Chains[s.HomeChainSelector].CCIPHome == nil {
		return errors.New("CCIPHome contract does not exist")
	}
	if state.Chains[s.HomeChainSelector].CapabilityRegistry == nil {
		return errors.New("CapabilityRegistry contract does not exist")
	}

	if e.OCRSecrets.IsEmpty() {
		return errors.New("OCR secrets must be set")
	}

	return nil
}

// AddDonAndSetCandidateChangesetConfig is a separate config struct
// because the validation is slightly different from SetCandidateChangesetConfig.
// In particular, we check to make sure we don't already have a DON for the chain.
type AddDonAndSetCandidateChangesetConfig struct {
	SetCandidateConfigBase `json:"setCandidateConfigBase"`

	// Only set one plugin at a time while you are adding the DON for the first time.
	// For subsequent SetCandidate call use SetCandidateChangeset as that fetches the already added DONID and sets the candidate.
	PluginInfo SetCandidatePluginInfo `json:"pluginInfo"`

	// WARNING: Do not use if calling this changeset in isolation
	DonIDOverride uint32 `json:"donIdOverride"`
}

func (a AddDonAndSetCandidateChangesetConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	if err := a.SetCandidateConfigBase.Validate(e, state); err != nil {
		return err
	}

	if err := a.PluginInfo.Validate(e, state, a.HomeChainSelector, a.FeedChainSelector); err != nil {
		return fmt.Errorf("validate plugin info %s: %w", a.PluginInfo.String(), err)
	}
	for chainSelector := range a.PluginInfo.OCRConfigPerRemoteChainSelector {
		// check if a DON already exists for this chain
		donID, err := internal.DonIDForChain(
			state.Chains[a.HomeChainSelector].CapabilityRegistry,
			state.Chains[a.HomeChainSelector].CCIPHome,
			chainSelector,
		)
		if err != nil {
			return fmt.Errorf("fetch don id for chain: %w", err)
		}
		// if don already exists use SetCandidateChangeset instead
		if donID != 0 {
			return fmt.Errorf("don already exists in CR for chain %d, it has id %d", chainSelector, donID)
		}
	}

	return nil
}

// AddDonAndSetCandidateChangeset adds new DON for destination to home chain
// and sets the plugin config as candidateConfig for the don.
//
// This is the first step to creating a CCIP DON and must be executed before any
// other changesets (SetCandidateChangeset, PromoteCandidateChangeset)
// can be executed.
//
// Note that these operations must be done together because the createDON call
// in the capability registry calls the capability config contract, so we must
// provide suitable calldata for CCIPHome.
// AddDonAndSetCandidateChangeset is not idempotent, if AddDON is called more than once for the same chain,
// it will throw an error because the DON would already exist for that chain.
func AddDonAndSetCandidateChangeset(
	e cldf.Environment,
	cfg AddDonAndSetCandidateChangesetConfig,
) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	err = cfg.Validate(e, state)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("%w: %w", cldf.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("get node info: %w", err)
	}
	nonBootstraps := nodes.NonBootstraps()

	var expectedDonID uint32
	if cfg.DonIDOverride != 0 {
		expectedDonID = cfg.DonIDOverride
	} else {
		expectedDonID, err = state.Chains[cfg.HomeChainSelector].CapabilityRegistry.GetNextDONId(&bind.CallOpts{
			Context: e.GetContext(),
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("get next don id: %w", err)
		}
	}

	dons := make([]ccipseqs.DONAddition, len(cfg.PluginInfo.OCRConfigPerRemoteChainSelector))
	i := 0
	for chainSelector, params := range cfg.PluginInfo.OCRConfigPerRemoteChainSelector {
		offRampAddress, err := state.GetOffRampAddressBytes(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
			state.Chains[cfg.HomeChainSelector].CCIPHome,
			e.OCRSecrets,
			offRampAddress,
			chainSelector,
			nonBootstraps,
			state.Chains[cfg.HomeChainSelector].RMNHome.Address(),
			params.OCRParameters,
			params.CommitOffChainConfig,
			params.ExecuteOffChainConfig,
			cfg.PluginInfo.SkipChainConfigValidation,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		pluginOCR3Config, ok := newDONArgs[cfg.PluginInfo.PluginType]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("missing plugin %s in ocr3Configs",
				cfg.PluginInfo.PluginType.String())
		}

		dons[i] = ccipseqs.DONAddition{
			ExpectedID:       expectedDonID,
			PluginConfig:     pluginOCR3Config,
			PeerIDs:          nonBootstraps.PeerIDs(),
			F:                nonBootstraps.DefaultF(),
			IsPublic:         false,
			AcceptsWorkflows: false,
		}
		i++
		expectedDonID++
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.AddDONAndSetCandidateSequence,
		ccipseqs.DONSequenceDeps{
			HomeChain: e.BlockChains.EVMChains()[cfg.HomeChainSelector],
		},
		ccipseqs.AddDONAndSetCandidateSequenceInput{
			CapabilitiesRegistry: state.Chains[cfg.HomeChainSelector].CapabilityRegistry.Address(),
			NoSend:               cfg.MCMS != nil,
			DONs:                 dons,
		},
	)
	return opsutil.AddEVMCallSequenceToCSOutput(e, cldf.ChangesetOutput{}, report, err, state.EVMMCMSStateByChain(), cfg.MCMS, fmt.Sprintf("addDON and setCandidates for %s plugin on new chains", cfg.PluginInfo.PluginType.String()))
}

type SetCandidateChangesetConfig struct {
	SetCandidateConfigBase `json:"setCandidateConfigBase"`

	PluginInfo []SetCandidatePluginInfo `json:"pluginInfo"`

	// WARNING: Do not use if calling this changeset in isolation
	DonIDOverrides map[uint64]uint32 `json:"donIdOverrides"`
}

func (s SetCandidateChangesetConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) (map[uint64]uint32, error) {
	err := s.SetCandidateConfigBase.Validate(e, state)
	if err != nil {
		return nil, err
	}

	chainToDonIDs := make(map[uint64]uint32)
	for _, plugin := range s.PluginInfo {
		if err := plugin.Validate(e, state, s.HomeChainSelector, s.FeedChainSelector); err != nil {
			return nil, fmt.Errorf("validate plugin info %s: %w", plugin.String(), err)
		}
		for chainSelector := range plugin.OCRConfigPerRemoteChainSelector {
			if donIDOverride, ok := s.DonIDOverrides[chainSelector]; ok {
				chainToDonIDs[chainSelector] = donIDOverride
				continue
			}
			donID, err := internal.DonIDForChain(
				state.Chains[s.HomeChainSelector].CapabilityRegistry,
				state.Chains[s.HomeChainSelector].CCIPHome,
				chainSelector,
			)
			if err != nil {
				return nil, fmt.Errorf("fetch don id for chain: %w", err)
			}
			// if don doesn't exist use AddDonAndSetCandidateChangeset instead
			if donID == 0 {
				return nil, fmt.Errorf("don doesn't exist in CR for chain %d", chainSelector)
			}
			chainToDonIDs[chainSelector] = donID
		}
	}
	return chainToDonIDs, nil
}

// SetCandidateChangeset generates a proposal to call setCandidate on the CCIPHome through the capability registry.
// A DON must exist in order to use this changeset effectively, i.e AddDonAndSetCandidateChangeset must be called first.
func SetCandidateChangeset(
	e cldf.Environment,
	cfg SetCandidateChangesetConfig,
) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainToDonIDs, err := cfg.Validate(e, state)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("%w: %w", cldf.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("get node info: %w", err)
	}
	nonBootstraps := nodes.NonBootstraps()

	pluginInfos := make([]string, 0)
	dons := make([]ccipseqs.DONUpdate, 0)
	for _, plugin := range cfg.PluginInfo {
		pluginInfos = append(pluginInfos, plugin.String())
		for chainSelector, params := range plugin.OCRConfigPerRemoteChainSelector {
			offRampAddress, err := state.GetOffRampAddressBytes(chainSelector)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
				state.Chains[cfg.HomeChainSelector].CCIPHome,
				e.OCRSecrets,
				offRampAddress,
				chainSelector,
				nodes.NonBootstraps(),
				state.Chains[cfg.HomeChainSelector].RMNHome.Address(),
				params.OCRParameters,
				params.CommitOffChainConfig,
				params.ExecuteOffChainConfig,
				plugin.SkipChainConfigValidation,
			)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			config, ok := newDONArgs[plugin.PluginType]
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("missing %s plugin in ocr3Configs", plugin.PluginType.String())
			}

			// get the current candidate config
			donID := chainToDonIDs[chainSelector]
			existingDigest, err := state.Chains[cfg.HomeChainSelector].CCIPHome.GetCandidateDigest(&bind.CallOpts{
				Context: e.GetContext(),
			}, donID, config.PluginType)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("get candidate digest from ccipHome: %w", err)
			}
			if existingDigest != [32]byte{} {
				e.Logger.Warnw("Overwriting existing candidate config", "digest", existingDigest, "donID", donID, "pluginType", config.PluginType)
			}

			dons = append(dons, ccipseqs.DONUpdate{
				ID:             donID,
				PluginConfig:   config,
				PeerIDs:        nonBootstraps.PeerIDs(),
				F:              nonBootstraps.DefaultF(),
				IsPublic:       false,
				ExistingDigest: existingDigest,
			})
		}
	}
	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.SetCandidateSequence,
		ccipseqs.DONSequenceDeps{
			HomeChain: e.BlockChains.EVMChains()[cfg.HomeChainSelector],
		},
		ccipseqs.SetCandidateSequenceInput{
			CapabilitiesRegistry: state.Chains[cfg.HomeChainSelector].CapabilityRegistry.Address(),
			NoSend:               cfg.MCMS != nil,
			DONs:                 dons,
		},
	)
	return opsutil.AddEVMCallSequenceToCSOutput(e, cldf.ChangesetOutput{}, report, err, state.EVMMCMSStateByChain(), cfg.MCMS, fmt.Sprintf("setCandidates for plugins: %v", pluginInfos))
}

type RevokeCandidateChangesetConfig struct {
	HomeChainSelector uint64 `json:"homeChainSelector"`

	// RemoteChainSelector is the chain selector whose candidate config we want to revoke.
	RemoteChainSelector uint64           `json:"remoteChainSelector"`
	PluginType          types.PluginType `json:"pluginType"`

	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *proposalutils.TimelockConfig `json:"mcms,omitempty"`
}

func (r RevokeCandidateChangesetConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) (donID uint32, err error) {
	if err := cldf.IsValidChainSelector(r.HomeChainSelector); err != nil {
		return 0, fmt.Errorf("home chain selector invalid: %w", err)
	}
	if err := cldf.IsValidChainSelector(r.RemoteChainSelector); err != nil {
		return 0, fmt.Errorf("don chain selector invalid: %w", err)
	}
	if len(e.NodeIDs) == 0 {
		return 0, errors.New("NodeIDs must be set")
	}
	if state.Chains[r.HomeChainSelector].CCIPHome == nil {
		return 0, errors.New("CCIPHome contract does not exist")
	}
	if state.Chains[r.HomeChainSelector].CapabilityRegistry == nil {
		return 0, errors.New("CapabilityRegistry contract does not exist")
	}
	homeChainState, exists := state.Chains[r.HomeChainSelector]
	if !exists {
		return 0, fmt.Errorf("home chain %d does not exist", r.HomeChainSelector)
	}
	if err := commoncs.ValidateOwnership(e.GetContext(), r.MCMS != nil, e.BlockChains.EVMChains()[r.HomeChainSelector].DeployerKey.From, homeChainState.Timelock.Address(), homeChainState.CapabilityRegistry); err != nil {
		return 0, err
	}

	// check that the don exists for this chain
	donID, err = internal.DonIDForChain(
		state.Chains[r.HomeChainSelector].CapabilityRegistry,
		state.Chains[r.HomeChainSelector].CCIPHome,
		r.RemoteChainSelector,
	)
	if err != nil {
		return 0, fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return 0, fmt.Errorf("don doesn't exist in CR for chain %d", r.RemoteChainSelector)
	}

	// check that candidate digest is not zero - this is enforced onchain.
	candidateDigest, err := state.Chains[r.HomeChainSelector].CCIPHome.GetCandidateDigest(nil, donID, uint8(r.PluginType))
	if err != nil {
		return 0, fmt.Errorf("fetching candidate digest from cciphome: %w", err)
	}
	if candidateDigest == [32]byte{} {
		return 0, errors.New("candidate config digest is zero, can't revoke it")
	}

	return donID, nil
}

func RevokeCandidateChangeset(e cldf.Environment, cfg RevokeCandidateChangesetConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	donID, err := cfg.Validate(e, state)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("%w: %w", cldf.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("fetch nodes info: %w", err)
	}

	txOpts := e.BlockChains.EVMChains()[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = cldf.SimTransactOpts()
	}

	homeChain := e.BlockChains.EVMChains()[cfg.HomeChainSelector]
	ops, err := revokeCandidateOps(
		txOpts,
		homeChain,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		state.Chains[cfg.HomeChainSelector].CCIPHome,
		nodes.NonBootstraps(),
		donID,
		uint8(cfg.PluginType),
		cfg.MCMS != nil,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("revoke candidate ops: %w", err)
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}

	timelocks := map[uint64]string{cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address().Hex()}
	inspectors := map[uint64]mcmssdk.Inspector{cfg.HomeChainSelector: mcmsevmsdk.NewInspector(e.BlockChains.EVMChains()[cfg.HomeChainSelector].Client)}
	batches := []mcmstypes.BatchOperation{{ChainSelector: mcmstypes.ChainSelector(cfg.HomeChainSelector), Transactions: ops}}

	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMS, nil)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build mcm addresses per chain: %w", err)
	}
	prop, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		batches,
		fmt.Sprintf("revokeCandidate for don %d", cfg.RemoteChainSelector),
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*prop}}, nil
}

func revokeCandidateOps(
	txOpts *bind.TransactOpts,
	homeChain cldf_evm.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	nodes deployment.Nodes,
	donID uint32,
	pluginType uint8,
	mcmsEnabled bool,
) ([]mcmstypes.Transaction, error) {
	if donID == 0 {
		return nil, errors.New("donID is zero")
	}

	candidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, pluginType)
	if err != nil {
		return nil, fmt.Errorf("fetching candidate digest from cciphome: %w", err)
	}

	encodedRevokeCandidateCall, err := internal.CCIPHomeABI.Pack(
		"revokeCandidate",
		donID,
		pluginType,
		candidateDigest,
	)
	if err != nil {
		return nil, fmt.Errorf("pack set candidate call: %w", err)
	}

	updateDonTx, err := capReg.UpdateDON(
		txOpts,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: shared.CCIPCapabilityID,
				Config:       encodedRevokeCandidateCall,
			},
		},
		false, // isPublic
		nodes.DefaultF(),
	)

	// note: error check is handled below
	if !mcmsEnabled {
		_, err = cldf.ConfirmIfNoErrorWithABI(
			homeChain, updateDonTx,
			capabilities_registry.CapabilitiesRegistryABI, err)
		if err != nil {
			return nil, fmt.Errorf("error confirming UpdateDON call in revoke candidate (don: %d; ptype: %s): %w",
				donID, types.PluginType(pluginType).String(), err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to call UpdateDON in revoke candidate (don: %d; ptype: %s): %w",
			donID, types.PluginType(pluginType).String(), err)
	}

	tx, err := proposalutils.TransactionForChain(homeChain.Selector, capReg.Address().Hex(), updateDonTx.Data(),
		big.NewInt(0), string(shared.CapabilitiesRegistry), []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateDON mcms tx in revoke candidate (don: %d; ptype: %s): %w",
			donID, types.PluginType(pluginType).String(), err)
	}

	return []mcmstypes.Transaction{tx}, nil
}

type ChainConfig struct {
	Readers              [][32]byte              `json:"readers"`
	FChain               uint8                   `json:"fChain"`
	EncodableChainConfig chainconfig.ChainConfig `json:"encodableChainConfig"`
}

type UpdateChainConfigConfig struct {
	HomeChainSelector  uint64                        `json:"homeChainSelector"`
	RemoteChainRemoves []uint64                      `json:"remoteChainRemoves"`
	RemoteChainAdds    map[uint64]ChainConfig        `json:"remoteChainAdds"`
	MCMS               *proposalutils.TimelockConfig `json:"mcms,omitempty"`
}

func (c UpdateChainConfigConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	if err := cldf.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("home chain selector invalid: %w", err)
	}
	if len(c.RemoteChainRemoves) == 0 && len(c.RemoteChainAdds) == 0 {
		return errors.New("no chain adds or removes")
	}
	homeChainState, exists := state.Chains[c.HomeChainSelector]
	if !exists {
		return fmt.Errorf("home chain %d does not exist", c.HomeChainSelector)
	}
	if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, e.BlockChains.EVMChains()[c.HomeChainSelector].DeployerKey.From, homeChainState.Timelock.Address(), homeChainState.CCIPHome); err != nil {
		return err
	}
	for _, remove := range c.RemoteChainRemoves {
		if err := cldf.IsValidChainSelector(remove); err != nil {
			return fmt.Errorf("chain remove selector invalid: %w", err)
		}
		if _, ok := state.SupportedChains()[remove]; !ok {
			return fmt.Errorf("chain to remove %d is not supported", remove)
		}
	}
	for add, ccfg := range c.RemoteChainAdds {
		if err := cldf.IsValidChainSelector(add); err != nil {
			return fmt.Errorf("chain remove selector invalid: %w", err)
		}
		if _, ok := state.SupportedChains()[add]; !ok {
			return fmt.Errorf("chain to add %d is not supported", add)
		}
		if ccfg.FChain == 0 {
			return fmt.Errorf("fChain must be set for selector %d", add)
		}
		if len(ccfg.Readers) == 0 {
			return fmt.Errorf("readers must be set for selector %d", add)
		}
		if err := ccfg.EncodableChainConfig.Validate(); err != nil {
			return fmt.Errorf("invalid chain config for selector %d: %w", add, err)
		}
	}
	return nil
}

func UpdateChainConfigChangeset(e cldf.Environment, cfg UpdateChainConfigConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("%w: %w", cldf.ErrInvalidConfig, err)
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Create mapping of all removals to check if we are removing and re-adding chains
	removes := make(map[uint64]struct{}, len(cfg.RemoteChainRemoves))
	for _, chain := range cfg.RemoteChainRemoves {
		removes[chain] = struct{}{}
	}

	var adds []ccip_home.CCIPHomeChainConfigArgs
	for chain, ccfg := range cfg.RemoteChainAdds {
		encodedChainConfig, err := chainconfig.EncodeChainConfig(ccfg.EncodableChainConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("encoding chain config: %w", err)
		}
		chainConfig := ccip_home.CCIPHomeChainConfig{
			Readers: ccfg.Readers,
			FChain:  ccfg.FChain,
			Config:  encodedChainConfig,
		}
		existingCfg, err := state.Chains[cfg.HomeChainSelector].CCIPHome.GetChainConfig(nil, chain)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("get chain config for selector %d: %w", chain, err)
		}
		// Don't add chain configs again, unless we are removing and re-adding it.
		if _, ok := removes[chain]; !ok && isChainConfigEqual(existingCfg, chainConfig) {
			e.Logger.Infow("Chain config already exists, not applying again",
				"addedChain", chain,
				"chainConfig", chainConfig,
			)
			continue
		}
		adds = append(adds, ccip_home.CCIPHomeChainConfigArgs{
			ChainSelector: chain,
			ChainConfig:   chainConfig,
		})
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.ApplyChainConfigUpdatesSequence,
		ccipseqs.DONSequenceDeps{
			HomeChain: e.BlockChains.EVMChains()[cfg.HomeChainSelector],
		},
		ccipseqs.ApplyChainConfigUpdatesSequenceInput{
			CCIPHome:           state.Chains[cfg.HomeChainSelector].CCIPHome.Address(),
			NoSend:             cfg.MCMS != nil,
			RemoteChainAdds:    adds,
			RemoteChainRemoves: cfg.RemoteChainRemoves,
			BatchSize:          4, // Conservative batch size to avoid exceeding gas limits (TODO: Make this configurable)
		},
	)
	e.Logger.Infof("Proposed chain config update on chain %d removes %v, adds %+v", cfg.HomeChainSelector, cfg.RemoteChainRemoves, cfg.RemoteChainAdds)
	return opsutil.AddEVMCallSequenceToCSOutput(e, cldf.ChangesetOutput{}, report, err, state.EVMMCMSStateByChain(), cfg.MCMS, "Update chain configs on CCIPHome")
}

func isChainConfigEqual(a, b ccip_home.CCIPHomeChainConfig) bool {
	mapReader := make(map[[32]byte]struct{})
	for i := range a.Readers {
		mapReader[a.Readers[i]] = struct{}{}
	}
	for i := range b.Readers {
		if _, ok := mapReader[b.Readers[i]]; !ok {
			return false
		}
	}
	return bytes.Equal(a.Config, b.Config) &&
		a.FChain == b.FChain
}

type DeployDonIDClaimerConfig struct{}

func deployDonIDClaimerChangesetLogic(e cldf.Environment, _ DeployDonIDClaimerConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return cldf.ChangesetOutput{}, err
	}

	ab := cldf.NewMemoryAddressBook()
	homeChainSel, err := state.HomeChainSelector()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get HomeChainSelector: %w", err)
	}

	chain := e.BlockChains.EVMChains()[homeChainSel]
	err = deployDonIDClaimerContract(e, ab, state, chain)
	if err != nil {
		e.Logger.Errorw("Failed to deploy donIDClaimer contract", "err", err, "addressBook", ab)
		return cldf.ChangesetOutput{
			AddressBook: ab,
		}, fmt.Errorf("failed to deploy donIDClaimer contract: %w", err)
	}
	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func deployDonIDClaimerContract(e cldf.Environment, ab cldf.AddressBook, state stateview.CCIPOnChainState, chain cldf_evm.Chain) error {
	chainState, chainExists := state.Chains[chain.Selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}

	if state.Chains[chain.Selector].DonIDClaimer == nil {
		_, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*don_id_claimer.DonIDClaimer] {
				donIDClaimerAddr, tx2, donIDClaimerC, err2 := don_id_claimer.DeployDonIDClaimer(
					chain.DeployerKey,
					chain.Client,
					chainState.CapabilityRegistry.Address(),
				)
				return cldf.ContractDeploy[*don_id_claimer.DonIDClaimer]{
					Address: donIDClaimerAddr, Contract: donIDClaimerC, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.DonIDClaimer, deployment.Version1_6_1), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy donIDClaimer contract", "chain", chain.String(), "err", err)
			return err
		}
	} else {
		e.Logger.Infow("DonIDClaimer already deployed", "chain", chain.String(), "addr", chainState.DonIDClaimer.Address)
	}

	return nil
}

func deployDonIDClaimerPrecondition(e cldf.Environment, _ DeployDonIDClaimerConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	homeChainSel, err := state.HomeChainSelector()
	if err != nil {
		return fmt.Errorf("failed to get homeChainSelector state: %w", err)
	}

	return donIDClaimerValidationHelper(state, homeChainSel)
}

type DonIDClaimerOffSetConfig struct {
	OffSet uint32 `json:"offset"`
}

func donIDClaimerOffSetChangesetLogic(e cldf.Environment, cfg DonIDClaimerOffSetConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return cldf.ChangesetOutput{}, err
	}

	homeChainSel, err := state.HomeChainSelector()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get HomeChainSelector: %w", err)
	}

	// perform the offset operation
	donIDClaimer := state.Chains[homeChainSel].DonIDClaimer

	txOpts := e.BlockChains.EVMChains()[homeChainSel].DeployerKey
	txOpts.Context = e.GetContext()

	tx, err := donIDClaimer.SyncNextDONIdWithOffset(txOpts, cfg.OffSet)
	if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[homeChainSel], tx, don_id_claimer.DonIDClaimerABI, err); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error apply offset to donIDClaimer for chain %d: %w", homeChainSel, err)
	}

	return cldf.ChangesetOutput{}, err
}

func donIDClaimerOffSetChangesetPrecondition(e cldf.Environment, c DonIDClaimerOffSetConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return err
	}

	homeChainSel, err := state.HomeChainSelector()
	if err != nil {
		return fmt.Errorf("failed to get homeChainSelector state: %w", err)
	}

	// check the donIDClaimer contract exist
	if state.Chains[homeChainSel].DonIDClaimer == nil {
		return errors.New("donIDClaimer contract does not exist")
	}

	err = donIDClaimerValidationHelper(state, homeChainSel)
	if err != nil {
		return err
	}

	txOpts := e.BlockChains.EVMChains()[homeChainSel].DeployerKey
	// ensure deployer key is authorized
	isAuthorizedDeployer, err := state.Chains[homeChainSel].DonIDClaimer.IsAuthorizedDeployer(&bind.CallOpts{
		Context: e.GetContext(),
	}, txOpts.From)
	if err != nil {
		return fmt.Errorf("failed to run IsAuthorizedDeployed on home chain for donIDClaimer: %w", err)
	}

	if !isAuthorizedDeployer {
		return fmt.Errorf("deployerKey %v is not authorized deployer on donIDClaimer. ", txOpts.From.String())
	}

	return nil
}

func donIDClaimerValidationHelper(state stateview.CCIPOnChainState, homeChainSelector uint64) error {
	if err := cldf.IsValidChainSelector(homeChainSelector); err != nil {
		return fmt.Errorf("home chain selector invalid: %w", err)
	}

	_, exists := state.Chains[homeChainSelector]
	if !exists {
		return fmt.Errorf("home chain %d does not exist", homeChainSelector)
	}

	if state.Chains[homeChainSelector].CapabilityRegistry == nil {
		return errors.New("capabilityRegistry contract does not exist")
	}

	return nil
}
