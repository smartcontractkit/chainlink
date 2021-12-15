package static

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Masterminds/semver/v3"
	uuid "github.com/satori/go.uuid"
)

// Version is the version of application
// Must be either "unset" or valid semver
var Version = "unset"

// Sha string "unset"
var Sha = "unset"

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
		return
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

const (
	EvmMaxInFlightTransactionsWarningLabel = `WARNING: If this happens a lot, you may need to increase ETH_MAX_IN_FLIGHT_TRANSACTIONS to boost your node's transaction throughput, however you do this at your own risk. You MUST first ensure your ethereum node is configured not to ever evict local transactions that exceed this number otherwise the node can get permanently stuck`
	EvmMaxQueuedTransactionsLabel          = `WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse`
	EthNodeConnectivityProblemLabel        = `WARNING: If this happens a lot, it may be a sign that your eth node has a connectivity problem, and your transactions are not making it to any miners`
)
