package registrysyncer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type capabilitiesRegistryDONInfo struct {
	Id                       uint32                                        `json:"id"`
	ConfigCount              uint32                                        `json:"configCount"`
	F                        uint8                                         `json:"f"`
	IsPublic                 bool                                          `json:"isPublic"`
	AcceptsWorkflows         bool                                          `json:"acceptsWorkflows"`
	NodeP2PIds               []string                                      `json:"nodeP2PIds"`
	CapabilityConfigurations []capabilitiesRegistryCapabilityConfiguration `json:"capabilityConfigurations"`
}

type capabilitiesRegistryCapabilityConfiguration struct {
	CapabilityId string `json:"capabilityId"`
	Config       string `json:"config"`
}

type capabilitiesRegistryNodeInfo struct {
	NodeOperatorId      uint32   `json:"nodeOperatorId"`
	ConfigCount         uint32   `json:"configCount"`
	WorkflowDONId       uint32   `json:"workflowDONId"`
	Signer              string   `json:"signer"`
	P2pId               string   `json:"p2pId"`
	HashedCapabilityIds []string `json:"hashedCapabilityIds"`
	CapabilitiesDONIds  []string `json:"capabilitiesDONIds"`
}

type capabilitiesRegistryCapabilityInfo struct {
	HashedId              string `json:"hashedId"`
	LabelledName          string `json:"labelledName"`
	Version               string `json:"version"`
	CapabilityType        uint8  `json:"capabilityType"`
	ResponseType          uint8  `json:"responseType"`
	ConfigurationContract string `json:"configurationContract"`
	IsDeprecated          bool   `json:"isDeprecated"`
}

func (t *State) MarshalJSON() ([]byte, error) {
	idsToDONs := make(map[string]capabilitiesRegistryDONInfo)
	for k, v := range t.IDsToDONs {
		nodeP2PIds := make([]string, len(v.NodeP2PIds))
		for i, id := range v.NodeP2PIds {
			nodeP2PIds[i] = hex.EncodeToString(id[:])
		}
		configs := make([]capabilitiesRegistryCapabilityConfiguration, len(v.CapabilityConfigurations))
		for i, c := range v.CapabilityConfigurations {
			configs[i] = capabilitiesRegistryCapabilityConfiguration{
				CapabilityId: hex.EncodeToString(c.CapabilityId[:]),
				Config:       hex.EncodeToString(c.Config),
			}
		}
		idsToDONs[fmt.Sprintf("%d", k)] = capabilitiesRegistryDONInfo{
			Id:                       v.Id,
			ConfigCount:              v.ConfigCount,
			F:                        v.F,
			IsPublic:                 v.IsPublic,
			AcceptsWorkflows:         v.AcceptsWorkflows,
			NodeP2PIds:               nodeP2PIds,
			CapabilityConfigurations: configs,
		}
	}

	idsToNodes := make(map[string]capabilitiesRegistryNodeInfo)
	for k, v := range t.IDsToNodes {
		hashedCapabilityIds := make([]string, len(v.HashedCapabilityIds))
		for i, id := range v.HashedCapabilityIds {
			hashedCapabilityIds[i] = hex.EncodeToString(id[:])
		}
		capabilitiesDONIds := make([]string, len(v.CapabilitiesDONIds))
		for i, id := range v.CapabilitiesDONIds {
			capabilitiesDONIds[i] = id.String()
		}
		idsToNodes[fmt.Sprintf("%x", k[:])] = capabilitiesRegistryNodeInfo{
			NodeOperatorId:      v.NodeOperatorId,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              hex.EncodeToString(v.Signer[:]),
			P2pId:               hex.EncodeToString(v.P2pId[:]),
			HashedCapabilityIds: hashedCapabilityIds,
			CapabilitiesDONIds:  capabilitiesDONIds,
		}
	}

	idsToCapabilities := make(map[string]capabilitiesRegistryCapabilityInfo)
	for k, v := range t.IDsToCapabilities {
		idsToCapabilities[fmt.Sprintf("%x", k[:])] = capabilitiesRegistryCapabilityInfo{
			HashedId:              hex.EncodeToString(v.HashedId[:]),
			LabelledName:          v.LabelledName,
			Version:               v.Version,
			CapabilityType:        v.CapabilityType,
			ResponseType:          v.ResponseType,
			ConfigurationContract: v.ConfigurationContract.Hex(),
			IsDeprecated:          v.IsDeprecated,
		}
	}

	return json.Marshal(&struct {
		IDsToDONs         map[string]capabilitiesRegistryDONInfo
		IDsToNodes        map[string]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]capabilitiesRegistryCapabilityInfo
	}{
		IDsToDONs:         idsToDONs,
		IDsToNodes:        idsToNodes,
		IDsToCapabilities: idsToCapabilities,
	})
}

func (t *State) UnmarshalJSON(data []byte) error {
	temp := struct {
		IDsToDONs         map[string]capabilitiesRegistryDONInfo
		IDsToNodes        map[string]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]capabilitiesRegistryCapabilityInfo
	}{
		IDsToDONs:         make(map[string]capabilitiesRegistryDONInfo),
		IDsToNodes:        make(map[string]capabilitiesRegistryNodeInfo),
		IDsToCapabilities: make(map[string]capabilitiesRegistryCapabilityInfo),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	t.IDsToDONs = make(map[DonID]kcr.CapabilitiesRegistryDONInfo)
	for k, v := range temp.IDsToDONs {
		id, err := strconv.ParseUint(k, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse DON ID: %w", err)
		}
		nodeP2PIds := make([][32]byte, len(v.NodeP2PIds))
		for i, p2pid := range v.NodeP2PIds {
			b, err2 := hex.DecodeString(p2pid)
			if err2 != nil {
				return fmt.Errorf("failed to decode nodeP2PId: %w", err2)
			}
			copy(nodeP2PIds[i][:], b[:32])
		}
		configs := make([]kcr.CapabilitiesRegistryCapabilityConfiguration, len(v.CapabilityConfigurations))
		for i, c := range v.CapabilityConfigurations {
			capabilityId, err2 := hex.DecodeString(c.CapabilityId)
			if err2 != nil {
				return fmt.Errorf("failed to decode capabilityId: %w", err2)
			}
			config, err2 := hex.DecodeString(c.Config)
			if err2 != nil {
				return fmt.Errorf("failed to decode capability config: %w", err2)
			}
			configs[i] = kcr.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: to32Byte(capabilityId),
				Config:       config,
			}
		}
		t.IDsToDONs[DonID(id)] = kcr.CapabilitiesRegistryDONInfo{
			Id:                       v.Id,
			ConfigCount:              v.ConfigCount,
			F:                        v.F,
			IsPublic:                 v.IsPublic,
			AcceptsWorkflows:         v.AcceptsWorkflows,
			NodeP2PIds:               nodeP2PIds,
			CapabilityConfigurations: configs,
		}
	}

	t.IDsToNodes = make(map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo)
	for k, v := range temp.IDsToNodes {
		key, err := hex.DecodeString(k)
		if err != nil {
			return fmt.Errorf("failed to decode node key: %w", err)
		}
		var peerID p2ptypes.PeerID
		copy(peerID[:], key[:32])

		hashedCapabilityIds := make([][32]byte, len(v.HashedCapabilityIds))
		for i, id := range v.HashedCapabilityIds {
			b, err2 := hex.DecodeString(id)
			if err2 != nil {
				return fmt.Errorf("failed to decode hashedCapabilityId: %w", err2)
			}
			copy(hashedCapabilityIds[i][:], b[:32])
		}

		capabilitiesDONIds := make([]*big.Int, len(v.CapabilitiesDONIds))
		for i, id := range v.CapabilitiesDONIds {
			bigInt := new(big.Int)
			bigInt.SetString(id, 10)
			capabilitiesDONIds[i] = bigInt
		}
		signer, err := hex.DecodeString(v.Signer)
		if err != nil {
			return fmt.Errorf("failed to decode signer: %w", err)
		}
		p2pId, err := hex.DecodeString(v.P2pId)
		if err != nil {
			return fmt.Errorf("failed to decode p2pId: %w", err)
		}
		t.IDsToNodes[peerID] = kcr.CapabilitiesRegistryNodeInfo{
			NodeOperatorId:      v.NodeOperatorId,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              to32Byte(signer),
			P2pId:               to32Byte(p2pId),
			HashedCapabilityIds: hashedCapabilityIds,
			CapabilitiesDONIds:  capabilitiesDONIds,
		}
	}

	t.IDsToCapabilities = make(map[HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo)
	for k, v := range temp.IDsToCapabilities {
		id, err := hex.DecodeString(k)
		if err != nil {
			return fmt.Errorf("failed to decode capability ID: %w", err)
		}
		var hashedId HashedCapabilityID
		copy(hashedId[:], id[:32])

		t.IDsToCapabilities[hashedId] = kcr.CapabilitiesRegistryCapabilityInfo{
			HashedId:              hashedId,
			LabelledName:          v.LabelledName,
			Version:               v.Version,
			CapabilityType:        v.CapabilityType,
			ResponseType:          v.ResponseType,
			ConfigurationContract: common.HexToAddress(v.ConfigurationContract),
			IsDeprecated:          v.IsDeprecated,
		}
	}

	return nil
}

type syncerORM struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

func newORM(ds sqlutil.DataSource, lggr logger.Logger) syncerORM {
	namedLogger := lggr.Named("RegistrySyncerORM")
	return syncerORM{
		ds:   ds,
		lggr: namedLogger,
	}
}

func (orm syncerORM) addState(ctx context.Context, state State) error {
	stateJSON, err := state.MarshalJSON()
	if err != nil {
		return err
	}
	hash := sha256.Sum256(stateJSON)
	_, err = orm.ds.ExecContext(
		ctx,
		`INSERT INTO registry_syncer_states (data, data_hash) VALUES ($1, $2) ON CONFLICT (data_hash) DO NOTHING`,
		stateJSON, fmt.Sprintf("%x", hash[:]),
	)
	return err
}

func (orm syncerORM) latestState(ctx context.Context) (*State, error) {
	var state State
	var stateJSON string
	err := orm.ds.GetContext(ctx, &stateJSON, `SELECT data FROM registry_syncer_states ORDER BY created_at DESC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	err = state.UnmarshalJSON([]byte(stateJSON))
	if err != nil {
		return nil, err
	}
	return &state, nil
}
