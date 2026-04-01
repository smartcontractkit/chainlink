package evm

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	fqv2ops "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/fee_quoter"
	fqv2seq "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/sequences"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	viewshared "github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

// --- FeeQuoter Helpers ---

func isEthereumChain(selector uint64) bool {
	return selector == chain_selectors.ETHEREUM_MAINNET.Selector ||
		selector == chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector ||
		selector == chain_selectors.ETHEREUM_TESTNET_HOODI.Selector
}

// expectedNetworkFeeUSDCents: Ethereum involvement → 50, otherwise → 10
func expectedNetworkFeeUSDCents(srcSel, destSel uint64) uint32 {
	if isEthereumChain(destSel) || isEthereumChain(srcSel) {
		return 50
	}
	return 10
}

// expectedDefaultTokenFeeUSDCents: →ETH=150, ETH→=50, →SOL=35, other=25
func expectedDefaultTokenFeeUSDCents(srcSel, destSel uint64) uint16 {
	if isEthereumChain(destSel) {
		return 150
	}
	if isEthereumChain(srcSel) {
		return 50
	}
	destFamily, _ := chain_selectors.GetSelectorFamily(destSel)
	if destFamily == chain_selectors.FamilySolana {
		return 35
	}
	return 25
}

func getFeeTokensV2(callOpts *bind.CallOpts, backend bind.ContractBackend, addr common.Address) ([]common.Address, error) {
	parsed, err := abi.JSON(strings.NewReader(fqv2ops.FeeQuoterABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse FeeQuoter v2.0 ABI: %w", err)
	}
	bc := bind.NewBoundContract(addr, parsed, backend, backend, backend)
	var out []any
	if err := bc.Call(callOpts, &out, "getFeeTokens"); err != nil {
		return nil, fmt.Errorf("failed to call getFeeTokens on FeeQuoter v2.0 %s: %w", addr.Hex(), err)
	}
	return *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address), nil
}

func (c CCIPChainState) validateFeeTokenSuperset(
	callOpts *bind.CallOpts,
	fqAddr string,
	feeTokens []common.Address,
) error {
	if c.PriceRegistry == nil {
		return nil
	}
	legacyFeeTokens, err := c.PriceRegistry.GetFeeTokens(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get fee tokens from v1.5 PriceRegistry: %w", err)
	}
	feeTokenSet := make(map[common.Address]bool, len(feeTokens))
	for _, ft := range feeTokens {
		feeTokenSet[ft] = true
	}
	var errs []error
	for _, legacyFT := range legacyFeeTokens {
		if !feeTokenSet[legacyFT] {
			errs = append(errs, fmt.Errorf("FeeQuoter %s missing fee token %s from v1.5 PriceRegistry",
				fqAddr, legacyFT.Hex()))
		}
	}
	return errors.Join(errs...)
}

// --- FeeQuoter Validation ---

// ValidateFeeQuoter validates all FeeQuoter contracts (v1.6 and/or v2.0) for a chain
func (c CCIPChainState) ValidateFeeQuoter(
	e cldf.Environment,
	sourceChainSel uint64,
	connectedChains []uint64,
	fqV2 *fqv2ops.FeeQuoterContract,
	backend bind.ContractBackend,
) error {
	if c.FeeQuoter == nil && fqV2 == nil {
		return errors.New("no FeeQuoter contract (v1.6 or v2.0) found in the state")
	}
	callOpts := &bind.CallOpts{Context: e.GetContext()}
	var errs []error

	// v1.6 static config checks
	var v16FeeTokens []common.Address
	v16LaneReady := false
	if c.FeeQuoter != nil {
		fqAddr := c.FeeQuoter.Address().Hex()
		switch {
		case c.FeeQuoterVersion == nil:
			errs = append(errs, fmt.Errorf("FeeQuoter %s: version not set, cannot perform lane-level validation", fqAddr))
		case c.FeeQuoterVersion.Major() != 1:
			if fqV2 == nil {
				errs = append(errs, fmt.Errorf("FeeQuoter %s: unsupported version %s for lane-level validation",
					fqAddr, c.FeeQuoterVersion.String()))
			} else {
				e.Logger.Debugw("Skipping FeeQuoter v1.6 validation for non-v1 contract",
					"chain", sourceChainSel,
					"feeQuoter", fqAddr,
					"version", c.FeeQuoterVersion.String())
			}
		default:
			e.Logger.Debugw("Validating FeeQuoter v1.6", "chain", sourceChainSel, "feeQuoter", fqAddr, "connectedChains", len(connectedChains))
			staticConfig, err := c.FeeQuoter.GetStaticConfig(callOpts)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get static config for FeeQuoter %s: %w", fqAddr, err))
			} else {
				linktokenAddr, err := c.LinkTokenAddress()
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to get link token address from state: %w", err))
				} else if staticConfig.LinkToken != linktokenAddr {
					errs = append(errs, fmt.Errorf("FeeQuoter %s LinkToken mismatch: expected %s, got %s",
						fqAddr, linktokenAddr.Hex(), staticConfig.LinkToken.Hex()))
				}
				if staticConfig.TokenPriceStalenessThreshold == 0 {
					errs = append(errs, fmt.Errorf("FeeQuoter %s: TokenPriceStalenessThreshold is 0", fqAddr))
				}
			}
			feeTokens, err := c.FeeQuoter.GetFeeTokens(callOpts)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get fee tokens from FeeQuoter %s: %w", fqAddr, err))
			} else {
				v16FeeTokens = feeTokens
				v16LaneReady = len(v16FeeTokens) > 0
			}
		}
	}

	// v2.0 static config + owner checks
	v20Ready := fqV2 != nil
	if fqV2 != nil {
		fqAddr := fqV2.Address().Hex()
		e.Logger.Debugw("Validating FeeQuoter v2.0", "chain", sourceChainSel, "feeQuoterV2", fqAddr, "connectedChains", len(connectedChains))
		staticConfig, err := fqV2.GetStaticConfig(callOpts)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get static config for FeeQuoter v2.0 %s: %w", fqAddr, err))
			v20Ready = false
		} else {
			linktokenAddr, err := c.LinkTokenAddress()
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get link token address from state: %w", err))
			} else if staticConfig.LinkToken != linktokenAddr {
				errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s LinkToken mismatch: expected %s, got %s",
					fqAddr, linktokenAddr.Hex(), staticConfig.LinkToken.Hex()))
			}
		}
		owner, err := fqV2.Owner(callOpts)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get owner from FeeQuoter v2.0 %s: %w", fqAddr, err))
		} else if c.Timelock != nil && owner != c.Timelock.Address() {
			errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s not owned by Timelock %s, actual owner: %s",
				fqAddr, c.Timelock.Address().Hex(), owner.Hex()))
		}
	}

	var effectiveFqV2 *fqv2ops.FeeQuoterContract
	if v20Ready {
		effectiveFqV2 = fqV2
	}

	// Fee token validation (version-aware)
	if err := c.validateAllFeeTokenConfigs(callOpts, v16FeeTokens, effectiveFqV2, backend); err != nil {
		errs = append(errs, err)
	}

	if len(connectedChains) == 0 {
		return errors.Join(errs...)
	}

	// Dest chain config validation (version-aware).
	var laneV16FeeTokens []common.Address
	if v16LaneReady {
		laneV16FeeTokens = v16FeeTokens
	}
	if err := c.validateAllDestChainConfigs(callOpts, sourceChainSel, connectedChains, laneV16FeeTokens, effectiveFqV2); err != nil {
		errs = append(errs, err)
	}

	// Token transfer fee validation (version-aware)
	if err := c.validateAllTokenTransferFeeConfigs(e, callOpts, connectedChains, effectiveFqV2); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// validateAllFeeTokenConfigs validates fee tokens for all present FeeQuoter versions
// v1.6: non-empty, v1.5 PriceRegistry superset, premium multiplier per token, v1.5 cross-check
// v2.0: non-empty, v1.5 PriceRegistry superset
func (c CCIPChainState) validateAllFeeTokenConfigs(
	callOpts *bind.CallOpts,
	v16FeeTokens []common.Address,
	fqV2 *fqv2ops.FeeQuoterContract,
	backend bind.ContractBackend,
) error {
	var errs []error

	// v1.6 fee token checks
	if c.FeeQuoter != nil && v16FeeTokens != nil {
		fqAddr := c.FeeQuoter.Address().Hex()
		if len(v16FeeTokens) == 0 {
			errs = append(errs, fmt.Errorf("FeeQuoter %s has no fee tokens configured", fqAddr))
		}
		if err := c.validateFeeTokenSuperset(callOpts, fqAddr, v16FeeTokens); err != nil {
			errs = append(errs, err)
		}

		// Premium multiplier validation + v1.5 cross-check
		var anyLegacyOnRamp *evm_2_evm_onramp.EVM2EVMOnRamp
		if c.EVM2EVMOnRamp != nil {
			for _, onRamp := range c.EVM2EVMOnRamp {
				if onRamp != nil {
					anyLegacyOnRamp = onRamp
					break
				}
			}
		}
		for _, feeToken := range v16FeeTokens {
			premium, err := c.FeeQuoter.GetPremiumMultiplierWeiPerEth(callOpts, feeToken)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get PremiumMultiplierWeiPerEth for token %s on FeeQuoter %s: %w",
					feeToken.Hex(), fqAddr, err))
				continue
			}
			if anyLegacyOnRamp != nil {
				// Cross-check against v1.5 first — legacy mismatch is the most actionable signal.
				legacyFeeTokenCfg, err := anyLegacyOnRamp.GetFeeTokenConfig(callOpts, feeToken)
				if err == nil && legacyFeeTokenCfg.Enabled && premium != legacyFeeTokenCfg.PremiumMultiplierWeiPerEth {
					errs = append(errs, fmt.Errorf("FeeQuoter %s PremiumMultiplierWeiPerEth mismatch for fee token %s: "+
						"v1.6 has %d, v1.5 OnRamp had %d",
						fqAddr, feeToken.Hex(), premium, legacyFeeTokenCfg.PremiumMultiplierWeiPerEth))
				}
			} else if premium == 0 {
				// No legacy to compare — flag zero as standalone issue.
				errs = append(errs, fmt.Errorf("FeeQuoter %s PremiumMultiplierWeiPerEth is 0 for fee token %s",
					fqAddr, feeToken.Hex()))
			}
		}
	}

	// v2.0 fee token checks
	if fqV2 != nil {
		fqAddr := fqV2.Address().Hex()
		feeTokens, err := getFeeTokensV2(callOpts, backend, fqV2.Address())
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get fee tokens from FeeQuoter v2.0 %s: %w", fqAddr, err))
		} else {
			if len(feeTokens) == 0 {
				errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s has no fee tokens configured", fqAddr))
			}
			if err := c.validateFeeTokenSuperset(callOpts, fqAddr, feeTokens); err != nil {
				errs = append(errs, err)
			}

			// v1.6 -> v2.0 fee token subset check: all v1.6 fee tokens must exist in v2.0.
			if v16FeeTokens != nil {
				v16Set := make(map[common.Address]bool, len(v16FeeTokens))
				for _, ft := range v16FeeTokens {
					v16Set[ft] = true
				}
				v20Set := make(map[common.Address]bool, len(feeTokens))
				for _, ft := range feeTokens {
					v20Set[ft] = true
				}
				for _, ft := range v16FeeTokens {
					if !v20Set[ft] {
						errs = append(errs, fmt.Errorf("FeeQuoter v1.6 has fee token %s not present in v2.0 FeeQuoter %s",
							ft.Hex(), fqAddr))
					}
				}
			}
		}
	}

	return errors.Join(errs...)
}

// validateAllDestChainConfigs validates dest chain configs across v1.5, v1.6, and v2.0
// Cross-version field counts: v1.5↔v1.6 (11), v1.5↔v2.0 (6+NetworkFeeUSDCents), v1.6↔v2.0 (10)
func (c CCIPChainState) validateAllDestChainConfigs(
	callOpts *bind.CallOpts,
	sourceChainSel uint64,
	connectedChains []uint64,
	v16FeeTokens []common.Address,
	fqV2 *fqv2ops.FeeQuoterContract,
) error {
	var errs []error

	for _, destChainSel := range connectedChains {
		var v16Cfg *fee_quoter.FeeQuoterDestChainConfig
		var v20Cfg *fqv2ops.DestChainConfig
		var legacyCfg *evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig

		if c.FeeQuoter != nil && v16FeeTokens != nil {
			cfg, err := c.FeeQuoter.GetDestChainConfig(callOpts, destChainSel)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get FeeQuoter v1.6 dest chain config for chain %d: %w", destChainSel, err))
			} else {
				v16Cfg = &cfg
			}
		}
		if fqV2 != nil {
			cfg, err := fqV2.GetDestChainConfig(callOpts, destChainSel)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get FeeQuoter v2.0 dest chain config for chain %d: %w", destChainSel, err))
			} else {
				v20Cfg = &cfg
			}
		}
		if legacyOnRamp := c.EVM2EVMOnRamp[destChainSel]; legacyOnRamp != nil {
			cfg, err := legacyOnRamp.GetDynamicConfig(callOpts)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get v1.5 OnRamp dynamic config for dest chain %d: %w", destChainSel, err))
			} else {
				legacyCfg = &cfg
			}
		}

		v16Enabled := v16Cfg != nil && v16Cfg.IsEnabled
		v20Enabled := v20Cfg != nil && v20Cfg.IsEnabled

		// Skip v1.6 checks when lane is enabled only in v2.0.
		if v16Cfg != nil && (v16Enabled || !v20Enabled) {
			if err := c.validateV16DestChainConfig(callOpts, sourceChainSel, destChainSel, *v16Cfg, legacyCfg); err != nil {
				errs = append(errs, err)
			}
		}
		if v20Cfg != nil {
			v16ForV20 := v16Cfg
			if !v16Enabled && v20Enabled {
				v16ForV20 = nil
			}
			if err := c.validateV20DestChainConfig(callOpts, sourceChainSel, destChainSel, *v20Cfg, v16ForV20, legacyCfg, fqV2); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

// validateV16DestChainConfig validates a single v1.6 dest chain config
func (c CCIPChainState) validateV16DestChainConfig(
	callOpts *bind.CallOpts,
	sourceChainSel, destChainSel uint64,
	destCfg fee_quoter.FeeQuoterDestChainConfig,
	legacyCfg *evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig,
) error {
	header := fmt.Sprintf("FeeQuoter v1.6 %s dest chain %d", c.FeeQuoter.Address().Hex(), destChainSel)

	if !destCfg.IsEnabled {
		return groupErrors(header, []error{errors.New("not enabled — dest chain not configured, skipping field checks")})
	}

	var errs []error

	// Cross-version field mapping: v1.6 <-> v1.5
	if legacyCfg != nil {
		if err := compareFieldChecks("v1.5<->v1.6", []fieldCheck{
			{"MaxNumberOfTokensPerMsg", uint64(destCfg.MaxNumberOfTokensPerMsg), uint64(legacyCfg.MaxNumberOfTokensPerMsg)},
			{"DestGasOverhead", uint64(destCfg.DestGasOverhead), uint64(legacyCfg.DestGasOverhead)},
			{"DestDataAvailabilityOverheadGas", uint64(destCfg.DestDataAvailabilityOverheadGas), uint64(legacyCfg.DestDataAvailabilityOverheadGas)},
			{"DestGasPerDataAvailabilityByte", uint64(destCfg.DestGasPerDataAvailabilityByte), uint64(legacyCfg.DestGasPerDataAvailabilityByte)},
			{"DestDataAvailabilityMultiplierBps", uint64(destCfg.DestDataAvailabilityMultiplierBps), uint64(legacyCfg.DestDataAvailabilityMultiplierBps)},
			{"MaxDataBytes", uint64(destCfg.MaxDataBytes), uint64(legacyCfg.MaxDataBytes)},
			{"MaxPerMsgGasLimit", uint64(destCfg.MaxPerMsgGasLimit), uint64(legacyCfg.MaxPerMsgGasLimit)},
			{"DefaultTokenDestGasOverhead", uint64(destCfg.DefaultTokenDestGasOverhead), uint64(legacyCfg.DefaultTokenDestGasOverhead)},
			{"EnforceOutOfOrder", destCfg.EnforceOutOfOrder, legacyCfg.EnforceOutOfOrder},
			{"DestGasPerPayloadByteBase", uint64(destCfg.DestGasPerPayloadByteBase), uint64(uint8(legacyCfg.DestGasPerPayloadByte))}, //nolint:gosec // G115: intentional v1.5 uint16->uint8 truncation during migration
		}); err != nil {
			errs = append(errs, err)
		}
	}

	// v1.6 field checks (always applied)
	if destCfg.ChainFamilySelector == [4]byte{} {
		errs = append(errs, errors.New("ChainFamilySelector is empty"))
	}
	// Compare staleness against v1.5 only when a legacy lane exists.
	if legacyCfg != nil && c.PriceRegistry != nil {
		st, stErr := c.PriceRegistry.GetStalenessThreshold(callOpts)
		if stErr != nil {
			errs = append(errs, fmt.Errorf("failed to get staleness threshold from v1.5 PriceRegistry: %w", stErr))
		} else if !st.IsUint64() || st.Uint64() > uint64(^uint32(0)) {
			errs = append(errs, fmt.Errorf("v1.5 PriceRegistry StalenessThreshold %s overflows uint32", st.String()))
		} else if want := uint32(st.Uint64()); want != destCfg.GasPriceStalenessThreshold { //nolint:gosec // G115: safe, verified <= MaxUint32
			errs = append(errs, fmt.Errorf("v1.5<->v1.6:\n  GasPriceStalenessThreshold: got=%d, want=%d",
				destCfg.GasPriceStalenessThreshold, want))
		}
	} else if destCfg.GasPriceStalenessThreshold == 0 {
		errs = append(errs, errors.New("GasPriceStalenessThreshold is 0"))
	}
	if err := compareFieldChecks("deploy-constants", []fieldCheck{
		{"DestGasPerPayloadByteHigh", uint64(destCfg.DestGasPerPayloadByteHigh), uint64(ccipevm.CalldataGasPerByteHigh)},
		{"DestGasPerPayloadByteThreshold", uint64(destCfg.DestGasPerPayloadByteThreshold), uint64(ccipevm.CalldataGasPerByteThreshold)},
		{"DefaultTxGasLimit", uint64(destCfg.DefaultTxGasLimit), uint64(200_000)},
		{"GasMultiplierWeiPerEth", destCfg.GasMultiplierWeiPerEth, uint64(11e17)},
	}); err != nil {
		errs = append(errs, err)
	}
	if err := compareFieldChecks("topology", []fieldCheck{
		{"NetworkFeeUSDCents", uint64(destCfg.NetworkFeeUSDCents), uint64(expectedNetworkFeeUSDCents(sourceChainSel, destChainSel))},
		{"DefaultTokenFeeUSDCents", uint64(destCfg.DefaultTokenFeeUSDCents), uint64(expectedDefaultTokenFeeUSDCents(sourceChainSel, destChainSel))},
	}); err != nil {
		errs = append(errs, err)
	}

	return groupErrors(header, errs)
}

// validateV20DestChainConfig validates a single v2.0 dest chain config
func (c CCIPChainState) validateV20DestChainConfig(
	callOpts *bind.CallOpts,
	sourceChainSel, destChainSel uint64,
	destCfgV2 fqv2ops.DestChainConfig,
	v16Cfg *fee_quoter.FeeQuoterDestChainConfig,
	legacyCfg *evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig,
	fqV2 *fqv2ops.FeeQuoterContract,
) error {
	header := fmt.Sprintf("FeeQuoter v2.0 %s dest chain %d", fqV2.Address().Hex(), destChainSel)

	if !destCfgV2.IsEnabled {
		return groupErrors(header, []error{errors.New("not enabled — dest chain not configured, skipping field checks")})
	}

	var errs []error

	// v2.0 business-rule fields
	if destCfgV2.ChainFamilySelector == [4]byte{} {
		errs = append(errs, errors.New("ChainFamilySelector is empty"))
	}
	if err := compareFieldChecks("deploy-constants", []fieldCheck{
		{"LinkFeeMultiplierPercent", uint64(destCfgV2.LinkFeeMultiplierPercent), uint64(fqv2seq.LinkFeeMultiplierPercent)},
		{"DefaultTxGasLimit", uint64(destCfgV2.DefaultTxGasLimit), uint64(200_000)},
	}); err != nil {
		errs = append(errs, err)
	}
	if err := compareFieldChecks("topology", []fieldCheck{
		{"NetworkFeeUSDCents", uint64(destCfgV2.NetworkFeeUSDCents), uint64(expectedNetworkFeeUSDCents(sourceChainSel, destChainSel))},
		{"DefaultTokenFeeUSDCents", uint64(destCfgV2.DefaultTokenFeeUSDCents), uint64(expectedDefaultTokenFeeUSDCents(sourceChainSel, destChainSel))},
	}); err != nil {
		errs = append(errs, err)
	}

	// Cross-version field mapping: v1.6 <-> v2.0 (v2.0 is got, v1.6 is want)
	if v16Cfg != nil {
		if err := compareFieldChecks("v1.6<->v2.0", []fieldCheck{
			{"IsEnabled", destCfgV2.IsEnabled, v16Cfg.IsEnabled},
			{"MaxDataBytes", uint64(destCfgV2.MaxDataBytes), uint64(v16Cfg.MaxDataBytes)},
			{"MaxPerMsgGasLimit", uint64(destCfgV2.MaxPerMsgGasLimit), uint64(v16Cfg.MaxPerMsgGasLimit)},
			{"DestGasOverhead", uint64(destCfgV2.DestGasOverhead), uint64(v16Cfg.DestGasOverhead)},
			{"DestGasPerPayloadByteBase", uint64(destCfgV2.DestGasPerPayloadByteBase), uint64(v16Cfg.DestGasPerPayloadByteBase)},
			{"ChainFamilySelector", destCfgV2.ChainFamilySelector, v16Cfg.ChainFamilySelector},
			{"DefaultTokenDestGasOverhead", uint64(destCfgV2.DefaultTokenDestGasOverhead), uint64(v16Cfg.DefaultTokenDestGasOverhead)},
			{"DefaultTxGasLimit", uint64(destCfgV2.DefaultTxGasLimit), uint64(v16Cfg.DefaultTxGasLimit)},
		}); err != nil {
			errs = append(errs, err)
		}
	}

	// Cross-version field mapping: v2.0 <-> v1.5
	if legacyCfg != nil {
		v15Checks := []fieldCheck{
			{"DestGasOverhead", uint64(destCfgV2.DestGasOverhead), uint64(legacyCfg.DestGasOverhead)},
			{"MaxDataBytes", uint64(destCfgV2.MaxDataBytes), uint64(legacyCfg.MaxDataBytes)},
			{"MaxPerMsgGasLimit", uint64(destCfgV2.MaxPerMsgGasLimit), uint64(legacyCfg.MaxPerMsgGasLimit)},
			{"DefaultTokenDestGasOverhead", uint64(destCfgV2.DefaultTokenDestGasOverhead), uint64(legacyCfg.DefaultTokenDestGasOverhead)},
			{"DestGasPerPayloadByteBase", uint64(destCfgV2.DestGasPerPayloadByteBase), uint64(uint8(legacyCfg.DestGasPerPayloadByte))}, //nolint:gosec // G115: intentional v1.5 uint16->uint8 truncation during migration
		}
		if err := compareFieldChecks("v1.5<->v2.0", v15Checks); err != nil {
			errs = append(errs, err)
		}
	}

	return groupErrors(header, errs)
}

// validateAllTokenTransferFeeConfigs validates token transfer fees across v1.5, v1.6, and v2.0
// Cross-version field counts: v1.5↔v1.6 (5), v1.6→v2.0 (FeeUSDCents, DestGasOverhead, DestBytesOverhead; DeciBps+MaxFeeUSDCents dropped)
func (c CCIPChainState) validateAllTokenTransferFeeConfigs(
	e cldf.Environment,
	callOpts *bind.CallOpts,
	connectedChains []uint64,
	fqV2 *fqv2ops.FeeQuoterContract,
) error {
	if c.FeeQuoter == nil && fqV2 == nil {
		return nil
	}
	if c.TokenAdminRegistry == nil {
		return errors.New("no TokenAdminRegistry contract found, cannot validate token transfer fee configs")
	}

	allTokens, err := viewshared.GetSupportedTokens(c.TokenAdminRegistry)
	if err != nil {
		return fmt.Errorf("failed to get configured tokens from TokenAdminRegistry: %w", err)
	}

	addrToSymbol := make(map[common.Address]string)
	if symbolMap, symErr := c.TokenAddressBySymbol(); symErr == nil {
		for symbol, addr := range symbolMap {
			addrToSymbol[addr] = string(symbol)
		}
	}

	e.Logger.Debugw("Validating TokenTransferFeeConfigs", "tokens", len(allTokens), "connectedChains", len(connectedChains))
	var mu sync.Mutex
	var errs []error
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)
	for _, tokenAddr := range allTokens {
		token := tokenAddr
		tokenLabel := token.Hex()
		if sym, ok := addrToSymbol[token]; ok {
			tokenLabel = fmt.Sprintf("%s (%s)", sym, token.Hex())
		}
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			var tokenErrs []error
			for _, destChainSel := range connectedChains {
				if err := c.validateTokenTransferFee(callOpts, destChainSel, token, tokenLabel, fqV2); err != nil {
					tokenErrs = append(tokenErrs, err)
				}
			}
			if len(tokenErrs) > 0 {
				mu.Lock()
				errs = append(errs, tokenErrs...)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	return errors.Join(errs...)
}

// validateTokenTransferFee validates a single token+dest pair across all present FeeQuoter versions
func (c CCIPChainState) validateTokenTransferFee(
	callOpts *bind.CallOpts,
	destChainSel uint64,
	token common.Address,
	tokenLabel string,
	fqV2 *fqv2ops.FeeQuoterContract,
) error {
	var errs []error

	var v16Cfg *fee_quoter.FeeQuoterTokenTransferFeeConfig
	v16Enabled := false
	if c.FeeQuoter != nil {
		cfg, err := c.FeeQuoter.GetTokenTransferFeeConfig(callOpts, destChainSel, token)
		if err == nil {
			v16Enabled = cfg.IsEnabled
			if cfg.IsEnabled {
				v16Cfg = &cfg
			}
		}
	}

	var v20Cfg *fqv2ops.TokenTransferFeeConfig
	v20Enabled := false
	if fqV2 != nil {
		cfg, err := fqV2.GetTokenTransferFeeConfig(callOpts, destChainSel, token)
		if err == nil {
			v20Enabled = cfg.IsEnabled
			if cfg.IsEnabled {
				v20Cfg = &cfg
			}
		}
	}

	// Cross-version IsEnabled consistency
	if c.FeeQuoter != nil && fqV2 != nil && v16Enabled != v20Enabled {
		errs = append(errs, fmt.Errorf("IsEnabled mismatch: v1.6=%v, v2.0=%v", v16Enabled, v20Enabled))
	}

	if v16Cfg == nil && v20Cfg == nil {
		header := fmt.Sprintf("token %s dest %d", tokenLabel, destChainSel)
		return groupErrors(header, errs)
	}

	// v1.6 invariants + v1.5 cross-check
	if v16Cfg != nil {
		if v16Cfg.MinFeeUSDCents >= v16Cfg.MaxFeeUSDCents {
			errs = append(errs, fmt.Errorf("v1.6: MinFeeUSDCents (%d) must be less than MaxFeeUSDCents (%d)",
				v16Cfg.MinFeeUSDCents, v16Cfg.MaxFeeUSDCents))
		}
		if v16Cfg.DestBytesOverhead < globals.CCIPLockOrBurnV1RetBytes {
			errs = append(errs, fmt.Errorf("v1.6: DestBytesOverhead (%d) must be at least %d",
				v16Cfg.DestBytesOverhead, globals.CCIPLockOrBurnV1RetBytes))
		}

		// v1.6 <-> v1.5 field mapping
		if legacyOnRamp := c.EVM2EVMOnRamp[destChainSel]; legacyOnRamp != nil {
			legacyTTF, legacyErr := legacyOnRamp.GetTokenTransferFeeConfig(callOpts, token)
			if legacyErr == nil && legacyTTF.IsEnabled {
				if err := compareFieldChecks("v1.5<->v1.6", []fieldCheck{
					{"MinFeeUSDCents", uint64(v16Cfg.MinFeeUSDCents), uint64(legacyTTF.MinFeeUSDCents)},
					{"MaxFeeUSDCents", uint64(v16Cfg.MaxFeeUSDCents), uint64(legacyTTF.MaxFeeUSDCents)},
					{"DeciBps", uint64(v16Cfg.DeciBps), uint64(legacyTTF.DeciBps)},
					{"DestGasOverhead", uint64(v16Cfg.DestGasOverhead), uint64(legacyTTF.DestGasOverhead)},
					{"DestBytesOverhead", uint64(v16Cfg.DestBytesOverhead), uint64(legacyTTF.DestBytesOverhead)},
				}); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	// v2.0 invariants + cross-checks
	if v20Cfg != nil {
		if v20Cfg.DestBytesOverhead < globals.CCIPLockOrBurnV1RetBytes {
			errs = append(errs, fmt.Errorf("v2.0: DestBytesOverhead (%d) must be at least %d",
				v20Cfg.DestBytesOverhead, globals.CCIPLockOrBurnV1RetBytes))
		}

		// v2.0 <-> v1.6 field mapping
		if v16Cfg != nil {
			if v20Cfg.FeeUSDCents != v16Cfg.MinFeeUSDCents {
				errs = append(errs, fmt.Errorf("v2.0: FeeUSDCents (%d) != v1.6 MinFeeUSDCents (%d)",
					v20Cfg.FeeUSDCents, v16Cfg.MinFeeUSDCents))
			}
			if err := compareFieldChecks("v1.6<->v2.0", []fieldCheck{
				{"DestGasOverhead", uint64(v20Cfg.DestGasOverhead), uint64(v16Cfg.DestGasOverhead)},
				{"DestBytesOverhead", uint64(v20Cfg.DestBytesOverhead), uint64(v16Cfg.DestBytesOverhead)},
			}); err != nil {
				errs = append(errs, err)
			}
		}
	}

	header := fmt.Sprintf("token %s dest %d", tokenLabel, destChainSel)
	return groupErrors(header, errs)
}
