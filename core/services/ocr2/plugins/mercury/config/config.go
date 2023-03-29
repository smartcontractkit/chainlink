// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type PluginConfig struct {
	RawServerURL string              `json:"serverURL" toml:"serverURL"`
	ServerPubKey utils.PlainHexBytes `json:"serverPubKey" toml:"serverPubKey"`
}

func ValidatePluginConfig(config PluginConfig) (merr error) {
	if config.RawServerURL == "" {
		merr = errors.New("Mercury: ServerURL must be specified")
	} else {
		var normalizedURI string
		if schemeRegexp.MatchString(config.RawServerURL) {
			normalizedURI = config.RawServerURL
		} else {
			normalizedURI = fmt.Sprintf("wss://%s", config.RawServerURL)
		}
		uri, err := url.ParseRequestURI(normalizedURI)
		if err != nil {
			merr = pkgerrors.Wrap(err, "Mercury: invalid value for ServerURL")
		} else if !(uri.Scheme == "" || uri.Scheme == "wss") {
			merr = pkgerrors.Errorf(`Mercury: invalid scheme specified for MercuryServer, got: %q (scheme: %q) but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`, config.RawServerURL, uri.Scheme)
		}
	}
	if len(config.ServerPubKey) != 32 {
		merr = errors.Join(merr, errors.New("Mercury: ServerPubKey is required and must be a 32-byte hex string"))
	}
	return merr
}

var schemeRegexp = regexp.MustCompile(`^(.*)://`)
var wssRegexp = regexp.MustCompile(`^wss://`)

func (p PluginConfig) ServerURL() string {
	return wssRegexp.ReplaceAllString(p.RawServerURL, "")
}
