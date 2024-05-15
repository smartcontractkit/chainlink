// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/ethereum/go-ethereum/common"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type PluginConfig struct {
	RawServerURL string              `json:"serverURL" toml:"serverURL"`
	ServerPubKey utils.PlainHexBytes `json:"serverPubKey" toml:"serverPubKey"`

	ChannelDefinitionsContractAddress   common.Address `json:"channelDefinitionsContractAddress" toml:"channelDefinitionsContractAddress"`
	ChannelDefinitionsContractFromBlock int64          `json:"channelDefinitionsContractFromBlock" toml:"channelDefinitionsContractFromBlock"`

	// NOTE: ChannelDefinitions is an override.
	// If Channe}lDefinitions is specified, values for
	// ChannelDefinitionsContractAddress and
	// ChannelDefinitionsContractFromBlock will be ignored
	ChannelDefinitions string `json:"channelDefinitions" toml:"channelDefinitions"`

	// BenchmarkMode is a flag to enable benchmarking mode. In this mode, the
	// transmitter will not transmit anything at all and instead emit
	// logs/metrics.
	BenchmarkMode bool `json:"benchmarkMode" toml:"benchmarkMode"`

	// KeyBundleIDs maps supported keys to their respective bundle IDs
	// Key must match llo's ReportFormat
	KeyBundleIDs map[string]string `json:"keyBundleIDs" toml:"keyBundleIDs"`
}

func (p PluginConfig) Validate() (merr error) {
	if p.RawServerURL == "" {
		merr = errors.New("llo: ServerURL must be specified")
	} else {
		var normalizedURI string
		if schemeRegexp.MatchString(p.RawServerURL) {
			normalizedURI = p.RawServerURL
		} else {
			normalizedURI = fmt.Sprintf("wss://%s", p.RawServerURL)
		}
		uri, err := url.ParseRequestURI(normalizedURI)
		if err != nil {
			merr = fmt.Errorf("llo: invalid value for ServerURL: %w", err)
		} else if uri.Scheme != "wss" {
			merr = fmt.Errorf(`llo: invalid scheme specified for MercuryServer, got: %q (scheme: %q) but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`, p.RawServerURL, uri.Scheme)
		}
	}

	if p.ChannelDefinitions != "" {
		if p.ChannelDefinitionsContractAddress != (common.Address{}) {
			merr = errors.Join(merr, errors.New("llo: ChannelDefinitionsContractAddress is not allowed if ChannelDefinitions is specified"))
		}
		if p.ChannelDefinitionsContractFromBlock != 0 {
			merr = errors.Join(merr, errors.New("llo: ChannelDefinitionsContractFromBlock is not allowed if ChannelDefinitions is specified"))
		}
		var cd llotypes.ChannelDefinitions
		if err := json.Unmarshal([]byte(p.ChannelDefinitions), &cd); err != nil {
			merr = errors.Join(merr, fmt.Errorf("channelDefinitions is invalid JSON: %w", err))
		}
	} else {
		if p.ChannelDefinitionsContractAddress == (common.Address{}) {
			merr = errors.Join(merr, errors.New("llo: ChannelDefinitionsContractAddress is required if ChannelDefinitions is not specified"))
		}
	}

	if len(p.ServerPubKey) != 32 {
		merr = errors.Join(merr, errors.New("llo: ServerPubKey is required and must be a 32-byte hex string"))
	}

	merr = errors.Join(merr, validateKeyBundleIDs(p.KeyBundleIDs))

	return merr
}

func validateKeyBundleIDs(keyBundleIDs map[string]string) error {
	for k, v := range keyBundleIDs {
		if k == "" {
			return errors.New("llo: KeyBundleIDs: key must not be empty")
		}
		if v == "" {
			return errors.New("llo: KeyBundleIDs: value must not be empty")
		}
		if _, err := llotypes.ReportFormatFromString(k); err != nil {
			return fmt.Errorf("llo: KeyBundleIDs: key must be a recognized report format, got: %s (err: %w)", k, err)
		}
		if !chaintype.IsSupportedChainType(chaintype.ChainType(k)) {
			return fmt.Errorf("llo: KeyBundleIDs: key must be a supported chain type, got: %s", k)
		}
	}
	return nil
}

var schemeRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)
var wssRegexp = regexp.MustCompile(`^wss://`)

func (p PluginConfig) ServerURL() string {
	return wssRegexp.ReplaceAllString(p.RawServerURL, "")
}
