package chainreader_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/chainreader"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestChainReader(t *testing.T) {
	const anyConfig = "config :D"
	wrappedConfig, err := values.NewMap(map[string]interface{}{"anything": anyConfig})
	require.NoError(t, err)

	const anyWorkflowId = "workflowId"
	anyArguments, err := values.NewMap(map[string]interface{}{"key": "value"})
	require.NoError(t, err)

	factory := mockFactory{t, anyConfig}
	reader := chainreader.NewChainReader(factory, "evm", "1337")

	t.Run("Test Info returns the correct structure", func(t *testing.T) {
		actual := reader.Info()
		require.NotNil(t, actual)

		expected, err := capabilities.NewCapabilityInfo(
			"chainreader-evm-1337",
			capabilities.CapabilityTypeAction,
			"Reads from evm 1337",
			"v0.0.1",
		)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Execute makes callbacks with registered capabilities", func(t *testing.T) {
		ctx := testutils.Context(t)

		require.NoError(t, reader.RegisterWorkflow(ctx, anyWorkflowId))
		t.Skip("Not written yet")
	})

	t.Run("Execute proxies errors from the chain reader", func(t *testing.T) {
		t.Skip("Not written yet")
	})

	t.Run("Execute returns an error if the workflow was never registered", func(t *testing.T) {
		t.Skip("Not written yet")
	})

	t.Run("Execute returns an error if the workflow was deregistered", func(t *testing.T) {
		t.Skip("Not written yet")
	})
}

type mockReader struct {
}

var _ types.ChainReader = &mockReader{}

func (m *mockReader) GetLatestValue(ctx context.Context, contractName string, method string, params, returnVal any) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	//TODO implement me
	panic("implement me")
}

type mockFactory struct {
	t              *testing.T
	expectedConfig string
}

func (m mockFactory) NewChainReader(config string) (types.ChainReader, error) {
	assert.Equal(m.t, m.expectedConfig, config)
	return &mockReader{}, nil
}
