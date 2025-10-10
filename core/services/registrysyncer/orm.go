package registrysyncer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

type capabilitiesRegistryNodeInfo struct {
	NodeOperatorId      uint32         `json:"nodeOperatorId"`
	ConfigCount         uint32         `json:"configCount"`
	WorkflowDONId       uint32         `json:"workflowDONId"`
	Signer              types.PeerID   `json:"signer"`
	P2pId               types.PeerID   `json:"p2pId"`
	EncryptionPublicKey [32]byte       `json:"encryptionPublicKey"`
	HashedCapabilityIds []types.PeerID `json:"hashedCapabilityIds"`
	CapabilitiesDONIds  []string       `json:"capabilitiesDONIds"`
}

func (l *LocalRegistry) MarshalJSON() ([]byte, error) {
	idsToNodes := make(map[types.PeerID]capabilitiesRegistryNodeInfo)
	for k, v := range l.IDsToNodes {
		hashedCapabilityIDs := make([]types.PeerID, len(v.HashedCapabilityIDs))
		for i, id := range v.HashedCapabilityIDs {
			hashedCapabilityIDs[i] = types.PeerID(id[:])
		}
		capabilitiesDONIds := make([]string, len(v.CapabilitiesDONIds))
		for i, id := range v.CapabilitiesDONIds {
			capabilitiesDONIds[i] = id.String()
		}
		idsToNodes[k] = capabilitiesRegistryNodeInfo{
			NodeOperatorId:      v.NodeOperatorID,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              types.PeerID(v.Signer[:]),
			P2pId:               types.PeerID(v.P2pID[:]),
			EncryptionPublicKey: v.EncryptionPublicKey,
			HashedCapabilityIds: hashedCapabilityIDs,
			CapabilitiesDONIds:  capabilitiesDONIds,
		}
	}

	b, err := json.Marshal(&struct {
		IDsToDONs         map[DonID]DON
		IDsToNodes        map[types.PeerID]capabilitiesRegistryNodeInfo
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
		IDsToNodes        map[types.PeerID]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]Capability
	}{
		IDsToDONs:         make(map[DonID]DON),
		IDsToNodes:        make(map[types.PeerID]capabilitiesRegistryNodeInfo),
		IDsToCapabilities: make(map[string]Capability),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	l.IDsToDONs = temp.IDsToDONs

	l.IDsToNodes = make(map[types.PeerID]NodeInfo)
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
		l.IDsToNodes[peerID] = NodeInfo{
			NodeOperatorID:      v.NodeOperatorId,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              v.Signer,
			P2pID:               v.P2pId,
			EncryptionPublicKey: v.EncryptionPublicKey,
			HashedCapabilityIDs: hashedCapabilityIds,
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
	namedLogger := logger.Named(lggr, "RegistrySyncerORM")
	return orm{
		ds:   ds,
		lggr: namedLogger,
	}
}

func (orm orm) AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error {
	orm.lggr.Debugw("Adding local registry to DB...")
	return sqlutil.TransactDataSource(ctx, orm.ds, nil, func(tx sqlutil.DataSource) error {
		localRegistryJSON, err := localRegistry.MarshalJSON()
		if err != nil {
			return err
		}
		hash := sha256.Sum256(localRegistryJSON)
		// update if and only if the hash does not match the latest value
		r, err := tx.ExecContext(
			ctx,
			`INSERT INTO registry_syncer_states (data, data_hash) 
            SELECT $1, $2 
            WHERE $2 NOT IN (
                SELECT data_hash FROM registry_syncer_states 
                ORDER BY id DESC LIMIT 1
            )`,
			localRegistryJSON, hex.EncodeToString(hash[:]),
		)
		if err != nil {
			return fmt.Errorf("failed to insert into registry_syncer: %w", err)
		}

		n, _ := r.RowsAffected()
		if n != 0 {
			id, _ := r.LastInsertId()
			orm.lggr.Debugw("Inserted new local registry", "id", id, "hash", hex.EncodeToString(hash[:]), "registry", localRegistry)
		} else {
			orm.lggr.Debugw("No rows affected, local registry updated. ", "hash", hex.EncodeToString(hash[:]))
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
	hash := sha256.Sum256([]byte(localRegistryJSON))
	err = localRegistry.UnmarshalJSON([]byte(localRegistryJSON))
	if err != nil {
		return nil, err
	}
	orm.lggr.Debugw("Fetched latest local registry from DB", "hash", hex.EncodeToString(hash[:]), "registry", localRegistry)

	return &localRegistry, nil
}
