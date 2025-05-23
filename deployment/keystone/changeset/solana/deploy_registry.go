package solana

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

type DeployRequest = struct {
	ChainSel    uint64
	Qualifier   string
	Labels      *datastore.LabelSet
	BuildConfig *BuildSolanaConfig
}
