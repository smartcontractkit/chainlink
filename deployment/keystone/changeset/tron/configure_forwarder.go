package tron

import (
	"context"
	"errors"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSetV2[*ConfigureForwarderRequest] = ConfigureForwarder{}

type ConfigureForwarder struct{}

func (cs ConfigureForwarder) VerifyPreconditions(env cldf.Environment, req *ConfigureForwarderRequest) error {
	if len(req.WFNodeIDs) == 0 {
		return errors.New("WFNodeIDs must not be empty")
	}
	return nil
}

type ConfigureForwarderRequest struct {
	WFDonName        string
	WFNodeIDs        []string
	RegistryChainSel uint64
	Chains           map[uint64]struct{}
	TriggerOptions   *cldf_tron.TriggerOptions
}

func (cs ConfigureForwarder) Apply(env cldf.Environment, req *ConfigureForwarderRequest) (cldf.ChangesetOutput, error) {
	wfDon, err := internal.NewRegisteredDon(env, internal.RegisteredDonConfig{
		NodeIDs:          req.WFNodeIDs,
		Name:             req.WFDonName,
		RegistryChainSel: req.RegistryChainSel,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create registered don: %w", err)
	}

	return cldf.ChangesetOutput{}, configureForwarderContracts(env, req, wfDon)
}

func configureForwarderContracts(env cldf.Environment, req *ConfigureForwarderRequest, wfdon *internal.RegisteredDon) error {
	tronChains := env.BlockChains.TronChains()
	contractSetsResp, err := LoadTronOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to get contract sets: %w", err)
	}

	for _, chain := range tronChains {
		if _, shouldInclude := req.Chains[chain.Selector]; len(req.Chains) > 0 && !shouldInclude {
			continue
		}
		// get the forwarder contract for the chain
		contracts, ok := contractSetsResp.Chains[chain.Selector]
		if !ok {
			return fmt.Errorf("failed to get contract set for chain %d", chain.Selector)
		}
		err := configureForwarder(env.Logger, chain, contracts.Forwarder, []internal.RegisteredDon{*wfdon}, req.TriggerOptions)
		if err != nil {
			return fmt.Errorf("failed to configure forwarder for chain selector %d: %w", chain.Selector, err)
		}
	}

	return nil
}

// determineTronChainFamily checks what chain family the DON nodes are configured with
// Returns FamilyTron for native Tron configs, FamilyEVM for EVM configs with ChainType='tron'
func determineTronChainFamily(dn internal.RegisteredDon) string {
	hasTronFamily := false

	// Check all nodes to see what chain families they support
	for _, node := range dn.Nodes {
		for details := range node.SelToOCRConfig {
			if family, err := chainsel.GetSelectorFamily(details.ChainSelector); err == nil {
				if family == chainsel.FamilyTron {
					hasTronFamily = true
				}
			}
		}
	}

	// Prefer native Tron if available, fall back to EVM
	if hasTronFamily {
		return chainsel.FamilyTron
	}
	return chainsel.FamilyEVM
}

func configureForwarder(lggr logger.Logger, chain cldf_tron.Chain, fwdrAddress address.Address, dons []internal.RegisteredDon, triggerOpts *cldf_tron.TriggerOptions) error {
	if fwdrAddress == nil {
		return errors.New("nil forwarder contract")
	}

	for _, dn := range dons {
		if !dn.Info.AcceptsWorkflows {
			continue
		}
		ver := dn.Info.ConfigCount // note config count on the don info is the version on the forwarder

		// Check which chain family is available for backward compatibility
		// Nodes might be configured as native Tron or as EVM chains with ChainType='tron'
		chainFamily := determineTronChainFamily(dn)
		signers := dn.Signers(chainFamily)

		txInfo, err := chain.TriggerContractAndConfirm(context.Background(), fwdrAddress, "setConfig(uint32,uint32,uint8,address[])", []any{"uint32", dn.Info.Id, "uint32", ver, "uint8", dn.Info.F, "address[]", signers}, triggerOpts)
		if err != nil {
			return fmt.Errorf("failed to setConfig for donId %d, err: %w", dn.Info.Id, err)
		}

		lggr.Infof("Configured donId %d, txInfo: %+v", dn.Info.Id, txInfo)
	}

	return nil
}
