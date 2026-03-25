package evm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
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

// --- Helpers ---

type fieldCheck struct {
	name string
	got  any
	want any
}

func compareFieldChecks(section string, checks []fieldCheck) error {
	var lines []string
	for _, chk := range checks {
		if chk.got != chk.want {
			lines = append(lines, fmt.Sprintf("%s: got=%v, want=%v", chk.name, chk.got, chk.want))
		}
	}
	if len(lines) == 0 {
		return nil
	}
	return fmt.Errorf("%s:\n  %s", section, strings.Join(lines, "\n  "))
}

// groupErrors wraps sub-errors under a single header line to avoid repeating context.
func groupErrors(header string, errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var lines []string
	for _, e := range errs {
		for _, line := range strings.Split(e.Error(), "\n") {
			if line != "" {
				lines = append(lines, line)
			}
		}
	}
	return fmt.Errorf("%s:\n  %s", header, strings.Join(lines, "\n  "))
}

// --- Ownership ---

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

	for _, ct := range []struct {
		name     string
		contract ownableContract
		present  bool
	}{
		{"FeeQuoter", c.FeeQuoter, c.FeeQuoter != nil},
		{"NonceManager", c.NonceManager, c.NonceManager != nil},
		{"OnRamp", c.OnRamp, c.OnRamp != nil},
		{"OffRamp", c.OffRamp, c.OffRamp != nil},
		{"Router", c.Router, c.Router != nil},
	} {
		if !ct.present {
			continue
		}
		if err := checkOwnership(callOpts, ct.name, ct.contract, timelockAddr); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// V16ActiveChainSelectors returns chain selectors with an active or candidate v1.6 DON config in CCIPHome
func (c CCIPChainState) V16ActiveChainSelectors(ctx context.Context) (map[uint64]bool, error) {
	if c.CCIPHome == nil {
		return nil, errors.New("no CCIPHome contract found in the state")
	}
	if c.CapabilityRegistry == nil {
		return nil, errors.New("no CapabilityRegistry contract found in the state")
	}
	ccipDons, err := shared.GetCCIPDonsFromCapRegistry(ctx, c.CapabilityRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed to get CCIP DONs from capability registry: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx}
	active := make(map[uint64]bool, len(ccipDons))
	for _, don := range ccipDons {
		commitConfigs, err := c.CCIPHome.GetAllConfigs(callOpts, don.Id, uint8(types.PluginTypeCCIPCommit))
		if err != nil {
			continue
		}
		execConfigs, err := c.CCIPHome.GetAllConfigs(callOpts, don.Id, uint8(types.PluginTypeCCIPExec))
		if err != nil {
			continue
		}
		if chainSel := chainSelFromConfigs(commitConfigs, execConfigs); chainSel != 0 {
			active[chainSel] = true
		}
	}
	return active, nil
}

// chainSelFromConfigs extracts the chain selector from CCIPHome configs,
// falling back through active→candidate for both commit and exec
func chainSelFromConfigs(commit, exec ccip_home.GetAllConfigs) uint64 {
	sel := commit.ActiveConfig.Config.ChainSelector
	if sel == 0 {
		sel = commit.CandidateConfig.Config.ChainSelector
	}
	if sel == 0 {
		sel = exec.ActiveConfig.Config.ChainSelector
		if sel == 0 {
			sel = exec.CandidateConfig.Config.ChainSelector
		}
	}
	return sel
}
