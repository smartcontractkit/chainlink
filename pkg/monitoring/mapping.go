package monitoring

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/pb"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/protobuf/proto"
)

func MakeTransmissionMapping(
	envelope Envelope,
	chainConfig ChainConfig,
	feedConfig FeedConfig,
) (map[string]interface{}, error) {
	data := []byte{}
	if envelope.LatestAnswer != nil {
		data = envelope.LatestAnswer.Bytes()
	}
	out := map[string]interface{}{
		"block_number": uint64ToBeBytes(envelope.BlockNumber),
		"answer": map[string]interface{}{
			"data":      data,
			"timestamp": envelope.LatestTimestamp.Unix(),
			"config_digest": map[string]interface{}{
				"string": base64.StdEncoding.EncodeToString(envelope.ConfigDigest[:]),
			},
			"epoch": map[string]interface{}{
				"long": int64(envelope.Epoch),
			},
			"round": map[string]interface{}{
				"int": int32(envelope.Round),
			},
		},
		"chain_config": map[string]interface{}{
			"link.chain.ocr2.chain_config": chainConfig.ToMapping(),
		},
		// Deprecated in favour of chain_config.
		"solana_chain_config": map[string]interface{}{
			"network_name": "",
			"network_id":   "",
			"chain_id":     "",
		},
		"feed_config": feedConfig.ToMapping(),
		"link_balance": map[string]interface{}{
			"bytes": uint64ToBeBytes(envelope.LinkBalance),
		},
	}
	return out, nil
}

func MakeConfigSetSimplifiedMapping(
	envelope Envelope,
	feedConfig FeedConfig,
) (map[string]interface{}, error) {
	offchainConfig, err := parseOffchainConfig(envelope.ContractConfig.OffchainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OffchainConfig blob from the program state: %w", err)
	}
	signers, err := json.Marshal(envelope.ContractConfig.Signers)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signers: %w", err)
	}
	transmitters, err := json.Marshal(envelope.ContractConfig.Transmitters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transmitters: %w", err)
	}
	s, err := json.Marshal(int32ArrToInt64Arr(offchainConfig.S))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schedule: %w", err)
	}
	oracles, err := createConfigSetSimplifiedOracles(offchainConfig.OffchainPublicKeys, offchainConfig.PeerIds, envelope.ContractConfig.Transmitters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode oracle set: %w", err)
	}
	out := map[string]interface{}{
		"config_digest":      base64.StdEncoding.EncodeToString(envelope.ConfigDigest[:]),
		"block_number":       uint64ToBeBytes(envelope.BlockNumber),
		"signers":            string(signers),
		"transmitters":       string(transmitters),
		"f":                  int32(envelope.ContractConfig.F),
		"delta_progress":     uint64ToBeBytes(offchainConfig.DeltaProgressNanoseconds),
		"delta_resend":       uint64ToBeBytes(offchainConfig.DeltaResendNanoseconds),
		"delta_round":        uint64ToBeBytes(offchainConfig.DeltaRoundNanoseconds),
		"delta_grace":        uint64ToBeBytes(offchainConfig.DeltaGraceNanoseconds),
		"delta_stage":        uint64ToBeBytes(offchainConfig.DeltaStageNanoseconds),
		"r_max":              int64(offchainConfig.RMax),
		"s":                  string(s),
		"oracles":            string(oracles),
		"feed_state_account": feedConfig.GetContractAddress(),
	}
	return out, nil
}

// Helpers

func uint64ToBeBytes(input uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, input)
	return buf
}

func parseOffchainConfig(buf []byte) (*pb.OffchainConfigProto, error) {
	config := &pb.OffchainConfigProto{}
	err := proto.Unmarshal(buf, config)
	return config, err
}

func int32ArrToInt64Arr(xs []uint32) []int64 {
	out := make([]int64, len(xs))
	for i, x := range xs {
		out[i] = int64(x)
	}
	return out
}

func createConfigSetSimplifiedOracles(offchainPublicKeys [][]byte, peerIDs []string, transmitters []types.Account) ([]byte, error) {
	if len(offchainPublicKeys) != len(peerIDs) && len(transmitters) != len(peerIDs) {
		return nil, fmt.Errorf("length missmatch len(offchainPublicKeys)=%d , len(transmitters)=%d, len(peerIDs)=%d", len(offchainPublicKeys), len(transmitters), len(peerIDs))
	}
	out := make([]interface{}, len(transmitters))
	for i := 0; i < len(transmitters); i++ {
		out[i] = map[string]interface{}{
			"transmitter":         transmitters[i],
			"peer_id":             peerIDs[i],
			"offchain_public_key": offchainPublicKeys[i],
		}
	}
	s, err := json.Marshal(out)
	return s, err
}
