package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	Config                       = Var("CL_CONFIG")
	DatabaseAllowSimplePasswords = Var("CL_DATABASE_ALLOW_SIMPLE_PASSWORDS")
	DatabaseURL                  = Secret("CL_DATABASE_URL")
	DatabaseBackupURL            = Secret("CL_DATABASE_BACKUP_URL")
	PasswordKeystore             = Secret("CL_PASSWORD_KEYSTORE")
	PasswordVRF                  = Secret("CL_PASSWORD_VRF")
	PyroscopeAuthToken           = Secret("CL_PYROSCOPE_AUTH_TOKEN")
	PrometheusAuthToken          = Secret("CL_PROMETHEUS_AUTH_TOKEN")
	ThresholdKeyShare            = Secret("CL_THRESHOLD_KEY_SHARE")
	// Migrations env vars
	EVMChainIDNotNullMigration0195 = "CL_EVM_CHAINID_NOT_NULL_MIGRATION_0195"
	CustomDefaults                 = Var("CL_CHAIN_DEFAULTS")
)

// LOOPP commands and vars
var (
	MedianPlugin   = NewPlugin("median")
	MercuryPlugin  = NewPlugin("mercury")
	SolanaPlugin   = NewPlugin("solana")
	StarknetPlugin = NewPlugin("starknet")
	// PrometheusDiscoveryHostName is the externally accessible hostname
	// published by the node in the `/discovery` endpoint. Generally, it is expected to match
	// the public hostname of node.
	// Cluster step up like kubernetes may need to set this explicitly to ensure
	// that Prometheus can discovery LOOPps.
	// In house we observed that the resolved value of os.Hostname was not accessible to
	// outside of the given pod
	PrometheusDiscoveryHostName = Var("CL_PROMETHEUS_DISCOVERY_HOSTNAME")
	// LOOPPHostName is the hostname used for HTTP communication between the
	// node and LOOPps. In most cases this does not need to be set explicitly.
	LOOPPHostName = Var("CL_LOOPP_HOSTNAME")
	// Work around for Solana LOOPPs configured with zero values.
	MinOCR2MaxDurationQuery = Var("CL_MIN_OCR2_MAX_DURATION_QUERY")
	// PipelineOvertime is an undocumented escape hatch for overriding the default padding in pipeline executions.
	PipelineOvertime = Var("CL_PIPELINE_OVERTIME")
)

type Var string

func (e Var) Get() string { return os.Getenv(string(e)) }

// Lookup wraps [os.LookupEnv]
func (e Var) Lookup() (string, bool) { return os.LookupEnv(string(e)) }

func (e Var) IsTrue() bool { return strings.ToLower(e.Get()) == "true" }

type Secret string

func (e Secret) Get() models.Secret { return models.Secret(os.Getenv(string(e))) }

type Plugin struct {
	Cmd Var
	Env Var
}

func NewPlugin(kind string) Plugin {
	kind = strings.ToUpper(kind)
	return Plugin{
		Cmd: Var(fmt.Sprintf("CL_%s_CMD", kind)),
		Env: Var(fmt.Sprintf("CL_%s_ENV", kind)),
	}
}
