// Package ref formats manifest versions for consumers (go install, docker tags, etc.).
package ref

import (
	"strings"
)

// ForInstall prepends "v" to semver-looking versions; SHA and other pins pass through.
func ForInstall(version string) string {
	if semverLike(version) {
		return "v" + version
	}
	return version
}

// ForConsumer is an alias for ForInstall (docker tags, golangci-lint-action, etc.).
func ForConsumer(version string) string {
	return ForInstall(version)
}

func semverLike(version string) bool {
	if version == "" {
		return false
	}
	dot := strings.IndexByte(version, '.')
	if dot <= 0 {
		return false
	}
	if version[0] < '0' || version[0] > '9' {
		return false
	}
	if dot+1 >= len(version) {
		return false
	}
	return version[dot+1] >= '0' && version[dot+1] <= '9'
}
