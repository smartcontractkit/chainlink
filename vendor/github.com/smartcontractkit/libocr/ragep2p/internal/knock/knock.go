package knock

import (
	"crypto/ed25519"
	"fmt"

	"github.com/smartcontractkit/libocr/ragep2p/types"
)

const domainSeparator = "ragep2p 1.0.0 knock knock"
const version = byte(0x02) // 0x2 rather than 0x1 to prevent accidental connection with previous OCR networking

// knock = version (1 byte) || pk (ed25519.PublicKeySize) || sig (ed25519.SignatureSize)
const KnockSize = 1 + ed25519.PublicKeySize + ed25519.SignatureSize

func messageToSign(peerID types.PeerID) []byte {
	msg := make([]byte, 0, len(domainSeparator)+len(peerID))
	msg = append(msg, []byte(domainSeparator)...)
	msg = append(msg, peerID[:]...)
	return msg
}

// Builds a knock message based on the PeerID of the node being dialed (other),
// the dialing node (self). The dialing node's secretKey is used for signing the
// message and must correspond to self.
//
// Returns a knock message exactly KnockSize bytes long.
func BuildKnock(other types.PeerID, self types.PeerID, secretKey ed25519.PrivateKey) []byte {
	msg := messageToSign(other)
	sig := ed25519.Sign(secretKey, msg)

	knock := make([]byte, 0, KnockSize)
	knock = append(knock, version)
	knock = append(knock, self[:]...)
	knock = append(knock, sig...)
	return knock
}

var (
	ErrSizeMismatch     = fmt.Errorf("knock size mismatch")
	ErrFromSelfDial     = fmt.Errorf("knock from self-dial")
	ErrInvalidSignature = fmt.Errorf("knock has invalid signature")
)

// Verifies a knock message allegedly destined to self. If the message is valid,
// returns the PeerId of the sender. Otherwise returns nil and an error.
func VerifyKnock(self types.PeerID, knock []byte) (*types.PeerID, error) {
	if len(knock) != KnockSize {
		return nil, fmt.Errorf("knock hash wrong length %v, expected %v", len(knock), KnockSize)
	}

	if knock[0] != version {
		return nil, fmt.Errorf("knock has wrong version %v, expected %v", knock[0], version)
	}
	knock = knock[1:]

	var other types.PeerID
	if len(other) != copy(other[:], knock[:ed25519.PublicKeySize]) {
		return nil, ErrSizeMismatch
	}

	if other == self {
		return nil, ErrFromSelfDial
	}

	msg := messageToSign(self)
	sig := knock[ed25519.PublicKeySize:]
	if !ed25519.Verify(ed25519.PublicKey(other[:]), msg, sig) {
		return nil, ErrInvalidSignature
	}

	return &other, nil
}
