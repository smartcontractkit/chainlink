package generic

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type RelayGetter interface {
	GetIDToRelayerMap() map[types.RelayID]loop.Relayer
}

type RelayerSet struct {
	wrappedRelayers map[types.RelayID]core.Relayer
}

func NewRelayerSet(relayGetter RelayGetter, externalJobID uuid.UUID, jobID int32, isNew bool) (*RelayerSet, error) {
	wrappedRelayers := map[types.RelayID]core.Relayer{}

	relayers := relayGetter.GetIDToRelayerMap()

	for id, relayer := range relayers {
		wrappedRelayers[id] = relayerWrapper{Relayer: relayer, ExternalJobID: externalJobID, JobID: jobID, New: isNew}
	}

	return &RelayerSet{wrappedRelayers: wrappedRelayers}, nil
}

func (r *RelayerSet) Get(_ context.Context, id types.RelayID) (core.Relayer, error) {
	if relayer, ok := r.wrappedRelayers[id]; ok {
		return relayer, nil
	}

	return nil, fmt.Errorf("relayer with id %s not found", id)
}

func (r *RelayerSet) List(_ context.Context, relayIDs ...types.RelayID) (map[types.RelayID]core.Relayer, error) {
	if len(relayIDs) == 0 {
		return r.wrappedRelayers, nil
	}

	filterer := map[types.RelayID]bool{}
	for _, id := range relayIDs {
		filterer[id] = true
	}

	result := map[types.RelayID]core.Relayer{}
	for id, relayer := range r.wrappedRelayers {
		if _, ok := filterer[id]; ok {
			result[id] = relayer
		}
	}

	return result, nil
}

type relayerWrapper struct {
	loop.Relayer
	ExternalJobID uuid.UUID
	JobID         int32
	New           bool // Whether this is a first time job add.
}

func (r relayerWrapper) NewPluginProvider(ctx context.Context, rargs core.RelayArgs, pargs core.PluginArgs) (core.PluginProvider, error) {
	relayArgs := types.RelayArgs{
		ExternalJobID:      r.ExternalJobID,
		JobID:              r.JobID,
		ContractID:         rargs.ContractID,
		New:                r.New,
		RelayConfig:        rargs.RelayConfig,
		ProviderType:       rargs.ProviderType,
		MercuryCredentials: rargs.MercuryCredentials,
	}

	return r.Relayer.NewPluginProvider(ctx, relayArgs, types.PluginArgs(pargs))
}
