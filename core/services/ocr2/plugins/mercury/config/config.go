// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type PluginConfig struct {
	ServerHost   string              `json:"serverHost"`
	ServerPubKey utils.PlainHexBytes `json:"serverPubKey"`
}

func ValidatePluginConfig(config PluginConfig) (err error) {
	if config.ServerHost == "" {
		err = errors.New("Mercury: ServerHost must be specified")
	} else if !utils.IsHostnamePort(config.ServerHost) {
		// FIXME: actually according to James this is also valid: localhost:4040/some/foo
		err = errors.Errorf(`Mercury: invalid value specified for MercuryServer, got :%s but expected value in the form of "address:port" e.g. "192.0.2.2:4242"`, config.ServerHost)
	}
	if len(config.ServerPubKey) != 32 {
		err = multierr.Combine(err, errors.New("Mercury: ServerPubKey is required and must be a 32-byte hex string"))
	}
	return err
}
