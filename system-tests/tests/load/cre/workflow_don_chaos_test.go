package cre

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	frameworkGrafana "github.com/smartcontractkit/chainlink-testing-framework/framework/grafana"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
)

func Ptr[T any](value T) *T { return &value }

func a(ns, text string, dashboardUIDs []string, from, to *time.Time) frameworkGrafana.Annotation {
	a := frameworkGrafana.Annotation{
		Text:         fmt.Sprintf("Namespace: %s, Test: %s", ns, text),
		StartTime:    from,
		Tags:         []string{"chaos"},
		DashboardUID: dashboardUIDs,
	}
	if !to.IsZero() {
		a.EndTime = to
	}
	return a
}

// prepareChaos creates a namespace scoped chaos runner and Grafana client
func prepareChaos(t *testing.T) (*havoc.NamespaceScopedChaosRunner, *frameworkGrafana.Client, error) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	c, err := havoc.NewChaosMeshClient()
	if err != nil {
		t.Error("Failed to create chaos mesh client", err)
	}
	gURL := os.Getenv("GRAFANA_URL")
	gToken := os.Getenv("GRAFANA_TOKEN")
	if gURL == "" || gToken == "" {
		return nil, nil, errors.New("GRAFANA_URL or GRAFANA_TOKEN environment variables not set")
	}
	return havoc.NewNamespaceRunner(l, c, false), frameworkGrafana.NewGrafanaClient(gURL, gToken), nil
}

func TestChaos(t *testing.T) {
	in, err := framework.Load[TestConfigLoadTest](t)
	require.NoError(t, err, "couldn't load test config")
	runChaosSuite(t, in)
}

// runChaosSuite runs chaos suite for Assets, Workflow and Writer NodeSets
func runChaosSuite(t *testing.T, testConfig *TestConfigLoadTest) {
	cr, gc, err := prepareChaos(t)
	require.NoError(t, err)
	cribCfg := testConfig.Infra.CRIB
	chaosCfg := testConfig.Chaos

	testDuration, err := time.ParseDuration(testConfig.Duration)
	require.NoError(t, err, "could not parse test duration")

	rpcLatency, err := time.ParseDuration(chaosCfg.Latency)
	require.NoError(t, err, "could not parse chaos latency")
	rpcJitter, err := time.ParseDuration(chaosCfg.Jitter)
	require.NoError(t, err, "could not parse chaos jitter")

	waitDur, err := time.ParseDuration(chaosCfg.WaitBeforeStart)
	require.NoError(t, err, "could not parse chaos wait time")
	expFullDur, err := time.ParseDuration(chaosCfg.ExperimentFullInterval)
	require.NoError(t, err, "could not parse chaos experiment full interval")
	expInjectDur, err := time.ParseDuration(chaosCfg.ExperimentInjectionInterval)
	require.NoError(t, err, "could not parse chaos experiment injection interval")

	reorgFunc := func(rpcs []*blockchain.Node, blocks int) {
		for _, n := range rpcs {
			t.Logf("Reorg: %d", blocks)
			r := rpc.New(n.ExternalHTTPUrl, nil)
			tcName := fmt.Sprintf("%s-%d-blocks", n.ExternalHTTPUrl, blocks)
			t.Run(tcName, func(t *testing.T) {
				n := time.Now()
				err := r.GethSetHead(blocks)
				if err != nil {
					t.Error("Failed to set block head on Geth", err)
				}
				time.Sleep(expFullDur)
				_, _, err = gc.Annotate(a(cribCfg.Namespace, tcName, chaosCfg.DashboardUIDs, Ptr(n), Ptr(time.Now())))
				if err != nil {
					t.Error("Failed to annotate grafana with chaos labels", err)
				}
			})
		}
	}

	type testCase struct {
		name     string
		run      func(t *testing.T, r []*blockchain.Node)
		validate func(t *testing.T)
	}
	var testCases []testCase

	switch chaosCfg.Mode {
	case "clean":
		framework.L.Info().Msg("Skipping chaos experiments")
		return
	case "reorg":
		testCases = []testCase{
			{
				name: "Chain reorgs below finality",
				run: func(t *testing.T, r []*blockchain.Node) {
					reorgFunc(r, 50)
				},
				validate: func(t *testing.T) {},
			},
		}
	case "rpc":
		testCases = []testCase{
			{
				name: "Realistic RPC Latency",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodDelay(context.Background(),
						havoc.PodDelayCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/name",
							LabelValues:       []string{"geth-1337"},
							Latency:           rpcLatency,
							Jitter:            rpcJitter,
							Correlation:       "0",
							InjectionDuration: testDuration,
						})
					require.NoError(t, err)
				},
			},
		}
	case "full":
		testCases = []testCase{
			// blockchain
			{
				name: "Fail EVM Chain",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"geth-1337"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "400ms+200ms jitter for EVM Chain",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodDelay(context.Background(),
						havoc.PodDelayCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"geth-1337"},
							Latency:           400 * time.Millisecond,
							Jitter:            200 * time.Millisecond,
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "30% corrupt for EVM Chain",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodCorrupt(context.Background(),
						havoc.PodCorruptCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"geth-1337"},
							Corrupt:           "30",
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "30% loss for EVM Chain",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodLoss(context.Background(),
						havoc.PodLossCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"geth-1337"},
							Loss:              "30",
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition EVM Chain from 2 Assets nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"geth-1337"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"asset-0", "asset-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition EVM Chain from 2 Workflow nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"geth-1337"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"workflow-0", "workflow-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition EVM Chain from 2 Writer nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"geth-1337"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"writer-0", "writer-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},

			// node sets

			// assets
			{
				name: "Fail Assets Node",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"asset-0"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "400ms+200ms jitter for 2 DBs of Assets Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodDelay(context.Background(),
						havoc.PodDelayCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"asset-db-0", "asset-db-1"},
							Latency:           400 * time.Millisecond,
							Jitter:            200 * time.Millisecond,
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "400ms+200ms jitter for 2 Assets Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodDelay(context.Background(),
						havoc.PodDelayCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"asset-0", "asset-1"},
							Latency:           400 * time.Millisecond,
							Jitter:            200 * time.Millisecond,
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "30% corrupt for 2 Assets Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodCorrupt(context.Background(),
						havoc.PodCorruptCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"asset-0", "asset-1"},
							Corrupt:           "30",
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "30% loss for 2 Assets Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodLoss(context.Background(),
						havoc.PodLossCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"asset-0", "asset-1"},
							Loss:              "30",
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition 2 Assets nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"asset-0"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"asset-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Fail 2 Assets Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"asset-0", "asset-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition 2 Assets Nodes <> 2 Assets Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"asset-0", "asset-1"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"asset-2", "asset-3"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},

			// workflow
			{
				name: "Fail Workflow Node",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"workflow-0"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "400ms+200ms jitter for 2 DBs of Workflow Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodDelay(context.Background(),
						havoc.PodDelayCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"workflow-db-0", "workflow-db-1"},
							Latency:           400 * time.Millisecond,
							Jitter:            200 * time.Millisecond,
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "400ms+200ms jitter for 2 Workflow Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodDelay(context.Background(),
						havoc.PodDelayCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"workflow-0", "workflow-1"},
							Latency:           400 * time.Millisecond,
							Jitter:            200 * time.Millisecond,
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "30% corrupt for 2 Workflow Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodCorrupt(context.Background(),
						havoc.PodCorruptCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"workflow-0", "workflow-1"},
							Corrupt:           "30",
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "30% loss for 2 Workflow Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodLoss(context.Background(),
						havoc.PodLossCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"workflow-0", "workflow-1"},
							Loss:              "30",
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition 2 Workflow nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"workflow-0"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"workflow-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Fail 2 Workflow Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"workflow-0", "workflow-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition 2 Workflow Nodes <> 2 Workflow Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"workflow-0", "workflow-1"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"workflow-2", "workflow-3"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},

			// writer
			{
				name: "Fail Writer Node",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"writer-0"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "400ms+200ms jitter for 2 DBs of Writer Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodDelay(context.Background(),
						havoc.PodDelayCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"writer-db-0", "writer-db-1"},
							Latency:           400 * time.Millisecond,
							Jitter:            200 * time.Millisecond,
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "400ms+200ms jitter for 2 Writer Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodDelay(context.Background(),
						havoc.PodDelayCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"writer-0", "writer-1"},
							Latency:           400 * time.Millisecond,
							Jitter:            200 * time.Millisecond,
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "30% corrupt for 2 Writer Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodCorrupt(context.Background(),
						havoc.PodCorruptCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"writer-0", "writer-1"},
							Corrupt:           "30",
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "30% loss for 2 Writer Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodLoss(context.Background(),
						havoc.PodLossCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"writer-0", "writer-1"},
							Loss:              "30",
							Correlation:       "0",
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition 2 Writer nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"writer-0"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"writer-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Fail 2 Writer Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"writer-0", "writer-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Partition 2 Writer Nodes <> 2 Writer Nodes",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"writer-0", "writer-1"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"writer-2", "writer-3"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},

			// node set <> node set
			{
				name: "2 Assets Nodes <> 2 Workflow Nodes partition",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"assets-0", "assets-1"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"workflow-0", "workflow-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "2 Workflow Nodes <> 2 Writer Nodes partition",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodPartition(context.Background(),
						havoc.PodPartitionCfg{
							Namespace:         cribCfg.Namespace,
							LabelFromKey:      "app.kubernetes.io/instance",
							LabelFromValues:   []string{"workflow-0", "workflow-1"},
							LabelToKey:        "app.kubernetes.io/instance",
							LabelToValues:     []string{"writer-0", "writer-1"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},

			// job distributor
			{
				name: "Fail Job Distributor",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"job-distributor"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},
			{
				name: "Fail Job Distributor DB",
				run: func(t *testing.T, r []*blockchain.Node) {
					_, err := cr.RunPodFail(context.Background(),
						havoc.PodFailCfg{
							Namespace:         cribCfg.Namespace,
							LabelKey:          "app.kubernetes.io/instance",
							LabelValues:       []string{"postgres"},
							InjectionDuration: expInjectDur,
						})
					assert.NoError(t, err)
				},
				validate: func(t *testing.T) {},
			},

			// EA experiments
			// TODO: there is no EA deployed atm
		}
	default:
		panic("unknown chaos test type, can be 'rpc' or 'full', check [chaos] in TOML config")
	}

	t.Logf("Starting chaos tests in %s", waitDur)
	time.Sleep(waitDur)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			n := time.Now()
			// TODO: load test is not reliable right now and in 50% hands on JD calls,
			// TODO: when it'd be stable we should remove these static URLs with real output of a test
			testConfig.Blockchains[0].Out = &blockchain.Output{
				Family: "evm",
				Nodes: []*blockchain.Node{
					{
						ExternalWSUrl:   "wss://crib-df-cre-chaos-test-geth-1337-ws.main.stage.cldev.sh",
						ExternalHTTPUrl: "https://crib-df-cre-chaos-test-geth-1337-http.main.stage.cldev.sh",
					},
				},
			}
			testCase.run(t, testConfig.Blockchains[0].Out.Nodes)
			time.Sleep(expFullDur)
			_, _, err := gc.Annotate(a(cribCfg.Namespace, testCase.name, chaosCfg.DashboardUIDs, Ptr(n), Ptr(time.Now())))
			if err != nil {
				t.Error("Failed to annotate grafana with chaos labels", err)
			}
			if testCase.validate != nil {
				testCase.validate(t)
			}
		})
	}
}
