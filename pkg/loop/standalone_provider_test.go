package loop_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
)

func TestRegisterStandAloneProvider(t *testing.T) {
	s := grpc.NewServer()

	p := test.StaticPluginProvider{}
	err := loop.RegisterStandAloneProvider(s, p, "some-type-we-do-not-support")
	require.ErrorContains(t, err, "stand alone provider only supports median")

	err = loop.RegisterStandAloneProvider(s, p, "median")
	require.ErrorContains(t, err, "expected median provider got")

	stopCh := newStopCh(t)
	pr := newPluginRelayerExec(t, stopCh)
	mp := newMedianProvider(t, pr)
	err = loop.RegisterStandAloneProvider(s, mp, "median")
	require.NoError(t, err)
}
