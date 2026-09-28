package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"

	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ connector.Signer = ocr2KeySigner{}

// NewOCR2KeySigner wraps the node's OCR2 key bundle as a vault envelope signer
// so responses are signed with the same secp256k1 key that signs OCR3 reports.
// The recovered address matches the node's registry Signer entry, letting the
// gateway (and clients) verify envelope signatures against the publicly known
// DON member signer set.
func NewOCR2KeySigner(kb ocr2key.KeyBundle) (connector.Signer, error) {
	if kb == nil {
		return nil, errors.New("OCR2 key bundle is nil")
	}
	if kb.ChainType() != corekeys.EVM {
		return nil, fmt.Errorf("vault envelope signing requires an EVM OCR2 key bundle, got chain type %q", kb.ChainType())
	}
	return ocr2KeySigner{kb: kb}, nil
}

type ocr2KeySigner struct {
	kb ocr2key.KeyBundle
}

// Sign returns an EIP-191 style signature over the flattened data, matching the
// recovery convention of utils.GetSignersEthAddress (used by the gateway to
// verify envelope signatures).
func (s ocr2KeySigner) Sign(_ context.Context, data ...[]byte) ([]byte, error) {
	msg := gwcommon.Flatten(data...)
	hash := utils.GenerateEthPrefixedMsgHash(msg)
	sig, err := s.kb.SignBlob(hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign vault response digest: %w", err)
	}
	return sig, nil
}
