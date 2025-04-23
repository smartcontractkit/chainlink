package deployment

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type OnchainClient = deployment.OnchainClient

var MaybeDataErr = deployment.MaybeDataErr

// RPC Config

type (
	URLSchemePreference = deployment.URLSchemePreference
	RPC                 = deployment.RPC
	RPCConfig           = deployment.RPCConfig
)

var (
	URLSchemePreferenceFromString = deployment.URLSchemePreferenceFromString
	URLSchemePreferenceNone       = deployment.URLSchemePreferenceNone
	URLSchemePreferenceWS         = deployment.URLSchemePreferenceWS
	URLSchemePreferenceHTTP       = deployment.URLSchemePreferenceHTTP
)

// MultiClient

type (
	MultiClient = deployment.MultiClient
	RetryConfig = deployment.RetryConfig
)

var (
	NewMultiClient              = deployment.NewMultiClient
	RPCDefaultRetryAttempts     = deployment.RPCDefaultRetryAttempts
	RPCDefaultRetryDelay        = deployment.RPCDefaultRetryDelay
	RPCDefaultDialRetryAttempts = deployment.RPCDefaultDialRetryAttempts
	RPCDefaultDialRetryDelay    = deployment.RPCDefaultDialRetryDelay
)
