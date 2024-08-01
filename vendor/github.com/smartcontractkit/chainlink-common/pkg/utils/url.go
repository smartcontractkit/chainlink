package utils

import "github.com/smartcontractkit/chainlink-common/pkg/config"

// Deprecated: use config.URL
type URL = config.URL

// Deprecated: use config.ParseURL
func ParseURL(s string) (*URL, error) { return config.ParseURL(s) }

// Deprecated: use config.MustParseURL
func MustParseURL(s string) *URL { return config.MustParseURL(s) }
