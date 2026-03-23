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

// ValidateNonceManager checks NonceManager previous ramps against v1.5 contracts.
func (c CCIPChainState) ValidateNonceManager(
	e cldf.Environment,
	selector uint64,
	connectedChains []uint64,
) error {
	if c.NonceManager == nil {
		return errors.New("no NonceManager contract found in the state")
	}
	e.Logger.Debugw("Validating NonceManager", "chain", selector, "nonceManager", c.NonceManager.Address().Hex(), "connectedChains", len(connectedChains))
	callOpts := &bind.CallOpts{Context: e.GetContext()}
	var errs []error

	for _, remoteChainSel := range connectedChains {
		if remoteChainSel == selector {
			continue
		}
		previousRamps, err := c.NonceManager.GetPreviousRamps(callOpts, remoteChainSel)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get previous ramps for remote chain %d on chain %d: %w",
				remoteChainSel, selector, err))
			continue
		}
		if c.EVM2EVMOnRamp != nil && c.EVM2EVMOnRamp[remoteChainSel] != nil {
			expectedOnRamp := c.EVM2EVMOnRamp[remoteChainSel].Address()
			if previousRamps.PrevOnRamp != expectedOnRamp {
				errs = append(errs, fmt.Errorf("NonceManager %s PrevOnRamp mismatch for remote chain %d on chain %d: expected %s, got %s",
					c.NonceManager.Address().Hex(), remoteChainSel, selector,
					expectedOnRamp.Hex(), previousRamps.PrevOnRamp.Hex()))
			}
		}
		if c.EVM2EVMOffRamp != nil && c.EVM2EVMOffRamp[remoteChainSel] != nil {
			expectedOffRamp := c.EVM2EVMOffRamp[remoteChainSel].Address()
			if previousRamps.PrevOffRamp != expectedOffRamp {
				errs = append(errs, fmt.Errorf("NonceManager %s PrevOffRamp mismatch for remote chain %d on chain %d: expected %s, got %s",
					c.NonceManager.Address().Hex(), remoteChainSel, selector,
					expectedOffRamp.Hex(), previousRamps.PrevOffRamp.Hex()))
			}
		}
	}
	return errors.Join(errs...)
}

// ValidateRMNProxy checks that RMNProxy.GetARM() returns the RMNRemote address.
func (c CCIPChainState) ValidateRMNProxy(e cldf.Environment) error {
	if c.RMNProxy == nil {
		return errors.New("no RMNProxy contract found in the state")
	}
	if c.RMNRemote == nil {
		return errors.New("no RMNRemote contract found for RMNProxy validation")
	}
	callOpts := &bind.CallOpts{Context: e.GetContext()}
	armAddr, err := c.RMNProxy.GetARM(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get ARM from RMNProxy %s: %w", c.RMNProxy.Address().Hex(), err)
	}
	if armAddr != c.RMNRemote.Address() {
		return fmt.Errorf("RMNProxy %s GetARM mismatch: expected RMNRemote %s, got %s",
			c.RMNProxy.Address().Hex(), c.RMNRemote.Address().Hex(), armAddr.Hex())
	}
	return nil
}

func isEthereumChain(selector uint64) bool {
	return selector == chain_selectors.ETHEREUM_MAINNET.Selector ||
		selector == chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector
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

// ValidateFeeQuoter performs chain-level and lane-level validation.
func (c CCIPChainState) ValidateFeeQuoter(
	e cldf.Environment,
	sourceChainSel uint64,
	connectedChains []uint64,
	fqV2 *fqv2ops.FeeQuoterContract,
) error {
	if c.FeeQuoter == nil {
		return errors.New("no FeeQuoter contract found in the state")
	}
	callOpts := &bind.CallOpts{Context: e.GetContext()}
	fqAddr := c.FeeQuoter.Address().Hex()
	e.Logger.Debugw("Validating FeeQuoter", "chain", sourceChainSel, "feeQuoter", fqAddr, "connectedChains", len(connectedChains))
	var errs []error

	staticConfig, err := c.FeeQuoter.GetStaticConfig(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get static config for FeeQuoter %s: %w", fqAddr, err)
	}
	linktokenAddr, err := c.LinkTokenAddress()
	if err != nil {
		return fmt.Errorf("failed to get link token address from state: %w", err)
	}
	if staticConfig.LinkToken != linktokenAddr {
		errs = append(errs, fmt.Errorf("FeeQuoter %s LinkToken mismatch: expected %s, got %s",
			fqAddr, linktokenAddr.Hex(), staticConfig.LinkToken.Hex()))
	}

	feeTokens, err := c.FeeQuoter.GetFeeTokens(callOpts)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to get fee tokens from FeeQuoter %s: %w", fqAddr, err))
		return errors.Join(errs...)
	}

	if err := c.validateFeeTokenConfigs(callOpts, fqAddr, feeTokens); err != nil {
		errs = append(errs, err)
	}

	if len(connectedChains) == 0 {
		// No lanes wired yet — skip lane-level validation (valid during early deployment)
		return errors.Join(errs...)
	}

	if c.FeeQuoterVersion == nil {
		errs = append(errs, fmt.Errorf("FeeQuoter %s: version not set, cannot perform lane-level validation", fqAddr))
		return errors.Join(errs...)
	}
	switch c.FeeQuoterVersion.Major() {
	case 1:
		if err := c.validateDestChainConfigs(callOpts, fqAddr, sourceChainSel, connectedChains, feeTokens); err != nil {
			errs = append(errs, err)
		}
		if err := c.validateTokenTransferFeeConfigs(e, callOpts, fqAddr, connectedChains, fqV2); err != nil {
			errs = append(errs, err)
		}
	case 2:
		// TODO: implement FeeQuoter 2.0 lane-level validation
	default:
		errs = append(errs, fmt.Errorf("FeeQuoter %s: unsupported version %s for lane-level validation",
			fqAddr, c.FeeQuoterVersion.String()))
	}

	return errors.Join(errs...)
}

// validateDestChainConfigs cross-checks against v1.5 OnRamp if migrated, or validates defaults if fresh
func (c CCIPChainState) validateDestChainConfigs(
	callOpts *bind.CallOpts,
	fqAddr string,
	sourceChainSel uint64,
	connectedChains []uint64,
	feeTokens []common.Address,
) error {
	var errs []error

	for _, destChainSel := range connectedChains {
		destCfg, err := c.FeeQuoter.GetDestChainConfig(callOpts, destChainSel)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get FeeQuoter dest chain config for chain %d: %w", destChainSel, err))
			continue
		}
		if !destCfg.IsEnabled {
			errs = append(errs, fmt.Errorf("FeeQuoter %s dest chain config not enabled for chain %d", fqAddr, destChainSel))
		}

		legacyOnRamp := c.EVM2EVMOnRamp[destChainSel]
		if legacyOnRamp != nil {
			if err := c.validateFeeQuoterAgainstLegacyOnRamp(callOpts, fqAddr, destChainSel, destCfg, legacyOnRamp); err != nil {
				errs = append(errs, err)
			}
			// GasMultiplierWeiPerEth moved from per-token to per-dest in v1.6
			for _, ft := range feeTokens {
				legacyFTCfg, err := legacyOnRamp.GetFeeTokenConfig(callOpts, ft)
				if err != nil || !legacyFTCfg.Enabled {
					continue
				}
				if destCfg.GasMultiplierWeiPerEth != legacyFTCfg.GasMultiplierWeiPerEth {
					errs = append(errs, fmt.Errorf("FeeQuoter %s GasMultiplierWeiPerEth mismatch for dest chain %d: "+
						"v1.6=%d, v1.5 FeeTokenConfig=%d",
						fqAddr, destChainSel, destCfg.GasMultiplierWeiPerEth, legacyFTCfg.GasMultiplierWeiPerEth))
				}
				break
			}
		} else {
			if err := validateFeeQuoterDestCfgDefaults(fqAddr, sourceChainSel, destChainSel, destCfg); err != nil {
				errs = append(errs, err)
			}
		}

		if destCfg.ChainFamilySelector == [4]byte{} {
			errs = append(errs, fmt.Errorf("FeeQuoter %s ChainFamilySelector is empty for dest chain %d", fqAddr, destChainSel))
		}
		if destCfg.GasPriceStalenessThreshold == 0 {
			errs = append(errs, fmt.Errorf("FeeQuoter %s GasPriceStalenessThreshold is 0 for dest chain %d", fqAddr, destChainSel))
		}
		for _, chk := range []struct {
			name     string
			got      uint64
			expected uint64
		}{
			{"DestGasPerPayloadByteHigh", uint64(destCfg.DestGasPerPayloadByteHigh), uint64(ccipevm.CalldataGasPerByteHigh)},
			{"DestGasPerPayloadByteThreshold", uint64(destCfg.DestGasPerPayloadByteThreshold), uint64(ccipevm.CalldataGasPerByteThreshold)},
			{"DefaultTxGasLimit", uint64(destCfg.DefaultTxGasLimit), 200_000},
			{"NetworkFeeUSDCents", uint64(destCfg.NetworkFeeUSDCents), uint64(expectedNetworkFeeUSDCents(sourceChainSel, destChainSel))},
		} {
			if chk.got != chk.expected {
				errs = append(errs, fmt.Errorf("FeeQuoter %s %s mismatch for dest chain %d: expected %d, got %d",
					fqAddr, chk.name, destChainSel, chk.expected, chk.got))
			}
		}

		destFamily, _ := chain_selectors.GetSelectorFamily(destChainSel)
		if destFamily != chain_selectors.FamilyEVM && !destCfg.EnforceOutOfOrder {
			errs = append(errs, fmt.Errorf("FeeQuoter %s EnforceOutOfOrder must be true for non-EVM dest chain %d (family %s)",
				fqAddr, destChainSel, destFamily))
		}
	}

	return errors.Join(errs...)
}

// validateFeeTokenConfigs checks fee token presence, v1.5 PriceRegistry superset, and premium multipliers
func (c CCIPChainState) validateFeeTokenConfigs(
	callOpts *bind.CallOpts,
	fqAddr string,
	feeTokens []common.Address,
) error {
	var errs []error

	if len(feeTokens) == 0 {
		errs = append(errs, fmt.Errorf("FeeQuoter %s has no fee tokens configured", fqAddr))
	}
	if c.PriceRegistry != nil {
		legacyFeeTokens, err := c.PriceRegistry.GetFeeTokens(callOpts)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get fee tokens from v1.5 PriceRegistry: %w", err))
		} else {
			feeTokenSet := make(map[common.Address]bool, len(feeTokens))
			for _, ft := range feeTokens {
				feeTokenSet[ft] = true
			}
			for _, legacyFT := range legacyFeeTokens {
				if !feeTokenSet[legacyFT] {
					errs = append(errs, fmt.Errorf("FeeQuoter %s missing fee token %s from v1.5 PriceRegistry",
						fqAddr, legacyFT.Hex()))
				}
			}
		}
	}

	var anyLegacyOnRamp *evm_2_evm_onramp.EVM2EVMOnRamp
	if c.EVM2EVMOnRamp != nil {
		for _, onRamp := range c.EVM2EVMOnRamp {
			if onRamp != nil {
				anyLegacyOnRamp = onRamp
				break
			}
		}
	}

	for _, feeToken := range feeTokens {
		premium, err := c.FeeQuoter.GetPremiumMultiplierWeiPerEth(callOpts, feeToken)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get PremiumMultiplierWeiPerEth for token %s on FeeQuoter %s: %w",
				feeToken.Hex(), fqAddr, err))
			continue
		}
		if premium == 0 {
			errs = append(errs, fmt.Errorf("FeeQuoter %s PremiumMultiplierWeiPerEth is 0 for fee token %s",
				fqAddr, feeToken.Hex()))
		}
		if anyLegacyOnRamp != nil {
			legacyFeeTokenCfg, err := anyLegacyOnRamp.GetFeeTokenConfig(callOpts, feeToken)
			if err == nil && legacyFeeTokenCfg.Enabled && premium != legacyFeeTokenCfg.PremiumMultiplierWeiPerEth {
				errs = append(errs, fmt.Errorf("FeeQuoter %s PremiumMultiplierWeiPerEth mismatch for fee token %s: "+
					"v1.6 has %d, v1.5 OnRamp had %d",
					fqAddr, feeToken.Hex(), premium, legacyFeeTokenCfg.PremiumMultiplierWeiPerEth))
			}
		}
	}

	return errors.Join(errs...)
}

// validateTokenTransferFeeConfigs checks per-token-per-dest fee invariants and v1.5 cross-checks.
// When fqV2 is non-nil, also validates v2.0 token transfer fees in the same pass to avoid duplicate RPC calls.
func (c CCIPChainState) validateTokenTransferFeeConfigs(
	e cldf.Environment,
	callOpts *bind.CallOpts,
	fqAddr string,
	connectedChains []uint64,
	fqV2 *fqv2ops.FeeQuoterContract,
) error {
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
	sem := make(chan struct{}, 20) // max 20 concurrent RPC calls
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
				ttfCfg, err := c.FeeQuoter.GetTokenTransferFeeConfig(callOpts, destChainSel, token)
				if err != nil {
					continue
				}
				if !ttfCfg.IsEnabled {
					continue
				}
				if ttfCfg.MinFeeUSDCents >= ttfCfg.MaxFeeUSDCents {
					tokenErrs = append(tokenErrs, fmt.Errorf("FeeQuoter %s TokenTransferFeeConfig for token %s to dest chain %d: "+
						"MinFeeUSDCents (%d) must be less than MaxFeeUSDCents (%d)",
						fqAddr, tokenLabel, destChainSel,
						ttfCfg.MinFeeUSDCents, ttfCfg.MaxFeeUSDCents))
				}
				if ttfCfg.DestBytesOverhead < globals.CCIPLockOrBurnV1RetBytes {
					tokenErrs = append(tokenErrs, fmt.Errorf("FeeQuoter %s TokenTransferFeeConfig for token %s to dest chain %d: "+
						"DestBytesOverhead (%d) must be at least %d",
						fqAddr, tokenLabel, destChainSel,
						ttfCfg.DestBytesOverhead, globals.CCIPLockOrBurnV1RetBytes))
				}

				// v1.5 legacy cross-check
				legacyOnRamp := c.EVM2EVMOnRamp[destChainSel]
				if legacyOnRamp != nil {
					legacyTTF, legacyErr := legacyOnRamp.GetTokenTransferFeeConfig(callOpts, token)
					if legacyErr == nil && legacyTTF.IsEnabled {
						for _, chk := range []struct {
							name   string
							v16Val uint64
							v15Val uint64
						}{
							{"MinFeeUSDCents", uint64(ttfCfg.MinFeeUSDCents), uint64(legacyTTF.MinFeeUSDCents)},
							{"MaxFeeUSDCents", uint64(ttfCfg.MaxFeeUSDCents), uint64(legacyTTF.MaxFeeUSDCents)},
							{"DeciBps", uint64(ttfCfg.DeciBps), uint64(legacyTTF.DeciBps)},
							{"DestGasOverhead", uint64(ttfCfg.DestGasOverhead), uint64(legacyTTF.DestGasOverhead)},
							{"DestBytesOverhead", uint64(ttfCfg.DestBytesOverhead), uint64(legacyTTF.DestBytesOverhead)},
						} {
							if chk.v16Val != chk.v15Val {
								tokenErrs = append(tokenErrs, fmt.Errorf("FeeQuoter %s TokenTransferFeeConfig for token %s to dest %d: "+
									"%s mismatch: v1.6=%d, v1.5=%d",
									fqAddr, tokenLabel, destChainSel, chk.name, chk.v16Val, chk.v15Val))
							}
						}
					}
				}

				// v2.0 cross-check (reuses v1.6 ttfCfg already fetched above)
				if fqV2 != nil {
					ttfCfgV2, v2Err := fqV2.GetTokenTransferFeeConfig(callOpts, destChainSel, token)
					if v2Err == nil && ttfCfgV2.IsEnabled {
						fqV2Addr := fqV2.Address().Hex()
						if ttfCfgV2.DestBytesOverhead < globals.CCIPLockOrBurnV1RetBytes {
							tokenErrs = append(tokenErrs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest chain %d: "+
								"DestBytesOverhead (%d) must be at least %d",
								fqV2Addr, tokenLabel, destChainSel,
								ttfCfgV2.DestBytesOverhead, globals.CCIPLockOrBurnV1RetBytes))
						}
						if ttfCfgV2.FeeUSDCents != ttfCfg.MinFeeUSDCents {
							tokenErrs = append(tokenErrs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest %d: "+
								"FeeUSDCents (%d) != v1.6 MinFeeUSDCents (%d)",
								fqV2Addr, tokenLabel, destChainSel, ttfCfgV2.FeeUSDCents, ttfCfg.MinFeeUSDCents))
						}
						if ttfCfg.DeciBps > 0 {
							tokenErrs = append(tokenErrs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest %d: "+
								"v1.6 DeciBps=%d is non-zero but DeciBps is removed in v2.0 (percentage fee lost)",
								fqV2Addr, tokenLabel, destChainSel, ttfCfg.DeciBps))
						}
						if ttfCfg.MaxFeeUSDCents > ttfCfg.MinFeeUSDCents {
							tokenErrs = append(tokenErrs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest %d: "+
								"v1.6 MaxFeeUSDCents (%d) > MinFeeUSDCents (%d) — fee cap is not present in v2.0",
								fqV2Addr, tokenLabel, destChainSel, ttfCfg.MaxFeeUSDCents, ttfCfg.MinFeeUSDCents))
						}
						for _, chk := range []struct {
							name   string
							v20Val uint64
							v16Val uint64
						}{
							{"DestGasOverhead", uint64(ttfCfgV2.DestGasOverhead), uint64(ttfCfg.DestGasOverhead)},
							{"DestBytesOverhead", uint64(ttfCfgV2.DestBytesOverhead), uint64(ttfCfg.DestBytesOverhead)},
						} {
							if chk.v20Val != chk.v16Val {
								tokenErrs = append(tokenErrs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest %d: "+
									"%s mismatch: v2.0=%d, v1.6=%d",
									fqV2Addr, tokenLabel, destChainSel, chk.name, chk.v20Val, chk.v16Val))
							}
						}
					}
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

// validateFeeQuoterAgainstLegacyOnRamp cross-checks v1.6 dest chain config against v1.5 OnRamp DynamicConfig
func (c CCIPChainState) validateFeeQuoterAgainstLegacyOnRamp(
	callOpts *bind.CallOpts,
	fqAddr string,
	destChainSel uint64,
	destCfg fee_quoter.FeeQuoterDestChainConfig,
	legacyOnRamp *evm_2_evm_onramp.EVM2EVMOnRamp,
) error {
	legacyCfg, err := legacyOnRamp.GetDynamicConfig(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get v1.5 OnRamp dynamic config for dest chain %d: %w", destChainSel, err)
	}
	var errs []error

	for _, chk := range []struct {
		name string
		v16  any
		v15  any
	}{
		{"MaxNumberOfTokensPerMsg", uint64(destCfg.MaxNumberOfTokensPerMsg), uint64(legacyCfg.MaxNumberOfTokensPerMsg)},
		{"DestGasOverhead", uint64(destCfg.DestGasOverhead), uint64(legacyCfg.DestGasOverhead)},
		{"DestDataAvailabilityOverheadGas", uint64(destCfg.DestDataAvailabilityOverheadGas), uint64(legacyCfg.DestDataAvailabilityOverheadGas)},
		{"DestGasPerDataAvailabilityByte", uint64(destCfg.DestGasPerDataAvailabilityByte), uint64(legacyCfg.DestGasPerDataAvailabilityByte)},
		{"DestDataAvailabilityMultiplierBps", uint64(destCfg.DestDataAvailabilityMultiplierBps), uint64(legacyCfg.DestDataAvailabilityMultiplierBps)},
		{"MaxDataBytes", uint64(destCfg.MaxDataBytes), uint64(legacyCfg.MaxDataBytes)},
		{"MaxPerMsgGasLimit", uint64(destCfg.MaxPerMsgGasLimit), uint64(legacyCfg.MaxPerMsgGasLimit)},
		{"DefaultTokenDestGasOverhead", uint64(destCfg.DefaultTokenDestGasOverhead), uint64(legacyCfg.DefaultTokenDestGasOverhead)},
		{"DefaultTokenFeeUSDCents", uint64(destCfg.DefaultTokenFeeUSDCents), uint64(legacyCfg.DefaultTokenFeeUSDCents)},
		{"EnforceOutOfOrder", destCfg.EnforceOutOfOrder, legacyCfg.EnforceOutOfOrder},
		{"DestGasPerPayloadByteBase", uint64(destCfg.DestGasPerPayloadByteBase), uint64(uint8(legacyCfg.DestGasPerPayloadByte))}, //nolint:gosec // G115: intentional v1.5 uint16→uint8 truncation during migration
	} {
		if chk.v16 != chk.v15 {
			errs = append(errs, fmt.Errorf("FeeQuoter %s %s mismatch for dest chain %d: v1.6=%v, v1.5=%v",
				fqAddr, chk.name, destChainSel, chk.v16, chk.v15))
		}
	}

	return errors.Join(errs...)
}

// validateFeeQuoterDestCfgDefaults checks dest config against known defaults.
func validateFeeQuoterDestCfgDefaults(
	fqAddr string,
	sourceChainSel uint64,
	destChainSel uint64,
	destCfg fee_quoter.FeeQuoterDestChainConfig,
) error {
	var errs []error

	for _, chk := range []struct {
		name     string
		got      uint64
		expected uint64
	}{
		{"MaxNumberOfTokensPerMsg", uint64(destCfg.MaxNumberOfTokensPerMsg), 10},
		{"MaxDataBytes", uint64(destCfg.MaxDataBytes), 30_000},
		{"MaxPerMsgGasLimit", uint64(destCfg.MaxPerMsgGasLimit), 3_000_000},
		{"DestGasOverhead", uint64(destCfg.DestGasOverhead), uint64(ccipevm.DestGasOverhead)},
		{"DestGasPerPayloadByteBase", uint64(destCfg.DestGasPerPayloadByteBase), uint64(ccipevm.CalldataGasPerByteBase)},
		{"DefaultTokenDestGasOverhead", uint64(destCfg.DefaultTokenDestGasOverhead), 90_000},
		{"DestDataAvailabilityOverheadGas", uint64(destCfg.DestDataAvailabilityOverheadGas), 100},
		{"DestGasPerDataAvailabilityByte", uint64(destCfg.DestGasPerDataAvailabilityByte), 16},
		{"DestDataAvailabilityMultiplierBps", uint64(destCfg.DestDataAvailabilityMultiplierBps), 1},
		{"GasMultiplierWeiPerEth", destCfg.GasMultiplierWeiPerEth, uint64(11e17)},
		{"DefaultTokenFeeUSDCents", uint64(destCfg.DefaultTokenFeeUSDCents), uint64(expectedDefaultTokenFeeUSDCents(sourceChainSel, destChainSel))},
	} {
		if chk.got != chk.expected {
			errs = append(errs, fmt.Errorf("FeeQuoter %s %s mismatch for dest chain %d: expected %d, got %d",
				fqAddr, chk.name, destChainSel, chk.expected, chk.got))
		}
	}

	return errors.Join(errs...)
}

type ownableContract interface {
	Owner(opts *bind.CallOpts) (common.Address, error)
	Address() common.Address
}

func checkOwnership(callOpts *bind.CallOpts, name string, contract ownableContract, expectedOwner common.Address) error {
	owner, err := contract.Owner(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get %s owner: %w", name, err)
	}
	if owner != expectedOwner {
		return fmt.Errorf("%s %s not owned by expected owner %s, actual owner: %s",
			name, contract.Address().Hex(), expectedOwner.Hex(), owner.Hex())
	}
	return nil
}

// ValidateContractOwnership checks CCIP contracts are owned by the MCMS Timelock.
func (c CCIPChainState) ValidateContractOwnership(e cldf.Environment) error {
	if c.Timelock == nil {
		return errors.New("timelock not found in state, cannot validate ownership")
	}
	timelockAddr := c.Timelock.Address()
	callOpts := &bind.CallOpts{Context: e.GetContext()}
	var errs []error

	if c.FeeQuoter != nil {
		if err := checkOwnership(callOpts, "FeeQuoter", c.FeeQuoter, timelockAddr); err != nil {
			errs = append(errs, err)
		}
	}
	if c.NonceManager != nil {
		if err := checkOwnership(callOpts, "NonceManager", c.NonceManager, timelockAddr); err != nil {
			errs = append(errs, err)
		}
	}
	/* if c.RMNRemote != nil {
		if err := checkOwnership(callOpts, "RMNRemote", c.RMNRemote, timelockAddr); err != nil {
			errs = append(errs, err)
		}
	} */
	if c.OnRamp != nil {
		if err := checkOwnership(callOpts, "OnRamp", c.OnRamp, timelockAddr); err != nil {
			errs = append(errs, err)
		}
	}
	if c.OffRamp != nil {
		if err := checkOwnership(callOpts, "OffRamp", c.OffRamp, timelockAddr); err != nil {
			errs = append(errs, err)
		}
	}
	if c.Router != nil {
		if err := checkOwnership(callOpts, "Router", c.Router, timelockAddr); err != nil {
			errs = append(errs, err)
		}
	}

	/* if c.ProposerMcm != nil {
		if err := checkOwnership(callOpts, "ProposerMcm", c.ProposerMcm, c.ProposerMcm.Address()); err != nil {
			errs = append(errs, err)
		}
	}
	if c.CancellerMcm != nil {
		if err := checkOwnership(callOpts, "CancellerMcm", c.CancellerMcm, c.CancellerMcm.Address()); err != nil {
			errs = append(errs, err)
		}
	}
	if c.BypasserMcm != nil {
		if err := checkOwnership(callOpts, "BypasserMcm", c.BypasserMcm, c.BypasserMcm.Address()); err != nil {
			errs = append(errs, err)
		}
	} */

	return errors.Join(errs...)
}

// — FeeQuoter v2.0 validation —

// normalizedDestChainConfig holds DestChainConfig fields shared between v1.6 and v2.0.
type normalizedDestChainConfig struct {
	IsEnabled                   bool
	MaxDataBytes                uint64
	MaxPerMsgGasLimit           uint64
	DestGasOverhead             uint64
	DestGasPerPayloadByteBase   uint64
	ChainFamilySelector         [4]byte
	DefaultTokenFeeUSDCents     uint64
	DefaultTokenDestGasOverhead uint64
	DefaultTxGasLimit           uint64
	NetworkFeeUSDCents          uint64
}

func normalizeFromV16FQ(cfg fee_quoter.FeeQuoterDestChainConfig) normalizedDestChainConfig {
	return normalizedDestChainConfig{
		IsEnabled:                   cfg.IsEnabled,
		MaxDataBytes:                uint64(cfg.MaxDataBytes),
		MaxPerMsgGasLimit:           uint64(cfg.MaxPerMsgGasLimit),
		DestGasOverhead:             uint64(cfg.DestGasOverhead),
		DestGasPerPayloadByteBase:   uint64(cfg.DestGasPerPayloadByteBase),
		ChainFamilySelector:         cfg.ChainFamilySelector,
		DefaultTokenFeeUSDCents:     uint64(cfg.DefaultTokenFeeUSDCents),
		DefaultTokenDestGasOverhead: uint64(cfg.DefaultTokenDestGasOverhead),
		DefaultTxGasLimit:           uint64(cfg.DefaultTxGasLimit),
		NetworkFeeUSDCents:          uint64(cfg.NetworkFeeUSDCents),
	}
}

func normalizeFromV20FQ(cfg fqv2ops.DestChainConfig) normalizedDestChainConfig {
	return normalizedDestChainConfig{
		IsEnabled:                   cfg.IsEnabled,
		MaxDataBytes:                uint64(cfg.MaxDataBytes),
		MaxPerMsgGasLimit:           uint64(cfg.MaxPerMsgGasLimit),
		DestGasOverhead:             uint64(cfg.DestGasOverhead),
		DestGasPerPayloadByteBase:   uint64(cfg.DestGasPerPayloadByteBase),
		ChainFamilySelector:         cfg.ChainFamilySelector,
		DefaultTokenFeeUSDCents:     uint64(cfg.DefaultTokenFeeUSDCents),
		DefaultTokenDestGasOverhead: uint64(cfg.DefaultTokenDestGasOverhead),
		DefaultTxGasLimit:           uint64(cfg.DefaultTxGasLimit),
		NetworkFeeUSDCents:          uint64(cfg.NetworkFeeUSDCents),
	}
}

// compareNormalizedDestConfigs errors on any field mismatch between a and b.
func compareNormalizedDestConfigs(fqAddr string, destChainSel uint64, label string, a, b normalizedDestChainConfig) error {
	var errs []error
	if a.IsEnabled != b.IsEnabled {
		errs = append(errs, fmt.Errorf("FeeQuoter %s dest chain %d IsEnabled mismatch (%s): %v vs %v",
			fqAddr, destChainSel, label, a.IsEnabled, b.IsEnabled))
	}
	for _, chk := range []struct {
		name string
		aVal uint64
		bVal uint64
	}{
		{"MaxDataBytes", a.MaxDataBytes, b.MaxDataBytes},
		{"MaxPerMsgGasLimit", a.MaxPerMsgGasLimit, b.MaxPerMsgGasLimit},
		{"DestGasOverhead", a.DestGasOverhead, b.DestGasOverhead},
		{"DestGasPerPayloadByteBase", a.DestGasPerPayloadByteBase, b.DestGasPerPayloadByteBase},
		{"DefaultTokenFeeUSDCents", a.DefaultTokenFeeUSDCents, b.DefaultTokenFeeUSDCents},
		{"DefaultTokenDestGasOverhead", a.DefaultTokenDestGasOverhead, b.DefaultTokenDestGasOverhead},
		{"DefaultTxGasLimit", a.DefaultTxGasLimit, b.DefaultTxGasLimit},
		{"NetworkFeeUSDCents", a.NetworkFeeUSDCents, b.NetworkFeeUSDCents},
	} {
		if chk.aVal != chk.bVal {
			errs = append(errs, fmt.Errorf("FeeQuoter %s dest chain %d %s mismatch (%s): %d vs %d",
				fqAddr, destChainSel, chk.name, label, chk.aVal, chk.bVal))
		}
	}
	if a.ChainFamilySelector != b.ChainFamilySelector {
		errs = append(errs, fmt.Errorf("FeeQuoter %s dest chain %d ChainFamilySelector mismatch (%s): %x vs %x",
			fqAddr, destChainSel, label, a.ChainFamilySelector, b.ChainFamilySelector))
	}
	return errors.Join(errs...)
}

// getFeeTokensV2 calls getFeeTokens on a FeeQuoter 2.0 via raw ABI call
// (the operations-gen wrapper doesn't expose GetFeeTokens).
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

// ValidateFeeQuoterV2 validates a FeeQuoter v2.0 deployment against on-chain state.
func (c CCIPChainState) ValidateFeeQuoterV2(
	e cldf.Environment,
	sourceChainSel uint64,
	connectedChains []uint64,
	fqV2 *fqv2ops.FeeQuoterContract,
	backend bind.ContractBackend,
) error {
	callOpts := &bind.CallOpts{Context: e.GetContext()}
	fqAddr := fqV2.Address().Hex()
	e.Logger.Debugw("Validating FeeQuoter v2.0", "chain", sourceChainSel, "feeQuoterV2", fqAddr, "connectedChains", len(connectedChains))
	var errs []error

	staticConfig, err := fqV2.GetStaticConfig(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get static config for FeeQuoter v2.0 %s: %w", fqAddr, err)
	}
	linktokenAddr, err := c.LinkTokenAddress()
	if err != nil {
		return fmt.Errorf("failed to get link token address from state: %w", err)
	}
	if staticConfig.LinkToken != linktokenAddr {
		errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s LinkToken mismatch: expected %s, got %s",
			fqAddr, linktokenAddr.Hex(), staticConfig.LinkToken.Hex()))
	}

	owner, err := fqV2.Owner(callOpts)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to get owner from FeeQuoter v2.0 %s: %w", fqAddr, err))
	} else if c.Timelock != nil && owner != c.Timelock.Address() {
		errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s not owned by Timelock %s, actual owner: %s",
			fqAddr, c.Timelock.Address().Hex(), owner.Hex()))
	}

	if len(connectedChains) == 0 {
		return errors.Join(errs...)
	}

	if err := c.validateFeeTokenConfigsV20(callOpts, fqAddr, backend, fqV2.Address()); err != nil {
		errs = append(errs, err)
	}
	if err := c.validateDestChainConfigsV20(callOpts, fqAddr, sourceChainSel, connectedChains, fqV2); err != nil {
		errs = append(errs, err)
	}
	// When v1.6 FQ exists, token transfer fees are validated in the combined pass (ValidateFeeQuoter).
	// When v1.6 FQ is absent, run standalone v2.0 token fee validation.
	if c.FeeQuoter == nil {
		if err := c.validateTokenTransferFeeConfigsV20(callOpts, fqAddr, connectedChains, fqV2); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// validateFeeTokenConfigsV20 checks fee token presence and v1.5 PriceRegistry superset for v2.0.
func (c CCIPChainState) validateFeeTokenConfigsV20(
	callOpts *bind.CallOpts,
	fqAddr string,
	backend bind.ContractBackend,
	addr common.Address,
) error {
	feeTokens, err := getFeeTokensV2(callOpts, backend, addr)
	if err != nil {
		return fmt.Errorf("failed to get fee tokens from FeeQuoter v2.0 %s: %w", fqAddr, err)
	}

	var errs []error
	if len(feeTokens) == 0 {
		errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s has no fee tokens configured", fqAddr))
	}

	if c.PriceRegistry != nil {
		legacyFeeTokens, err := c.PriceRegistry.GetFeeTokens(callOpts)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get fee tokens from v1.5 PriceRegistry for v2.0 FQ check: %w", err))
		} else {
			feeTokenSet := make(map[common.Address]bool, len(feeTokens))
			for _, ft := range feeTokens {
				feeTokenSet[ft] = true
			}
			for _, legacyFT := range legacyFeeTokens {
				if !feeTokenSet[legacyFT] {
					errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s missing fee token %s from v1.5 PriceRegistry",
						fqAddr, legacyFT.Hex()))
				}
			}
		}
	}

	return errors.Join(errs...)
}

// validateDestChainConfigsV20 validates v2.0 dest chain configs against v1.6 and v1.5 state.
// Case B (c.FeeQuoter != nil): cross-checks v1.6↔v2.0 and v1.5↔v2.0 for each dest.
// Case C (c.FeeQuoter == nil): cross-checks v1.5↔v2.0 directly where a v1.5 OnRamp exists.
func (c CCIPChainState) validateDestChainConfigsV20(
	callOpts *bind.CallOpts,
	fqAddr string,
	sourceChainSel uint64,
	connectedChains []uint64,
	fqV2 *fqv2ops.FeeQuoterContract,
) error {
	var errs []error

	for _, destChainSel := range connectedChains {
		destCfgV2, err := fqV2.GetDestChainConfig(callOpts, destChainSel)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get FeeQuoter v2.0 dest chain config for chain %d: %w", destChainSel, err))
			continue
		}
		if !destCfgV2.IsEnabled {
			errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s dest chain config not enabled for chain %d", fqAddr, destChainSel))
		}

		// v2.0-specific fields
		if destCfgV2.LinkFeeMultiplierPercent != fqv2seq.LinkFeeMultiplierPercent {
			errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s dest chain %d LinkFeeMultiplierPercent: expected %d, got %d",
				fqAddr, destChainSel, fqv2seq.LinkFeeMultiplierPercent, destCfgV2.LinkFeeMultiplierPercent))
		}
		expectedNetworkFee := uint16(expectedNetworkFeeUSDCents(sourceChainSel, destChainSel)) //nolint:gosec // G115: max value is 50, always fits uint16
		if destCfgV2.NetworkFeeUSDCents != expectedNetworkFee {
			errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s dest chain %d NetworkFeeUSDCents: expected %d, got %d",
				fqAddr, destChainSel, expectedNetworkFee, destCfgV2.NetworkFeeUSDCents))
		}
		if destCfgV2.DefaultTxGasLimit != 200_000 {
			errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s dest chain %d DefaultTxGasLimit: expected 200000, got %d",
				fqAddr, destChainSel, destCfgV2.DefaultTxGasLimit))
		}

		normV2 := normalizeFromV20FQ(destCfgV2)
		_ = normV2 // used in v1.6↔v2.0 comparison below

		if c.FeeQuoter != nil {
			destCfgV16, err := c.FeeQuoter.GetDestChainConfig(callOpts, destChainSel)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get FeeQuoter v1.6 dest chain config for chain %d: %w", destChainSel, err))
			} else {
				normV16 := normalizeFromV16FQ(destCfgV16)
				if err := compareNormalizedDestConfigs(fqAddr, destChainSel, "v1.6↔v2.0", normV16, normV2); err != nil {
					errs = append(errs, err)
				}
			}
		}

		if legacyOnRamp := c.EVM2EVMOnRamp[destChainSel]; legacyOnRamp != nil {
			if err := compareV15OnRampWithV20DestCfg(callOpts, fqAddr, destChainSel, destCfgV2, legacyOnRamp); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

// compareV15OnRampWithV20DestCfg compares v1.5 OnRamp DynamicConfig with v2.0 DestChainConfig.
func compareV15OnRampWithV20DestCfg(
	callOpts *bind.CallOpts,
	fqAddr string,
	destChainSel uint64,
	destCfgV2 fqv2ops.DestChainConfig,
	legacyOnRamp *evm_2_evm_onramp.EVM2EVMOnRamp,
) error {
	legacyCfg, err := legacyOnRamp.GetDynamicConfig(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get v1.5 OnRamp dynamic config for dest chain %d (v2.0 cross-check): %w", destChainSel, err)
	}
	var errs []error
	for _, chk := range []struct {
		name string
		v20  any
		v15  any
	}{
		{"DestGasOverhead", uint64(destCfgV2.DestGasOverhead), uint64(legacyCfg.DestGasOverhead)},
		{"MaxDataBytes", uint64(destCfgV2.MaxDataBytes), uint64(legacyCfg.MaxDataBytes)},
		{"MaxPerMsgGasLimit", uint64(destCfgV2.MaxPerMsgGasLimit), uint64(legacyCfg.MaxPerMsgGasLimit)},
		{"DefaultTokenDestGasOverhead", uint64(destCfgV2.DefaultTokenDestGasOverhead), uint64(legacyCfg.DefaultTokenDestGasOverhead)},
		{"DefaultTokenFeeUSDCents", uint64(destCfgV2.DefaultTokenFeeUSDCents), uint64(legacyCfg.DefaultTokenFeeUSDCents)},
		{"DestGasPerPayloadByteBase", uint64(destCfgV2.DestGasPerPayloadByteBase), uint64(uint8(legacyCfg.DestGasPerPayloadByte))}, //nolint:gosec // G115: intentional v1.5 uint16→uint8 truncation during migration
	} {
		if chk.v20 != chk.v15 {
			errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s %s mismatch for dest chain %d: v2.0=%v, v1.5=%v",
				fqAddr, chk.name, destChainSel, chk.v20, chk.v15))
		}
	}
	return errors.Join(errs...)
}

// validateTokenTransferFeeConfigsV20 validates per-token per-dest fee configs for v2.0.
func (c CCIPChainState) validateTokenTransferFeeConfigsV20(
	callOpts *bind.CallOpts,
	fqAddr string,
	connectedChains []uint64,
	fqV2 *fqv2ops.FeeQuoterContract,
) error {
	if c.TokenAdminRegistry == nil {
		return errors.New("no TokenAdminRegistry contract found, cannot validate v2.0 token transfer fee configs")
	}

	allTokens, err := viewshared.GetSupportedTokens(c.TokenAdminRegistry)
	if err != nil {
		return fmt.Errorf("failed to get configured tokens from TokenAdminRegistry for v2.0 validation: %w", err)
	}

	addrToSymbol := make(map[common.Address]string)
	if symbolMap, symErr := c.TokenAddressBySymbol(); symErr == nil {
		for symbol, addr := range symbolMap {
			addrToSymbol[addr] = string(symbol)
		}
	}

	var errs []error
	for _, tokenAddr := range allTokens {
		tokenLabel := tokenAddr.Hex()
		if sym, ok := addrToSymbol[tokenAddr]; ok {
			tokenLabel = fmt.Sprintf("%s (%s)", sym, tokenAddr.Hex())
		}

		for _, destChainSel := range connectedChains {
			ttfCfgV2, err := fqV2.GetTokenTransferFeeConfig(callOpts, destChainSel, tokenAddr)
			if err != nil {
				continue
			}
			if !ttfCfgV2.IsEnabled {
				continue
			}
			if ttfCfgV2.DestBytesOverhead < globals.CCIPLockOrBurnV1RetBytes {
				errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest chain %d: "+
					"DestBytesOverhead (%d) must be at least %d",
					fqAddr, tokenLabel, destChainSel,
					ttfCfgV2.DestBytesOverhead, globals.CCIPLockOrBurnV1RetBytes))
			}

			if c.FeeQuoter == nil {
				continue
			}

			ttfCfgV16, err := c.FeeQuoter.GetTokenTransferFeeConfig(callOpts, destChainSel, tokenAddr)
			if err != nil || !ttfCfgV16.IsEnabled {
				continue
			}

			// FeeUSDCents replaces MinFeeUSDCents — values must match
			if ttfCfgV2.FeeUSDCents != ttfCfgV16.MinFeeUSDCents {
				errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest %d: "+
					"FeeUSDCents (%d) != v1.6 MinFeeUSDCents (%d)",
					fqAddr, tokenLabel, destChainSel, ttfCfgV2.FeeUSDCents, ttfCfgV16.MinFeeUSDCents))
			}
			// DeciBps removed in v2.0; a non-zero v1.6 value means percentage fee was active and is now lost
			if ttfCfgV16.DeciBps > 0 {
				errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest %d: "+
					"v1.6 DeciBps=%d is non-zero but DeciBps is removed in v2.0 (percentage fee lost)",
					fqAddr, tokenLabel, destChainSel, ttfCfgV16.DeciBps))
			}
			// MaxFeeUSDCents removed in v2.0; a non-trivial cap means fee ceiling would be lost
			if ttfCfgV16.MaxFeeUSDCents > ttfCfgV16.MinFeeUSDCents {
				errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest %d: "+
					"v1.6 MaxFeeUSDCents (%d) > MinFeeUSDCents (%d) — fee cap is not present in v2.0",
					fqAddr, tokenLabel, destChainSel, ttfCfgV16.MaxFeeUSDCents, ttfCfgV16.MinFeeUSDCents))
			}
			// Overhead fields must match
			for _, chk := range []struct {
				name   string
				v20Val uint64
				v16Val uint64
			}{
				{"DestGasOverhead", uint64(ttfCfgV2.DestGasOverhead), uint64(ttfCfgV16.DestGasOverhead)},
				{"DestBytesOverhead", uint64(ttfCfgV2.DestBytesOverhead), uint64(ttfCfgV16.DestBytesOverhead)},
			} {
				if chk.v20Val != chk.v16Val {
					errs = append(errs, fmt.Errorf("FeeQuoter v2.0 %s TokenTransferFeeConfig for token %s to dest %d: "+
						"%s mismatch: v2.0=%d, v1.6=%d",
						fqAddr, tokenLabel, destChainSel, chk.name, chk.v20Val, chk.v16Val))
				}
			}
		}
	}

	return errors.Join(errs...)
}
