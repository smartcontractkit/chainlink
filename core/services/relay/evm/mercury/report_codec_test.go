package mercury

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReportCodec(t *testing.T) {
	r := ReportCodec{}

	t.Run("BuildReport constructs a report", func(t *testing.T) {
		paos := []median.ParsedAttributedObservation{
			{
	Timestamp       uint32
	Value           *big.Int
	JuelsPerFeeCoin *big.Int
	Observer        commontypes.OracleID
			},
			{
			},
			{},
			{},
		}
		r.BuildReport(pao)
	})

	t.Run("MedianFromReport gets the median", func(t *testing.T) {
		b, err := hexutil.Decode(sampleReportHex)
		require.NoError(t, err)
		sampleReport := types.Report(b)
		median, err := r.MedianFromReport(sampleReport)
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(69000), median)
	})

	t.Run("Report has len less than or equal to MaxReportLength", func(t *testing.T) {
	})
}
