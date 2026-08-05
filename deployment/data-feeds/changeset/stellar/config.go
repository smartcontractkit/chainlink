package stellar

import "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

// Datastore contract types for the two DF contracts.
const (
	CacheContract datastore.ContractType = "DataFeedsCache"
	ProxyContract datastore.ContractType = "DataFeedsProxy"
)
