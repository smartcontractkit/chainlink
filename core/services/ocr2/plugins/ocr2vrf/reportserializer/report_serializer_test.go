package reportserializer_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/types"

	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/reportserializer"
)

func Test_Serialize_Deserialize(t *testing.T) {
	altbn128Suite := &altbn_128.PairingSuite{}
	reportSerializer := reportserializer.NewReportSerializer(altbn128Suite.G1())

	unserializedReport := types.AbstractReport{
		JuelsPerFeeCoin:   big.NewInt(10),
		RecentBlockHeight: 100,
		RecentBlockHash:   common.HexToHash("0x002"),
		Outputs: []types.AbstractVRFOutput{{
			BlockHeight:       10,
			ConfirmationDelay: 20,
			Callbacks: []types.AbstractCostedCallbackRequest{{
				RequestID:      1,
				NumWords:       2,
				Requester:      common.HexToAddress("0x03"),
				Arguments:      []byte{4},
				SubscriptionID: big.NewInt(5),
				GasAllowance:   big.NewInt(6),
				Price:          big.NewInt(7),
				GasPrice:       big.NewInt(0),
				WeiPerUnitLink: big.NewInt(0),
			}},
		}},
	}
	r, err := reportSerializer.SerializeReport(unserializedReport)
	require.NoError(t, err)
	require.Equal(t, uint(len(r)), reportSerializer.ReportLength(unserializedReport))
	// TODO: Add deserialization after this point to verify.
}
