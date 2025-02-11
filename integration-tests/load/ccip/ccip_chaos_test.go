package ccip

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/load-tests/ccip/template"
)

func TestK8sChaos(t *testing.T) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	c, err := havoc.NewChaosMeshClient()
	require.NoError(t, err)

	config, err := tc.GetConfig([]string{"Chaos"}, tc.CCIP)
	require.NoError(t, err)
	cfg := config.CCIP.Chaos

	testCases := []struct {
		name     string
		run      func(t *testing.T)
		validate func(t *testing.T)
	}{
		// pod failures
		{
			name: "Fail src chain",
			run: func(t *testing.T) {
				src, err := template.PodFail(c, l, template.PodFailCfg{
					Namespace:         cfg.Namespace,
					LabelKey:          "instance",
					LabelValues:       []string{"geth-1337"},
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail dst chain",
			run: func(t *testing.T) {
				dst, err := template.PodFail(c, l, template.PodFailCfg{
					Namespace:         cfg.Namespace,
					LabelKey:          "instance",
					LabelValues:       []string{"geth-2337"},
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				dst.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail one node",
			run: func(t *testing.T) {
				node1, err := template.PodFail(c, l, template.PodFailCfg{
					Namespace:         cfg.Namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"ccip-0"},
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				node1.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail two nodes",
			run: func(t *testing.T) {
				node1, err := template.PodFail(c, l, template.PodFailCfg{
					Namespace:         cfg.Namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"ccip-0", "ccip-1"},
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				node1.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		// network delay
		{
			name: "Slow src chain",
			run: func(t *testing.T) {
				src, err := template.PodDelay(c, l, template.PodDelayCfg{
					Namespace:         cfg.Namespace,
					LabelKey:          "instance",
					LabelValues:       []string{"geth-1337"},
					Latency:           200 * time.Millisecond,
					Jitter:            200 * time.Millisecond,
					Correlation:       "0",
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Slow dst chain",
			run: func(t *testing.T) {
				src, err := template.PodDelay(c, l, template.PodDelayCfg{
					Namespace:         cfg.Namespace,
					LabelKey:          "instance",
					LabelValues:       []string{"geth-2337"},
					Latency:           200 * time.Millisecond,
					Jitter:            200 * time.Millisecond,
					Correlation:       "0",
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "One slow node",
			run: func(t *testing.T) {
				src, err := template.PodDelay(c, l, template.PodDelayCfg{
					Namespace:         cfg.Namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"ccip-0"},
					Latency:           200 * time.Millisecond,
					Jitter:            200 * time.Millisecond,
					Correlation:       "0",
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Two slow nodes",
			run: func(t *testing.T) {
				src, err := template.PodDelay(c, l, template.PodDelayCfg{
					Namespace:         cfg.Namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"ccip-0", "ccip-1"},
					Latency:           200 * time.Millisecond,
					Jitter:            200 * time.Millisecond,
					Correlation:       "0",
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "One node partition",
			run: func(t *testing.T) {
				src, err := template.PodPartition(c, l, template.PodPartitionCfg{
					Namespace:         cfg.Namespace,
					LabelFromKey:      "app.kubernetes.io/instance",
					LabelFromValues:   []string{"ccip-0"},
					LabelToKey:        "app.kubernetes.io/instance",
					LabelToValues:     []string{"ccip-1", "ccip-2", "ccip-3"},
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Two nodes partition",
			run: func(t *testing.T) {
				src, err := template.PodPartition(c, l, template.PodPartitionCfg{
					Namespace:         cfg.Namespace,
					LabelFromKey:      "app.kubernetes.io/instance",
					LabelFromValues:   []string{"ccip-0", "ccip-1"},
					LabelToKey:        "app.kubernetes.io/instance",
					LabelToValues:     []string{"ccip-2", "ccip-3"},
					InjectionDuration: cfg.GetExperimentInjectionInterval(),
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Reorg src chain below finality",
			run: func(t *testing.T) {
				r := rpc.New(cfg.SrcChainURL, nil)
				err := r.GethSetHead(cfg.ReorgDepthBelowFinality)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Reorg dst chain below finality",
			run: func(t *testing.T) {
				r := rpc.New(cfg.DstChainURL, nil)
				err := r.GethSetHead(cfg.ReorgDepthBelowFinality)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Reorg src chain above finality",
			run: func(t *testing.T) {
				r := rpc.New(cfg.SrcChainURL, nil)
				err := r.GethSetHead(cfg.ReorgDepthAboveFinality)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Reorg dst chain above finality",
			run: func(t *testing.T) {
				r := rpc.New(cfg.DstChainURL, nil)
				err := r.GethSetHead(cfg.ReorgDepthAboveFinality)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
	}

	// Start WASP load test here, apply average load profile that you expect in production!
	// Configure timeouts and validate all the test cases until the test ends
	// or start the load test manually

	// Run test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.run(t)
			time.Sleep(cfg.GetExperimentInterval())
			testCase.validate(t)
		})
	}
}
