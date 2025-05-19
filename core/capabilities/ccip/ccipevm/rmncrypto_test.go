package ccipevm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func Test_VerifyRmnReportSignatures(t *testing.T) {
	// NOTE: The following test data (public keys, signatures, ...) are shared from the RMN team.

	onchainRmnRemoteAddr := common.HexToAddress("0x7821bcd6944457d17c631157efeb0c621baa76eb")

	rmnHomeContractConfigDigestHex := "0x785936570d1c7422ef30b7da5555ad2f175fa2dd97a2429a2e71d1e07c94e060"
	rmnHomeContractConfigDigest := common.FromHex(rmnHomeContractConfigDigestHex)
	require.Len(t, rmnHomeContractConfigDigest, 32)
	var rmnHomeContractConfigDigest32 [32]byte
	copy(rmnHomeContractConfigDigest32[:], rmnHomeContractConfigDigest)

	rootHex := "0x48e688aefc20a04fdec6b8ff19df358fd532455659dcf529797cda358e9e5205"
	root := common.FromHex(rootHex)
	require.Len(t, root, 32)
	var root32 [32]byte
	copy(root32[:], root)

	onRampAddr := common.HexToAddress("0x6662cb20464f4be557262693bea0409f068397ed")

	destChainEvmID := int64(4083663998511321420)

	reportData := cciptypes.RMNReport{
		ReportVersionDigest:         cciptypes.Bytes32(crypto.Keccak256Hash([]byte("RMN_V1_6_ANY2EVM_REPORT"))),
		DestChainID:                 cciptypes.NewBigIntFromInt64(destChainEvmID),
		DestChainSelector:           5266174733271469989,
		RmnRemoteContractAddress:    common.HexToAddress("0x3d015cec4411357eff4ea5f009a581cc519f75d3").Bytes(),
		OfframpAddress:              common.HexToAddress("0xc5cdb7711a478058023373b8ae9e7421925140f8").Bytes(),
		RmnHomeContractConfigDigest: rmnHomeContractConfigDigest32,
		LaneUpdates: []cciptypes.RMNLaneUpdate{
			{
				SourceChainSelector: 8258882951688608272,
				OnRampAddress:       onRampAddr.Bytes(),
				MinSeqNr:            9018980618932210108,
				MaxSeqNr:            8239368306600774074,
				MerkleRoot:          root32,
			},
		},
	}

	rmnCrypto := NewEVMRMNCrypto(logger.Test(t))

	r, _ := cciptypes.NewBytes32FromString("0x89546b4652d0377062a398e413344e4da6034ae877c437d0efe0e5246b70a9a1")
	s, _ := cciptypes.NewBytes32FromString("0x95eef2d24d856ccac3886db8f4aebea60684ed73942392692908fed79a679b4e")

	err := rmnCrypto.VerifyReportSignatures(
		t.Context(),
		[]cciptypes.RMNECDSASignature{{R: r, S: s}},
		reportData,
		[]cciptypes.UnknownAddress{onchainRmnRemoteAddr.Bytes()},
	)
	assert.NoError(t, err)
}
