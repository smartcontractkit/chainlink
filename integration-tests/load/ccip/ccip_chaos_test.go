package ccip

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	framework "github.com/smartcontractkit/chainlink-testing-framework/framework/grafana"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
	"github.com/smartcontractkit/chainlink/load-tests/ccip/template"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func Ptr[T any](value T) *T { return &value }

func a(ns, text string, from, to *time.Time) framework.Annotation {
	a := framework.Annotation{
		Text:         fmt.Sprintf("Namespace: %s, Test: %s", ns, text),
		StartTime:    from,
		Tags:         []string{"chaos"},
		DashboardUID: []string{"WaspDebug"},
	}
	if !to.IsZero() {
		a.EndTime = to
	}
	return a
}

func TestK8sChaos(t *testing.T) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	c, err := havoc.NewChaosMeshClient()
	require.NoError(t, err)

	config, err := tc.GetConfig([]string{"Chaos"}, tc.CCIP)
	require.NoError(t, err)
	cfg := config.CCIP.Chaos
	cr := template.NewChaosRunner(l, c)

	gc := framework.NewGrafanaClient(os.Getenv("GRAFANA_URL"), os.Getenv("GRAFANA_TOKEN"))

	testCases := []struct {
		name     string
		run      func(t *testing.T)
		validate func(t *testing.T)
	}{
		// pod failures
		{
			name: "Fail src chain",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					template.PodFailCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "instance",
						LabelValues:       []string{"geth-1337"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail dst chain",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					template.PodFailCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "instance",
						LabelValues:       []string{"geth-2337"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail one CL node",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					template.PodFailCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"ccip-0"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail two CL nodes",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					template.PodFailCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"ccip-0", "ccip-1"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail one CL node DB",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					template.PodFailCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"chainlink-don-db-0"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail one RMN node",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					template.PodFailCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"rmn-0"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail two RMN nodes",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					template.PodFailCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"rmn-0", "rmn-1"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// network delay
		{
			name: "Slow src chain",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					template.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "instance",
						LabelValues:       []string{"geth-1337"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Slow dst chain",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					template.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "instance",
						LabelValues:       []string{"geth-2337"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "One slow CL node",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					template.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"ccip-0"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Two slow CL nodes",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					template.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"ccip-0", "ccip-1"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "One slow CL node DB",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					template.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"chainlink-don-db-0"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "One slow RMN node",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					template.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"rmn-0"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Two slow RMN nodes",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					template.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"rmn-0", "rmn-1"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// network partition
		{
			name: "CL node <> CL nodes partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					template.PodPartitionCfg{
						Namespace:         cfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"ccip-1", "ccip-2", "ccip-3"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "2 CL nodes <> 2 CL nodes partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					template.PodPartitionCfg{
						Namespace:         cfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0", "ccip-1"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"ccip-2", "ccip-3"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "CL node <> DB partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					template.PodPartitionCfg{
						Namespace:         cfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"chainlink-don-db-0"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "RMN node <> RMN node",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					template.PodPartitionCfg{
						Namespace:         cfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"rmn-0"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"rmn-1", "rmn-2", "rmn-3"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "2 RMN nodes <> 2 RMN nodes partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					template.PodPartitionCfg{
						Namespace:         cfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"rmn-0", "rmn-1"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"rmn-2", "rmn-3"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "2 CL nodes <> 2 RMN nodes partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					template.PodPartitionCfg{
						Namespace:         cfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0", "ccip-1"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"rmn-2", "rmn-3"},
						InjectionDuration: cfg.GetExperimentInjectionInterval(),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// reorgs
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
			n := time.Now()
			testCase.run(t)
			time.Sleep(cfg.GetExperimentInterval())
			_, _, err = gc.Annotate(a(cfg.Namespace, testCase.name, Ptr(n), Ptr(time.Now())))
			testCase.validate(t)
		})
	}
}
