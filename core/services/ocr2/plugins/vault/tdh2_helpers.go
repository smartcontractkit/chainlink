package vault

import (
	"fmt"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
)

func tdh2ToTDH2EasyPK(pk *tdh2.PublicKey) (*tdh2easy.PublicKey, error) {
	tdh2PubKeyBytes, err := pk.Marshal()
	if err != nil {
		return nil, fmt.Errorf("could not marshal tdh2 public key: %w", err)
	}
	publicKey := &tdh2easy.PublicKey{}
	err = publicKey.Unmarshal(tdh2PubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal to tdh2easy public key: %w", err)
	}
	return publicKey, nil
}

func tdh2ToTDH2EasyKS(pk *tdh2.PrivateShare) (*tdh2easy.PrivateShare, error) {
	tdh2PrivKeyShareBytes, err := pk.Marshal()
	if err != nil {
		return nil, fmt.Errorf("could not marshal tdh2 private key share: %w", err)
	}
	privateKeyShare := &tdh2easy.PrivateShare{}
	err = privateKeyShare.Unmarshal(tdh2PrivKeyShareBytes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal to tdh2easy private key share: %w", err)
	}
	return privateKeyShare, nil
}
