package ocrcommon_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

var _ ocrtypes.OnchainKeyring = (*fakeOnchainKeyring)(nil)

var (
	account      ocrtypes.Account = "Test-Account"
	configDigest                  = ocrtypes.ConfigDigest([]byte("kKfYauxXBMjuP5EuuyacN6BwCfKJnP6d"))
	seqNr        uint64           = 11
	rwi                           = ocr3types.ReportWithInfo[[]byte]{
		Report: []byte("report"),
		Info:   []byte("info"),
	}
	signatures = []types.AttributedOnchainSignature{{
		Signature: []byte("signature1"),
		Signer:    1,
	}, {
		Signature: []byte("signature2"),
		Signer:    2,
	}}
	pubKey             = ocrtypes.OnchainPublicKey("pub-key")
	maxSignatureLength = 12
	sigs               = []byte("some-signatures")
)

type fakeOnchainKeyring struct {
}

func (f fakeOnchainKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return pubKey
}

func (f fakeOnchainKeyring) Sign(rc ocrtypes.ReportContext, r ocrtypes.Report) (signature []byte, err error) {
	if !reflect.DeepEqual(rc.ConfigDigest, configDigest) {
		return nil, fmt.Errorf("expected configDigest %v but got %v", configDigest, rc.ReportTimestamp.ConfigDigest)
	}

	if rc.Epoch != uint32(seqNr) {
		return nil, fmt.Errorf("expected Epoch %v but got %v", seqNr, rc.Epoch)
	}

	if rc.Round != 0 {
		return nil, fmt.Errorf("expected Round %v but got %v", 0, rc.Round)
	}

	if !reflect.DeepEqual(r, rwi.Report) {
		return nil, fmt.Errorf("expected Report %v but got %v", rwi.Report, r)
	}
	return nil, nil
}

func (f fakeOnchainKeyring) Verify(pk ocrtypes.OnchainPublicKey, rc ocrtypes.ReportContext, r ocrtypes.Report, signature []byte) bool {
	if !reflect.DeepEqual(pk, pubKey) {
		return false
	}

	if !reflect.DeepEqual(rc.ConfigDigest, configDigest) {
		return false
	}

	if rc.Epoch != uint32(seqNr) {
		return false
	}

	if rc.Round != 0 {
		return false
	}

	if !reflect.DeepEqual(r, rwi.Report) {
		return false
	}

	if !reflect.DeepEqual(signature, sigs) {
		return false
	}

	return true
}

func (f fakeOnchainKeyring) MaxSignatureLength() int {
	return maxSignatureLength
}

func TestOCR3OnchainKeyringAdapter(t *testing.T) {
	kr := ocrcommon.NewOCR3OnchainKeyringAdapter(fakeOnchainKeyring{})

	_, err := kr.Sign(configDigest, seqNr, rwi)
	require.NoError(t, err)
	require.True(t, kr.Verify(pubKey, configDigest, seqNr, rwi, sigs))

	require.Equal(t, pubKey, kr.PublicKey())
	require.Equal(t, maxSignatureLength, kr.MaxSignatureLength())
}

var _ ocrtypes.ContractTransmitter = (*fakeContractTransmitter)(nil)

type fakeContractTransmitter struct {
}

func (f fakeContractTransmitter) Transmit(ctx context.Context, rc ocrtypes.ReportContext, report ocrtypes.Report, s []ocrtypes.AttributedOnchainSignature) error {
	if !reflect.DeepEqual(report, rwi.Report) {
		return fmt.Errorf("expected Report %v but got %v", rwi.Report, report)
	}

	if !reflect.DeepEqual(s, signatures) {
		return fmt.Errorf("expected signatures %v but got %v", signatures, s)
	}

	if !reflect.DeepEqual(rc.ConfigDigest, configDigest) {
		return fmt.Errorf("expected configDigest %v but got %v", configDigest, rc.ReportTimestamp.ConfigDigest)
	}

	if rc.Epoch != uint32(seqNr) {
		return fmt.Errorf("expected Epoch %v but got %v", seqNr, rc.Epoch)
	}

	if rc.Round != 0 {
		return fmt.Errorf("expected Round %v but got %v", 0, rc.Round)
	}

	return nil
}

func (f fakeContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (configDigest ocrtypes.ConfigDigest, epoch uint32, err error) {
	panic("not implemented")
}

func (f fakeContractTransmitter) FromAccount() (ocrtypes.Account, error) {
	return account, nil
}

func TestContractTransmitter(t *testing.T) {
	ct := ocrcommon.NewOCR3ContractTransmitterAdapter(fakeContractTransmitter{})

	require.NoError(t, ct.Transmit(context.Background(), configDigest, seqNr, rwi, signatures))

	a, err := ct.FromAccount()
	require.NoError(t, err)
	require.Equal(t, a, account)
}
