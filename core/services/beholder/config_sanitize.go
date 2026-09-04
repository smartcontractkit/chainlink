package beholder

import (
	"fmt"
	"net/url"

	gotoml "github.com/pelletier/go-toml/v2"
)

const (
	redactedUsername    = "userxxx"
	redactedPassword    = "passwordxxx"
	queryPlaceholder    = "<QUERY>"
	fragmentPlaceholder = "<FRAGMENT>"
)

// SanitizeConfigTOML removes credentials and query/fragment secrets from URL
// values in a TOML config string. It is used to sanitize the node configuration
// before it is emitted through the Beholder message emitter.
//
// For every string value that parses as a URL with a non-empty scheme and host,
// the following transformations are applied:
//   - usernames are replaced with "userxxx" (if present)
//   - passwords are replaced with "passwordxxx" (if present)
//   - query strings are replaced with "?<QUERY>" (if present)
//   - fragments are replaced with "#<FRAGMENT>" (if present)
//   - paths, ports, and schemes are left untouched
//
// Clean URLs (without userinfo, query, or fragment) and non-URL strings are
// returned unchanged. The TOML is decoded to a generic map, sanitized, and
// re-encoded, so keys may be reordered alphabetically. The resulting TOML is
// semantically identical to the input except for the sanitized URL values.
func SanitizeConfigTOML(in string) (string, error) {
	var m map[string]any
	if err := gotoml.Unmarshal([]byte(in), &m); err != nil {
		return "", fmt.Errorf("decode config TOML: %w", err)
	}
	sanitizeValue(m)
	out, err := gotoml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("encode config TOML: %w", err)
	}
	return string(out), nil
}

func sanitizeValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		for k, sub := range val {
			val[k] = sanitizeValue(sub)
		}
		return val
	case []any:
		for i, sub := range val {
			val[i] = sanitizeValue(sub)
		}
		return val
	case string:
		return sanitizeURLString(val)
	default:
		return v
	}
}

func sanitizeURLString(s string) string {
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return s
	}

	changed := false
	hasFragment := u.Fragment != "" || u.RawFragment != ""

	if u.User != nil {
		username := u.User.Username()
		_, hasPass := u.User.Password()
		switch {
		case username != "" && hasPass:
			u.User = url.UserPassword(redactedUsername, redactedPassword)
		case username != "":
			u.User = url.User(redactedUsername)
		case hasPass:
			u.User = url.UserPassword("", redactedPassword)
		default:
			u.User = url.User("")
		}
		changed = true
	}

	if u.RawQuery != "" || u.ForceQuery {
		u.RawQuery = queryPlaceholder
		u.ForceQuery = false
		changed = true
	}

	if hasFragment {
		u.Fragment = ""
		u.RawFragment = ""
		changed = true
	}

	if !changed {
		return s
	}

	out := u.String()
	if hasFragment {
		out += "#" + fragmentPlaceholder
	}
	return out
}
