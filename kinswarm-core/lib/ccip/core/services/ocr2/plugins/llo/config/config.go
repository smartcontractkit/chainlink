// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"

	"github.com/ethereum/go-ethereum/common"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	mercuryconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type PluginConfig struct {
	ChannelDefinitionsContractAddress   common.Address `json:"channelDefinitionsContractAddress" toml:"channelDefinitionsContractAddress"`
	ChannelDefinitionsContractFromBlock int64          `json:"channelDefinitionsContractFromBlock" toml:"channelDefinitionsContractFromBlock"`

	// NOTE: ChannelDefinitions is an override.
	// If ChannelDefinitions is specified, values for
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

	DonID uint32 `json:"donID" toml:"donID"`

	// Mercury servers
	Servers map[string]utils.PlainHexBytes `json:"servers" toml:"servers"`
}

func (p *PluginConfig) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p PluginConfig) GetServers() (servers []mercuryconfig.Server) {
	for url, pubKey := range p.Servers {
		servers = append(servers, mercuryconfig.Server{URL: wssRegexp.ReplaceAllString(url, ""), PubKey: pubKey})
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].URL < servers[j].URL
	})
	return
}

func (p PluginConfig) Validate() (merr error) {
	if p.DonID == 0 {
		merr = errors.Join(merr, errors.New("llo: DonID must be specified and not zero"))
	}

	if len(p.Servers) == 0 {
		merr = errors.Join(merr, errors.New("llo: At least one Mercury server must be specified"))
	} else {
		for serverName, serverPubKey := range p.Servers {
			if err := validateURL(serverName); err != nil {
				merr = errors.Join(merr, fmt.Errorf("llo: invalid value for ServerURL: %w", err))
			}
			if len(serverPubKey) != 32 {
				merr = errors.Join(merr, errors.New("llo: ServerPubKey must be a 32-byte hex string"))
			}
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
		// TODO: Verify Opts format here?
		// MERC-3524
	} else {
		if p.ChannelDefinitionsContractAddress == (common.Address{}) {
			merr = errors.Join(merr, errors.New("llo: ChannelDefinitionsContractAddress is required if ChannelDefinitions is not specified"))
		}
	}

	merr = errors.Join(merr, validateKeyBundleIDs(p.KeyBundleIDs))

	return merr
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
		return fmt.Errorf(`llo: invalid value for ServerURL, got: %q`, rawServerURL)
	}
	if uri.Scheme != "wss" {
		return fmt.Errorf(`llo: invalid scheme specified for MercuryServer, got: %q (scheme: %q) but expected a websocket url e.g. "192.0.2.2:4242" or "wss://192.0.2.2:4242"`, rawServerURL, uri.Scheme)
	}
	return nil
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
