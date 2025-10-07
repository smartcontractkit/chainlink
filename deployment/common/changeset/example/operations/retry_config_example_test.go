package example

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

func TestDisableRetryExampleChangeset(t *testing.T) {
	env, err := environment.New(t.Context())
	require.NoError(t, err)

	changesetInput := operations.EmptyInput{}
	_, err = DisableRetryExampleChangeset{}.Apply(*env, changesetInput)
	require.ErrorContains(t, err, "operation failed")
}

func TestUpdateInputExampleChangeset(t *testing.T) {

	env, err := environment.New(t.Context())
	require.NoError(t, err)

	changesetInput := operations.EmptyInput{}
	_, err = UpdateInputExampleChangeset{}.Apply(*env, changesetInput)
	require.NoError(t, err)
}
