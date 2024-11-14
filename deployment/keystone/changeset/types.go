package changeset

import (
	"encoding/hex"
	"errors"
	"fmt"

	v1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func NewP2PSignerEncFromJD(ccfg *v1.ChainConfig, pubkeyStr string) (*P2PSignerEnc, error) {
	if ccfg == nil {
		return nil, errors.New("nil ocr2config")
	}
	var pubkey [32]byte
	if _, err := hex.Decode(pubkey[:], []byte(pubkeyStr)); err != nil {
		return nil, fmt.Errorf("failed to decode pubkey %s: %w", pubkey, err)
	}
	ocfg := ccfg.Ocr2Config
	p2p := p2pkey.PeerID{}
	if err := p2p.UnmarshalString(ocfg.P2PKeyBundle.PeerId); err != nil {
		return nil, fmt.Errorf("failed to unmarshal peer id %s: %w", ocfg.P2PKeyBundle.PeerId, err)
	}

	signer := ocfg.OcrKeyBundle.OnchainSigningAddress
	if len(signer) != 40 {
		return nil, fmt.Errorf("invalid onchain signing address %s", ocfg.OcrKeyBundle.OnchainSigningAddress)
	}
	signerB, err := hex.DecodeString(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signer %s: %w", signer, err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return &P2PSignerEnc{
		Signer:              sigb,
		P2PKey:              p2p,
		EncryptionPublicKey: pubkey, // TODO. no current way to get this from the node itself (and therefore not in clo or jd)
	}, nil
}
