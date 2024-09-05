package verifier

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
)

func Test_Verifier(t *testing.T) {
	t.Parallel()

	signedReportBinary := hexutil.MustDecode(`0x0006e1dde86b8a12add45546a14ea7e5efd10b67a373c6f4c41ecfa17d0005350000000000000000000000000000000000000000000000000000000000000201000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000022000000000000000000000000000000000000000000000000000000000000002800001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012000034c9214519c942ad0aa84a3dd31870e6efe8b3fcab4e176c5226879b26c77000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000669150aa0000000000000000000000000000000000001504e1e6c380271bb8b129ac8f7c0000000000000000000000000000000000001504e1e6c380271bb8b129ac8f7c00000000000000000000000000000000000000000000000000000000669150ab0000000000000000000000000000000000000000000000000000002482116240000000000000000000000000000000000000000000000000000000247625a04000000000000000000000000000000000000000000000000000000024880743400000000000000000000000000000000000000000000000000000000000000002710ac21df88ab70c8822b68be53d7bed65c82ffc9204c1d7ccf3c6c4048b3ca2cafb26e7bbd8f13fe626c946baa5ffcb444319c4229b945ea65d0c99c21978a100000000000000000000000000000000000000000000000000000000000000022c07843f17aa3ecd55f52e99e889906f825f49e4ddfa9c74ca487dd4ff101cc636108a5323be838e658dffa1be67bd91e99f68c4bf86936b76c5d8193b707597`)
	m := make(map[string]interface{})
	err := mercury.PayloadTypes.UnpackIntoMap(m, signedReportBinary)
	require.NoError(t, err)

	signedReport := SignedReport{
		RawRs:         m["rawRs"].([][32]byte),
		RawSs:         m["rawSs"].([][32]byte),
		RawVs:         m["rawVs"].([32]byte),
		ReportContext: m["reportContext"].([3][32]byte),
		Report:        m["report"].([]byte),
	}

	f := uint8(1)

	v := NewVerifier()

	t.Run("Verify errors with unauthorized signers", func(t *testing.T) {
		_, err := v.Verify(signedReport, f, []common.Address{})
		require.Error(t, err)
		assert.EqualError(t, err, "verification failed: node unauthorized\nsigner 0x3fc9FaA15d71EeD614e5322bd9554Fb35cC381d2 not in list of authorized nodes\nverification failed: node unauthorized\nsigner 0xBa6534da0E49c71cD9d0292203F1524876f33E23 not in list of authorized nodes")
	})

	t.Run("Verify succeeds with authorized signers", func(t *testing.T) {
		signers, err := v.Verify(signedReport, f, []common.Address{
			common.HexToAddress("0xde25e5b4005f611e356ce203900da4e37d72d58f"),
			common.HexToAddress("0x256431d41cf0d944f5877bc6c93846a9829dfc03"),
			common.HexToAddress("0x3fc9faa15d71eed614e5322bd9554fb35cc381d2"),
			common.HexToAddress("0xba6534da0e49c71cd9d0292203f1524876f33e23"),
		})
		require.NoError(t, err)
		assert.Equal(t, []common.Address{
			common.HexToAddress("0x3fc9faa15d71eed614e5322bd9554fb35cc381d2"),
			common.HexToAddress("0xBa6534da0E49c71cD9d0292203F1524876f33E23"),
		}, signers)
	})

	t.Run("Verify fails if report has been tampered with", func(t *testing.T) {
		badReport := signedReport
		badReport.Report = []byte{0x0011}
		_, err := v.Verify(badReport, f, []common.Address{
			common.HexToAddress("0xde25e5b4005f611e356ce203900da4e37d72d58f"),
			common.HexToAddress("0x256431d41cf0d944f5877bc6c93846a9829dfc03"),
			common.HexToAddress("0x3fc9faa15d71eed614e5322bd9554fb35cc381d2"),
			common.HexToAddress("0xba6534da0e49c71cd9d0292203f1524876f33e23"),
		})

		require.Error(t, err)
	})

	t.Run("Verify fails if rawVs has been changed", func(t *testing.T) {
		badReport := signedReport
		badReport.RawVs = [32]byte{0x0011}
		_, err := v.Verify(badReport, f, []common.Address{
			common.HexToAddress("0xde25e5b4005f611e356ce203900da4e37d72d58f"),
			common.HexToAddress("0x256431d41cf0d944f5877bc6c93846a9829dfc03"),
			common.HexToAddress("0x3fc9faa15d71eed614e5322bd9554fb35cc381d2"),
			common.HexToAddress("0xba6534da0e49c71cd9d0292203f1524876f33e23"),
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to recover signature: invalid signature recovery id")
	})
}
