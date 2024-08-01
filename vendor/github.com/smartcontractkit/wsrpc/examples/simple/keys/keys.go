package keys

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"

	"github.com/smartcontractkit/wsrpc/credentials"
)

type Client struct {
	ID                 string
	Name               string
	PubKey             string
	PrivKey            string
	RegisteredOnServer bool
}

// Server Client keys in hex

const ServerPubKey = "3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808"
const ServerPrivKey = "c1afd224cec2ff6066746bf9b7cdf7f9f4694ab7ef2ca1692ff923a30df203483b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808"

var Clients = []Client{
	{
		ID:                 "1",
		Name:               "Alice",
		PubKey:             "0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9",
		PrivKey:            "5ae73174dfe4ae293e25ba6845e76d9819e5f5432ff8820d5c996be292abf14f0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9",
		RegisteredOnServer: true,
	},
	{
		ID:                 "2",
		Name:               "Bob",
		PubKey:             "9a36f1819c60970b0cb16585cacf35ed824b48bc8ac24980bf615bc8d7b9661c",
		PrivKey:            "1b6093915ce64fa5ca15147808da47ccb64f50a16841edde12fd797c3e52d8169a36f1819c60970b0cb16585cacf35ed824b48bc8ac24980bf615bc8d7b9661c",
		RegisteredOnServer: true,
	},
	// This user is not registered on the server
	{
		ID:                 "3",
		Name:               "Charlie",
		PubKey:             "235750320bb723760add5969b3c51342e829542d399f82b07a0458092c3960af",
		PrivKey:            "84206ddc52ad6f83569ed409829f194db1d3bbb65e7c04db5ca098d1f020ee47235750320bb723760add5969b3c51342e829542d399f82b07a0458092c3960af",
		RegisteredOnServer: false,
	},
}

func FromHex(keyHex string) []byte {
	privKey := make([]byte, hex.DecodedLen(len(keyHex)))
	hex.Decode(privKey, []byte(keyHex))

	return privKey
}

// ToStaticSizedBytes convert bytes to a statically sized byte array of the
// of ed25519.PublicKeySize
func ToStaticSizedBytes(b []byte) (credentials.StaticSizedPublicKey, error) {
	var sb credentials.StaticSizedPublicKey

	if ed25519.PublicKeySize != copy(sb[:], b) {
		return sb, errors.New("copying public key failed")
	}

	return sb, nil
}
