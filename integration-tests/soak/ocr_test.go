package soak

import (
	"fmt"
	"testing"

	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"time"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/smartcontractkit/chainlink-testing-framework/havoc"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestOCRv1Soak(t *testing.T) {
	config, err := tc.GetConfig([]string{"Soak"}, tc.OCR)
	require.NoError(t, err, "Error getting config")
	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config)
	require.NoError(t, err, "Error creating OCR soak test")
	executeOCRSoakTest(t, ocrSoakTest, &config)
}

func TestOCRv2Soak(t *testing.T) {
	config, err := tc.GetConfig([]string{"Soak"}, tc.OCR2)
	require.NoError(t, err, "Error getting config")

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config)
	require.NoError(t, err, "Error creating OCR soak test")
	executeOCRSoakTest(t, ocrSoakTest, &config)
}

func TestOCRSoak_GethReorgBelowFinality_FinalityTagDisabled(t *testing.T) {
	config, err := tc.GetConfig([]string{t.Name()}, tc.OCR)
	require.NoError(t, err, "Error getting config")
	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config)
	require.NoError(t, err, "Error creating OCR soak test")
	executeOCRSoakTest(t, ocrSoakTest, &config)
}

func TestOCRSoak_GethReorgBelowFinality_FinalityTagEnabled(t *testing.T) {
	config, err := tc.GetConfig([]string{t.Name()}, tc.OCR)
	require.NoError(t, err, "Error getting config")
	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config)
	require.NoError(t, err, "Error creating OCR soak test")
	executeOCRSoakTest(t, ocrSoakTest, &config)
}

func TestOCRSoak_GasSpike(t *testing.T) {
	config, err := tc.GetConfig([]string{t.Name()}, tc.OCR)
	require.NoError(t, err, "Error getting config")
	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config)
	require.NoError(t, err, "Error creating OCR soak test")
	executeOCRSoakTest(t, ocrSoakTest, &config)
}

// TestOCRSoak_ChangeBlockGasLimit changes next block gas limit and sets it to percentage of last gasUsed in previous block creating congestion
func TestOCRSoak_ChangeBlockGasLimit(t *testing.T) {
	config, err := tc.GetConfig([]string{t.Name()}, tc.OCR)
	require.NoError(t, err, "Error getting config")
	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config)
	require.NoError(t, err, "Error creating OCR soak test")
	executeOCRSoakTest(t, ocrSoakTest, &config)
}

// TestOCRSoak_RPCDownForAllCLNodes simulates a network chaos by bringing down network to RPC node for all Chainlink Nodes
func TestOCRSoak_RPCDownForAllCLNodes(t *testing.T) {
	config, err := tc.GetConfig([]string{t.Name()}, tc.OCR)
	require.NoError(t, err, "Error getting config")
	require.True(t, config.Network.IsSimulatedGethSelected(), "This test requires simulated geth")

	namespace := fmt.Sprintf("%s-%s", "soak-ocr-simulated-geth", uuid.NewString()[0:5])
	chaos, err := gethNetworkDownChaos(GethNetworkDownChaosOpts{
		DelayCreate: time.Minute * 2,
		Duration:    time.Minute * 5,
		Name:        "ocr-soak-test-geth-network-down-all-nodes",
		Description: "Geth Network Down For All Chainlink Nodes",
		Namespace:   namespace,
		TargetSelector: v1alpha1.PodSelector{
			Selector: v1alpha1.PodSelectorSpec{
				GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
					Namespaces:     []string{namespace},
					LabelSelectors: map[string]string{"app": "chainlink-0"},
				},
			},
			Mode: v1alpha1.AllMode,
		},
	})
	require.NoError(t, err, "Error creating chaos")
	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config,
		testsetups.WithNamespace(namespace),
		testsetups.WithChaos([]*havoc.Chaos{chaos}),
	)
	require.NoError(t, err, "Error creating OCR soak test")
	executeOCRSoakTest(t, ocrSoakTest, &config)
}

// TestOCRSoak_RPCDownForAllCLNodes simulates a network chaos by bringing down network to RPC node for 50% of Chainlink Nodes
func TestOCRSoak_RPCDownForHalfCLNodes(t *testing.T) {
	config, err := tc.GetConfig([]string{t.Name()}, tc.OCR)
	require.NoError(t, err, "Error getting config")
	require.True(t, config.Network.IsSimulatedGethSelected(), "This test requires simulated geth")

	namespace := fmt.Sprintf("%s-%s", "soak-ocr-simulated-geth", uuid.NewString()[0:5])
	chaos, err := gethNetworkDownChaos(GethNetworkDownChaosOpts{
		DelayCreate: time.Minute * 2,
		Duration:    time.Minute * 5,
		Name:        "ocr-soak-test-geth-network-down-half-nodes",
		Description: "Geth Network Down For 50 Percent Of Chainlink Nodes",
		Namespace:   namespace,
		TargetSelector: v1alpha1.PodSelector{
			Selector: v1alpha1.PodSelectorSpec{
				GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
					Namespaces:     []string{namespace},
					LabelSelectors: map[string]string{"app": "chainlink-0"},
				},
			},
			Mode:  v1alpha1.FixedPercentMode,
			Value: "50",
		},
	})
	require.NoError(t, err, "Error creating chaos")
	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config,
		testsetups.WithNamespace(namespace),
		testsetups.WithChaos([]*havoc.Chaos{chaos}),
	)
	require.NoError(t, err, "Error creating OCR soak test")
	executeOCRSoakTest(t, ocrSoakTest, &config)
}

func executeOCRSoakTest(t *testing.T, test *testsetups.OCRSoakTest, config *tc.TestConfig) {
	l := logging.GetTestLogger(t)

	// validate Seth config before anything else, but only for live networks (simulated will fail, since there's no chain started yet)
	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	if !network.Simulated {
		_, err := seth_utils.GetChainClient(config, network)
		require.NoError(t, err, "Error creating seth client")
	}

	l.Info().Str("test", t.Name()).Msg("Starting OCR soak test")
	if !test.Interrupted() {
		test.DeployEnvironment(config)
	}

	if test.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(test.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		} else {
			err := test.Environment().Client.RemoveNamespace(test.Environment().Cfg.Namespace)
			if err != nil {
				l.Error().Err(err).Msg("Error removing namespace")
			}
		}
	})
	if test.Interrupted() {
		err := test.LoadState()
		require.NoError(t, err, "Error loading state")
		test.Resume()
	} else {
		test.Setup(config)
		test.Run()
	}
}

type GethNetworkDownChaosOpts struct {
	Name           string
	Namespace      string
	Description    string
	TargetSelector v1alpha1.PodSelector
	DelayCreate    time.Duration
	Duration       time.Duration
}

func gethNetworkDownChaos(opts GethNetworkDownChaosOpts) (*havoc.Chaos, error) {
	k8sClient, err := havoc.NewChaosMeshClient()
	if err != nil {
		return nil, err
	}
	return havoc.NewChaos(havoc.ChaosOpts{
		Description: opts.Description,
		DelayCreate: opts.DelayCreate,
		Object: &v1alpha1.NetworkChaos{
			TypeMeta: metav1.TypeMeta{
				Kind:       "NetworkChaos",
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      opts.Name,
				Namespace: opts.Namespace,
			},
			Spec: v1alpha1.NetworkChaosSpec{
				Action: v1alpha1.LossAction,
				PodSelector: v1alpha1.PodSelector{
					Mode: v1alpha1.AllMode,
					Selector: v1alpha1.PodSelectorSpec{
						GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
							Namespaces:     []string{opts.Namespace},
							LabelSelectors: map[string]string{"app": "geth"},
						},
					},
				},
				Duration:  ptr.Ptr(opts.Duration.String()),
				Direction: v1alpha1.Both,
				Target:    &opts.TargetSelector,
				TcParameter: v1alpha1.TcParameter{
					Loss: &v1alpha1.LossSpec{
						Loss: "100",
					},
				},
			},
		},
		Client: k8sClient,
		Logger: &havoc.Logger,
	})

}
