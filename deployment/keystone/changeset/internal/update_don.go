package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"google.golang.org/protobuf/proto"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"

	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

// CapabilityConfig is a struct that holds a capability and its configuration
type CapabilityConfig struct {
	Capability kcr.CapabilitiesRegistryCapability
	Config     []byte // this is the marshalled proto config. if nil, a default config is used
}

type UpdateDonRequest struct {
	Registry *kcr.CapabilitiesRegistry
	Chain    deployment.Chain

	P2PIDs            []p2pkey.PeerID    // this is the unique identifier for the don
	CapabilityConfigs []CapabilityConfig // if Config subfield is nil, a default config is used
}

func (r *UpdateDonRequest) appendNodeCapabilitiesRequest() *AppendNodeCapabilitiesRequest {
	out := &AppendNodeCapabilitiesRequest{
		Chain:             r.Chain,
		Registry:          r.Registry,
		P2pToCapabilities: make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability),
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
	if r.Registry == nil {
		return fmt.Errorf("registry is required")
	}
	if len(r.P2PIDs) == 0 {
		return fmt.Errorf("p2pIDs is required")
	}
	return nil
}

type UpdateDonResponse struct {
	DonInfo kcr.CapabilitiesRegistryDONInfo
}

func UpdateDon(lggr logger.Logger, req *UpdateDonRequest) (*UpdateDonResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	getDonsResp, err := req.Registry.GetDONs(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Dons: %w", err)
	}

	don, err := lookupDonByPeerIDs(getDonsResp, req.P2PIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup don by p2pIDs: %w", err)
	}
	cfgs, err := computeConfigs(req.Registry, req.CapabilityConfigs, don)
	if err != nil {
		return nil, fmt.Errorf("failed to compute configs: %w", err)
	}

	_, err = AppendNodeCapabilitiesImpl(lggr, req.appendNodeCapabilitiesRequest())
	if err != nil {
		return nil, fmt.Errorf("failed to append node capabilities: %w", err)
	}

	tx, err := req.Registry.UpdateDON(req.Chain.DeployerKey, don.Id, don.NodeP2PIds, cfgs, don.IsPublic, don.F)
	if err != nil {
		err = kslib.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call UpdateDON: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm UpdateDON transaction %s: %w", tx.Hash().String(), err)
	}
	out := don
	out.CapabilityConfigurations = cfgs
	return &UpdateDonResponse{DonInfo: out}, nil
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

func computeConfigs(registry *kcr.CapabilitiesRegistry, caps []CapabilityConfig, donInfo kcr.CapabilitiesRegistryDONInfo) ([]kcr.CapabilitiesRegistryCapabilityConfiguration, error) {
	out := make([]kcr.CapabilitiesRegistryCapabilityConfiguration, len(caps))
	for i, cap := range caps {
		out[i] = kcr.CapabilitiesRegistryCapabilityConfiguration{}
		id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, cap.Capability.LabelledName, cap.Capability.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to get capability id: %w", err)
		}
		out[i].CapabilityId = id
		if out[i].Config == nil {
			c := kslib.DefaultCapConfig(cap.Capability.CapabilityType, int(donInfo.F))
			cb, err := proto.Marshal(c)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal capability config for %v: %w", c, err)
			}
			out[i].Config = cb
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
