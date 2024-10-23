package ccipevm

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_encoding_utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	ctx := tests.Context(t)

	rmnCrypto := NewEVMRMNCrypto(logger.Test(t))

	r, _ := cciptypes.NewBytes32FromString("0x89546b4652d0377062a398e413344e4da6034ae877c437d0efe0e5246b70a9a1")
	s, _ := cciptypes.NewBytes32FromString("0x95eef2d24d856ccac3886db8f4aebea60684ed73942392692908fed79a679b4e")

	err := rmnCrypto.VerifyReportSignatures(
		ctx,
		[]cciptypes.RMNECDSASignature{{R: r, S: s}},
		reportData,
		[]cciptypes.UnknownAddress{onchainRmnRemoteAddr.Bytes()},
	)
	assert.NoError(t, err)
}

func TestIt2(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString("llGUN4Pb+Bk1pg6Y8hip2bWyiCP7Iii72RMg1jL6z1MAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAkhAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAs1ZzPiKSRPIAAAAAAAAAAAAAAABTcPeMavLanPZkI4Kjp1+dWuycwQAAAAAAAAAAAAAAAK1dV62bsX003ruIVmqy9duHnMRvAAu+bkQ/7qY52YeMrEtOmPWbDBhp3WAWOtdvdvuly38AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAu5jSVHvcbRgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARxh/vej3RU5Q0GcEQEDExbte3o9dZE8NMvoYoAz9ZJPAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAB6IIihv8nYHFU2iuFowsAlcMuBTw==")
	require.NoError(t, err)
	t.Logf("b: %v", b)
}

func TestIt(t *testing.T) {

	const (
		_remoteRpc     = "http://localhost:52066"
		_rmnRemoteAddr = "0x5370f78c6af2da9cf6642382a3a75f9d5aec9cc1"

		_onchainPubKey = "0x1fd10272c19d6588d3527dc960aee3e201f1f1c2"

		_onRampAddr  = "0x0000000000000000000000007a2088a1bfc9d81c55368ae168c2c02570cb814f"
		_offRampAddr = "0xad5d57ad9bb17d34debb88566ab2f5db879cc46f"

		_merkleRoot     = "1c61fef7a3dd153943419c1101031316ed7b7a3d75913c34cbe8628033f5924f"
		_sourceChainSel = 3379446385462418246

		_r = "db4ced6758edc4f40d9b770fff4e07b91878275afbb37fe9903a6ef884f08eda"
		_s = "641a141153b1542f4e95269cca761109d912ba6c449e25ef47e453253eb90e8d"
	)

	cl, err := ethclient.Dial(_remoteRpc)
	require.NoError(t, err)

	rmnRemoteClient, err := rmn_remote.NewRMNRemote(common.HexToAddress(_rmnRemoteAddr), cl)
	require.NoError(t, err)

	cfg, err := rmnRemoteClient.GetVersionedConfig(nil)
	require.NoError(t, err)

	js, _ := json.MarshalIndent(cfg, " ", " ")
	fmt.Println(string(js))

	onRampAddr := common.HexToAddress(_onRampAddr)
	merkleRoot := common.Hex2Bytes(_merkleRoot)

	b, err := rmnRemoteClient.GetVerifyPreimage(
		nil,
		common.HexToAddress(_offRampAddr),
		[]rmn_remote.InternalMerkleRoot{
			{
				SourceChainSelector: _sourceChainSel,
				OnRampAddress:       common.LeftPadBytes(onRampAddr.Bytes(), 32),
				MinSeqNr:            1,
				MaxSeqNr:            1,
				MerkleRoot:          [32]byte(merkleRoot),
			},
		},
		[]rmn_remote.IRMNRemoteSignature{
			{
				R: [32]byte(common.Hex2Bytes(_r)),
				S: [32]byte(common.Hex2Bytes(_s)),
			},
		},
		big.NewInt(0),
	)

	// fetch evm chain id
	evmChainID, err := cl.ChainID(context.Background())
	require.NoError(t, err)
	t.Logf("EVM ChainID: %v", evmChainID)

	localChainSel, err := rmnRemoteClient.GetLocalChainSelector(nil)
	require.NoError(t, err)
	t.Logf("LocalChainSelector: %v", localChainSel)

	require.NoError(t, err)
	t.Logf("VerifyPreimage: %v", b)

	offChainPreImage := []byte{150, 81, 148, 55, 131, 219, 248, 25, 53, 166, 14, 152, 242, 24, 169, 217, 181, 178, 136, 35, 251, 34, 40, 187, 217, 19, 32, 214, 50, 250, 207, 83, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 33, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 179, 86, 115, 62, 34, 146, 68, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 83, 112, 247, 140, 106, 242, 218, 156, 246, 100, 35, 130, 163, 167, 95, 157, 90, 236, 156, 193, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 173, 93, 87, 173, 155, 177, 125, 52, 222, 187, 136, 86, 106, 178, 245, 219, 135, 156, 196, 111, 0, 11, 190, 110, 68, 63, 238, 166, 57, 217, 135, 140, 172, 75, 78, 152, 245, 155, 12, 24, 105, 221, 96, 22, 58, 215, 111, 118, 251, 165, 203, 127, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 46, 230, 52, 149, 30, 247, 27, 70, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 28, 97, 254, 247, 163, 221, 21, 57, 67, 65, 156, 17, 1, 3, 19, 22, 237, 123, 122, 61, 117, 145, 60, 52, 203, 232, 98, 128, 51, 245, 146, 79, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 122, 32, 136, 161, 191, 201, 216, 28, 85, 54, 138, 225, 104, 194, 192, 37, 112, 203, 129, 79}
	t.Logf("offChainPreImage: %v", hex.EncodeToString(offChainPreImage))

	contractPreImage := []byte{150, 81, 148, 55, 131, 219, 248, 25, 53, 166, 14, 152, 242, 24, 169, 217, 181, 178, 136, 35, 251, 34, 40, 187, 217, 19, 32, 214, 50, 250, 207, 83, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 33, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 179, 86, 115, 62, 34, 146, 68, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 241, 125, 93, 205, 169, 207, 37, 136, 156, 236, 154, 229, 97, 11, 15, 185, 114, 95, 101, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 173, 93, 87, 173, 155, 177, 125, 52, 222, 187, 136, 86, 106, 178, 245, 219, 135, 156, 196, 111, 0, 11, 190, 110, 68, 63, 238, 166, 57, 217, 135, 140, 172, 75, 78, 152, 245, 155, 12, 24, 105, 221, 96, 22, 58, 215, 111, 118, 251, 165, 203, 127, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 46, 230, 52, 149, 30, 247, 27, 70, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 28, 97, 254, 247, 163, 221, 21, 57, 67, 65, 156, 17, 1, 3, 19, 22, 237, 123, 122, 61, 117, 145, 60, 52, 203, 232, 98, 128, 51, 245, 146, 79, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 122, 32, 136, 161, 191, 201, 216, 28, 85, 54, 138, 225, 104, 194, 192, 37, 112, 203, 129, 79}
	t.Log("contractPreImage: ", hex.EncodeToString(contractPreImage))
}

func Test_Stuff(t *testing.T) {
	tabi, err := ccip_encoding_utils.EncodingUtilsMetaData.GetAbi()
	require.NoError(t, err)

	// offChainPreImage
	reportStructA := "0x9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf5300000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000921000000000000000000000000000000000000000000000000b356733e229244f20000000000000000000000005370f78c6af2da9cf6642382a3a75f9d5aec9cc1000000000000000000000000ad5d57ad9bb17d34debb88566ab2f5db879cc46f000bbe6e443feea639d9878cac4b4e98f59b0c1869dd60163ad76f76fba5cb7f00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000002ee634951ef71b4600000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000011c61fef7a3dd153943419c1101031316ed7b7a3d75913c34cbe8628033f5924f00000000000000000000000000000000000000000000000000000000000000200000000000000000000000007a2088a1bfc9d81c55368ae168c2c02570cb814f"
	reportDataA := hexutil.MustDecode(reportStructA)
	unpackedA, err := tabi.Methods["exposeRmnReport"].Inputs.Unpack(reportDataA)
	require.NoError(t, err)

	// contractPreImage
	reportStructB := "0x9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf5300000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000921000000000000000000000000000000000000000000000000b356733e229244f20000000000000000000000000cf17d5dcda9cf25889cec9ae5610b0fb9725f65000000000000000000000000ad5d57ad9bb17d34debb88566ab2f5db879cc46f000bbe6e443feea639d9878cac4b4e98f59b0c1869dd60163ad76f76fba5cb7f00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000002ee634951ef71b4600000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000011c61fef7a3dd153943419c1101031316ed7b7a3d75913c34cbe8628033f5924f00000000000000000000000000000000000000000000000000000000000000200000000000000000000000007a2088a1bfc9d81c55368ae168c2c02570cb814f"
	reportDataB := hexutil.MustDecode(reportStructB)
	unpackedB, err := tabi.Methods["exposeRmnReport"].Inputs.Unpack(reportDataB)
	require.NoError(t, err)

	assert.Equal(t, unpackedA, unpackedB)
	t.Log(unpackedA, unpackedB)
}
