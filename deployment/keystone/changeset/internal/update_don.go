package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// CapabilityConfig is a struct that holds a capability and its configuration
type CapabilityConfig struct {
	Capability kcr.CapabilitiesRegistryCapability
	Config     []byte // this is the marshalled proto config. if nil, a default config is used
}

type UpdateDonRequest struct {
	Chain                cldf_evm.Chain
	CapabilitiesRegistry *kcr.CapabilitiesRegistry

	P2PIDs            []p2pkey.PeerID    // this is the unique identifier for the don
	CapabilityConfigs []CapabilityConfig // if Config subfield is nil, a default config is used

	UseMCMS bool
}

func (r *UpdateDonRequest) AppendNodeCapabilitiesRequest() *AppendNodeCapabilitiesRequest {
	out := &AppendNodeCapabilitiesRequest{
		Chain:                r.Chain,
		CapabilitiesRegistry: r.CapabilitiesRegistry,
		P2pToCapabilities:    make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability),
		UseMCMS:              r.UseMCMS,
	}
	for _, p2pid := range r.P2PIDs {
		if _, exists := out.P2pToCapabilities[p2pid]; !exists {
			out.P2pToCapabilities[p2pid] = make([]kcr.CapabilitiesRegistryCapability, 0)
		}
		for _, cc := range r.CapabilityConfigs {
			out.P2pToCapabilities[p2pid] = append(out.P2pToCapabilities[p2pid], cc.Capability)
		}
	}
	return out
}

func (r *UpdateDonRequest) Validate() error {
	if r.CapabilitiesRegistry == nil {
		return errors.New("registry is required")
	}
	if len(r.P2PIDs) == 0 {
		return errors.New("p2pIDs is required")
	}
	return nil
}

type UpdateDonResponse struct {
	DonInfo kcr.CapabilitiesRegistryDONInfo
	Ops     *types.BatchOperation
}

func UpdateDon(_ logger.Logger, req *UpdateDonRequest) (*UpdateDonResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	registry := req.CapabilitiesRegistry
	getDonsResp, err := registry.GetDONs(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Dons: %w", err)
	}

	don, err := lookupDonByPeerIDs(getDonsResp, req.P2PIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup don by p2pIDs: %w", err)
	}

	if don.AcceptsWorkflows {
		// TODO: CRE-277 ensure forwarders are support the next DON version
		// https://github.com/smartcontractkit/chainlink/blob/4fc61bb156fe57bfd939b836c02c413ad1209ebb/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L812
		// and
		// https://github.com/smartcontractkit/chainlink/blob/4fc61bb156fe57bfd939b836c02c413ad1209ebb/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L274
		return nil, fmt.Errorf("refusing to update workflow don %d at config version %d because we cannot validate that all forwarder contracts are ready to accept the new configure version", don.Id, don.ConfigCount)
	}
	cfgs, err := computeConfigs(registry, req.CapabilityConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to compute configs: %w", err)
	}

	txOpts := req.Chain.DeployerKey
	if req.UseMCMS {
		txOpts = cldf.SimTransactOpts()
	}
	tx, err := registry.UpdateDON(txOpts, don.Id, don.NodeP2PIds, cfgs, don.IsPublic, don.F)
	if err != nil {
		err = cldf.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call UpdateDON: %w", err)
	}
	var ops types.BatchOperation
	if !req.UseMCMS {
		_, err = req.Chain.Confirm(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm UpdateDON transaction %s: %w", tx.Hash().String(), err)
		}
	} else {
		ops, err = proposalutils.BatchOperationForChain(req.Chain.Selector, registry.Address().Hex(), tx.Data(), big.NewInt(0), string(CapabilitiesRegistry), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create batch operation: %w", err)
		}
	}

	out := don
	out.CapabilityConfigurations = cfgs
	return &UpdateDonResponse{DonInfo: out, Ops: &ops}, nil
}

func PeerIDsToBytes(p2pIDs []p2pkey.PeerID) [][32]byte {
	out := make([][32]byte, len(p2pIDs))
	for i, p2pID := range p2pIDs {
		out[i] = p2pID
	}
	return out
}

func BytesToPeerIDs(p2pIDs [][32]byte) []p2pkey.PeerID {
	out := make([]p2pkey.PeerID, len(p2pIDs))
	for i, p2pID := range p2pIDs {
		out[i] = p2pID
	}
	return out
}

func computeConfigs(registry *kcr.CapabilitiesRegistry, capCfgs []CapabilityConfig) ([]kcr.CapabilitiesRegistryCapabilityConfiguration, error) {
	out := make([]kcr.CapabilitiesRegistryCapabilityConfiguration, len(capCfgs))
	for i, capCfg := range capCfgs {
		out[i] = kcr.CapabilitiesRegistryCapabilityConfiguration{}
		id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, capCfg.Capability.LabelledName, capCfg.Capability.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to get capability id: %w", err)
		}
		out[i].CapabilityId = id
		out[i].Config = capCfg.Config
		if out[i].Config == nil {
			return nil, fmt.Errorf("config is required for capability %s", capCfg.Capability.LabelledName)
		}
	}
	return out, nil
}

func SortedHash(p2pids [][32]byte) string {
	sha256Hash := sha256.New()
	sort.Slice(p2pids, func(i, j int) bool {
		return bytes.Compare(p2pids[i][:], p2pids[j][:]) < 0
	})
	for _, id := range p2pids {
		sha256Hash.Write(id[:])
	}
	return hex.EncodeToString(sha256Hash.Sum(nil))
}

func lookupDonByPeerIDs(donResp []kcr.CapabilitiesRegistryDONInfo, wanted []p2pkey.PeerID) (kcr.CapabilitiesRegistryDONInfo, error) {
	var don kcr.CapabilitiesRegistryDONInfo
	wantedDonID := SortedHash(PeerIDsToBytes(wanted))
	found := false
	for i, di := range donResp {
		gotID := SortedHash(di.NodeP2PIds)
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

func verboseDonNotFound(donResp []kcr.CapabilitiesRegistryDONInfo, wanted []p2pkey.PeerID) error {
	type debugDonInfo struct {
		OnchainID  uint32
		P2PIDsHash string
		Want       []p2pkey.PeerID
		Got        []p2pkey.PeerID
	}
	debugIds := make([]debugDonInfo, len(donResp))
	for i, di := range donResp {
		debugIds[i] = debugDonInfo{
			OnchainID:  di.Id,
			P2PIDsHash: SortedHash(di.NodeP2PIds),
			Want:       wanted,
			Got:        BytesToPeerIDs(di.NodeP2PIds),
		}
	}
	wantedID := SortedHash(PeerIDsToBytes(wanted))
	b, err2 := json.Marshal(debugIds)
	if err2 == nil {
		return fmt.Errorf("don not found by p2pIDs %s in %s", wantedID, b)
	}
	return fmt.Errorf("don not found by p2pIDs %s in %v", wantedID, debugIds)
}
