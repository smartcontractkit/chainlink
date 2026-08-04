package stellar

import "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

const (
	// CacheContract is the datastore contract type for the DataFeedsCache Soroban contract.
	CacheContract datastore.ContractType = "DataFeedsCache"
	// ProxyContract is the datastore contract type for the DataFeedsProxy Soroban contract.
	ProxyContract datastore.ContractType = "DataFeedsProxy"
)
