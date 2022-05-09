package static

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Masterminds/semver/v3"
	uuid "github.com/satori/go.uuid"
)

// Version is the version of application
// Must be either "unset" or valid semver
var Version = "unset"

// Sha string "unset"
var Sha = "unset"

// InitTime holds the initial start timestamp, based on this package's init func.
var InitTime time.Time

const (
	// ExternalInitiatorAccessKeyHeader is the header name for the access key
	// used by external initiators to authenticate
	ExternalInitiatorAccessKeyHeader = "X-Chainlink-EA-AccessKey"
	// ExternalInitiatorSecretHeader is the header name for the secret used by
	// external initiators to authenticate
	ExternalInitiatorSecretHeader = "X-Chainlink-EA-Secret"
)

func init() {
	InitTime = time.Now()

	checkVersion()
}

func checkVersion() {
	if Version == "unset" {
		if os.Getenv("CHAINLINK_DEV") == "true" {
			return
		}
		log.Println(`Version was unset but CHAINLINK_DEV was not set to "true". Chainlink should be built with static.Version set to a valid semver for production builds.`)
	} else if _, err := semver.NewVersion(Version); err != nil {
		panic(fmt.Sprintf("Version invalid: %q is not valid semver", Version))
	}
}

func buildPrettyVersion() string {
	if Version == "unset" {
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
