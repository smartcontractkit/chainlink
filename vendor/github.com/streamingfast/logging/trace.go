package logging

import (
	"os"
	"regexp"
	"strings"
)

// IsTraceEnabled receives the a short name value (usually the app name) and the fully qualified
// package identifier (i.e. `github.com/dfuse-io/logging/subpackage`) and determines if tracing should
// be enabled for such values.
//
// To take the decision, this function inspects the `TRACE` environment variable. If the variable
// is **not** set, tracing is disabled. If the value is either `true`, `*` or `.*`, then tracing
// is enabled.
//
// For other values, we split the `TRACE` value using `,` separator. For each split element, if
// the element matches directly the short name, tracing is enabled and if the element as a Regexp
// object matches (partially, not fully) the `packageID`, tracing is enabled.
//
// The method also supports deny elements. If you prefix the element with `-`, it means disable
// trace for this match. This is applied after all allow element has been processed, so it's possible
// to enabled except a specific package (i.e. `TRACE=.*,-github.com/specific/package`).
//
// In all other cases, tracing is disabled.
//
// Deprecated: Define your logger and `Tracer` directly with `var zlog, tracer = logging.PackageLogger(<shortName>, "...")`
// instead of separately, `tracer.Enabled()` can then be used to determine if tracing should be enabled (can be enable dynamically).
func IsTraceEnabled(shortName string, packageID string) bool {
	trace := os.Getenv("TRACE")
	if trace == "" {
		return false
	}

	if trace == "true" || trace == "TRUE" || trace == "*" || trace == ".*" {
		return true
	}

	isEnabled := false
	for _, filter := range strings.Split(trace, ",") {
		if logFilter(filter).isAllowed(shortName, packageID) {
			isEnabled = true
			break
		}
	}

	// Now if it's denied, it should not be enabled
	for _, filter := range strings.Split(trace, ",") {
		if logFilter(filter).isDenied(shortName, packageID) {
			isEnabled = false
			break
		}
	}

	return isEnabled
}

type logFilter string

func (l logFilter) isAllowed(name string, packageID string) bool {
	if len(l) == 0 || l[0] == '-' {
		return false
	}

	return string(l) == name || regexp.MustCompile(string(l)).MatchString(packageID)
}

func (l logFilter) isDenied(name string, packageID string) bool {
	if len(l) == 0 || l[0] != '-' {
		return false
	}

	query := string(l[1:])
	if len(query) == 0 {
		return false
	}

	return matchPackage(query, name, packageID)
}

func matchPackage(query string, name string, packageID string) bool {
	if query == name || query == packageID {
		return true
	}

	regex, err := regexp.Compile(query)
	if (err == nil && regex.MatchString(packageID)) || (err != nil && query == packageID) {
		return true
	}

	return false
}
