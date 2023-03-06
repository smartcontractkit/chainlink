package protocol

import (
	"bytes"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type AttributedObservation struct {
	Observation types.Observation
	Observer    commontypes.OracleID
}

type AttestedReportOne struct {
	Skip      bool
	Report    types.Report
	Signature []byte
}

func MakeAttestedReportOneSkip() AttestedReportOne {
	return AttestedReportOne{true, nil, nil}
}

func MakeAttestedReportOneNoskip(
	repctx types.ReportContext,
	report types.Report,
	signer func(types.ReportContext, types.Report) ([]byte, error),
) (AttestedReportOne, error) {
	sig, err := signer(repctx, report)
	if err != nil {
		return AttestedReportOne{}, fmt.Errorf("error while signing in MakeAttestedReportOneNoskip: %w", err)
	}

	return AttestedReportOne{false, report, sig}, nil
}

func (rep AttestedReportOne) EqualExceptSignature(rep2 AttestedReportOne) bool {
	return rep.Skip == rep2.Skip && bytes.Equal(rep.Report, rep2.Report)
}

// Verify is used by the leader to check the signature a process attaches to its
// report message (the c.Sig value.)
func (aro *AttestedReportOne) Verify(contractSigner types.OnchainKeyring, publicKey types.OnchainPublicKey, repctx types.ReportContext) (err error) {
	if aro.Skip {
		if len(aro.Report) != 0 || len(aro.Signature) != 0 {
			return fmt.Errorf("AttestedReportOne with Skip=true has non-empty Report or Signature")
		}
	} else {
		ok := contractSigner.Verify(publicKey, repctx, aro.Report, aro.Signature)
		if !ok {
			return fmt.Errorf("failed to verify signature on AttestedReportOne")
		}
	}
	return nil
}

type AttestedReportMany struct {
	Report               types.Report
	AttributedSignatures []types.AttributedOnchainSignature
}

func (rep *AttestedReportMany) VerifySignatures(
	numSignatures int,
	onchainKeyring types.OnchainKeyring,
	oracleIdentities []config.OracleIdentity,
	repctx types.ReportContext,
) error {
	if numSignatures != len(rep.AttributedSignatures) {
		return fmt.Errorf("wrong number of signatures, expected %v and got %v", numSignatures, len(rep.AttributedSignatures))
	}
	seen := make(map[commontypes.OracleID]bool)
	for i, sig := range rep.AttributedSignatures {
		if seen[sig.Signer] {
			return fmt.Errorf("duplicate Signature by %v", sig.Signer)
		}
		seen[sig.Signer] = true
		if !(0 <= int(sig.Signer) && int(sig.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", sig.Signer)
		}
		if !onchainKeyring.Verify(oracleIdentities[sig.Signer].OnchainPublicKey, repctx, rep.Report, sig.Signature) {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify", i, sig.Signer, oracleIdentities[sig.Signer].OnchainPublicKey)
		}
	}
	return nil
}
