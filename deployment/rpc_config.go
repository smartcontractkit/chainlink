package deployment

import (
	"errors"
	"fmt"
	"strings"
)

type URLSchemePreference int

const (
	URLSchemePreferenceNone URLSchemePreference = iota
	URLSchemePreferenceWS
	URLSchemePreferenceHTTP
)

func URLSchemePreferenceFromString(s string) (URLSchemePreference, error) {
	switch strings.ToLower(s) {
	case "none":
		return URLSchemePreferenceNone, nil
	case "ws":
		return URLSchemePreferenceWS, nil
	case "http":
		return URLSchemePreferenceHTTP, nil
	default:
		return URLSchemePreferenceNone, fmt.Errorf("invalid URLSchemePreference: %s", s)
	}
}

func (u *URLSchemePreference) UnmarshalText(text []byte) error {
	preference, err := URLSchemePreferenceFromString(string(text))
	if err != nil {
		return err
	}
	*u = preference

	return nil
}

type RPC struct {
	Name               string
	WSURL              string
	HTTPURL            string
	PreferredURLScheme URLSchemePreference
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

// RPCConfig is a configuration for a chain.
// It contains a chain selector and a list of RPCs
type RPCConfig struct {
	ChainSelector uint64
	RPCs          []RPC
}
