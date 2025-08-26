//nolint:nolintlint
package tdh2shim

import (
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault/sanmarinodkg/dkg"
)

// Shim for extracting the public TDH2 public key from a DKG result.
// Note: in this demo code, a direct Unmarshal of the value returned by result.MasterPublicKey() is sufficient,
// but the consumer of this code must not rely on this behavior, and instead use the provided shim.
func TDH2PublicKeyFromDKGResult(result dkg.Result) (*tdh2easy.PublicKey, error) {
	mpkBytes, err := result.MasterPublicKey()
	if err != nil {
		return nil, err
	}
	pk := new(tdh2easy.PublicKey)
	err = pk.Unmarshal(mpkBytes)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

// Shim for extracting the private TDH2 share from a DKG result. Requires the secret key of the participant for
// decrypting the master secret key share.
// Note: in this demo code, a direct Unmarshal of the value returned by result.MasterSecretKeyShare(...) is sufficient,
// but the consumer of this code must not rely on this behavior, and instead use the provided shim.
func TDH2PrivateShareFromDKGResult(result dkg.Result, sk dkg.ParticipantSecretKey) (*tdh2easy.PrivateShare, error) {
	mskShare, err := result.MasterSecretKeyShare(sk)
	if err != nil {
		return nil, err
	}
	ps := new(tdh2easy.PrivateShare)
	if err := ps.Unmarshal(mskShare); err != nil {
		return nil, err
	}
	return ps, nil
}
