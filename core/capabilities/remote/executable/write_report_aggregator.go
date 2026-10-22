package executable

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"maps"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	reportmetadata "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// writeReportSignatureAggregator builds WriteReport payloads from the signatures of all
// requests received for a report, so that a single node providing invalid signatures cannot
// make the report fail onchain verification. It is used together with
// writeReportExcludeSignaturesHasher, which groups requests that differ only in their signatures.
//
// Signatures are verified with the onchain keyring of the report's chain family against the
// signers of the calling DON's consensus OCR3 config, where the signer index is the SignerId.
// Verification is deferred until the request quorum is reached, so requests that never reach
// a quorum cost no signature verification
//
// Whenever the signatures cannot be verified (unknown signer set, unexpected report header, or
// fewer than F+1 valid signatures once the request quorum is reached), the request is executed
// with the payload of the latest request, so the aggregation never fails a request that would
// succeed without it.
type writeReportSignatureAggregator struct {
	// workflowDONSigners holds, per workflow DON, the onchain public keys of each oracle by chain family.
	workflowDONSigners map[uint32][]map[string]ocrtypes.OnchainPublicKey
	verifiers          map[writeReportFamily]*ocrcommon.OCR3OnchainKeyringAdapter
	lggr               logger.Logger
}

var _ types.RequestAccumulatorFactory = (*writeReportSignatureAggregator)(nil)

var writeReportFamilyChainTypes = map[writeReportFamily]corekeys.ChainType{
	writeReportFamilyEVM:     corekeys.EVM,
	writeReportFamilySolana:  corekeys.Solana,
	writeReportFamilyAptos:   corekeys.Aptos,
	writeReportFamilyStellar: corekeys.Stellar,
}

// NewWriteReportSignatureAggregator returns a RequestAccumulatorFactory for WriteReport methods.
// workflowDONSigners maps workflow DON IDs to the multichain-encoded signers of their consensus
// OCR3 config, ordered by oracle index.
func NewWriteReportSignatureAggregator(workflowDONSigners map[uint32][]ocrtypes.OnchainPublicKey, lggr logger.Logger) (types.RequestAccumulatorFactory, error) {
	lggr = logger.Named(lggr, "WriteReportSignatureAggregator")

	// Keyrings verify against the given public key only, so a random key bundle can verify
	// signatures of any signer.
	verifiers := make(map[writeReportFamily]*ocrcommon.OCR3OnchainKeyringAdapter, len(writeReportFamilyChainTypes))
	for family, chainType := range writeReportFamilyChainTypes {
		kb, err := ocr2key.New(chainType)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s keyring: %w", chainType, err)
		}
		verifiers[family] = ocrcommon.NewOCR3OnchainKeyringAdapter(kb)
	}

	signers := make(map[uint32][]map[string]ocrtypes.OnchainPublicKey, len(workflowDONSigners))
	for donID, donSigners := range workflowDONSigners {
		keys := make([]map[string]ocrtypes.OnchainPublicKey, len(donSigners))
		for i, s := range donSigners {
			decoded, err := ocrcommon.UnmarshalMultichainPublicKey(s)
			if err != nil {
				lggr.Warnw("failed to decode workflow DON signer, its signatures will be dropped", "workflowDonID", donID, "signerID", i, "err", err)
				continue
			}
			keys[i] = decoded
		}
		signers[donID] = keys
	}

	return &writeReportSignatureAggregator{
		workflowDONSigners: signers,
		verifiers:          verifiers,
		lggr:               lggr,
	}, nil
}

func (a *writeReportSignatureAggregator) NewAccumulator(_ context.Context, capabilityID string, callingDon commoncap.DON) (types.RequestAccumulator, error) {
	family, err := getWriteReportFamily(capabilityID)
	if err != nil {
		return nil, err
	}

	acc := &writeReportAccumulator{
		family:     family,
		callingDon: callingDon,
		verifier:   a.verifiers[family],
		candidates: map[uint32]map[string]*signatureCandidate{},
		lggr:       logger.With(a.lggr, "capabilityID", capabilityID, "callingDonID", callingDon.ID),
	}

	donSigners := a.workflowDONSigners[callingDon.ID]
	if len(donSigners) == 0 {
		acc.useFallback("missing_signers", "no signers known for calling DON")
		return acc, nil
	}
	chainType := string(writeReportFamilyChainTypes[family])
	acc.signers = make([]ocrtypes.OnchainPublicKey, len(donSigners))
	found := false
	for i, keys := range donSigners {
		if key, ok := keys[chainType]; ok {
			acc.signers[i] = key
			found = true
		}
	}
	if !found {
		acc.useFallback("missing_family_signers", "no signers of the report chain family known for calling DON", "chainType", chainType)
	}
	return acc, nil
}

type signatureCandidate struct {
	signature []byte
	senders   []p2ptypes.PeerID
}

type writeReportAccumulator struct {
	family     writeReportFamily
	callingDon commoncap.DON
	verifier   *ocrcommon.OCR3OnchainKeyringAdapter
	// signers holds the public keys of the report chain family by SignerId, nil if unknown.
	signers []ocrtypes.OnchainPublicKey
	lggr    logger.Logger

	latestPayload []byte
	received      bool

	// base is the latest parsed request; its signatures are replaced by the aggregated ones.
	base          commoncap.CapabilityRequest
	baseWrReq     writeReportRequest
	headerChecked bool

	// candidates hold the unverified signatures received
	candidates map[uint32]map[string]*signatureCandidate
	// aggregated is the payload built from verified signatures once computed
	aggregated []byte

	fallback bool
}

var _ types.RequestAccumulator = (*writeReportAccumulator)(nil)

func (a *writeReportAccumulator) useFallback(reason, msg string, keyvals ...any) {
	a.candidates = nil
	a.fallback = true
	a.lggr.Warnw(msg+"; executing with the payload of the latest request", append(keyvals, "reason", reason)...)
}

func (a *writeReportAccumulator) Add(sender p2ptypes.PeerID, payload []byte) error {
	if a.fallback {
		a.latestPayload = payload
		a.received = true
		return nil
	}
	if a.aggregated != nil {
		return nil
	}

	req, err := pb.UnmarshalCapabilityRequest(payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal capability request: %w", err)
	}
	if req.Payload == nil {
		return errors.New("capability request payload is nil")
	}
	wrReq, err := unmarshalWriteReportRequest(a.family, req.Payload)
	if err != nil {
		return err
	}
	report := wrReq.GetReport()

	a.latestPayload = payload
	a.received = true

	// The report header is covered by the request hash, so it is identical for all requests.
	if !a.headerChecked {
		a.headerChecked = true
		a.checkReportHeader(report)
		if a.fallback {
			return nil
		}
	}

	a.base = req
	a.baseWrReq = wrReq
	a.addCandidates(sender, report.Sigs)
	return nil
}

func (a *writeReportAccumulator) checkReportHeader(report *sdk.ReportResponse) {
	if len(report.ConfigDigest) != len(ocrtypes.ConfigDigest{}) {
		a.useFallback("invalid_config_digest", "invalid report config digest length", "length", len(report.ConfigDigest))
		return
	}

	md, _, err := reportmetadata.Decode(report.RawReport)
	if err != nil {
		a.useFallback("undecodable_report_metadata", "failed to decode report metadata", "err", err)
		return
	}
	if md.DONID != a.callingDon.ID {
		a.useFallback("don_id_mismatch", "report DON ID does not match calling DON ID, signer set cannot be verified",
			"reportDonID", md.DONID)
		return
	}
	if md.DONConfigVersion != a.callingDon.ConfigVersion {
		a.useFallback("config_version_mismatch", "report DON config version does not match calling DON config version, signer set cannot be verified",
			"reportDonConfigVersion", md.DONConfigVersion, "callingDonConfigVersion", a.callingDon.ConfigVersion)
	}
}

// addCandidates buffers the signatures of a request for verification once the quorum is reached.
// A request contributes at most one signature per signer, which bounds the verification work a
// malicious sender  can cause
func (a *writeReportAccumulator) addCandidates(sender p2ptypes.PeerID, sigs []*sdk.AttributedSignature) {
	if len(sigs) > len(a.signers) {
		a.lggr.Warnw("dropping report signatures of request with more signatures than signers",
			"sender", sender.String(), "signatures", len(sigs), "signers", len(a.signers))
	}
	var invalid []string
	seen := make(map[uint32]struct{}, len(sigs))
	for _, sig := range sigs {
		if sig == nil {
			invalid = append(invalid, "nil signature")
			continue
		}

		if int(sig.SignerId) >= len(a.signers) || a.signers[sig.SignerId] == nil {
			invalid = append(invalid, fmt.Sprintf("signer %d: unknown signer", sig.SignerId))
			continue
		}

		if _, ok := seen[sig.SignerId]; ok {
			invalid = append(invalid, fmt.Sprintf("signer %d: duplicate signer", sig.SignerId))
			continue
		}
		seen[sig.SignerId] = struct{}{}

		bySig, ok := a.candidates[sig.SignerId]
		if !ok {
			bySig = map[string]*signatureCandidate{}
			a.candidates[sig.SignerId] = bySig
		}
		c, ok := bySig[string(sig.Signature)]
		if !ok {
			c = &signatureCandidate{signature: bytes.Clone(sig.Signature)}
			bySig[string(sig.Signature)] = c
		}
		c.senders = append(c.senders, sender)
	}

	if len(invalid) > 0 {
		a.lggr.Warnw("dropping invalid report signatures", "sender", sender.String(), "invalid", invalid)
	}
}

// Payload is called once the request quorum is reached. Since at least one of the quorum
// requests comes from an honest node, fewer than F+1 valid signatures means the signatures
// could not be verified rather than an attack, so the latest payload is used instead.
func (a *writeReportAccumulator) Payload() ([]byte, bool, error) {
	if a.aggregated != nil {
		return a.aggregated, true, nil
	}
	if a.fallback || a.baseWrReq == nil {
		return a.latestPayload, a.received, nil
	}

	required := int(a.callingDon.F) + 1
	sigs := a.verifiedSignatures(required)
	if len(sigs) < required {
		a.useFallback("not_enough_valid_signatures", "not enough valid signatures in quorum requests",
			"validSignatures", len(sigs), "required", required)
		return a.latestPayload, a.received, nil
	}

	payload, err := a.aggregatedPayload(sigs)
	if err != nil {
		a.useFallback("aggregation_failed", "failed to build payload with aggregated signatures", "err", err)
		return a.latestPayload, a.received, nil
	}
	a.aggregated = payload
	a.candidates = nil
	return payload, true, nil
}

// verifiedSignatures returns up to required valid signatures of distinct signers. Signers are
// tried in ascending SignerId order and their candidates in signature byte order, so the result
// does not depend on the order in which requests arrived
func (a *writeReportAccumulator) verifiedSignatures(required int) []*sdk.AttributedSignature {
	report := a.baseWrReq.GetReport()
	digest := ocrtypes.ConfigDigest(report.ConfigDigest)
	reportWithInfo := ocr3types.ReportWithInfo[[]byte]{Report: report.RawReport}

	var invalid []string
	sigs := make([]*sdk.AttributedSignature, 0, required)
	usedSigners := map[string]struct{}{}
	for _, signerID := range slices.Sorted(maps.Keys(a.candidates)) {
		if len(sigs) == required {
			break
		}
		signer := a.signers[signerID]
		if _, ok := usedSigners[string(signer)]; ok {
			continue
		}
		bySig := a.candidates[signerID]
		for _, key := range slices.Sorted(maps.Keys(bySig)) {
			c := bySig[key]
			if a.isMalleable(c.signature) {
				invalid = append(invalid, fmt.Sprintf("signer %d from %v: malleable signature", signerID, c.senders))
				continue
			}
			if !a.verifier.Verify(signer, digest, report.SeqNr, reportWithInfo, c.signature) {
				invalid = append(invalid, fmt.Sprintf("signer %d from %v: invalid signature", signerID, c.senders))
				continue
			}
			sigs = append(sigs, &sdk.AttributedSignature{Signature: c.signature, SignerId: signerID})
			usedSigners[string(signer)] = struct{}{}
			break
		}
	}

	if len(invalid) > 0 {
		a.lggr.Warnw("dropping invalid report signatures", "invalid", invalid)
	}

	return sigs
}

var secp256k1HalfN = new(big.Int).Rsh(crypto.S256().Params().N, 1)

// isMalleable reports whether a secp256k1 signature has a high S value. Keyring verification
// accepts both S and N-S, but onchain verification may only accept the low S form produced by
// honest signers.
func (a *writeReportAccumulator) isMalleable(sig []byte) bool {
	if a.family != writeReportFamilyEVM && a.family != writeReportFamilySolana {
		return false
	}
	if len(sig) != crypto.SignatureLength {
		return false
	}
	return new(big.Int).SetBytes(sig[32:64]).Cmp(secp256k1HalfN) > 0
}

func (a *writeReportAccumulator) aggregatedPayload(sigs []*sdk.AttributedSignature) ([]byte, error) {
	wrReq, ok := proto.Clone(a.baseWrReq).(writeReportRequest)
	if !ok {
		return nil, errors.New("failed to clone WriteReportRequest")
	}
	wrReq.GetReport().Sigs = sigs

	wrAny, err := anypb.New(wrReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal WriteReportRequest to anypb: %w", err)
	}
	req := a.base
	req.Payload = wrAny
	payload, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability request: %w", err)
	}
	return payload, nil
}
