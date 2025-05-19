package llo

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func Test_ShouldRetireCache(t *testing.T) {
	lggr, observedLogs := logger.TestObserved(t, zapcore.DebugLevel)
	lp := &mockLogPoller{make([]logpoller.Log, 0), 0, nil, query.LimitAndSort{}}
	addr := common.Address{1}
	donID := uint32(1)
	donIDHash := DonIDToBytes32(donID)
	retiredConfigDigest := ocr2types.ConfigDigest{1, 2, 3, 4}

	lp.logs = append(lp.logs, logpoller.Log{Address: addr, Topics: pq.ByteaArray{PromoteStagingConfig[:], donIDHash[:], retiredConfigDigest[:]}, EventSig: PromoteStagingConfig, BlockNumber: 100, Data: makePromoteStagingConfigData(t, false)})

	src := newShouldRetireCache(lggr, lp, addr, donID)

	servicetest.Run(t, src)

	testutils.WaitForLogMessage(t, observedLogs, "markRetired: Got retired config digest")

	shouldRetire, err := src.ShouldRetire(retiredConfigDigest)
	require.NoError(t, err)
	assert.True(t, shouldRetire, "Should retire")
	shouldRetire, err = src.ShouldRetire(ocr2types.ConfigDigest{9})
	require.NoError(t, err)
	assert.False(t, shouldRetire, "Should not retire")
}

func makePromoteStagingConfigData(t *testing.T, isGreenProduction bool) []byte {
	event := configuratorABI.Events["PromoteStagingConfig"]
	data, err := event.Inputs.NonIndexed().Pack(isGreenProduction)
	require.NoError(t, err)
	return data
}
