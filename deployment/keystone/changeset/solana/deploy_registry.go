package solana

import "github.com/smartcontractkit/chainlink/deployment/helpers"

type DeployRequest = struct {
	ChainSel    uint64
	BuildConfig *helpers.BuildSolanaConfig
}
