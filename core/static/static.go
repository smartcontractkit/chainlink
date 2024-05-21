package static

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Version and Sha are set at compile time via build arguments.
var (
	// Version is the semantic version of the build or Unset.
	Version = Unset
	// Sha is the commit hash of the build or Unset.
	Sha = Unset
)

// InitTime holds the initial start timestamp.
var InitTime = time.Now()

const (
	// Unset is a sentinel value.
	Unset = "unset"
	// ExternalInitiatorAccessKeyHeader is the header name for the access key
	// used by external initiators to authenticate
	ExternalInitiatorAccessKeyHeader = "X-Chainlink-EA-AccessKey"
	// ExternalInitiatorSecretHeader is the header name for the secret used by
	// external initiators to authenticate
	ExternalInitiatorSecretHeader = "X-Chainlink-EA-Secret"
)

func buildPrettyVersion() string {
	if Version == Unset {
		return " "
	}
	return fmt.Sprintf(" %s ", Version)
}

// SetConsumerName sets a nicely formatted application_name on the
// database uri
func SetConsumerName(uri *url.URL, name string, id *uuid.UUID) {
	q := uri.Query()

	applicationName := fmt.Sprintf("Chainlink%s|%s", buildPrettyVersion(), name)
	if id != nil {
		applicationName += fmt.Sprintf("|%s", id.String())
	}
	if len(applicationName) > 63 {
		applicationName = applicationName[:63]
	}
	q.Set("application_name", applicationName)
	uri.RawQuery = q.Encode()
}

// Short returns a 7-character sha prefix and version, or Unset if blank.
func Short() (shaPre string, ver string) {
	return short(Sha, Version)
}

func short(sha, ver string) (string, string) {
	if sha == "" {
		sha = Unset
	} else if len(sha) > 7 {
		sha = sha[:7]
	}
	if ver == "" {
		ver = Unset
	}
	return sha, ver
}
