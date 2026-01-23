// This file overrides the signer extraction in keystone.go to use OnchainPublicKey
// for LLO tests. It's always compiled and checks an environment variable to determine
// whether to apply the override.
//
// Why this override is necessary:
//
// LLO (Low-Latency Oracle) reports are cryptographically signed by LLO DON nodes using
// their OnchainPublicKey (a 20-byte Ethereum address). The signature verification in
// SignedReportRemoteAggregator validates signatures against signer addresses stored in
// the capability registry.
//
// The default signer extraction in keystone.go uses OffchainPublicKey, which is correct
// for most OCR jobs. However, for LLO reports, we must use OnchainPublicKey because:
// 1. LLO reports are signed with OnchainPublicKey as part of the OCR protocol
// 2. The capability registry stores signers as OnchainPublicKey addresses
// 3. Signature verification compares recovered signer addresses against registry entries
//
// This override ensures that the capability registry uses the correct signer addresses
// for LLO signature verification, enabling the end-to-end test to work correctly.
//
// NOTE: Set USE_LLO_ONCHAIN_SIGNER=true when running LLO tests to enable this override.

package contracts

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
)

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
