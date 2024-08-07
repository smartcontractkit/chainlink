package registrysyncer

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

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

	b, err := json.Marshal(&struct {
		IDsToDONs         map[DonID]DON
		IDsToNodes        map[p2ptypes.PeerID]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]Capability
	}{
		IDsToDONs:         l.IDsToDONs,
		IDsToNodes:        idsToNodes,
		IDsToCapabilities: l.IDsToCapabilities,
	})
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func (l *LocalRegistry) UnmarshalJSON(data []byte) error {
	temp := struct {
		IDsToDONs         map[DonID]DON
		IDsToNodes        map[p2ptypes.PeerID]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]Capability
	}{
		IDsToDONs:         make(map[DonID]DON),
		IDsToNodes:        make(map[p2ptypes.PeerID]capabilitiesRegistryNodeInfo),
		IDsToCapabilities: make(map[string]Capability),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	l.IDsToDONs = temp.IDsToDONs

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

type ORM interface {
	AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error
	LatestLocalRegistry(ctx context.Context) (*LocalRegistry, error)
}

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource, lggr logger.Logger) orm {
	namedLogger := lggr.Named("RegistrySyncerORM")
	return orm{
		ds:   ds,
		lggr: namedLogger,
	}
}

func (orm orm) AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error {
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

func (orm orm) LatestLocalRegistry(ctx context.Context) (*LocalRegistry, error) {
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
