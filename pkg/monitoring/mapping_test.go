package monitoring

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapping(t *testing.T) {
	config, _, offchainConfig, _, err := generateContractConfig(31)
	require.NoError(t, err)
	envelope, err := generateEnvelope()
	require.NoError(t, err)
	envelope.ConfigDigest = config.ConfigDigest
	envelope.ContractConfig = config

	chainConfig := generateChainConfig()
	feedConfig := generateFeedConfig()

	t.Run("MakeTransmissionMapping", func(t *testing.T) {
		mapping, err := MakeTransmissionMapping(envelope, chainConfig, feedConfig)
		require.NoError(t, err)
		output := []byte{}
		serialized, err := transmissionCodec.BinaryFromNative(output, mapping)
		require.NoError(t, err)
		deserialized, _, err := transmissionCodec.NativeFromBinary(serialized)
		require.NoError(t, err)

		transmission, ok := deserialized.(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, transmission["block_number"], uint64ToBeBytes(envelope.BlockNumber))

		answer, ok := transmission["answer"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, answer["data"], envelope.LatestAnswer.Bytes())
		require.Equal(t, answer["timestamp"].(int64), envelope.LatestTimestamp.Unix())

		configDigest, ok := answer["config_digest"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, configDigest["string"].(string), base64.StdEncoding.EncodeToString(envelope.ConfigDigest[:]))

		epoch, ok := answer["epoch"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, epoch["long"].(int64), int64(envelope.Epoch))

		round, ok := answer["round"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, round["int"].(int32), int32(envelope.Round))

		// Deprecated in favour of chain_config
		chainConfigUnion, ok := transmission["chain_config"].(map[string]interface{})
		require.True(t, ok)
		decodedChainConfig, ok := chainConfigUnion["link.chain.ocr2.chain_config"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, decodedChainConfig["network_name"], chainConfig.GetNetworkName())
		require.Equal(t, decodedChainConfig["network_id"], chainConfig.GetNetworkID())
		require.Equal(t, decodedChainConfig["chain_id"], chainConfig.GetChainID())

		decodedFeedConfig, ok := transmission["feed_config"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, decodedFeedConfig, feedConfig.ToMapping())
	})

	t.Run("MakeSimplifiedConfigSetMapping", func(t *testing.T) {
		mapping, err := MakeConfigSetSimplifiedMapping(envelope, feedConfig)
		require.NoError(t, err)

		var output []byte
		serialized, err := configSetSimplifiedCodec.BinaryFromNative(output, mapping)
		require.NoError(t, err)
		deserialized, _, err := configSetSimplifiedCodec.NativeFromBinary(serialized)
		require.NoError(t, err)

		configSetSimplified, ok := deserialized.(map[string]interface{})
		require.True(t, ok)

		oracles, err := createConfigSetSimplifiedOracles(offchainConfig.OffchainPublicKeys, offchainConfig.PeerIds, config.Transmitters)
		require.NoError(t, err)

		require.Equal(t, configSetSimplified["config_digest"], base64.StdEncoding.EncodeToString(envelope.ConfigDigest[:]))
		require.Equal(t, configSetSimplified["block_number"], uint64ToBeBytes(envelope.BlockNumber))
		require.Equal(t, configSetSimplified["delta_progress"], uint64ToBeBytes(offchainConfig.DeltaProgressNanoseconds))
		require.Equal(t, configSetSimplified["delta_resend"], uint64ToBeBytes(offchainConfig.DeltaResendNanoseconds))
		require.Equal(t, configSetSimplified["delta_round"], uint64ToBeBytes(offchainConfig.DeltaRoundNanoseconds))
		require.Equal(t, configSetSimplified["delta_grace"], uint64ToBeBytes(offchainConfig.DeltaGraceNanoseconds))
		require.Equal(t, configSetSimplified["delta_stage"], uint64ToBeBytes(offchainConfig.DeltaStageNanoseconds))
		require.Equal(t, configSetSimplified["r_max"], int64(offchainConfig.RMax))
		require.Equal(t, configSetSimplified["f"], int32(config.F))
		require.Equal(t, configSetSimplified["signers"], jsonMarshalToString(t, config.Signers))
		require.Equal(t, configSetSimplified["transmitters"], jsonMarshalToString(t, config.Transmitters))
		require.Equal(t, configSetSimplified["s"], jsonMarshalToString(t, offchainConfig.S))
		require.Equal(t, configSetSimplified["oracles"], string(oracles))
		require.Equal(t, configSetSimplified["feed_state_account"], feedConfig.GetContractAddress())
	})

	t.Run("MakeSimplifiedConfigSetMapping works for an empty envelope", func(t *testing.T) {
		mapping, err := MakeConfigSetSimplifiedMapping(envelope, feedConfig)
		require.NoError(t, err)
		_, err = configSetSimplifiedCodec.BinaryFromNative(nil, mapping)
		require.NoError(t, err)
	})

	t.Run("MakeTransmissionMapping works for empty envelope", func(t *testing.T) {
		mapping, err := MakeTransmissionMapping(envelope, chainConfig, feedConfig)
		require.NoError(t, err)
		_, err = transmissionCodec.BinaryFromNative(nil, mapping)
		require.NoError(t, err)
	})
}

// Helpers

func jsonMarshalToString(t *testing.T, i interface{}) string {
	s, err := json.Marshal(i)
	require.NoError(t, err)
	return string(s)
}

func interfaceArrToUint32Arr(in []interface{}) []int64 {
	out := []int64{}
	for _, i := range in {
		out = append(out, i.(int64))
	}
	return out
}

func interfaceArrToBytesArr(in []interface{}) [][]byte {
	out := [][]byte{}
	for _, i := range in {
		out = append(out, i.([]byte))
	}
	return out
}

func interfaceArrToStringArr(in []interface{}) []string {
	out := []string{}
	for _, i := range in {
		out = append(out, i.(string))
	}
	return out
}
