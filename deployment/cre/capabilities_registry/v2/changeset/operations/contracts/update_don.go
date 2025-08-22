package contracts

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
)

type UpdateDONDeps struct {
	Env                  *cldf.Environment
	CapabilitiesRegistry *capabilities_registry_v2.CapabilitiesRegistry
}

type UpdateDONInput struct {
	ChainSelector uint64

	// P2PIDs are the peer ids that compose the don
	P2PIDs            []p2pkey.PeerID
	CapabilityConfigs []CapabilityConfig

	// DonName to update, this is required
	DonName string

	// F is the fault tolerance level
	// if omitted, the existing value fetched from the registry is used
	F uint8

	// IsPrivate indicates whether the DON is public or private
	// If omitted, the existing value fetched from the registry is used
	IsPrivate bool

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool
}

func (r *UpdateDONInput) Validate() error {
	if r.DonName == "" {
		return errors.New("must specify DONName")
	}

	return nil
}

type UpdateDONOutput struct {
	DonInfo capabilities_registry_v2.CapabilitiesRegistryDONInfo
}

// CapabilityConfig is a struct that holds a capability and its configuration
type CapabilityConfig struct {
	Capability capabilities_registry_v2.CapabilitiesRegistryCapability
	Config     []byte // this is the marshalled proto config. if nil, a default config is used
}

var UpdateDON = operations.NewOperation[UpdateDONInput, UpdateDONOutput, UpdateDONDeps](
	"update-don-op",
	semver.MustParse("1.0.0"),
	"Update DON in Capabilities Registry",
	func(b operations.Bundle, deps UpdateDONDeps, input UpdateDONInput) (UpdateDONOutput, error) {
		err := input.Validate()
		if err != nil {
			return UpdateDONOutput{}, err
		}

		registry := deps.CapabilitiesRegistry
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return UpdateDONOutput{}, cldf.ErrChainNotFound
		}

		var don capabilities_registry_v2.CapabilitiesRegistryDONInfo
		if input.DonName != "" {
			don, err = registry.GetDONByName(&bind.CallOpts{}, input.DonName)
		} else {
			getDonsResp, err := registry.GetDONs(&bind.CallOpts{})
			if err != nil {
				return UpdateDONOutput{}, fmt.Errorf("failed to get Dons: %w", err)
			}

			don, err = lookupDonByPeerIDs(getDonsResp, input.P2PIDs)
			if err != nil {
				return UpdateDONOutput{}, fmt.Errorf("failed to lookup don by p2pIDs: %w", err)
			}
		}

		if don.AcceptsWorkflows && !input.Force {
			// TODO: CRE-277 ensure forwarders are support the next DON version
			// https://github.com/smartcontractkit/chainlink/blob/4fc61bb156fe57bfd939b836c02c413ad1209ebb/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L812
			// and
			// https://github.com/smartcontractkit/chainlink/blob/4fc61bb156fe57bfd939b836c02c413ad1209ebb/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L274
			return UpdateDONOutput{}, fmt.Errorf("refusing to update workflow don %d at config version %d because we cannot validate that all forwarder contracts are ready to accept the new configure version", don.Id, don.ConfigCount)
		}

		cfgs, err := computeConfigs(input.CapabilityConfigs)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to compute configs: %w", err)
		}

		txOpts := chain.DeployerKey

		f := input.F
		if f == 0 {
			f = don.F
		}
		// this is implement as such to maintain backwards compatibility; the default (omitted) value of a bool is false
		var isPublic bool
		if input.IsPrivate {
			isPublic = false
		} else {
			isPublic = don.IsPublic
		}

		tx, err := registry.UpdateDON(txOpts, don.Id, capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams{
			Nodes:                    pkg.PeerIDsToBytes(input.P2PIDs),
			CapabilityConfigurations: cfgs,
			IsPublic:                 isPublic,
			F:                        f,
		})
		if err != nil {
			err = cldf.DecodeErr(kcr.CapabilitiesRegistryABI, err)
			return UpdateDONOutput{}, fmt.Errorf("failed to call UpdateDON: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to confirm UpdateDON transaction %s: %w", tx.Hash().String(), err)
		}

		ctx := b.GetContext()
		_, err = bind.WaitMined(ctx, chain.Client, tx)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to mine UpdateDON confirm transaction %s: %w", tx.Hash().String(), err)
		}

		return UpdateDONOutput{
			DonInfo: don,
		}, nil
	},
)

func lookupDonByPeerIDs(donResp []capabilities_registry_v2.CapabilitiesRegistryDONInfo, wanted []p2pkey.PeerID) (capabilities_registry_v2.CapabilitiesRegistryDONInfo, error) {
	var don capabilities_registry_v2.CapabilitiesRegistryDONInfo
	wantedDonID := pkg.SortedHash(pkg.PeerIDsToBytes(wanted))
	found := false
	for i, di := range donResp {
		gotID := pkg.SortedHash(di.NodeP2PIds)
		if gotID == wantedDonID {
			don = donResp[i]
			found = true
			break
		}
	}
	if !found {
		return don, verboseDonNotFound(donResp, wanted)
	}
	return don, nil
}

func verboseDonNotFound(donResp []capabilities_registry_v2.CapabilitiesRegistryDONInfo, wanted []p2pkey.PeerID) error {
	type debugDonInfo struct {
		OnchainID  uint32
		P2PIDsHash string
		Want       []p2pkey.PeerID
		Got        []p2pkey.PeerID
	}
	debugIDs := make([]debugDonInfo, len(donResp))
	for i, di := range donResp {
		debugIDs[i] = debugDonInfo{
			OnchainID:  di.Id,
			P2PIDsHash: pkg.SortedHash(di.NodeP2PIds),
			Want:       wanted,
			Got:        pkg.BytesToPeerIDs(di.NodeP2PIds),
		}
	}
	wantedID := pkg.SortedHash(pkg.PeerIDsToBytes(wanted))
	b, err2 := json.Marshal(debugIDs)
	if err2 == nil {
		return fmt.Errorf("don not found by p2pIDs %s in %s", wantedID, b)
	}
	return fmt.Errorf("don not found by p2pIDs %s in %v", wantedID, debugIDs)
}

func computeConfigs(capCfgs []CapabilityConfig) ([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, error) {
	out := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, len(capCfgs))
	for i, capCfg := range capCfgs {
		out[i] = capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{}
		out[i].CapabilityId = capCfg.Capability.CapabilityId
		out[i].Config = capCfg.Config
		if out[i].Config == nil {
			return nil, fmt.Errorf("config is required for capability %s", capCfg.Capability.CapabilityId)
		}
	}
	return out, nil
}
