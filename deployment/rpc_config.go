package deployment

import "github.com/pkg/errors"

type URLSchemePreference int

const (
	URLSchemePreferenceNone URLSchemePreference = iota
	URLSchemePreferenceWS
	URLSchemePreferenceHTTP
)

type RPC struct {
	Name               string
	WSURL              string
	HTTPURL            string
	PreferredURLScheme URLSchemePreference
}

type RPCConfig struct {
	ChainName string
	RPCs      []RPC
}

// ToEndpoint returns the correct endpoint based on the preferred URL scheme
// If the preferred URL scheme is not set, it will return the WS URL
// If the preferred URL scheme is set to WS, it will return the WS URL
// If the preferred URL scheme is set to HTTP, it will return the HTTP URL
func (r RPC) ToEndpoint() (string, error) {
	switch r.PreferredURLScheme {
	case URLSchemePreferenceNone, URLSchemePreferenceWS:
		return r.WSURL, nil
	case URLSchemePreferenceHTTP:
		return r.HTTPURL, nil
	default:
		return "", errors.New("Unknown URLSchemePreference")
	}
}
