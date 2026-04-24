package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/cron"
)

func TestSmoke(t *testing.T) {
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	pdConfig, err := products.LoadOutput[cron.Configurator](outputFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	cls, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		runs, err := cls[0].MustReadRunsByJob(pdConfig.Config[0].Out.JobID)
		require.NoError(c, err)
		require.GreaterOrEqual(c, len(runs.Data), 10)
		for _, j := range runs.Data {
			require.Equal(c, []interface{}{interface{}(nil)}, j.Attributes.Errors)
		}
	}, 2*time.Minute, 2*time.Second)
}
