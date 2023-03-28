package v2

import (
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var (
	EnvConfig = Env("CL_CONFIG")
	EnvDev    = Env("CL_DEV")

	EnvDatabaseAllowSimplePasswords = Env("CL_DATABASE_ALLOW_SIMPLE_PASSWORDS")
	EnvDatabaseURL                  = EnvSecret("CL_DATABASE_URL")
	EnvDatabaseBackupURL            = EnvSecret("CL_DATABASE_BACKUP_URL")
	EnvExplorerAccessKey            = EnvSecret("CL_EXPLORER_ACCESS_KEY")
	EnvExplorerSecret               = EnvSecret("CL_EXPLORER_SECRET")
	EnvPasswordKeystore             = EnvSecret("CL_PASSWORD_KEYSTORE")
	EnvPasswordVRF                  = EnvSecret("CL_PASSWORD_VRF")
	EnvPyroscopeAuthToken           = EnvSecret("CL_PYROSCOPE_AUTH_TOKEN")
	EnvPrometheusAuthToken          = EnvSecret("CL_PROMETHEUS_AUTH_TOKEN")
)

type Env string

func (e Env) Get() string { return os.Getenv(string(e)) }

func (e Env) IsTrue() bool { return strings.ToLower(e.Get()) == "true" }

type EnvSecret string

func (e EnvSecret) Get() models.Secret { return models.Secret(os.Getenv(string(e))) }
