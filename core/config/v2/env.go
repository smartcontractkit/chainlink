package v2

import (
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	EnvConfig = Env("CL_CONFIG")

	// LOOPP commands and vars
	EnvMedianPluginCmd   = Env("CL_MEDIAN_CMD")
	EnvSolanaPluginCmd   = Env("CL_SOLANA_CMD")
	EnvStarknetPluginCmd = Env("CL_STARKNET_CMD")
	// EnvPrometheusDiscoveryHostName is the externally accessible hostname
	// published by the node in the `/discovery` endpoint. Generally, it is expected to match
	// the public hostname of node.
	// Cluster step up like kubernetes may need to set this explicitly to ensure
	// that Prometheus can discovery LOOPps.
	// In house we observed that the resolved value of os.Hostname was not accessible to
	// outside of the given pod
	EnvPrometheusDiscoveryHostName = Env("CL_PROMETHEUS_DISCOVERY_HOSTNAME")
	// EnvLooopHostName is the hostname used for HTTP communication between the
	// node and LOOPps. In most cases this does not need to be set explicitly.
	EnvLooppHostName = Env("CL_LOOPP_HOSTNAME")

	EnvDatabaseAllowSimplePasswords = Env("CL_DATABASE_ALLOW_SIMPLE_PASSWORDS")
	EnvDatabaseURL                  = EnvSecret("CL_DATABASE_URL")
	EnvDatabaseBackupURL            = EnvSecret("CL_DATABASE_BACKUP_URL")
	EnvExplorerAccessKey            = EnvSecret("CL_EXPLORER_ACCESS_KEY")
	EnvExplorerSecret               = EnvSecret("CL_EXPLORER_SECRET")
	EnvPasswordKeystore             = EnvSecret("CL_PASSWORD_KEYSTORE")
	EnvPasswordVRF                  = EnvSecret("CL_PASSWORD_VRF")
	EnvPyroscopeAuthToken           = EnvSecret("CL_PYROSCOPE_AUTH_TOKEN")
	EnvPrometheusAuthToken          = EnvSecret("CL_PROMETHEUS_AUTH_TOKEN")
	EnvThresholdKeyShare            = EnvSecret("CL_THRESHOLD_KEY_SHARE")
)

type Env string

func (e Env) Get() string { return os.Getenv(string(e)) }

// Lookup wraps [os.LookupEnv]
func (e Env) Lookup() (string, bool) { return os.LookupEnv(string(e)) }

func (e Env) IsTrue() bool { return strings.ToLower(e.Get()) == "true" }

type EnvSecret string

func (e EnvSecret) Get() models.Secret { return models.Secret(os.Getenv(string(e))) }
