package vault

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func Test_ValidateSignatures_Valid(t *testing.T) {
	ctx, err := hex.DecodeString("000ec4f6a2ba011e909eccf64628855b848e08876a1edd938a1372a9e51adff100000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	sig1, err := hex.DecodeString("d1067844e2849b404d903730c4cae19f090d53a578a1e8dc16ecbdc0285c1f186599108abbe0073b78bc148a6504907474ed3a6881df917e6d142cff70acfb5900")
	require.NoError(t, err)
	sig2, err := hex.DecodeString("c7517c188d297093a6f602046fad7feafe19454ee9dc269b19c8e6c01268037d1f7b423eeecbc495dd2d9a65e106bc3eab849ddfd74a10cbd4ad50c7d953bd4b01")
	require.NoError(t, err)

	payload := []byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`)
	resp := vaulttypes.SignedOCRResponse{
		Payload: payload,
		Context: ctx,
		Signatures: [][]byte{
			sig1,
			sig2,
		},
	}
	allowedAddr := []common.Address{
		common.HexToAddress("0xd6da96fe596705b32bc3a0e11cdefad77feaad79"),
		common.HexToAddress("0x327aa349c9718cd36c877d1e90458fe1929768ad"),
		common.HexToAddress("0xe9bf394856d73402b30e160d0e05c847796f0e29"),
		common.HexToAddress("0xefd5bdb6c3256f04489a6ca32654d547297f48b9"),
	}

	err = vaulttypes.ValidateSignatures(&resp, allowedAddr, 1)
	require.NoError(t, err)
}

func Test_ValidateSignatures_InsufficientSignatures(t *testing.T) {
	ctx, err := hex.DecodeString("000ec4f6a2ba011e909eccf64628855b848e08876a1edd938a1372a9e51adff100000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	sig1, err := hex.DecodeString("d1067844e2849b404d903730c4cae19f090d53a578a1e8dc16ecbdc0285c1f186599108abbe0073b78bc148a6504907474ed3a6881df917e6d142cff70acfb5900")
	require.NoError(t, err)
	payload := []byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`)
	resp := vaulttypes.SignedOCRResponse{
		Payload: payload,
		Context: ctx,
		Signatures: [][]byte{
			sig1,
		},
	}
	allowedAddr := []common.Address{
		common.HexToAddress("0xd6da96fe596705b32bc3a0e11cdefad77feaad79"),
		common.HexToAddress("0x327aa349c9718cd36c877d1e90458fe1929768ad"),
		common.HexToAddress("0xe9bf394856d73402b30e160d0e05c847796f0e29"),
		common.HexToAddress("0xefd5bdb6c3256f04489a6ca32654d547297f48b9"),
	}

	err = vaulttypes.ValidateSignatures(&resp, allowedAddr, 2)
	require.ErrorContains(t, err, "not enough signatures: expected min 2, got 1")
}

func Test_ValidateSignatures_DoesntCountDuplicates(t *testing.T) {
	ctx, err := hex.DecodeString("000ec4f6a2ba011e909eccf64628855b848e08876a1edd938a1372a9e51adff100000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	sig1, err := hex.DecodeString("d1067844e2849b404d903730c4cae19f090d53a578a1e8dc16ecbdc0285c1f186599108abbe0073b78bc148a6504907474ed3a6881df917e6d142cff70acfb5900")
	require.NoError(t, err)
	payload := []byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`)
	resp := vaulttypes.SignedOCRResponse{
		Payload: payload,
		Context: ctx,
		Signatures: [][]byte{
			sig1,
			sig1,
		},
	}
	allowedAddr := []common.Address{
		common.HexToAddress("0xd6da96fe596705b32bc3a0e11cdefad77feaad79"),
		common.HexToAddress("0x327aa349c9718cd36c877d1e90458fe1929768ad"),
		common.HexToAddress("0xe9bf394856d73402b30e160d0e05c847796f0e29"),
		common.HexToAddress("0xefd5bdb6c3256f04489a6ca32654d547297f48b9"),
	}

	err = vaulttypes.ValidateSignatures(&resp, allowedAddr, 2)
	require.ErrorContains(t, err, "only 1 valid signatures, need at least 2")
}

func Test_ValidateSignatures_InvalidSignature(t *testing.T) {
	ctx, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	sig1, err := hex.DecodeString("d1067844e2849b404d903730c4cae19f090d53a578a1e8dc16ecbdc0285c1f186599108abbe0073b78bc148a6504907474ed3a6881df917e6d142cff70acfb5900")
	require.NoError(t, err)
	sig2, err := hex.DecodeString("c7517c188d297093a6f602046fad7feafe19454ee9dc269b19c8e6c01268037d1f7b423eeecbc495dd2d9a65e106bc3eab849ddfd74a10cbd4ad50c7d953bd4b01")
	require.NoError(t, err)
	payload := []byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`)
	resp := vaulttypes.SignedOCRResponse{
		Payload: payload,
		Context: ctx,
		Signatures: [][]byte{
			sig1,
			sig2,
		},
	}
	allowedAddr := []common.Address{
		common.HexToAddress("0xd6da96fe596705b32bc3a0e11cdefad77feaad79"),
		common.HexToAddress("0x327aa349c9718cd36c877d1e90458fe1929768ad"),
		common.HexToAddress("0xe9bf394856d73402b30e160d0e05c847796f0e29"),
		common.HexToAddress("0xefd5bdb6c3256f04489a6ca32654d547297f48b9"),
	}

	err = vaulttypes.ValidateSignatures(&resp, allowedAddr, 2)
	require.ErrorContains(t, err, "only 0 valid signatures, need at least 2")
}
