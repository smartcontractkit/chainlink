package ocr2config

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/crypto/curve25519"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// Serialized configs must be no larger than this (arbitrary bound, to prevent
// resource exhaustion attacks)
var maxSerializedOffchainConfigSize = 2_000_000

// offchainConfig contains the contents of the oracle Config objects
// which need to be serialized
type offchainConfig struct {
	DeltaProgress                           time.Duration
	DeltaResend                             time.Duration
	DeltaRound                              time.Duration
	DeltaGrace                              time.Duration
	DeltaStage                              time.Duration
	RMax                                    uint8
	S                                       []int
	OffchainPublicKeys                      []types.OffchainPublicKey
	PeerIDs                                 []string
	ReportingPluginConfig                   []byte
	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationReport                       time.Duration
	MaxDurationShouldAcceptFinalizedReport  time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration
	SharedSecretEncryptions                 config.SharedSecretEncryptions
}

func checkSize(serializedOffchainConfig []byte) error {
	if len(serializedOffchainConfig) <= maxSerializedOffchainConfigSize {
		return nil
	} else {
		return fmt.Errorf("OffchainConfig length is %d bytes which is greater than the max %d",
			len(serializedOffchainConfig),
			maxSerializedOffchainConfigSize,
		)
	}
}

// serialize returns a binary serialization of o
func (o offchainConfig) serialize() []byte {
	offchainConfigProto := enprotoOffchainConfig(o)
	rv, err := proto.Marshal(&offchainConfigProto)
	if err != nil {
		panic(err)
	}
	if err := checkSize(rv); err != nil {
		panic(err.Error())
	}
	return rv
}

func deserializeOffchainConfig(
	b []byte,
) (offchainConfig, error) {
	if err := checkSize(b); err != nil {
		return offchainConfig{}, err
	}

	offchainConfigPB := OffchainConfigProto{}
	if err := proto.Unmarshal(b, &offchainConfigPB); err != nil {
		return offchainConfig{}, fmt.Errorf("could not unmarshal ContractConfig.OffchainConfig protobuf: %w", err)
	}

	return deprotoOffchainConfig(&offchainConfigPB)
}

func deprotoOffchainConfig(
	offchainConfigProto *OffchainConfigProto,
) (offchainConfig, error) {
	S := make([]int, 0, len(offchainConfigProto.GetS()))
	for _, elem := range offchainConfigProto.GetS() {
		S = append(S, int(elem))
	}

	offchainPublicKeys := make([]types.OffchainPublicKey, 0, len(offchainConfigProto.GetOffchainPublicKeys()))
	for _, ocpkRaw := range offchainConfigProto.GetOffchainPublicKeys() {
		var ocpk types.OffchainPublicKey
		if len(ocpkRaw) != len(ocpk) {
			return offchainConfig{}, fmt.Errorf("invalid offchain public key: %x", ocpkRaw)
		}
		copy(ocpk[:], ocpkRaw)
		offchainPublicKeys = append(offchainPublicKeys, ocpk)
	}

	sharedSecretEncryptions, err := deprotoSharedSecretEncryptions(offchainConfigProto.GetSharedSecretEncryptions())
	if err != nil {
		return offchainConfig{}, fmt.Errorf("could not unmarshal shared protobuf: %w", err)
	}

	return offchainConfig{
		time.Duration(offchainConfigProto.GetDeltaProgressNanoseconds()),
		time.Duration(offchainConfigProto.GetDeltaResendNanoseconds()),
		time.Duration(offchainConfigProto.GetDeltaRoundNanoseconds()),
		time.Duration(offchainConfigProto.GetDeltaGraceNanoseconds()),
		time.Duration(offchainConfigProto.GetDeltaStageNanoseconds()),
		uint8(offchainConfigProto.GetRMax()),
		S,
		offchainPublicKeys,
		offchainConfigProto.GetPeerIds(),
		offchainConfigProto.GetReportingPluginConfig(),
		time.Duration(offchainConfigProto.GetMaxDurationQueryNanoseconds()),
		time.Duration(offchainConfigProto.GetMaxDurationObservationNanoseconds()),
		time.Duration(offchainConfigProto.GetMaxDurationReportNanoseconds()),
		time.Duration(offchainConfigProto.GetMaxDurationShouldAcceptFinalizedReportNanoseconds()),
		time.Duration(offchainConfigProto.GetMaxDurationShouldTransmitAcceptedReportNanoseconds()),
		sharedSecretEncryptions,
	}, nil
}

func deprotoSharedSecretEncryptions(sharedSecretEncryptionsProto *SharedSecretEncryptionsProto) (config.SharedSecretEncryptions, error) {
	var diffieHellmanPoint [curve25519.PointSize]byte
	if len(diffieHellmanPoint) != len(sharedSecretEncryptionsProto.GetDiffieHellmanPoint()) {
		return config.SharedSecretEncryptions{}, fmt.Errorf("DiffieHellmanPoint has wrong length. Expected %v bytes, got %v bytes", len(diffieHellmanPoint), len(sharedSecretEncryptionsProto.GetDiffieHellmanPoint()))
	}
	copy(diffieHellmanPoint[:], sharedSecretEncryptionsProto.GetDiffieHellmanPoint())

	var sharedSecretHash common.Hash
	if len(sharedSecretHash) != len(sharedSecretEncryptionsProto.GetSharedSecretHash()) {
		return config.SharedSecretEncryptions{}, fmt.Errorf("sharedSecretHash has wrong length. Expected %v bytes, got %v bytes", len(sharedSecretHash), len(sharedSecretEncryptionsProto.GetSharedSecretHash()))
	}
	copy(sharedSecretHash[:], sharedSecretEncryptionsProto.GetSharedSecretHash())

	encryptions := make([]config.EncryptedSharedSecret, 0, len(sharedSecretEncryptionsProto.GetEncryptions()))
	for i, encryptionRaw := range sharedSecretEncryptionsProto.GetEncryptions() {
		var encryption config.EncryptedSharedSecret
		if len(encryption) != len(encryptionRaw) {
			return config.SharedSecretEncryptions{}, fmt.Errorf("Encryptions[%v] has wrong length. Expected %v bytes, got %v bytes", i, len(encryption), len(encryptionRaw))
		}
		copy(encryption[:], encryptionRaw)
		encryptions = append(encryptions, encryption)
	}

	return config.SharedSecretEncryptions{
		diffieHellmanPoint,
		sharedSecretHash,
		encryptions,
	}, nil
}

func enprotoOffchainConfig(o offchainConfig) OffchainConfigProto {
	s := make([]uint32, len(o.S))
	for i, d := range o.S {
		s[i] = uint32(d)
	}
	offchainPublicKeys := make([][]byte, 0, len(o.OffchainPublicKeys))
	for _, k := range o.OffchainPublicKeys {
		k := k // have to copy or we append the same key over and over
		offchainPublicKeys = append(offchainPublicKeys, k[:])
	}
	sharedSecretEncryptions := enprotoSharedSecretEncryptions(o.SharedSecretEncryptions)
	return OffchainConfigProto{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		uint64(o.DeltaProgress),
		uint64(o.DeltaResend),
		uint64(o.DeltaRound),
		uint64(o.DeltaGrace),
		uint64(o.DeltaStage),
		uint32(o.RMax),
		s,
		offchainPublicKeys,
		o.PeerIDs,
		o.ReportingPluginConfig,
		uint64(o.MaxDurationQuery),
		uint64(o.MaxDurationObservation),
		uint64(o.MaxDurationReport),
		uint64(o.MaxDurationShouldAcceptFinalizedReport),
		uint64(o.MaxDurationShouldTransmitAcceptedReport),
		&sharedSecretEncryptions,
	}
}

func enprotoSharedSecretEncryptions(e config.SharedSecretEncryptions) SharedSecretEncryptionsProto {
	encs := make([][]byte, 0, len(e.Encryptions))
	for _, enc := range e.Encryptions {
		enc := enc
		encs = append(encs, enc[:])
	}
	return SharedSecretEncryptionsProto{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		e.DiffieHellmanPoint[:],
		e.SharedSecretHash[:],
		encs,
	}
}
