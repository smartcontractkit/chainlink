package ocr2vrf_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/types"

	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf"
)

func Test_Serialize_Deserialize(t *testing.T) {
	altbn128Suite := &altbn_128.PairingSuite{}
	reportSerializer := ocr2vrf.ReportSerializer{
		G: altbn128Suite.G1(),
	}

	unserializedReport := types.AbstractReport{
		JulesPerFeeCoin:   big.NewInt(10),
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
				SubscriptionID: 5,
				GasAllowance:   big.NewInt(6),
				Price:          big.NewInt(7),
			}},
		}},
	}
	r, err := reportSerializer.SerializeReport(unserializedReport)
	require.NoError(t, err)

	report, err := reportSerializer.DeserializeReport(r)
	require.NoError(t, err)

	require.Equal(t, unserializedReport, types.AbstractReport{
		JulesPerFeeCoin:   report.JulesPerFeeCoin,
		RecentBlockHeight: report.RecentBlockHeight,
		RecentBlockHash:   common.Hash(report.RecentBlockHash),
		Outputs: []types.AbstractVRFOutput{{
			BlockHeight:       report.Outputs[0].BlockHeight,
			ConfirmationDelay: uint32(report.Outputs[0].ConfirmationDelay.Int64()),
			Callbacks: []types.AbstractCostedCallbackRequest{{
				RequestID:      report.Outputs[0].Callbacks[0].Callback.RequestID.Uint64(),
				NumWords:       report.Outputs[0].Callbacks[0].Callback.NumWords,
				Requester:      report.Outputs[0].Callbacks[0].Callback.Requester,
				Arguments:      report.Outputs[0].Callbacks[0].Callback.Arguments,
				SubscriptionID: report.Outputs[0].Callbacks[0].Callback.SubID,
				GasAllowance:   report.Outputs[0].Callbacks[0].Callback.GasAllowance,
				Price:          report.Outputs[0].Callbacks[0].Price,
			}},
		}},
	})
}

func Test_Serialize_Length(t *testing.T) {
	altbn128Suite := &altbn_128.PairingSuite{}
	reportSerializer := ocr2vrf.ReportSerializer{
		G: altbn128Suite.G1(),
	}

	unserializedReport := types.AbstractReport{
		JulesPerFeeCoin:   big.NewInt(10),
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
				SubscriptionID: 5,
				GasAllowance:   big.NewInt(6),
				Price:          big.NewInt(7),
			}},
		}},
	}
	r, err := reportSerializer.SerializeReport(unserializedReport)
	require.NoError(t, err)

	require.Equal(t, uint(len(r)), reportSerializer.ReportLength(unserializedReport))
}
