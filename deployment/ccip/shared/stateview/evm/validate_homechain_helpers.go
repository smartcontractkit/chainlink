package evm

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

const (
	pluginNameCommit = "commit"
	pluginNameExec   = "exec"
)

// ValidateActiveOCR3DigestMatchesOffRamp checks that the CCIPHome versioned config digest
// matches the OffRamp's LatestConfigDetails digest for the config's plugin type.
// It intentionally does not validate JD, RMNHome wiring, signature verification, or OffRamp FChain.
func ValidateActiveOCR3DigestMatchesOffRamp(
	e cldf.Environment,
	offRamp offramp.OffRampInterface,
	homeCfg ccip_home.CCIPHomeVersionedConfig,
) error {
	if homeCfg.ConfigDigest == [32]byte{} {
		return errors.New("active config digest is empty")
	}
	ocrConfig, err := offRamp.LatestConfigDetails(
		&bind.CallOpts{Context: e.GetContext()},
		homeCfg.Config.PluginType,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get plugin %d config for chain %d offRamp %s: %w",
			homeCfg.Config.PluginType, homeCfg.Config.ChainSelector, offRamp.Address().Hex(), err,
		)
	}

	return validateActiveOCR3Digest(offRamp, homeCfg, ocrConfig)
}

func validateActiveOCR3Digest(
	offRamp offramp.OffRampInterface,
	homeCfg ccip_home.CCIPHomeVersionedConfig,
	ocrConfig offramp.MultiOCR3BaseOCRConfig,
) error {
	if ocrConfig.ConfigInfo.ConfigDigest != homeCfg.ConfigDigest {
		pluginName := pluginNameExec
		if homeCfg.Config.PluginType == uint8(types.PluginTypeCCIPCommit) {
			pluginName = pluginNameCommit
		}

		return fmt.Errorf(
			"offRamp %s %s config digest mismatch with CCIPHome for chain %d: expected %x, got %x",
			offRamp.Address().Hex(), pluginName, homeCfg.Config.ChainSelector,
			homeCfg.ConfigDigest, ocrConfig.ConfigInfo.ConfigDigest,
		)
	}

	return nil
}

// OffRampOCR3ConfigDigest reads the OCR3 config digest from the chain's OffRamp binding.
func (c CCIPChainState) OffRampOCR3ConfigDigest(callOpts *bind.CallOpts, pluginType uint8) ([32]byte, error) {
	if c.OffRamp == nil {
		return [32]byte{}, errors.New("no OffRamp contract found in the state")
	}
	ocrConfig, err := c.OffRamp.LatestConfigDetails(callOpts, pluginType)
	if err != nil {
		return [32]byte{}, fmt.Errorf("OffRamp.LatestConfigDetails(plugin=%d): %w", pluginType, err)
	}

	return ocrConfig.ConfigInfo.ConfigDigest, nil
}

// ChainHasActiveCommitAndExecDigests reports whether CCIPHome has non-empty active commit and
// exec OCR3 configs for chainSel on the given DON.
func (c CCIPChainState) ChainHasActiveCommitAndExecDigests(ctx context.Context, donID uint32, chainSel uint64) error {
	if c.CCIPHome == nil {
		return errors.New("no CCIPHome contract found in the state")
	}
	callOpts := &bind.CallOpts{Context: ctx}
	hasCommitActive, hasExecActive := false, false

	for _, pt := range []struct {
		enum uint8
		name string
	}{
		{uint8(types.PluginTypeCCIPCommit), pluginNameCommit},
		{uint8(types.PluginTypeCCIPExec), pluginNameExec},
	} {
		configs, err := c.CCIPHome.GetAllConfigs(callOpts, donID, pt.enum)
		if err != nil {
			return fmt.Errorf("chain %d plugin %s: %w", chainSel, pt.name, err)
		}
		active := configs.ActiveConfig
		if active.ConfigDigest != [32]byte{} && active.Config.ChainSelector == chainSel {
			switch pt.enum {
			case uint8(types.PluginTypeCCIPCommit):
				hasCommitActive = true
			case uint8(types.PluginTypeCCIPExec):
				hasExecActive = true
			}
		}
	}
	if !hasCommitActive || !hasExecActive {
		return fmt.Errorf(
			"chain %d missing active OCR3 digest on CCIPHome (commit=%t exec=%t) — run SetChainConfigsAndCandidates and PromoteCandidatesAndSetOCR3 first",
			chainSel, hasCommitActive, hasExecActive,
		)
	}

	return nil
}

// DonIDForChainFromCCIPHome finds the DON ID for chainSel using CCIPHome OCR3 configs.
// Returns (0, nil) when no DON is registered for the chain.
func DonIDForChainFromCCIPHome(
	ctx context.Context,
	ccipHome *ccip_home.CCIPHome,
	dons []capabilities_registry.CapabilitiesRegistryDONInfo,
	chainSelector uint64,
) (uint32, error) {
	if ccipHome == nil {
		return 0, errors.New("ccipHome is nil")
	}
	callOpts := &bind.CallOpts{Context: ctx}
	var donIDs []uint32
	for _, don := range dons {
		if len(don.CapabilityConfigurations) != 1 ||
			don.CapabilityConfigurations[0].CapabilityId != shared.CCIPCapabilityID {
			continue
		}
		active, candidate, err := ActiveAndCandidateConfigs(callOpts, ccipHome, don.Id, uint8(types.PluginTypeCCIPCommit))
		if err != nil {
			return 0, err
		}
		if active.ConfigDigest == [32]byte{} && candidate.ConfigDigest == [32]byte{} {
			active, candidate, err = ActiveAndCandidateConfigs(callOpts, ccipHome, don.Id, uint8(types.PluginTypeCCIPExec))
			if err != nil {
				return 0, err
			}
		}
		if active.Config.ChainSelector == chainSelector || candidate.Config.ChainSelector == chainSelector {
			donIDs = append(donIDs, don.Id)
		}
	}
	if len(donIDs) > 1 {
		return 0, fmt.Errorf(
			"more than one DON found for (chain selector %d, ccip capability id %x) pair",
			chainSelector, shared.CCIPCapabilityID[:],
		)
	}
	if len(donIDs) == 0 {
		return 0, nil
	}

	return donIDs[0], nil
}

// DonIDForChain resolves the CCIP DON ID for chainSel from CCIPHome and the capability registry.
func (c CCIPChainState) DonIDForChain(ctx context.Context, chainSel uint64) (uint32, error) {
	if c.CCIPHome == nil {
		return 0, errors.New("no CCIPHome contract found in the state")
	}
	if c.CapabilityRegistry == nil {
		return 0, errors.New("no CapabilityRegistry contract found in the state")
	}
	ccipDons, err := shared.GetCCIPDonsFromCapRegistry(ctx, c.CapabilityRegistry)
	if err != nil {
		return 0, fmt.Errorf("get CCIP DONs from capability registry: %w", err)
	}

	return DonIDForChainFromCCIPHome(ctx, c.CCIPHome, ccipDons, chainSel)
}

// DonIDFor is like DonIDForChain but treats donID==0 as an error.
func (c CCIPChainState) DonIDFor(ctx context.Context, chainSel uint64) (uint32, error) {
	donID, err := c.DonIDForChain(ctx, chainSel)
	if err != nil {
		return 0, fmt.Errorf("chain %d: DonIDForChain: %w", chainSel, err)
	}
	if donID == 0 {
		return 0, fmt.Errorf("chain %d: no DON found in CCIPHome", chainSel)
	}

	return donID, nil
}

// ValidateChainHasActiveCommitAndExecDigests resolves the CCIP DON for chainSel and verifies
// CCIPHome has non-empty active commit and exec OCR3 configs for that chain.
func (c CCIPChainState) ValidateChainHasActiveCommitAndExecDigests(ctx context.Context, chainSel uint64) error {
	donID, err := c.DonIDFor(ctx, chainSel)
	if err != nil {
		return err
	}

	return c.ChainHasActiveCommitAndExecDigests(ctx, donID, chainSel)
}

// ActiveAndCandidateConfigs fetches the active and candidate CCIPHome OCR3 configs for
// the given DON and plugin type in a single call.
func ActiveAndCandidateConfigs(
	callOpts *bind.CallOpts,
	ccipHome *ccip_home.CCIPHome,
	donID uint32,
	pluginType uint8,
) (active, candidate ccip_home.CCIPHomeVersionedConfig, err error) {
	configs, err := ccipHome.GetAllConfigs(callOpts, donID, pluginType)
	if err != nil {
		return active, candidate, fmt.Errorf("CCIPHome.GetAllConfigs(donID=%d, plugin=%d): %w", donID, pluginType, err)
	}

	return configs.ActiveConfig, configs.CandidateConfig, nil
}

// VerifyActiveMatchesOffRamp asserts the active OCR3 digest on CCIPHome for (donID, pluginType)
// is non-empty and matches what the chain's OffRamp reports for the same plugin.
func (c CCIPChainState) VerifyActiveMatchesOffRamp(
	ctx context.Context,
	e cldf.Environment,
	ccipHome *ccip_home.CCIPHome,
	donID uint32,
	pluginType uint8,
) error {
	if ccipHome == nil {
		return errors.New("ccipHome is nil")
	}
	if c.OffRamp == nil {
		return errors.New("no OffRamp contract found in the state")
	}
	active, _, err := ActiveAndCandidateConfigs(&bind.CallOpts{Context: ctx}, ccipHome, donID, pluginType)
	if err != nil {
		return err
	}
	if active.ConfigDigest == [32]byte{} {
		return fmt.Errorf("donID %d plugin %d: active digest is empty", donID, pluginType)
	}

	return ValidateActiveOCR3DigestMatchesOffRamp(e, c.OffRamp, active)
}
