package example

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

func TestUnit_TerminalErrorExampleChangeset(t *testing.T) {
	env, err := environment.New(t.Context())
	require.NoError(t, err)

	changesetInput := operations.EmptyInput{}
	_, err = TerminalErrorExampleChangeset{}.Apply(*env, changesetInput)
	require.ErrorContains(t, err, "terminal error")
}
