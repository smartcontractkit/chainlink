package verifier

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	ErrVerificationFailed = errors.New("verification failed")

	ErrFailedUnmarshalPubkey          = fmt.Errorf("%w: failed to unmarshal pubkey", ErrVerificationFailed)
	ErrVerifyInvalidSignatureCount    = fmt.Errorf("%w: invalid signature count", ErrVerificationFailed)
	ErrVerifyMismatchedSignatureCount = fmt.Errorf("%w: mismatched signature count", ErrVerificationFailed)
	ErrVerifyInvalidSignature         = fmt.Errorf("%w: invalid signature", ErrVerificationFailed)
	ErrVerifySomeSignerUnauthorized   = fmt.Errorf("%w: node unauthorized", ErrVerificationFailed)
	ErrVerifyNonUniqueSignature       = fmt.Errorf("%w: signer has already signed", ErrVerificationFailed)
)

type SignedReport struct {
	RawRs         [][32]byte
	RawSs         [][32]byte
	RawVs         [32]byte
	ReportContext [3][32]byte
	Report        []byte
}

type Verifier interface {
	// Verify checks the report against its configuration, and then verifies signatures.
	// It replicates the Verifier contract's "verify" function for server side
	// report verification.
	// See also: contracts/src/v0.8/llo-feeds/Verifier.sol
	Verify(report SignedReport, f uint8, authorizedSigners []common.Address) (signers []common.Address, err error)
}

var _ Verifier = (*verifier)(nil)

type verifier struct{}

func NewVerifier() Verifier {
	return &verifier{}
}

func (v *verifier) Verify(sr SignedReport, f uint8, authorizedSigners []common.Address) (signers []common.Address, err error) {
	if len(sr.RawRs) != int(f+1) {
		return signers, fmt.Errorf("%w: expected the number of signatures (len(rs)) to equal the number of signatures required (f), but f=%d and len(rs)=%d", ErrVerifyInvalidSignatureCount, f+1, len(sr.RawRs))
	}
	if len(sr.RawRs) != len(sr.RawSs) {
		return signers, fmt.Errorf("%w: got %d rs and %d ss, expected equal", ErrVerifyMismatchedSignatureCount, len(sr.RawRs), len(sr.RawSs))
	}

	sigData := ReportToSigData(sr.ReportContext, sr.Report)

	signerMap := make(map[common.Address]bool)
	for _, signer := range authorizedSigners {
		signerMap[signer] = false
	}

	// Loop over every signature and collect errors. This wastes some CPU cycles, but we need to know everyone who
	// signed the report. Some risk mitigated by checking that the number of signatures matches the expected (F) earlier
	var verifyErrors error
	reportSigners := make([]common.Address, len(sr.RawRs)) // For logging + metrics, string for convenience
	for i := 0; i < len(sr.RawRs); i++ {
		sig := append(sr.RawRs[i][:], sr.RawSs[i][:]...)
		sig = append(sig, sr.RawVs[i]) // In the contract, you'll see vs+27. We don't do that here since geth adds +27 internally

		sigPubKey, err := crypto.Ecrecover(sigData, sig)
		if err != nil {
			verifyErrors = errors.Join(verifyErrors, fmt.Errorf("failed to recover signature: %w", err))
			continue
		}

		verified := crypto.VerifySignature(sigPubKey, sigData, sig[:64])
		if !verified {
			verifyErrors = errors.Join(verifyErrors, ErrVerifyInvalidSignature, fmt.Errorf("signature verification failed for pubKey: %x, sig: %x", sigPubKey, sig))
			continue
		}

		unmarshalledPub, err := crypto.UnmarshalPubkey(sigPubKey)
		if err != nil {
			verifyErrors = errors.Join(verifyErrors, ErrFailedUnmarshalPubkey, fmt.Errorf("public key=%x error=%w", sigPubKey, err))
			continue
		}

		address := crypto.PubkeyToAddress(*unmarshalledPub)
		reportSigners[i] = address
		encountered, authorized := signerMap[address]
		if !authorized {
			verifyErrors = errors.Join(verifyErrors, ErrVerifySomeSignerUnauthorized, fmt.Errorf("signer %s not in list of authorized nodes", address.String()))
			continue
		}
		if encountered {
			verifyErrors = errors.Join(verifyErrors, ErrVerifyNonUniqueSignature, fmt.Errorf("signer %s has already signed this report", address.String()))
			continue
		}
		signerMap[address] = true
		signers = append(signers, address)
	}
	return signers, verifyErrors
}

func ReportToSigData(reportCtx [3][32]byte, sr types.Report) []byte {
	sigData := crypto.Keccak256(sr)
	sigData = append(sigData, reportCtx[0][:]...)
	sigData = append(sigData, reportCtx[1][:]...)
	sigData = append(sigData, reportCtx[2][:]...)
	return crypto.Keccak256(sigData)
}
