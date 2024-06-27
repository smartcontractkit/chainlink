package encoding

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

func TestReportEncoder_EncodeExtract(t *testing.T) {
	encoder := reportEncoder{
		packer: NewAbiPacker(),
	}

	tests := []struct {
		name               string
		results            []ocr2keepers.CheckResult
		reportSize         int
		expectedFastGasWei int64
		expectedLinkNative int64
		expectedErr        error
	}{
		{
			"happy flow single",
			[]ocr2keepers.CheckResult{
				newResult(1, 1, core.GenUpkeepID(ocr2keepers.LogTrigger, "123"), 1, 1),
			},
			736,
			1,
			1,
			nil,
		},
		{
			"happy flow multiple",
			[]ocr2keepers.CheckResult{
				newResult(1, 1, core.GenUpkeepID(ocr2keepers.LogTrigger, "10"), 1, 1),
				newResult(1, 1, core.GenUpkeepID(ocr2keepers.ConditionTrigger, "20"), 1, 1),
				newResult(1, 1, core.GenUpkeepID(ocr2keepers.ConditionTrigger, "30"), 1, 1),
			},
			1312,
			3,
			3,
			nil,
		},
		{
			"happy flow highest block number first",
			[]ocr2keepers.CheckResult{
				newResult(1, 1, core.GenUpkeepID(ocr2keepers.ConditionTrigger, "30"), 1, 1),
				newResult(1, 1, core.GenUpkeepID(ocr2keepers.ConditionTrigger, "20"), 1, 1),
				newResult(1, 1, core.GenUpkeepID(ocr2keepers.LogTrigger, "10"), 1, 1),
			},
			1312,
			1000,
			2000,
			nil,
		},
		{
			"empty results",
			[]ocr2keepers.CheckResult{},
			0,
			0,
			0,
			ErrEmptyResults,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := encoder.Encode(tc.results...)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Nil(t, err)
			assert.Len(t, b, tc.reportSize)

			results, err := encoder.Extract(b)
			assert.Nil(t, err)
			assert.Len(t, results, len(tc.results))

			for i, r := range results {
				assert.Equal(t, r.UpkeepID, tc.results[i].UpkeepID)
				assert.Equal(t, r.WorkID, tc.results[i].WorkID)
				assert.Equal(t, r.Trigger, tc.results[i].Trigger)
			}
		})
	}
}

func newResult(block int64, checkBlock ocr2keepers.BlockNumber, id ocr2keepers.UpkeepIdentifier, fastGasWei, linkNative int64) ocr2keepers.CheckResult {
	tp := core.GetUpkeepType(id)

	trig := ocr2keepers.Trigger{
		BlockNumber: ocr2keepers.BlockNumber(block),
		BlockHash:   [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
	}

	if tp == ocr2keepers.LogTrigger {
		trig.LogTriggerExtension = &ocr2keepers.LogTriggerExtension{
			Index:     1,
			TxHash:    common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
			BlockHash: common.HexToHash("0xaaaaaaaa90123456789012345678901234567890123456789012345678901234"),
		}
	}

	payload, _ := core.NewUpkeepPayload(
		id.BigInt(),
		trig,
		[]byte{},
	)

	return ocr2keepers.CheckResult{
		UpkeepID:     id,
		Trigger:      payload.Trigger,
		WorkID:       payload.WorkID,
		Eligible:     true,
		GasAllocated: 100,
		PerformData:  []byte("data0"),
		FastGasWei:   big.NewInt(fastGasWei),
		LinkNative:   big.NewInt(linkNative),
	}
}
