package main

import (
	"testing"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre/testutils"
	"github.com/stretchr/testify/require"
)

func TestInitWorkflow(t *testing.T) {
	config := &Config{
		ChainSelector:   3379446385462418246,
		RegistryAddress: "0x5FbDB2315678afecb367f032d93F642f64180aa3",
	}
	runtime := testutils.NewRuntime(t, testutils.Secrets{})

	workflow, err := InitWorkflow(config, runtime.Logger(), nil)
	require.NoError(t, err)

	require.Len(t, workflow, 1)
	require.Equal(t, cron.Trigger(&cron.Config{}).CapabilityID(), workflow[0].CapabilityID())
}
