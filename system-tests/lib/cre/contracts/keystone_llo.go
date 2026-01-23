// This file overrides the signer extraction in keystone.go to use OnchainPublicKey
// for LLO tests. It's always compiled and checks an environment variable to determine
// whether to apply the override.

package contracts

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
)

// Override the signer extraction for LLO tests to use OnchainPublicKey.
// This is required because LLO reports are signed with OnchainPublicKey,
// and the signature verification checks against the signer address in the registry.
//
// This file overrides extractSignerAddressFromOCRConfig to use OnchainPublicKey
// instead of OffchainPublicKey when the USE_LLO_ONCHAIN_SIGNER environment variable is set.
//
// NOTE: Set USE_LLO_ONCHAIN_SIGNER=true when running LLO tests to enable this override.
func init() {
	if os.Getenv("USE_LLO_ONCHAIN_SIGNER") == "true" {
		extractSignerAddressFromOCRConfig = func(ocrCfg deployment.OCRConfig) [32]byte {
			var signer [32]byte
			if len(ocrCfg.OnchainPublicKey) >= 20 {
				// Use common.BytesToAddress to extract the 20-byte address (takes last 20 bytes)
				addr := common.BytesToAddress(ocrCfg.OnchainPublicKey)
				copy(signer[0:20], addr.Bytes())
			} else {
				panic(fmt.Errorf("OnchainPublicKey too short: %d bytes", len(ocrCfg.OnchainPublicKey)))
			}
			return signer
		}
	}
}
