package static

import (
	"fmt"
	"net/url"

	uuid "github.com/satori/go.uuid"
)

// Version the version of application
var Version = "unset"

// Sha string "unset"
var Sha = "unset"

// InstanceUUID is generated on startup and uniquely identifies this instance of Chainlink
var InstanceUUID uuid.UUID

const (
	// ExternalInitiatorAccessKeyHeader is the header name for the access key
	// used by external initiators to authenticate
	ExternalInitiatorAccessKeyHeader = "X-Chainlink-EA-AccessKey"
	// ExternalInitiatorSecretHeader is the header name for the secret used by
	// external initiators to authenticate
	ExternalInitiatorSecretHeader = "X-Chainlink-EA-Secret"
)

func init() {
	InstanceUUID = uuid.NewV4()
}

func buildPrettyVersion() string {
	if Version == "unset" {
		return " "
	}
	return fmt.Sprintf(" %s ", Version)
}

// SetConsumerName sets a nicely formatted application_name on the
// database uri
func SetConsumerName(uri *url.URL, name string) {
	q := uri.Query()

	applicationName := fmt.Sprintf("Chainlink%s| %s | %s", buildPrettyVersion(), name, InstanceUUID)
	if len(applicationName) > 63 {
		applicationName = applicationName[:63]
	}
	q.Set("application_name", applicationName)
	uri.RawQuery = q.Encode()
}
