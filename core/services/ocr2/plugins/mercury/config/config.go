// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type PluginConfig struct {
	URL          *models.URL         `json:"url"`
	ServerPubKey utils.PlainHexBytes `json:"serverPubKey"`
	ClientPubKey utils.PlainHexBytes `json:"clientPubKey"`
}

func ValidatePluginConfig(config PluginConfig) error {
	if config.URL == nil {
		return errors.New("Mercury URL must be specified")
	}
	// TODO: Where to validate that mercury definitely has a FeedID (and non-mercury has transmitter ID?)
	// if (config.FeedID == common.Hash{}) {
	//     return errors.New("FeedID must be specified and non-zero")
	// }
	if len(config.ServerPubKey) != 32 {
		return errors.New("ServerPubKey is required and must be a 32-byte hex string")
	}
	if len(config.ClientPubKey) != 32 {
		return errors.New("ClientPubKey is required and must be a 32-byte hex string")
	}
	return nil
}
