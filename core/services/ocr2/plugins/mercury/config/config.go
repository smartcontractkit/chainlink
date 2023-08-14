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

	"github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	decoder "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/encoding"
)

type PluginConfig struct {
	RawServerURL string              `json:"serverURL" toml:"serverURL"`
	ServerPubKey utils.PlainHexBytes `json:"serverPubKey" toml:"serverPubKey"`
	// InitialBlockNumber allows to set a custom "validFromBlockNumber" for
	// the first ever report in the case of a brand new feed, where the mercury
	// server does not have any previous reports. For a brand new feed, this
	// effectively sets the "first" validFromBlockNumber.
	InitialBlockNumber null.Int64 `json:"initialBlockNumber" toml:"initialBlockNumber"`

	LinkFeedID   *common.Hash `json:"linkFeedID" toml:"linkFeedID"`
	NativeFeedID *common.Hash `json:"nativeFeedID" toml:"nativeFeedID"`
}

func ValidatePluginConfig(config PluginConfig, feedId [32]byte) (merr error) {
	if config.RawServerURL == "" {
		merr = errors.New("mercury: ServerURL must be specified")
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
		} else if uri.Scheme != "wss" {
			merr = pkgerrors.Errorf(`Mercury: invalid scheme specified for MercuryServer, got: %q (scheme: %q) but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`, config.RawServerURL, uri.Scheme)
		}
	}

	if len(config.ServerPubKey) != 32 {
		merr = errors.Join(merr, errors.New("mercury: ServerPubKey is required and must be a 32-byte hex string"))
	}

	reportSchemaVersion, err := decoder.DecodeSchemaVersionFromFeedId(feedId)
	if err != nil {
		merr = errors.Join(merr, err)
	} else {
		switch reportSchemaVersion {
		case 1:
			if !config.InitialBlockNumber.Valid {
				merr = errors.Join(merr, errors.New("initialBlockNumber must be specified for v1 jobs"))
			}
			if config.LinkFeedID != nil {
				merr = errors.Join(merr, errors.New("linkFeedID may not be specified for v1 jobs"))
			}
			if config.NativeFeedID != nil {
				merr = errors.Join(merr, errors.New("nativeFeedID may not be specified for v1 jobs"))
			}
		case 2, 3:
			if config.LinkFeedID == nil {
				merr = errors.Join(merr, fmt.Errorf("linkFeedID must be specified for v%d jobs", reportSchemaVersion))
			}
			if config.NativeFeedID == nil {
				merr = errors.Join(merr, fmt.Errorf("nativeFeedID must be specified for v%d jobs", reportSchemaVersion))
			}
			if config.InitialBlockNumber.Valid {
				merr = errors.Join(merr, fmt.Errorf("initialBlockNumber may not be specified for v%d jobs", reportSchemaVersion))
			}
		default:
			merr = errors.Join(merr, fmt.Errorf("got unsupported schema version %d; supported versions are 1,2,3", reportSchemaVersion))
		}
	}

	return merr
}

var schemeRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)
var wssRegexp = regexp.MustCompile(`^wss://`)

func (p PluginConfig) ServerURL() string {
	return wssRegexp.ReplaceAllString(p.RawServerURL, "")
}
