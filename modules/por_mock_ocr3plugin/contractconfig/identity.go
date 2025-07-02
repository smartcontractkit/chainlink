package contractconfig

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/myname"
	"golang.org/x/crypto/curve25519"
)

func P2pPrivateKey(i int) ed25519.PrivateKey {
	return ed25519.NewKeyFromSeed([]byte(fmt.Sprintf("MontrealMontrealMontreal%8d", i)))
}

func OffchainPrivateKey(i int) ed25519.PrivateKey {
	return ed25519.NewKeyFromSeed([]byte(fmt.Sprintf("CanadaCanadaCanadaCanada%8d", i)))
}

func ConfigEncryptionPrivateKey(i int) [curve25519.ScalarSize]byte {
	var priv [curve25519.ScalarSize]byte
	copy(priv[:], []byte(fmt.Sprintf("Bonjour!Bonjour!Bonjour!%8d", i)))
	return priv
}

func OnchainPrivateKey(i int) ecdsa.PrivateKey {
	secret := new(big.Int)
	secret.SetBytes([]byte(fmt.Sprintf("AwesomAwesomAwesomAwesom%8d", i)))

	x, y := secp256k1.S256().ScalarBaseMult(secret.Bytes())
	return ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     x, Y: y,
		},
		D: secret,
	}
}

func TransmitterPrivateKey(i int) ecdsa.PrivateKey {
	secret := new(big.Int)
	secret.SetBytes(crypto.Keccak256([]byte(fmt.Sprintf("Poutine!Poutine!Poutine!%s%8d", myname.Name, i))))

	x, y := secp256k1.S256().ScalarBaseMult(secret.Bytes())
	return ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     x, Y: y,
		},
		D: secret,
	}
}

func TransmitterAddress(i int) common.Address {
	return crypto.PubkeyToAddress(TransmitterPrivateKey(i).PublicKey)
}

func GodPrivateKey() ecdsa.PrivateKey {
	secret := new(big.Int)
	secret.SetBytes(crypto.Keccak256([]byte("lakeshfk hadksjfhk hkjhsabdfkh bakshjdbf kahbdskf bo73yo47y23")))

	x, y := secp256k1.S256().ScalarBaseMult(secret.Bytes())
	return ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     x, Y: y,
		},
		D: secret,
	}
}

func GodAddress() common.Address {
	return crypto.PubkeyToAddress(GodPrivateKey().PublicKey)
}

func DestinationPrivateKey() ecdsa.PrivateKey {
	secret := new(big.Int)
	secret.SetBytes(crypto.Keccak256([]byte(fmt.Sprintf("destination address for %s", myname.Name))))

	x, y := secp256k1.S256().ScalarBaseMult(secret.Bytes())
	return ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     x, Y: y,
		},
		D: secret,
	}
}

func DestinationAddress() common.Address {
	return crypto.PubkeyToAddress(DestinationPrivateKey().PublicKey)
}

func offchainPublicKeyKeyFromPrivateKey(priv ed25519.PrivateKey) types.OffchainPublicKey {
	var result types.OffchainPublicKey
	copy(result[:], priv.Public().(ed25519.PublicKey))
	return result
}

func peerIDFromPrivateKey(priv ed25519.PrivateKey) string {
	peerID, err := ragetypes.PeerIDFromPrivateKey(priv)
	if err != nil {
		panic(err)
	}
	return peerID.String()
}

func accountFromPrivateKey(priv ecdsa.PrivateKey) types.Account {
	return types.Account(crypto.PubkeyToAddress(priv.PublicKey).Hex())
}

func OracleIdentity(i int) confighelper.OracleIdentityExtra {
	var configEncryptionPublicKey types.ConfigEncryptionPublicKey
	{
		scalar := ConfigEncryptionPrivateKey(i)
		curve25519.ScalarBaseMult((*[32]byte)(&configEncryptionPublicKey), &scalar)
	}

	return confighelper.OracleIdentityExtra{
		OracleIdentity: confighelper.OracleIdentity{
			OffchainPublicKey: offchainPublicKeyKeyFromPrivateKey(OffchainPrivateKey(i)),
			OnchainPublicKey:  crypto.PubkeyToAddress(OnchainPrivateKey(i).PublicKey).Bytes(),
			PeerID:            peerIDFromPrivateKey(P2pPrivateKey(i)),
			TransmitAccount:   accountFromPrivateKey(TransmitterPrivateKey(i)),
		},
		ConfigEncryptionPublicKey: configEncryptionPublicKey,
	}
}

func OracleIdentities(n int) []confighelper.OracleIdentityExtra {
	var result []confighelper.OracleIdentityExtra
	for i := 0; i < n; i++ {
		result = append(result, OracleIdentity(i))
	}
	return result
}
