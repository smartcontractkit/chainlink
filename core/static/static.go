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

//nolint
const (
	EvmMaxInFlightTransactionsWarningLabel   = `WARNING: If this happens a lot, you may need to increase ETH_MAX_IN_FLIGHT_TRANSACTIONS to boost your node's transaction throughput, however you do this at your own risk. You MUST first ensure your ethereum node is configured not to ever evict local transactions that exceed this number otherwise the node can get permanently stuck. See the documentation for more details: https://docs.chain.link/docs/configuration-variables/`
	EvmMaxQueuedTransactionsLabel            = `WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. Unless you are operating with very high throughput, this error is unlikely to be a problem with your Chainlink node configuration, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. It is recommended to run Chainlink with multiple primary and sendonly nodes for redundancy and to ensure fast and reliable transaction propagation. Increasing ETH_MAX_QUEUED_TRANSACTIONS will allow Chainlink to buffer more unsent transactions, but you should only do this if you need very high burst transmission rates. If you don't need very high burst throughput, increasing this limit is not the correct action to take here and will probably make things worse`
	EthNodeConnectivityProblemLabel          = `WARNING: If this happens a lot, it may be a sign that your eth node has a connectivity problem, and your transactions are not making it to any miners. It is recommended to run Chainlink with multiple primary and sendonly nodes for redundancy and to ensure fast and reliable transaction propagation.`
	EvmRPCTxFeeCapConfiguredIncorrectlyLabel = `WARNING: Gas price was rejected by the eth node for being too high. By default, go-ethereum (and clones) have a built-in upper limit for gas price. It is preferable to disable this and use Chainlink's internal gas limits instead. Your RPC node's RPCTxFeeCap needs to be increased (recommended configuration: --rpc.gascap=0 --rpc.txfeecap=0). If you want to limit Chainlink's max gas price, do so by setting ETH_MAX_GAS_PRICE_WEI on the Chainlink node instead. See the docs for more details: https://docs.chain.link/docs/configuration-variables/`
)
