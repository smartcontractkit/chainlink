package registrysyncer

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type capabilitiesDON struct {
	capabilities.DON
	CapabilityConfigurations map[string]capabilityConfiguration
}

type capabilityConfiguration struct {
	DefaultConfig       values.Map
	RemoteTriggerConfig capabilities.RemoteTriggerConfig
	RemoteTargetConfig  capabilities.RemoteTargetConfig
}

type capabilitiesRegistryNodeInfo struct {
	NodeOperatorId      uint32            `json:"nodeOperatorId"`
	ConfigCount         uint32            `json:"configCount"`
	WorkflowDONId       uint32            `json:"workflowDONId"`
	Signer              p2ptypes.PeerID   `json:"signer"`
	P2pId               p2ptypes.PeerID   `json:"p2pId"`
	HashedCapabilityIds []p2ptypes.PeerID `json:"hashedCapabilityIds"`
	CapabilitiesDONIds  []string          `json:"capabilitiesDONIds"`
}

func (l *LocalRegistry) MarshalJSON() ([]byte, error) {
	idsToDONs := make(map[DonID]capabilitiesDON)
	for donID, don := range l.IDsToDONs {
		capabilityConfigurations := make(map[string]capabilityConfiguration)
		for k, v := range don.CapabilityConfigurations {
			cfg := capabilityConfiguration{
				DefaultConfig:       *values.EmptyMap(),
				RemoteTriggerConfig: capabilities.RemoteTriggerConfig{},
				RemoteTargetConfig:  capabilities.RemoteTargetConfig{},
			}
			if v.DefaultConfig != nil {
				cfg.DefaultConfig = *v.DefaultConfig
			}
			if v.RemoteTriggerConfig != nil {
				cfg.RemoteTriggerConfig = *v.RemoteTriggerConfig
			}
			if v.RemoteTargetConfig != nil {
				cfg.RemoteTargetConfig = *v.RemoteTargetConfig
			}
			capabilityConfigurations[k] = cfg
		}
		idsToDONs[donID] = capabilitiesDON{
			DON:                      don.DON,
			CapabilityConfigurations: capabilityConfigurations,
		}
	}

	idsToNodes := make(map[p2ptypes.PeerID]capabilitiesRegistryNodeInfo)
	for k, v := range l.IDsToNodes {
		hashedCapabilityIds := make([]p2ptypes.PeerID, len(v.HashedCapabilityIds))
		for i, id := range v.HashedCapabilityIds {
			hashedCapabilityIds[i] = p2ptypes.PeerID(id[:])
		}
		capabilitiesDONIds := make([]string, len(v.CapabilitiesDONIds))
		for i, id := range v.CapabilitiesDONIds {
			capabilitiesDONIds[i] = id.String()
		}
		idsToNodes[k] = capabilitiesRegistryNodeInfo{
			NodeOperatorId:      v.NodeOperatorId,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              p2ptypes.PeerID(v.Signer[:]),
			P2pId:               p2ptypes.PeerID(v.P2pId[:]),
			HashedCapabilityIds: hashedCapabilityIds,
			CapabilitiesDONIds:  capabilitiesDONIds,
		}
	}

	return json.Marshal(&struct {
		IDsToDONs         map[DonID]capabilitiesDON
		IDsToNodes        map[p2ptypes.PeerID]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]Capability
	}{
		IDsToDONs:         idsToDONs,
		IDsToNodes:        idsToNodes,
		IDsToCapabilities: l.IDsToCapabilities,
	})
}

func (l *LocalRegistry) UnmarshalJSON(data []byte) error {
	temp := struct {
		IDsToDONs         map[DonID]capabilitiesDON
		IDsToNodes        map[p2ptypes.PeerID]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]Capability
	}{
		IDsToDONs:         make(map[DonID]capabilitiesDON),
		IDsToNodes:        make(map[p2ptypes.PeerID]capabilitiesRegistryNodeInfo),
		IDsToCapabilities: make(map[string]Capability),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	l.IDsToDONs = make(map[DonID]DON)
	for donID, v := range temp.IDsToDONs {
		capabilityConfigurations := make(map[string]capabilities.CapabilityConfiguration)
		for k, v := range v.CapabilityConfigurations {
			capabilityConfigurations[k] = capabilities.CapabilityConfiguration{
				DefaultConfig:       &v.DefaultConfig,
				RemoteTriggerConfig: &v.RemoteTriggerConfig,
				RemoteTargetConfig:  &v.RemoteTargetConfig,
			}
		}
		l.IDsToDONs[donID] = DON{
			DON:                      v.DON,
			CapabilityConfigurations: capabilityConfigurations,
		}
	}

	l.IDsToNodes = make(map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo)
	for peerID, v := range temp.IDsToNodes {
		hashedCapabilityIds := make([][32]byte, len(v.HashedCapabilityIds))
		for i, id := range v.HashedCapabilityIds {
			copy(hashedCapabilityIds[i][:], id[:])
		}

		capabilitiesDONIds := make([]*big.Int, len(v.CapabilitiesDONIds))
		for i, id := range v.CapabilitiesDONIds {
			bigInt := new(big.Int)
			bigInt.SetString(id, 10)
			capabilitiesDONIds[i] = bigInt
		}
		l.IDsToNodes[peerID] = kcr.CapabilitiesRegistryNodeInfo{
			NodeOperatorId:      v.NodeOperatorId,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              v.Signer,
			P2pId:               v.P2pId,
			HashedCapabilityIds: hashedCapabilityIds,
			CapabilitiesDONIds:  capabilitiesDONIds,
		}
	}

	l.IDsToCapabilities = temp.IDsToCapabilities

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
