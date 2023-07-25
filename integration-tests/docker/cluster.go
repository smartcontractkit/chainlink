package docker

import (
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/stretchr/testify/require"
	"testing"
)

/* These are the high level components you should reuse in your docker tests in other repos */

func NewChainlinkCluster(t *testing.T, nodes int) (*Environment, error) {
	lw, err := logwatch.NewLogWatch(t, nil)
	require.NoError(t, err)
	env, err := NewEnvironment(lw).
		WithContainer(NewGeth(nil)).
		WithContainer(NewMockServer(nil)).
		Start(true)
	require.NoError(t, err)
	gethComponent := env.Get("geth")[0].(*Geth)
	for i := 0; i < nodes; i++ {
		env.WithContainer(NewChainlink(NodeConfigOpts{
			EVM: NodeEVMSettings{
				HTTPURL: gethComponent.InternalHttpUrl,
				WSURL:   gethComponent.InternalWsUrl,
			}}))
	}
	env, err = env.Start(true)
	require.NoError(t, err)
	return env, nil
}
