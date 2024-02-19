package loop_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
)

func TestRegisterStandAloneProvider_Median(t *testing.T) {
	s := grpc.NewServer()

	p := test.StaticPluginProvider{}
	err := loop.RegisterStandAloneProvider(s, p, "some-type-we-do-not-support")
	require.ErrorContains(t, err, "unsupported stand alone provider")

	err = loop.RegisterStandAloneProvider(s, p, "median")
	require.ErrorContains(t, err, "expected median provider got")

	stopCh := newStopCh(t)
	pr := newPluginRelayerExec(t, false, stopCh)
	mp := newMedianProvider(t, pr)
	err = loop.RegisterStandAloneProvider(s, mp, "median")
	require.NoError(t, err)
}

func TestRegisterStandAloneProvider_GenericPlugin(t *testing.T) {
	s := grpc.NewServer()

	stopCh := newStopCh(t)
	pr := newPluginRelayerExec(t, false, stopCh)
	gp := newGenericPluginProvider(t, pr)
	err := loop.RegisterStandAloneProvider(s, gp, "plugin")
	require.NoError(t, err)
}
