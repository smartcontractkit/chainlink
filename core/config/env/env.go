package env

import (
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	Config = Env("CL_CONFIG")

	// LOOPP commands and vars
	MedianPluginCmd   = Env("CL_MEDIAN_CMD")
	SolanaPluginCmd   = Env("CL_SOLANA_CMD")
	StarknetPluginCmd = Env("CL_STARKNET_CMD")
	// PrometheusDiscoveryHostName is the externally accessible hostname
	// published by the node in the `/discovery` endpoint. Generally, it is expected to match
	// the public hostname of node.
	// Cluster step up like kubernetes may need to set this explicitly to ensure
	// that Prometheus can discovery LOOPps.
	// In house we observed that the resolved value of os.Hostname was not accessible to
	// outside of the given pod
	PrometheusDiscoveryHostName = Env("CL_PROMETHEUS_DISCOVERY_HOSTNAME")
	// EnvLooopHostName is the hostname used for HTTP communication between the
	// node and LOOPps. In most cases this does not need to be set explicitly.
	LooppHostName = Env("CL_LOOPP_HOSTNAME")

	DatabaseAllowSimplePasswords = Env("CL_DATABASE_ALLOW_SIMPLE_PASSWORDS")
	DatabaseURL                  = EnvSecret("CL_DATABASE_URL")
	DatabaseBackupURL            = EnvSecret("CL_DATABASE_BACKUP_URL")
	ExplorerAccessKey            = EnvSecret("CL_EXPLORER_ACCESS_KEY")
	ExplorerSecret               = EnvSecret("CL_EXPLORER_SECRET")
	PasswordKeystore             = EnvSecret("CL_PASSWORD_KEYSTORE")
	PasswordVRF                  = EnvSecret("CL_PASSWORD_VRF")
	PyroscopeAuthToken           = EnvSecret("CL_PYROSCOPE_AUTH_TOKEN")
	PrometheusAuthToken          = EnvSecret("CL_PROMETHEUS_AUTH_TOKEN")
	ThresholdKeyShare            = EnvSecret("CL_THRESHOLD_KEY_SHARE")
)

type Env string

func (e Env) Get() string { return os.Getenv(string(e)) }

// Lookup wraps [os.LookupEnv]
func (e Env) Lookup() (string, bool) { return os.LookupEnv(string(e)) }

func (e Env) IsTrue() bool { return strings.ToLower(e.Get()) == "true" }

type EnvSecret string

func (e EnvSecret) Get() models.Secret { return models.Secret(os.Getenv(string(e))) }
