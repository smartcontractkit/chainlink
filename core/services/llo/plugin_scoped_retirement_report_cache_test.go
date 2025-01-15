package llo

import (
	"errors"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

func FuzzPluginScopedRetirementReportCache_CheckAttestedRetirementReport(f *testing.F) {
	f.Add([]byte("not a protobuf"))
	f.Add([]byte{0x0a, 0x00})             // empty protobuf
	f.Add([]byte{0x0a, 0x02, 0x08, 0x01}) // invalid protobuf
	f.Add(([]byte)(nil))
	f.Add([]byte{})

	rrc := &mockRetirementReportCache{}
	v := &mockVerifier{}
	c := &mockCodec{}
	psrrc := NewPluginScopedRetirementReportCache(rrc, v, c)

	exampleDigest := ocr2types.ConfigDigest{1}

	f.Fuzz(func(t *testing.T, data []byte) {
		psrrc.CheckAttestedRetirementReport(exampleDigest, data) //nolint:errcheck // test that it doesn't panic, don't care about errors
	})
}

type mockRetirementReportCache struct {
	arr    []byte
	cfg    Config
	exists bool
}

func (m *mockRetirementReportCache) AttestedRetirementReport(digest ocr2types.ConfigDigest) ([]byte, bool) {
	return m.arr, m.exists
}
func (m *mockRetirementReportCache) Config(cd ocr2types.ConfigDigest) (Config, bool) {
	return m.cfg, m.exists
}

type mockVerifier struct {
	verify func(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo], signature []byte) bool
}

func (m *mockVerifier) Verify(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo], signature []byte) bool {
	return m.verify(key, digest, seqNr, r, signature)
}

type mockCodec struct {
	decode func([]byte) (datastreamsllo.RetirementReport, error)
}

func (m *mockCodec) Encode(datastreamsllo.RetirementReport) ([]byte, error) {
	panic("not implemented")
}
func (m *mockCodec) Decode(b []byte) (datastreamsllo.RetirementReport, error) {
	return m.decode(b)
}

func Test_PluginScopedRetirementReportCache(t *testing.T) {
	rrc := &mockRetirementReportCache{}
	v := &mockVerifier{}
	c := &mockCodec{}
	psrrc := NewPluginScopedRetirementReportCache(rrc, v, c)
	exampleDigest := ocr2types.ConfigDigest{1}
	exampleDigest2 := ocr2types.ConfigDigest{2}

	exampleUnattestedSerializedRetirementReport := []byte("foo example unattested retirement report")

	validArr := AttestedRetirementReport{
		RetirementReport: exampleUnattestedSerializedRetirementReport,
		SeqNr:            42,
		Sigs: []*AttributedOnchainSignature{
			{
				Signer:    0,
				Signature: []byte("bar0"),
			},
			{
				Signer:    1,
				Signature: []byte("bar1"),
			},
			{
				Signer:    2,
				Signature: []byte("bar2"),
			},
			{
				Signer:    3,
				Signature: []byte("bar3"),
			},
		},
	}
	serializedValidArr, err := proto.Marshal(&validArr)
	require.NoError(t, err)

	t.Run("CheckAttestedRetirementReport", func(t *testing.T) {
		t.Run("invalid", func(t *testing.T) {
			// config missing
			_, err := psrrc.CheckAttestedRetirementReport(exampleDigest, []byte("not valid"))
			assert.EqualError(t, err, "Verify failed; predecessor config not found for config digest 0100000000000000000000000000000000000000000000000000000000000000")

			rrc.cfg = Config{Digest: exampleDigest}
			rrc.exists = true

			// unmarshal failure
			_, err = psrrc.CheckAttestedRetirementReport(exampleDigest, []byte("not valid"))
			require.Error(t, err)
			assert.Contains(t, err.Error(), "Verify failed; failed to unmarshal protobuf: proto")

			// config is invalid (no signers)
			_, err = psrrc.CheckAttestedRetirementReport(exampleDigest, serializedValidArr)
			assert.EqualError(t, err, "Verify failed; attested report signer index out of bounds (got: 0, max: -1)")

			rrc.cfg = Config{Digest: exampleDigest, Signers: [][]byte{[]byte{0}, []byte{1}, []byte{2}, []byte{3}}, F: 1}

			// no valid sigs
			v.verify = func(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo], signature []byte) bool {
				return false
			}
			_, err = psrrc.CheckAttestedRetirementReport(exampleDigest, serializedValidArr)
			assert.EqualError(t, err, "Verify failed; not enough valid signatures (got: 0, need: 2)")

			// not enough valid sigs
			v.verify = func(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo], signature []byte) bool {
				return string(signature) == "bar0"
			}
			_, err = psrrc.CheckAttestedRetirementReport(exampleDigest, serializedValidArr)
			assert.EqualError(t, err, "Verify failed; not enough valid signatures (got: 1, need: 2)")

			// enough valid sigs, but codec decode fails
			v.verify = func(key types.OnchainPublicKey, digest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[llotypes.ReportInfo], signature []byte) bool {
				if string(signature) == "bar0" || string(signature) == "bar3" {
					return true
				}
				return false
			}
			c.decode = func([]byte) (datastreamsllo.RetirementReport, error) {
				return datastreamsllo.RetirementReport{}, errors.New("codec decode failed")
			}
			_, err = psrrc.CheckAttestedRetirementReport(exampleDigest, serializedValidArr)
			assert.EqualError(t, err, "Verify failed; failed to decode retirement report: codec decode failed")

			exampleRetirementReport := datastreamsllo.RetirementReport{ValidAfterSeconds: map[llotypes.ChannelID]uint32{
				0: 1,
			},
			}

			// enough valid sigs and codec decode succeeds
			c.decode = func(b []byte) (datastreamsllo.RetirementReport, error) {
				assert.Equal(t, exampleUnattestedSerializedRetirementReport, b)
				return exampleRetirementReport, nil
			}
			decoded, err := psrrc.CheckAttestedRetirementReport(exampleDigest, serializedValidArr)
			assert.NoError(t, err)
			assert.Equal(t, exampleRetirementReport, decoded)
		})
	})
	t.Run("AttestedRetirementReport", func(t *testing.T) {
		rrc.arr = []byte("foo")
		rrc.exists = true

		// exists
		arr, err := psrrc.AttestedRetirementReport(exampleDigest)
		assert.NoError(t, err)
		assert.Equal(t, rrc.arr, arr)

		rrc.exists = false

		// doesn't exist
		arr, err = psrrc.AttestedRetirementReport(exampleDigest2)
		assert.NoError(t, err)
		assert.Nil(t, arr)
	})
}
