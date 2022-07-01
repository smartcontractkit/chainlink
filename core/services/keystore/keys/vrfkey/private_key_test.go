package vrfkey

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestVRFKeys_PrivateKey(t *testing.T) {
	jsonKey := `{"PublicKey":"0xd2377bc6be8a2c5ce163e1867ee42ef111e320686f940a98e52e9c019ca0606800","vrf_key":{"address":"b94276ad4e5452732ec0cccf30ef7919b67844b6","crypto":{"cipher":"aes-128-ctr","ciphertext":"ff66d61d02dba54a61bab1ceb8414643f9e76b7351785d2959e2c8b50ee69a92","cipherparams":{"iv":"75705da271b11e330a27b8d593a3930c"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"efe5b372e4fe79d0af576a79d65a1ee35d0792d9c92b70107b5ada1817ea7c7b"},"mac":"e4d0bb08ffd004ab03aeaa42367acbd9bb814c6cfd981f5157503f54c30816e7"},"version":3}}`
	k, err := FromEncryptedJSON([]byte(jsonKey), "p4SsW0rD1!@#_")
	require.NoError(t, err)
	cryptoJSON, err := keystore.EncryptKey(k.toGethKey(), adulteratedPassword(testutils.Password), utils.FastScryptParams.N, utils.FastScryptParams.P)
	require.NoError(t, err)
	var gethKey gethKeyStruct
	err = json.Unmarshal(cryptoJSON, &gethKey)
	require.NoError(t, err)

	ek := EncryptedVRFKey{
		PublicKey: k.PublicKey,
		VRFKey:    gethKey,
	}

	pk, err := Decrypt(ek, testutils.Password)
	require.NoError(t, err)
	_, err = Decrypt(ek, "wrong-password")
	assert.Error(t, err)

	kv2 := pk.ToV2()

	assert.Equal(t, fmt.Sprintf("VRFKeyV2{PublicKey: %s}", kv2.PublicKey), kv2.String())
	assert.Equal(t, fmt.Sprintf("PrivateKey{k: <redacted>, PublicKey: %s}", pk.PublicKey), pk.String())
}
