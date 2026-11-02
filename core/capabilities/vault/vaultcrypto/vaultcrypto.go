// Package vaultcrypto holds the TDH2 share helpers used by the Vault DON and its consumers.
package vaultcrypto

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

// ErrInvalidRecipientKey matches EncryptShareBinary's errors for a malformed recipient key.
var ErrInvalidRecipientKey = errors.New("invalid recipient public key")

type invalidRecipientKeyError struct{ msg string }

func (e invalidRecipientKeyError) Error() string { return e.msg }

func (e invalidRecipientKeyError) Is(target error) bool {
	_, ok := target.(invalidRecipientKeyError)
	return ok || target == ErrInvalidRecipientKey
}

// GeneratePlaintextShare returns this node's decryption share after checking the ciphertext's label.
func GeneratePlaintextShare(publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare, encryptedSecret []byte, expectedLabel [32]byte) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("vault public key is not set")
	}
	if privateKeyShare == nil {
		return nil, errors.New("vault private key share is not set")
	}

	ct := &tdh2easy.Ciphertext{}
	if err := ct.UnmarshalVerify(encryptedSecret, publicKey); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ciphertext: %w", err)
	}

	if label := ct.Label(); label != expectedLabel {
		return nil, fmt.Errorf("failed to verify label on secret. error: secret label [%s] does not match workflow owner label [%s]",
			hex.EncodeToString(label[:]), hex.EncodeToString(expectedLabel[:]))
	}

	s, err := tdh2easy.Decrypt(ct, privateKeyShare)
	if err != nil {
		return nil, fmt.Errorf("could not generate decryption share: %w", err)
	}

	sb, err := s.Marshal()
	if err != nil {
		return nil, errors.New("could not marshal decryption share")
	}
	return sb, nil
}

// EncryptShareBinary seals the share for a hex encoded curve25519 public key.
func EncryptShareBinary(share []byte, recipientPublicKeyHex string) ([]byte, error) {
	publicKey, err := hex.DecodeString(recipientPublicKeyHex)
	if err != nil {
		return nil, invalidRecipientKeyError{"failed to convert public key to bytes: " + err.Error()}
	}

	if len(publicKey) != curve25519.PointSize {
		return nil, invalidRecipientKeyError{fmt.Sprintf("invalid public key size: expected %d bytes, got %d bytes", curve25519.PointSize, len(publicKey))}
	}

	pk := [curve25519.PointSize]byte(publicKey)
	encrypted, err := box.SealAnonymous(nil, share, &pk, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt decryption share: %w", err)
	}
	return encrypted, nil
}

// SharesForKey collects the shares for a recipient key from every matching entry,
// preferring an entry's binary shares over its legacy hex shares.
func SharesForKey(entries []*vaultcommon.EncryptedShares, recipientPublicKeyHex string) ([][]byte, error) {
	var out [][]byte
	for _, e := range entries {
		if e.GetEncryptionKey() != recipientPublicKeyHex {
			continue
		}
		if len(e.GetBinaryShares()) > 0 {
			out = append(out, e.GetBinaryShares()...)
			continue
		}
		for _, s := range e.GetShares() {
			b, err := hex.DecodeString(s)
			if err != nil {
				return nil, fmt.Errorf("failed to decode hex share: %w", err)
			}
			out = append(out, b)
		}
	}
	return out, nil
}

type decrypter interface {
	Decrypt(encrypted []byte) ([]byte, error)
}

// VerifyShares returns the valid shares, and an error for each rejected one.
// At most one share per node index is returned
func VerifyShares(ct *tdh2easy.Ciphertext, publicKey *tdh2easy.PublicKey, encryptedShares [][]byte, d decrypter) ([]*tdh2easy.DecryptionShare, []error) {
	valid := make([]*tdh2easy.DecryptionShare, 0, len(encryptedShares))
	seenIndexes := make(map[int]struct{}, len(encryptedShares))
	var rejected []error
	for i, es := range encryptedShares {
		b, err := d.Decrypt(es)
		if err != nil {
			rejected = append(rejected, fmt.Errorf("share %d: failed to decrypt: %w", i, err))
			continue
		}
		share := &tdh2easy.DecryptionShare{}
		if err := share.Unmarshal(b); err != nil {
			rejected = append(rejected, fmt.Errorf("share %d: failed to unmarshal: %w", i, err))
			continue
		}
		if err := tdh2easy.VerifyShare(ct, publicKey, share); err != nil {
			rejected = append(rejected, fmt.Errorf("share %d: failed to verify: %w", i, err))
			continue
		}
		if _, dup := seenIndexes[share.Index()]; dup {
			rejected = append(rejected, fmt.Errorf("share %d: duplicate share for node index %d", i, share.Index()))
			continue
		}
		seenIndexes[share.Index()] = struct{}{}
		valid = append(valid, share)
	}
	return valid, rejected
}

// DecryptSecret also returns the rejected shares' errors, for logging.
func DecryptSecret(encryptedSecret []byte, publicKey *tdh2easy.PublicKey, threshold int, encryptedShares [][]byte, d decrypter) ([]byte, []error, error) {
	ct := &tdh2easy.Ciphertext{}
	if err := ct.UnmarshalVerify(encryptedSecret, publicKey); err != nil {
		return nil, nil, errors.New("failed to unmarshal encrypted secret: " + err.Error())
	}

	valid, rejected := VerifyShares(ct, publicKey, encryptedShares, d)
	if len(valid) < threshold {
		return nil, rejected, fmt.Errorf("not enough decryption shares to decrypt the secret: have %d, need at least %d", len(valid), threshold)
	}

	// Aggregate's n is only an allocation hint.
	plaintext, err := tdh2easy.Aggregate(ct, valid, len(encryptedShares))
	if err != nil {
		return nil, rejected, errors.New("failed to aggregate decryption shares: " + err.Error())
	}
	return plaintext, rejected, nil
}
