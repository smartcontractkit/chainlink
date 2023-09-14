// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type PluginConfig struct {
	RawServerURL string              `json:"serverURL" toml:"serverURL"`
	ServerPubKey utils.PlainHexBytes `json:"serverPubKey" toml:"serverPubKey"`

	ChannelDefinitionsContractAddress   common.Address `json:"channelDefinitionsContractAddress" toml:"channelDefinitionsContractAddress"`
	ChannelDefinitionsContractFromBlock int64          `json:"channelDefinitionsContractFromBlock" toml:"channelDefinitionsContractFromBlock"`
}

func (p PluginConfig) Validate() (merr error) {
	if p.RawServerURL == "" {
		merr = errors.New("streams: ServerURL must be specified")
	} else {
		var normalizedURI string
		if schemeRegexp.MatchString(p.RawServerURL) {
			normalizedURI = p.RawServerURL
		} else {
			normalizedURI = fmt.Sprintf("wss://%s", p.RawServerURL)
		}
		uri, err := url.ParseRequestURI(normalizedURI)
		if err != nil {
			merr = pkgerrors.Wrap(err, "streams: invalid value for ServerURL")
		} else if uri.Scheme != "wss" {
			merr = pkgerrors.Errorf(`streams: invalid scheme specified for MercuryServer, got: %q (scheme: %q) but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`, p.RawServerURL, uri.Scheme)
		}
	}

	if len(p.ServerPubKey) != 32 {
		merr = errors.Join(merr, errors.New("streams: ServerPubKey is required and must be a 32-byte hex string"))
	}

	return merr
}

var schemeRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)
var wssRegexp = regexp.MustCompile(`^wss://`)

func (p PluginConfig) ServerURL() string {
	return wssRegexp.ReplaceAllString(p.RawServerURL, "")
}
