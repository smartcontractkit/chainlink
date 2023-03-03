package config

import (
	"os"
	"sync"

	"github.com/imdario/mergo"
	"github.com/rs/zerolog/log"
)

const (
	EnvVarPrefix = "TEST_"

	EnvVarNoManifestUpdate            = "NO_MANIFEST_UPDATE"
	EnvVarNoManifestUpdateDescription = "Skip updating manifest when connecting to the namespace"
	EnvVarNoManifestUpdateExample     = "false"

	EnvVarNamespace            = "ENV_NAMESPACE"
	EnvVarNamespaceDescription = "Namespace name to connect to"
	EnvVarNamespaceExample     = "chainlink-test-epic"

	EnvVarJobImage            = "ENV_JOB_IMAGE"
	EnvVarJobImageDescription = "Image to run as a job in k8s"
	EnvVarJobImageExample     = "795953128386.dkr.ecr.us-west-2.amazonaws.com/core-integration-tests:v1.0"

	EnvVarInsideK8s            = "ENV_INSIDE_K8S"
	EnvVarInsideK8sDescription = "Internal variable to turn forwarding strategy off inside k8s, do not use"
	EnvVarInsideK8sExample     = ""

	EnvVarCLImage            = "CHAINLINK_IMAGE"
	EnvVarCLImageDescription = "Chainlink image repository"
	EnvVarCLImageExample     = "public.ecr.aws/chainlink/chainlink"

	EnvVarCLTag            = "CHAINLINK_VERSION"
	EnvVarCLTagDescription = "Chainlink image tag"
	EnvVarCLTagExample     = "1.9.0"

	EnvVarUser            = "CHAINLINK_ENV_USER"
	EnvVarUserDescription = "Owner of an environment"
	EnvVarUserExample     = "Satoshi"

	EnvVarCLCommitSha            = "CHAINLINK_COMMIT_SHA"
	EnvVarCLCommitShaDescription = "The sha of the commit that you're running tests on. Mostly used for CI"
	EnvVarCLCommitShaExample     = "${{ github.sha }}"

	EnvVarTestTrigger            = "TEST_TRIGGERED_BY"
	EnvVarTestTriggerDescription = "How the test was triggered, either manual or CI."
	EnvVarTestTriggerExample     = "CI"

	EnvVarLogLevel            = "TEST_LOG_LEVEL"
	EnvVarLogLevelDescription = "Environment logging level"
	EnvVarLogLevelExample     = "info | debug | trace"

	EnvVarSelectedNetworks            = "SELECTED_NETWORKS"
	EnvVarSelectedNetworksDescription = "Networks to select for testing"
	EnvVarSelectedNetworksExample     = "SIMULATED"

	EnvVarDBURL            = "DATABASE_URL"
	EnvVarDBURLDescription = "DATABASE_URL needed for component test. This is only necessary if testhelper methods are imported from core"
	EnvVarDBURLExample     = "postgresql://postgres:node@localhost:5432/chainlink_test?sslmode=disable"

	EnvVarSlackKey            = "SLACK_API_KEY"
	EnvVarSlackKeyDescription = "The OAuth Slack API key to report tests results with"
	EnvVarSlackKeyExample     = "xoxb-example-key"

	EnvVarSlackChannel            = "SLACK_CHANNEL"
	EnvVarSlackChannelDescription = "The Slack code for the channel you want to send the notification to"
	EnvVarSlackChannelExample     = "C000000000"

	EnvVarSlackUser            = "SLACK_USER"
	EnvVarSlackUserDescription = "The Slack code for the user you want to notify"
	EnvVarSlackUserExample     = "U000000000"

	EnvVarPyroscopeServer      = "PYROSCOPE_SERVER"
	EnvVarPyroscopeEnvironment = "PYROSCOPE_ENVIRONMENT"
	EnvVarPyroscopeKey         = "PYROSCOPE_KEY"

	EnvVarToleration                 = "K8S_TOLERATION"
	EnvVarTolerationsUserDescription = "Node roles to tolerate"
	EnvVarTolerationsExample         = "foundations"

	EnvVarNodeSelector                = "K8S_NODE_SELECTOR"
	EnvVarNodeSelectorUserDescription = "Node role to deploy to"
	EnvVarNodeSelectorExample         = "foundations"

	EnvVarDetachRunner                = "DETACH_RUNNER"
	EnvVarDetachRunnerUserDescription = "Should we detach the remote runner after starting a test using it"
	EnvVarDetachRunnerExample         = "true"

	EnvVarEVMKeys                = "EVM_KEYS"
	EnvVarEVMKeysUserDescription = "The keys used to connect to the evm"
	EnvVarEVMKeysExample         = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

var (
	JSIIGlobalMu = &sync.Mutex{}
)

func MustMerge(targetVars interface{}, codeVars interface{}) {
	if err := mergo.Merge(targetVars, codeVars, mergo.WithOverride); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func MustEnvOverrideVersion(target interface{}) {
	image := os.Getenv(EnvVarCLImage)
	tag := os.Getenv(EnvVarCLTag)
	if image != "" && tag != "" {
		if err := mergo.Merge(target, map[string]interface{}{
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   image,
					"version": tag,
				},
			},
		}, mergo.WithOverride); err != nil {
			log.Fatal().Err(err).Send()
		}
	}
}
