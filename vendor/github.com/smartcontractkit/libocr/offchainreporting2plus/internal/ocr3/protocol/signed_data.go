package protocol

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Returns a byte slice whose first four bytes are the string "ocr3" and the rest
// of which is the sum returned by h. Used for domain separation vs ocr2, where
// we just directly sign sha256 hashes.
//
// Any signatures made with the OffchainKeyring should use ocr3DomainSeparatedSum!
func ocr3DomainSeparatedSum(h hash.Hash) []byte {
	result := make([]byte, 0, 4+32)
	result = append(result, []byte("ocr3")...)
	return h.Sum(result)
}

const signedObservationDomainSeparator = "ocr3 SignedObservation"

type SignedObservation struct {
	Observation types.Observation
	Signature   []byte
}

func MakeSignedObservation(
	ogid OutcomeGenerationID,
	seqNr uint64,
	query types.Query,
	observation types.Observation,
	signer func(msg []byte) (sig []byte, err error),
) (
	SignedObservation,
	error,
) {
	payload := signedObservationMsg(ogid, seqNr, query, observation)
	sig, err := signer(payload)
	if err != nil {
		return SignedObservation{}, err
	}
	return SignedObservation{observation, sig}, nil
}

func (so SignedObservation) Verify(ogid OutcomeGenerationID, seqNr uint64, query types.Query, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])
	// should never trigger since types.OffchainPublicKey is an array with length ed25519.PublicKeySize
	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedObservationMsg(ogid, seqNr, query, so.Observation), so.Signature)
	if !ok {
		return fmt.Errorf("SignedObservation has invalid signature")
	}

	return nil
}

func signedObservationMsg(ogid OutcomeGenerationID, seqNr uint64, query types.Query, observation types.Observation) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(signedObservationDomainSeparator))

	// ogid
	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	// seqNr
	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, seqNr)

	// query
	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	// observation
	_ = binary.Write(h, binary.BigEndian, uint64(len(observation)))
	_, _ = h.Write(observation)

	return ocr3DomainSeparatedSum(h)
}

type AttributedSignedObservation struct {
	SignedObservation SignedObservation
	Observer          commontypes.OracleID
}

type OutcomeInputsDigest [32]byte

func MakeOutcomeInputsDigest(
	ogid OutcomeGenerationID,
	previousOutcome ocr3types.Outcome,
	seqNr uint64,
	query types.Query,
	attributedObservations []types.AttributedObservation,
) OutcomeInputsDigest {
	h := sha256.New()

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, uint64(len(previousOutcome)))
	_, _ = h.Write(previousOutcome)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	_ = binary.Write(h, binary.BigEndian, uint64(len(attributedObservations)))
	for _, ao := range attributedObservations {

		_ = binary.Write(h, binary.BigEndian, uint64(len(ao.Observation)))
		_, _ = h.Write(ao.Observation)

		_ = binary.Write(h, binary.BigEndian, uint64(ao.Observer))
	}

	var result OutcomeInputsDigest
	h.Sum(result[:0])
	return result
}

type OutcomeDigest [32]byte

func MakeOutcomeDigest(outcome ocr3types.Outcome) OutcomeDigest {
	h := sha256.New()

	_, _ = h.Write(outcome)

	var result OutcomeDigest
	h.Sum(result[:0])
	return result
}

const prepareSignatureDomainSeparator = "ocr3 PrepareSignature"

type PrepareSignature []byte

func MakePrepareSignature(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
	signer func(msg []byte) ([]byte, error),
) (PrepareSignature, error) {
	return signer(prepareSignatureMsg(ogid, seqNr, outcomeInputsDigest, outcomeDigest))
}

func (sig PrepareSignature) Verify(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, prepareSignatureMsg(ogid, seqNr, outcomeInputsDigest, outcomeDigest), sig)
	if !ok {
		// Other less common causes include leader equivocation or actually invalid signatures.
		return fmt.Errorf("PrepareSignature failed to verify. This is commonly caused by non-determinism in the ReportingPlugin")
	}

	return nil
}

func prepareSignatureMsg(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(prepareSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(outcomeInputsDigest[:])

	_, _ = h.Write(outcomeDigest[:])

	return ocr3DomainSeparatedSum(h)
}

type AttributedPrepareSignature struct {
	Signature PrepareSignature
	Signer    commontypes.OracleID
}

const commitSignatureDomainSeparator = "ocr3 CommitSignature"

type CommitSignature []byte

func MakeCommitSignature(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
	signer func(msg []byte) ([]byte, error),
) (CommitSignature, error) {
	return signer(commitSignatureMsg(ogid, seqNr, outcomeDigest))
}

func (sig CommitSignature) Verify(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, commitSignatureMsg(ogid, seqNr, outcomeDigest), sig)
	if !ok {
		return fmt.Errorf("CommitSignature failed to verify")
	}

	return nil
}

func commitSignatureMsg(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(commitSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(outcomeDigest[:])

	return ocr3DomainSeparatedSum(h)
}

type AttributedCommitSignature struct {
	Signature CommitSignature
	Signer    commontypes.OracleID
}

type HighestCertifiedTimestamp struct {
	SeqNr                 uint64
	CommittedElsePrepared bool
}

func (t HighestCertifiedTimestamp) Less(t2 HighestCertifiedTimestamp) bool {
	return t.SeqNr < t2.SeqNr || t.SeqNr == t2.SeqNr && !t.CommittedElsePrepared && t2.CommittedElsePrepared
}

const signedHighestCertifiedTimestampDomainSeparator = "ocr3 SignedHighestCertifiedTimestamp"

type SignedHighestCertifiedTimestamp struct {
	HighestCertifiedTimestamp HighestCertifiedTimestamp
	Signature                 []byte
}

func MakeSignedHighestCertifiedTimestamp(
	ogid OutcomeGenerationID,
	highestCertifiedTimestamp HighestCertifiedTimestamp,
	signer func(msg []byte) ([]byte, error),
) (SignedHighestCertifiedTimestamp, error) {
	sig, err := signer(signedHighestCertifiedTimestampMsg(ogid, highestCertifiedTimestamp))
	if err != nil {
		return SignedHighestCertifiedTimestamp{}, err
	}

	return SignedHighestCertifiedTimestamp{
		highestCertifiedTimestamp,
		sig,
	}, nil
}

func (shct *SignedHighestCertifiedTimestamp) Verify(ogid OutcomeGenerationID, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedHighestCertifiedTimestampMsg(ogid, shct.HighestCertifiedTimestamp), shct.Signature)
	if !ok {
		return fmt.Errorf("SignedHighestCertifiedTimestamp signature failed to verify")
	}

	return nil
}

func signedHighestCertifiedTimestampMsg(
	ogid OutcomeGenerationID,
	highestCertifiedTimestamp HighestCertifiedTimestamp,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(signedHighestCertifiedTimestampDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, highestCertifiedTimestamp.SeqNr)

	var committedElsePreparedByte uint8
	if highestCertifiedTimestamp.CommittedElsePrepared {
		committedElsePreparedByte = 1
	} else {
		committedElsePreparedByte = 0
	}
	_, _ = h.Write([]byte{byte(committedElsePreparedByte)})

	return ocr3DomainSeparatedSum(h)
}

type AttributedSignedHighestCertifiedTimestamp struct {
	SignedHighestCertifiedTimestamp SignedHighestCertifiedTimestamp
	Signer                          commontypes.OracleID
}

type EpochStartProof struct {
	HighestCertified      CertifiedPrepareOrCommit
	HighestCertifiedProof []AttributedSignedHighestCertifiedTimestamp
}

func (qc *EpochStartProof) Verify(
	ogid OutcomeGenerationID,
	oracleIdentities []config.OracleIdentity,
	byzQuorumSize int,
) error {
	if byzQuorumSize != len(qc.HighestCertifiedProof) {
		return fmt.Errorf("wrong length of HighestCertifiedProof, expected %v for byz. quorum and got %v", byzQuorumSize, len(qc.HighestCertifiedProof))
	}

	maximumTimestamp := qc.HighestCertifiedProof[0].SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp

	seen := make(map[commontypes.OracleID]bool)
	for i, ashct := range qc.HighestCertifiedProof {
		if seen[ashct.Signer] {
			return fmt.Errorf("duplicate signature by %v", ashct.Signer)
		}
		seen[ashct.Signer] = true
		if !(0 <= int(ashct.Signer) && int(ashct.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", ashct.Signer)
		}
		if err := ashct.SignedHighestCertifiedTimestamp.Verify(ogid, oracleIdentities[ashct.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, ashct.Signer, oracleIdentities[ashct.Signer].OffchainPublicKey, err)
		}

		if maximumTimestamp.Less(ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp) {
			maximumTimestamp = ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp
		}
	}

	if qc.HighestCertified.Timestamp() != maximumTimestamp {
		return fmt.Errorf("mismatch between timestamp of HighestCertified (%v) and the max from HighestCertifiedProof (%v)", qc.HighestCertified.Timestamp(), maximumTimestamp)
	}

	if err := qc.HighestCertified.Verify(ogid.ConfigDigest, oracleIdentities, byzQuorumSize); err != nil {
		return fmt.Errorf("failed to verify HighestCertified: %w", err)
	}

	return nil
}

type CertifiedPrepareOrCommit interface {
	isCertifiedPrepareOrCommit()
	Epoch() uint64
	Timestamp() HighestCertifiedTimestamp
	IsGenesis() bool
	Verify(
		_ types.ConfigDigest,
		_ []config.OracleIdentity,
		byzQuorumSize int,
	) error
	CheckSize(n int, f int, limits ocr3types.ReportingPluginLimits, maxReportSigLen int) bool
}

var _ CertifiedPrepareOrCommit = &CertifiedPrepare{}

type CertifiedPrepare struct {
	PrepareEpoch             uint64
	SeqNr                    uint64
	OutcomeInputsDigest      OutcomeInputsDigest
	Outcome                  ocr3types.Outcome
	PrepareQuorumCertificate []AttributedPrepareSignature
}

func (hc *CertifiedPrepare) isCertifiedPrepareOrCommit() {}

func (hc *CertifiedPrepare) Epoch() uint64 {
	return uint64(hc.PrepareEpoch)
}

func (hc *CertifiedPrepare) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		hc.SeqNr,
		false,
	}
}

func (hc *CertifiedPrepare) IsGenesis() bool {
	return false
}

func (hc *CertifiedPrepare) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	byzQuorumSize int,
) error {
	if byzQuorumSize != len(hc.PrepareQuorumCertificate) {
		return fmt.Errorf("wrong number of signatures, expected %v for byz. quorum and got %v", byzQuorumSize, len(hc.PrepareQuorumCertificate))
	}

	ogid := OutcomeGenerationID{
		configDigest,
		hc.PrepareEpoch,
	}

	seen := make(map[commontypes.OracleID]bool)
	for i, aps := range hc.PrepareQuorumCertificate {
		if seen[aps.Signer] {
			return fmt.Errorf("duplicate signature by %v", aps.Signer)
		}
		seen[aps.Signer] = true
		if !(0 <= int(aps.Signer) && int(aps.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", aps.Signer)
		}
		if err := aps.Signature.Verify(ogid, hc.SeqNr, hc.OutcomeInputsDigest, MakeOutcomeDigest(hc.Outcome), oracleIdentities[aps.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, aps.Signer, oracleIdentities[aps.Signer].OffchainPublicKey, err)
		}
	}
	return nil
}

func (hc *CertifiedPrepare) CheckSize(n int, f int, limits ocr3types.ReportingPluginLimits, maxReportSigLen int) bool {
	if len(hc.Outcome) > limits.MaxOutcomeLength {
		return false
	}
	if len(hc.PrepareQuorumCertificate) != byzquorum.Size(n, f) {
		return false
	}
	for _, aps := range hc.PrepareQuorumCertificate {
		if len(aps.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return true
}

var _ CertifiedPrepareOrCommit = &CertifiedCommit{}

// The empty CertifiedCommit{} is the genesis value
type CertifiedCommit struct {
	CommitEpoch             uint64
	SeqNr                   uint64
	Outcome                 ocr3types.Outcome
	CommitQuorumCertificate []AttributedCommitSignature
}

func (hc *CertifiedCommit) isCertifiedPrepareOrCommit() {}

func (hc *CertifiedCommit) Epoch() uint64 {
	return uint64(hc.CommitEpoch)
}

func (hc *CertifiedCommit) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		hc.SeqNr,
		true,
	}
}

func (hc *CertifiedCommit) IsGenesis() bool {
	// We intentionally don't just compare with CertifiedCommit{}, because after
	// protobuf deserialization, we might end up with hc.Outcome = []byte{}
	return hc.CommitEpoch == 0 && hc.SeqNr == 0 && len(hc.Outcome) == 0 && len(hc.CommitQuorumCertificate) == 0
}

func (hc *CertifiedCommit) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	byzQuorumSize int,
) error {
	if hc.IsGenesis() {
		return nil
	}

	if byzQuorumSize != len(hc.CommitQuorumCertificate) {
		return fmt.Errorf("wrong number of signatures, expected %d for byz. quorum but got %d", byzQuorumSize, len(hc.CommitQuorumCertificate))
	}

	ogid := OutcomeGenerationID{
		configDigest,
		hc.CommitEpoch,
	}

	seen := make(map[commontypes.OracleID]bool)
	for i, acs := range hc.CommitQuorumCertificate {
		if seen[acs.Signer] {
			return fmt.Errorf("duplicate signature by %v", acs.Signer)
		}
		seen[acs.Signer] = true
		if !(0 <= int(acs.Signer) && int(acs.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", acs.Signer)
		}
		if err := acs.Signature.Verify(ogid, hc.SeqNr, MakeOutcomeDigest(hc.Outcome), oracleIdentities[acs.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, acs.Signer, oracleIdentities[acs.Signer].OffchainPublicKey, err)
		}
	}
	return nil
}

func (hc *CertifiedCommit) CheckSize(n int, f int, limits ocr3types.ReportingPluginLimits, maxReportSigLen int) bool {
	if hc.IsGenesis() {
		return true
	}

	if len(hc.Outcome) > limits.MaxOutcomeLength {
		return false
	}
	if len(hc.CommitQuorumCertificate) != byzquorum.Size(n, f) {
		return false
	}
	for _, acs := range hc.CommitQuorumCertificate {
		if len(acs.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return true
}
