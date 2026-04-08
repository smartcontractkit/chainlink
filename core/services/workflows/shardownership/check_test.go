package shardownership

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

type stubClient struct {
	resp *ringpb.GetWorkflowShardMappingResponse
	err  error
}

func (s *stubClient) GetWorkflowShardMapping(_ context.Context, _ []string) (*ringpb.GetWorkflowShardMappingResponse, error) {
	return s.resp, s.err
}

func (s *stubClient) ReportWorkflowTriggerRegistration(context.Context, *ringpb.ReportWorkflowTriggerRegistrationRequest) (*ringpb.ReportWorkflowTriggerRegistrationResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *stubClient) Close() error { return nil }

func TestCheckCommittedOwner(t *testing.T) {
	ctx := t.Context()
	wf := "abc123"

	t.Run("allow when mapped shard matches", func(t *testing.T) {
		c := &stubClient{resp: &ringpb.GetWorkflowShardMappingResponse{Mappings: map[string]uint32{wf: 2}}}
		v, _, err := CheckCommittedOwner(ctx, c, wf, 2)
		require.NoError(t, err)
		require.Equal(t, Allow, v)
	})

	t.Run("deny not owner when mapped to other shard", func(t *testing.T) {
		c := &stubClient{resp: &ringpb.GetWorkflowShardMappingResponse{Mappings: map[string]uint32{wf: 1}}}
		v, _, err := CheckCommittedOwner(ctx, c, wf, 2)
		require.NoError(t, err)
		require.Equal(t, DenyNotOwner, v)
	})

	t.Run("deny not owner when workflow missing from map", func(t *testing.T) {
		c := &stubClient{resp: &ringpb.GetWorkflowShardMappingResponse{Mappings: map[string]uint32{}}}
		v, _, err := CheckCommittedOwner(ctx, c, wf, 2)
		require.NoError(t, err)
		require.Equal(t, DenyNotOwner, v)
	})

	t.Run("deny orchestrator error", func(t *testing.T) {
		c := &stubClient{err: errors.New("rpc down")}
		v, _, err := CheckCommittedOwner(ctx, c, wf, 2)
		require.Error(t, err)
		require.Equal(t, DenyOrchestratorError, v)
	})
}
