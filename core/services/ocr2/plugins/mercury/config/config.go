// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/null"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type PluginConfig struct {
	// Must either specify details for single server OR multiple servers.
	// Specifying both is not valid.

	// Single mercury server
	// LEGACY: This is the old way of specifying a mercury server
	RawServerURL string              `json:"serverURL" toml:"serverURL"`
	ServerPubKey utils.PlainHexBytes `json:"serverPubKey" toml:"serverPubKey"`

	// Multi mercury servers
	// This is the preferred way to specify mercury server(s)
	Servers map[string]utils.PlainHexBytes `json:"servers" toml:"servers"`

	// InitialBlockNumber allows to set a custom "validFromBlockNumber" for
	// the first ever report in the case of a brand new feed, where the mercury
	// server does not have any previous reports. For a brand new feed, this
	// effectively sets the "first" validFromBlockNumber.
	InitialBlockNumber null.Int64 `json:"initialBlockNumber" toml:"initialBlockNumber"`

	LinkFeedID   *mercuryutils.FeedID `json:"linkFeedID" toml:"linkFeedID"`
	NativeFeedID *mercuryutils.FeedID `json:"nativeFeedID" toml:"nativeFeedID"`
}

func validateURL(rawServerURL string) error {
	var normalizedURI string
	if schemeRegexp.MatchString(rawServerURL) {
		normalizedURI = rawServerURL
	} else {
		normalizedURI = fmt.Sprintf("wss://%s", rawServerURL)
	}
	uri, err := url.ParseRequestURI(normalizedURI)
	if err != nil {
		return pkgerrors.Errorf(`Mercury: invalid value for ServerURL, got: %q`, rawServerURL)
	}
	if uri.Scheme != "wss" {
		return pkgerrors.Errorf(`Mercury: invalid scheme specified for MercuryServer, got: %q (scheme: %q) but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`, rawServerURL, uri.Scheme)
	}
	return nil
}

type Server struct {
	URL    string
	PubKey utils.PlainHexBytes
}

func (p PluginConfig) GetServers() (servers []Server) {
	if p.RawServerURL != "" {
		return []Server{{URL: wssRegexp.ReplaceAllString(p.RawServerURL, ""), PubKey: p.ServerPubKey}}
	}
	for url, pubKey := range p.Servers {
		servers = append(servers, Server{URL: wssRegexp.ReplaceAllString(url, ""), PubKey: pubKey})
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].URL < servers[j].URL
	})
	return
}

func ValidatePluginConfig(config PluginConfig, feedID mercuryutils.FeedID) (merr error) {
	if len(config.Servers) > 0 {
		if config.RawServerURL != "" || len(config.ServerPubKey) != 0 {
			merr = errors.Join(merr, errors.New("Mercury: Servers and RawServerURL/ServerPubKey may not be specified together"))
		} else {
			for serverName, serverPubKey := range config.Servers {
				if err := validateURL(serverName); err != nil {
					merr = errors.Join(merr, pkgerrors.Wrap(err, "Mercury: invalid value for ServerURL"))
				}
				if len(serverPubKey) != 32 {
					merr = errors.Join(merr, errors.New("Mercury: ServerPubKey must be a 32-byte hex string"))
				}
			}
		}
	} else if config.RawServerURL == "" {
		merr = errors.Join(merr, errors.New("Mercury: Servers must be specified"))
	} else {
		if err := validateURL(config.RawServerURL); err != nil {
			merr = errors.Join(merr, pkgerrors.Wrap(err, "Mercury: invalid value for ServerURL"))
		}
		if len(config.ServerPubKey) != 32 {
			merr = errors.Join(merr, errors.New("Mercury: If RawServerURL is specified, ServerPubKey is also required and must be a 32-byte hex string"))
		}
	}

	switch feedID.Version() {
	case 1:
		if config.LinkFeedID != nil {
			merr = errors.Join(merr, errors.New("linkFeedID may not be specified for v1 jobs"))
		}
		if config.NativeFeedID != nil {
			merr = errors.Join(merr, errors.New("nativeFeedID may not be specified for v1 jobs"))
		}
	case 2, 3, 4:
		if config.LinkFeedID == nil {
			merr = errors.Join(merr, fmt.Errorf("linkFeedID must be specified for v%d jobs", feedID.Version()))
		}
		if config.NativeFeedID == nil {
			merr = errors.Join(merr, fmt.Errorf("nativeFeedID must be specified for v%d jobs", feedID.Version()))
		}
		if config.InitialBlockNumber.Valid {
			merr = errors.Join(merr, fmt.Errorf("initialBlockNumber may not be specified for v%d jobs", feedID.Version()))
		}
	default:
		merr = errors.Join(merr, fmt.Errorf("got unsupported schema version %d; supported versions are 1,2,3,4", feedID.Version()))
	}

	return merr
}

var schemeRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)
var wssRegexp = regexp.MustCompile(`^wss://`)
