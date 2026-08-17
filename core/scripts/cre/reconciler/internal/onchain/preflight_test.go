package onchain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
)

func TestExpectedJDTypeLabel(t *testing.T) {
	t.Parallel()

	require.Equal(t, "bootstrap", expectedJDTypeLabel(domain.RoleBootstrap))
	require.Equal(t, "gateway", expectedJDTypeLabel(domain.RoleGateway))
	require.Equal(t, "plugin", expectedJDTypeLabel(domain.RoleStandard))
}

func TestValidateNodeLabels_SkipsWhenJDNotConfigured(t *testing.T) {
	t.Parallel()

	desired := &domain.DesiredState{} // JD.GRPC is empty
	cv := &domain.ChartValues{Nodes: []domain.ChartNodeInfo{{Name: "node-0"}}}
	state := &domain.StateFile{}

	err := ValidateNodeLabels(context.Background(), desired, cv, state)
	require.NoError(t, err)
}

func TestValidateNodeLabels_MissingDiscoveredCSAKey(t *testing.T) {
	t.Setenv(infra.JDAccessTokenEnv, "test-token")

	desired := &domain.DesiredState{JD: domain.JDConfig{GRPC: "grpc-jd.example.com:443"}}
	cv := &domain.ChartValues{Nodes: []domain.ChartNodeInfo{
		{Name: "node-0", NodeType: domain.RoleStandard},
	}}
	state := &domain.StateFile{} // no NodeRuntime entry for node-0

	err := ValidateNodeLabels(context.Background(), desired, cv, state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "node-0: no discovered CSA key")
}
