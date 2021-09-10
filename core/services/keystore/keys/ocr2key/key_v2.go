package ocr2key

import "encoding/hex"

type Raw []byte

func (raw Raw) Key() KeyBundle {
	return KeyBundle{}
}

func (raw Raw) String() string {
	return "<OCR Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

func NewV2() (KeyBundle, error) {
	return KeyBundle{}, nil
}

func (key KeyBundle) ID() string {
	return hex.EncodeToString(key.id[:])
}

func (key KeyBundle) Raw() Raw {
	return []byte{}
}

// SignOnChain returns an ethereum-style ECDSA secp256k1 signature on msg.
func (key KeyBundle) SignOnChain(msg []byte) ([]byte, error) {
	return []byte{}, nil
}

// SignOffChain returns an EdDSA-Ed25519 signature on msg.
func (key KeyBundle) SignOffChain(msg []byte) ([]byte, error) {
	return []byte{}, nil
}

// ConfigDiffieHellman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the OffChainEncryption key.)
// func (key KeyBundle) ConfigDiffieHellman(base *[curve25519.PointSize]byte) (
// 	*[curve25519.PointSize]byte, error,
// ) {
// 	return nil, nil
// }

// PublicKeyAddressOnChain returns public component of the keypair used in
// SignOnChain
// func (key KeyBundle) PublicKeyAddressOnChain() ocrtypes.OnChainSigningAddress {
// 	return ocrtypes.OnChainSigningAddress{}
// }

// PublicKeyOffChain returns the pbulic component of the keypair used in SignOffChain
// func (key KeyBundle) PublicKeyOffChain() ocrtypes.OffchainPublicKey {
// 	return ocrtypes.OffChainPublicKey(ed25519.PublicKey{})
// }

// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
// func (key KeyBundle) PublicKeyConfig() [curve25519.PointSize]byte {
// 	return [curve25519.PointSize]byte{}
// }

// func (key KeyBundle) GetID() string {
// 	return key.ID
// }

// func (key KeyBundle) String() string {
// 	return fmt.Sprintf("OCRKeyBundle{ID: %s}", key.ID())
// }

func (key KeyBundle) GoString() string {
	return key.String()
}
