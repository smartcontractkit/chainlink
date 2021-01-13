package static

// Version the version of application
var Version = "unset"

// Sha string "unset"
var Sha = "unset"

const (
	// ExternalInitiatorAccessKeyHeader is the header name for the access key
	// used by external initiators to authenticate
	ExternalInitiatorAccessKeyHeader = "X-Chainlink-EA-AccessKey"
	// ExternalInitiatorSecretHeader is the header name for the secret used by
	// external initiators to authenticate
	ExternalInitiatorSecretHeader = "X-Chainlink-EA-Secret"
)
