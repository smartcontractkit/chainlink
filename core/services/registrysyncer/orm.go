package registrysyncer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type capabilitiesRegistryDON struct {
	ID               uint32   `json:"id"`
	ConfigVersion    uint32   `json:"configVersion"`
	F                uint8    `json:"f"`
	IsPublic         bool     `json:"isPublic"`
	AcceptsWorkflows bool     `json:"acceptsWorkflows"`
	Members          []string `json:"members"`
}

type capabilitiesRegistryDONInfo struct {
	capabilitiesRegistryDON
	CapabilityConfigurations map[string]capabilities.CapabilityConfiguration `json:"capabilityConfigurations"`
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
	ID             string `json:"id"`
	CapabilityType int    `json:"capabilityType"`
}

func (l *LocalRegistry) MarshalJSON() ([]byte, error) {
	idsToDONs := make(map[string]capabilitiesRegistryDONInfo)
	for k, v := range l.IDsToDONs {
		members := make([]string, len(v.Members))
		for i, id := range v.Members {
			members[i] = hex.EncodeToString(id[:])
		}
		configs := make(map[string]capabilities.CapabilityConfiguration, len(v.CapabilityConfigurations))
		for i, c := range v.CapabilityConfigurations {
			configs[i] = capabilities.CapabilityConfiguration{
				DefaultConfig:       c.DefaultConfig,
				RemoteTriggerConfig: c.RemoteTriggerConfig,
			}
		}
		idsToDONs[fmt.Sprintf("%d", k)] = capabilitiesRegistryDONInfo{
			capabilitiesRegistryDON: capabilitiesRegistryDON{
				ID:               v.ID,
				ConfigVersion:    v.ConfigVersion,
				F:                v.F,
				IsPublic:         v.IsPublic,
				AcceptsWorkflows: v.AcceptsWorkflows,
				Members:          members,
			},
			CapabilityConfigurations: configs,
		}
	}

	idsToNodes := make(map[string]capabilitiesRegistryNodeInfo)
	for k, v := range l.IDsToNodes {
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
	for k, v := range l.IDsToCapabilities {
		idsToCapabilities[k] = capabilitiesRegistryCapabilityInfo{
			ID:             v.ID,
			CapabilityType: int(v.CapabilityType),
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

func (l *LocalRegistry) UnmarshalJSON(data []byte) error {
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

	l.IDsToDONs = make(map[DonID]DON)
	for k, v := range temp.IDsToDONs {
		id, err := strconv.ParseUint(k, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse DON ID: %w", err)
		}
		members := make([]types.PeerID, len(v.Members))
		for i, p2pid := range v.Members {
			b, err2 := hex.DecodeString(p2pid)
			if err2 != nil {
				return fmt.Errorf("failed to decode nodeP2PId: %w", err2)
			}
			copy(members[i][:], b[:32])
		}
		configs := make(map[string]capabilities.CapabilityConfiguration, len(v.CapabilityConfigurations))
		for i, c := range v.CapabilityConfigurations {
			configs[i] = capabilities.CapabilityConfiguration{
				DefaultConfig:       c.DefaultConfig,
				RemoteTriggerConfig: c.RemoteTriggerConfig,
			}
		}
		l.IDsToDONs[DonID(id)] = DON{
			DON: capabilities.DON{
				ID:               v.ID,
				ConfigVersion:    v.ConfigVersion,
				F:                v.F,
				IsPublic:         v.IsPublic,
				AcceptsWorkflows: v.AcceptsWorkflows,
				Members:          members,
			},
			CapabilityConfigurations: configs,
		}
	}

	l.IDsToNodes = make(map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo)
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
		l.IDsToNodes[peerID] = kcr.CapabilitiesRegistryNodeInfo{
			NodeOperatorId:      v.NodeOperatorId,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              to32Byte(signer),
			P2pId:               to32Byte(p2pId),
			HashedCapabilityIds: hashedCapabilityIds,
			CapabilitiesDONIds:  capabilitiesDONIds,
		}
	}

	l.IDsToCapabilities = make(map[string]Capability)
	for k, v := range temp.IDsToCapabilities {
		l.IDsToCapabilities[k] = Capability{
			ID:             k,
			CapabilityType: capabilities.CapabilityType(v.CapabilityType),
		}
	}

	return nil
}

func to32Byte(slice []byte) [32]byte {
	var b [32]byte
	copy(b[:], slice[:32])
	return b
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

func (orm syncerORM) addLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error {
	return sqlutil.TransactDataSource(ctx, orm.ds, nil, func(tx sqlutil.DataSource) error {
		localRegistryJSON, err := localRegistry.MarshalJSON()
		if err != nil {
			return err
		}
		hash := sha256.Sum256(localRegistryJSON)
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO registry_syncer_states (data, data_hash) VALUES ($1, $2) ON CONFLICT (data_hash) DO NOTHING`,
			localRegistryJSON, fmt.Sprintf("%x", hash[:]),
		)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `DELETE FROM registry_syncer_states
WHERE data_hash NOT IN (
    SELECT data_hash FROM registry_syncer_states
    ORDER BY id DESC
    LIMIT 10
);`)
		return err
	})
}

func (orm syncerORM) latestLocalRegistry(ctx context.Context) (*LocalRegistry, error) {
	var localRegistry LocalRegistry
	var localRegistryJSON string
	err := orm.ds.GetContext(ctx, &localRegistryJSON, `SELECT data FROM registry_syncer_states ORDER BY id DESC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	err = localRegistry.UnmarshalJSON([]byte(localRegistryJSON))
	if err != nil {
		return nil, err
	}
	return &localRegistry, nil
}
