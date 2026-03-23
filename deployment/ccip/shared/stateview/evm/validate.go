package evm

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	viewshared "github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

// ValidateNonceManager checks that NonceManager previous ramps point to the correct v1.5 contracts
func (c CCIPChainState) ValidateNonceManager(
	e cldf.Environment,
	selector uint64,
	connectedChains []uint64,
) error {
	if c.NonceManager == nil {
		return errors.New("no NonceManager contract found in the state")
	}
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

// ValidateRMNProxy checks that RMNProxy.GetARM() returns the RMNRemote address
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

// ValidateFeeQuoter performs chain-level checks and version-specific lane-level validation
// Migrated chains are cross-checked against v1.5, fresh chains against defaults
func (c CCIPChainState) ValidateFeeQuoter(
	e cldf.Environment,
	sourceChainSel uint64,
	connectedChains []uint64,
) error {
	if c.FeeQuoter == nil {
		return errors.New("no FeeQuoter contract found in the state")
	}
	callOpts := &bind.CallOpts{Context: e.GetContext()}
	fqAddr := c.FeeQuoter.Address().Hex()
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
		if err := c.validateTokenTransferFeeConfigs(callOpts, fqAddr, connectedChains); err != nil {
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
			// GasMultiplierWeiPerEth moved from per-token (v1.5 FeeTokenConfig) to per-dest (v1.6)
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

// validateTokenTransferFeeConfigs checks per-token-per-dest fee invariants and v1.5 cross-checks
func (c CCIPChainState) validateTokenTransferFeeConfigs(
	callOpts *bind.CallOpts,
	fqAddr string,
	connectedChains []uint64,
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

	var errs []error
	for _, tokenAddr := range allTokens {
		tokenLabel := tokenAddr.Hex()
		if sym, ok := addrToSymbol[tokenAddr]; ok {
			tokenLabel = fmt.Sprintf("%s (%s)", sym, tokenAddr.Hex())
		}

		for _, destChainSel := range connectedChains {
			ttfCfg, err := c.FeeQuoter.GetTokenTransferFeeConfig(callOpts, destChainSel, tokenAddr)
			if err != nil {
				continue
			}
			if !ttfCfg.IsEnabled {
				continue
			}
			if ttfCfg.MinFeeUSDCents >= ttfCfg.MaxFeeUSDCents {
				errs = append(errs, fmt.Errorf("FeeQuoter %s TokenTransferFeeConfig for token %s to dest chain %d: "+
					"MinFeeUSDCents (%d) must be less than MaxFeeUSDCents (%d)",
					fqAddr, tokenLabel, destChainSel,
					ttfCfg.MinFeeUSDCents, ttfCfg.MaxFeeUSDCents))
			}
			if ttfCfg.DestBytesOverhead < globals.CCIPLockOrBurnV1RetBytes {
				errs = append(errs, fmt.Errorf("FeeQuoter %s TokenTransferFeeConfig for token %s to dest chain %d: "+
					"DestBytesOverhead (%d) must be at least %d",
					fqAddr, tokenLabel, destChainSel,
					ttfCfg.DestBytesOverhead, globals.CCIPLockOrBurnV1RetBytes))
			}

			legacyOnRamp := c.EVM2EVMOnRamp[destChainSel]
			if legacyOnRamp == nil {
				continue
			}
			legacyTTF, err := legacyOnRamp.GetTokenTransferFeeConfig(callOpts, tokenAddr)
			if err != nil || !legacyTTF.IsEnabled {
				continue
			}
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
					errs = append(errs, fmt.Errorf("FeeQuoter %s TokenTransferFeeConfig for token %s to dest %d: "+
						"%s mismatch: v1.6=%d, v1.5=%d",
						fqAddr, tokenLabel, destChainSel, chk.name, chk.v16Val, chk.v15Val))
				}
			}
		}
	}

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

// validateFeeQuoterDestCfgDefaults checks fresh v1.6 dest config against known defaults
// (mirrors DefaultFeeQuoterDestChainConfig in cs_chain_contracts.go)
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

// ValidateContractOwnership checks all CCIP contracts are owned by the MCMS Timelock
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
	if c.RMNRemote != nil {
		if err := checkOwnership(callOpts, "RMNRemote", c.RMNRemote, timelockAddr); err != nil {
			errs = append(errs, err)
		}
	}
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

	if c.ProposerMcm != nil {
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
	}

	return errors.Join(errs...)
}
