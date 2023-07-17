package core

import (
	_ "embed"
)

//go:embed testdata/mergingsecretsdata/secrets-database.toml
var DatabaseSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-explorer.toml
var ExplorerSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-password.toml
var PasswordSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-pyroscope.toml
var PyroscopeSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-prometheus.toml
var PrometheusSecretsTOML string

//go:embed testdata/mergingsecretsdata/secrets-mercury-split-one.toml
var MercurySecretsTOMLSplitOne string

//go:embed testdata/mergingsecretsdata/secrets-mercury-split-two.toml
var MercurySecretsTOMLSplitTwo string

//go:embed testdata/mergingsecretsdata/secrets-threshold.toml
var ThresholdSecretsTOML string
