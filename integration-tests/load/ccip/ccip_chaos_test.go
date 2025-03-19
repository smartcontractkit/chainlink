package ccip

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	framework "github.com/smartcontractkit/chainlink-testing-framework/framework/grafana"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func Ptr[T any](value T) *T { return &value }

func a(ns, text string, dashboardUIDs []string, from, to *time.Time) framework.Annotation {
	a := framework.Annotation{
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

func prepareChaos(t *testing.T) (*ccip.Config, *havoc.NamespaceScopedChaosRunner, *framework.Client) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	c, err := havoc.NewChaosMeshClient()
	if err != nil {
		t.Error("Failed to create chaos mesh client", err)
	}

	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	if err != nil {
		t.Error("Failed to get config", err)
	}
	cfg := config.CCIP
	cr := havoc.NewNamespaceRunner(l, c, false)

	gc := framework.NewGrafanaClient(os.Getenv("GRAFANA_URL"), os.Getenv("GRAFANA_TOKEN"))
	return cfg, cr, gc
}

func runRealisticRPCLatencySuite(t *testing.T, testDuration, latency, jitter time.Duration) {
	config, cr, _ := prepareChaos(t)
	cfg := config.Chaos

	testCases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "Realistic RPC Latency",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					havoc.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"geth-1337", "geth-2337", "geth-90000001", "geth-90000002", "geth-90000003", "geth-90000004"},
						Latency:           latency,
						Jitter:            jitter,
						Correlation:       "0",
						InjectionDuration: testDuration,
					})
				if err != nil {
					t.Error("Failed to inject rpc latency", err)
				}
			},
		},
	}

	t.Logf("Starting chaos tests in %s", cfg.GetWaitBeforeStart().String())
	time.Sleep(cfg.GetWaitBeforeStart())

	// Run test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.run(t)
			time.Sleep(testDuration)
		})
	}
}

func TestChaos(t *testing.T) {
	runFullChaosSuite(t)
}

type cribNetworkConfig []struct {
	HTTPRPCs []struct {
		External string `json:"External"`
		Internal string `json:"Internal"`
	} `json:"HTTPRPCs"`
}

func readCRIBConfig(t *testing.T, cfg *ccip.Config) cribNetworkConfig {
	f, _ := os.ReadFile(*cfg.Load.CribEnvDirectory + "/ccip-v2-scripts-chains-details.json")
	var cribConfig cribNetworkConfig
	err := json.Unmarshal(f, &cribConfig)
	assert.NoError(t, err)
	return cribConfig
}

func runFullChaosSuite(t *testing.T) {
	config, cr, gc := prepareChaos(t)
	chaosCfg := config.Chaos
	cnc := readCRIBConfig(t, config)
	reorgFunc := func(cncs cribNetworkConfig, blocks int) {
		for _, cnc := range cncs {
			t.Logf("Reorg: %d", blocks)
			r := rpc.New(cnc.HTTPRPCs[0].External, nil)
			tcName := fmt.Sprintf("%s-%d-blocks", cnc.HTTPRPCs[0].External, blocks)
			t.Run(tcName, func(t *testing.T) {
				n := time.Now()
				err := r.GethSetHead(blocks)
				if err != nil {
					t.Error("Failed to set block head on geth", err)
				}
				time.Sleep(chaosCfg.GetExperimentInterval())
				_, _, err = gc.Annotate(a(chaosCfg.Namespace, tcName, chaosCfg.DashboardUIDs, Ptr(n), Ptr(time.Now())))
				if err != nil {
					t.Error("Failed to annotate grafana with chaos labels", err)
				}
			})
		}
	}

	testCases := []struct {
		name     string
		run      func(t *testing.T)
		validate func(t *testing.T)
	}{
		// reorgs
		{
			name: "Chain reorgs below finality",
			run: func(t *testing.T) {
				reorgFunc(cnc, 30)
			},
			validate: func(t *testing.T) {},
		},
		// pod failures
		{
			name: "Fail 3 chains",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					havoc.PodFailCfg{
						Namespace:         chaosCfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"geth-1337", "geth-2337", "geth-90000001"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail 3 CL nodes",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					havoc.PodFailCfg{
						Namespace:         chaosCfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"ccip-0", "ccip-1", "ccip-2"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail 3 CL node DB",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					havoc.PodFailCfg{
						Namespace:         chaosCfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"ccip-db-0", "ccip-db-1", "ccip-db-2"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// TODO: add RMNv2 then uncomment
		/**
		{
			name: "Fail two RMN nodes",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					havoc.PodFailCfg{
						Namespace:         chaosCfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"rmn-0", "rmn-1"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		*/
		// network delay
		{
			name: "Three slow chains",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					havoc.PodDelayCfg{
						Namespace:         chaosCfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"geth-1337", "geth-2337", "geth-90000001"},
						Latency:           400 * time.Millisecond,
						Jitter:            20 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Three slow CL nodes",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					havoc.PodDelayCfg{
						Namespace:         chaosCfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"ccip-0", "ccip-1", "ccip-2"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Three slow CL node DBs",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					havoc.PodDelayCfg{
						Namespace:         chaosCfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"ccip-db-0", "ccip-db-1", "ccip-db-2"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// TODO: add RMNv2 then uncomment
		/**
		{
			name: "Two slow RMN nodes",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					havoc.PodDelayCfg{
						Namespace:         chaosCfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"rmn-0", "rmn-1"},
						Latency:           200 * time.Millisecond,
						Jitter:            200 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		*/
		// network partition
		{
			name: "2 CL nodes <> 2 CL nodes partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					havoc.PodPartitionCfg{
						Namespace:         chaosCfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0", "ccip-1"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"ccip-2", "ccip-3"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "4 nodes partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					havoc.PodPartitionCfg{
						Namespace:         chaosCfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0", "ccip-1", "ccip-2", "ccip-3"},
						LabelToKey:        "app.kubernetes.io/name",
						LabelToValues:     []string{"ccip"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "CL node <> DB partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					havoc.PodPartitionCfg{
						Namespace:         chaosCfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"ccip-db-0"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// TODO: add RMNv2 then uncomment
		/**
		{
			name: "2 RMN nodes <> 2 RMN nodes partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					havoc.PodPartitionCfg{
						Namespace:         chaosCfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"rmn-0", "rmn-1"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"rmn-2", "rmn-3"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// TODO: add RMNv2 then uncomment
		{
			name: "2 CL nodes <> 2 RMN nodes partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					havoc.PodPartitionCfg{
						Namespace:         chaosCfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0", "ccip-1"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"rmn-2", "rmn-3"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},

		may be hard to recover run when all the testing has ended
		{
			name: "Chain reorgs above finality",
			run: func(t *testing.T) {
				reorgFunc(cnc, 250)
			},
			validate: func(t *testing.T) {},
		},

		{
			name: "8-8 CL nodes split brain",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					havoc.PodPartitionCfg{
						Namespace:         chaosCfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0", "ccip-1", "ccip-2", "ccip-3", "ccip-4", "ccip-5", "ccip-6", "ccip-7"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"ccip-8", "ccip-9", "ccip-10", "ccip-11", "ccip-12", "ccip-13", "ccip-14", "ccip-15"},
						InjectionDuration: chaosCfg.GetExperimentInjectionInterval(),
					})
				assert.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		*/
	}

	t.Logf("Starting chaos tests in %s", chaosCfg.GetWaitBeforeStart().String())
	time.Sleep(chaosCfg.GetWaitBeforeStart())

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			n := time.Now()
			testCase.run(t)
			time.Sleep(chaosCfg.GetExperimentInterval())
			_, _, err := gc.Annotate(a(chaosCfg.Namespace, testCase.name, chaosCfg.DashboardUIDs, Ptr(n), Ptr(time.Now())))
			if err != nil {
				t.Error("Failed to annotate grafana with chaos labels", err)
			}
			testCase.validate(t)
		})
	}
}
