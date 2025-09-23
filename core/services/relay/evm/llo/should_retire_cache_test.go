package llo

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
)

func Test_ShouldRetireCache(t *testing.T) {
	lp := &mockLogPoller{make([]logpoller.Log, 0), 0, nil, query.LimitAndSort{}}
	addr := common.Address{1}
	donID := uint32(1)
	donIDHash := DonIDToBytes32(donID)
	retiredConfigDigest := ocr2types.ConfigDigest{1, 2, 3, 4}

	log := logpoller.Log{Address: addr, Topics: pq.ByteaArray{PromoteStagingConfig[:], donIDHash[:], retiredConfigDigest[:]}, EventSig: PromoteStagingConfig, BlockNumber: 100, Data: makePromoteStagingConfigData(t, false)}
	lp.logs = append(lp.logs, log)

	src := newShouldRetireCache(logger.Test(t), lp, addr, donID)

	servicetest.Run(t, src)

	testutils.RequireEventually(t, func() bool {
		shouldRetire, err2 := src.ShouldRetire(retiredConfigDigest)
		require.NoError(t, err2)
		return shouldRetire
	})

	shouldRetire, err := src.ShouldRetire(ocr2types.ConfigDigest{9})
	require.NoError(t, err)
	assert.False(t, shouldRetire, "Should not retire")
}

func makePromoteStagingConfigData(t *testing.T, isGreenProduction bool) []byte {
	event := configuratorABI.Events["PromoteStagingConfig"]
	data, err := event.Inputs.NonIndexed().Pack(isGreenProduction)
	require.NoError(t, err)
	return data
}
