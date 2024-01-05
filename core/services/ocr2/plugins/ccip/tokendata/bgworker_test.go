package tokendata_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
)

func TestBackgroundWorker(t *testing.T) {
	ctx := testutils.Context(t)

	const numTokens = 100
	const numWorkers = 20
	const numMessages = 40
	const maxReaderLatencyMS = 200
	const percentOfTokensWithoutTokenData = 10

	tokens := make([]common.Address, numTokens)
	readers := make(map[common.Address]*tokendata.MockReader, numTokens)
	tokenDataReaders := make(map[common.Address]tokendata.Reader, numTokens)
	tokenData := make(map[common.Address][]byte)
	delays := make(map[common.Address]time.Duration)

	for i := range tokens {
		tokens[i] = utils.RandomAddress()
		readers[tokens[i]] = tokendata.NewMockReader(t)
		if rand.Intn(100) >= percentOfTokensWithoutTokenData {
			tokenDataReaders[tokens[i]] = readers[tokens[i]]
			tokenData[tokens[i]] = []byte(fmt.Sprintf("...token %x data...", tokens[i]))
		}

		// specify a random latency for the reader implementation
		readerLatency := rand.Intn(maxReaderLatencyMS)
		delays[tokens[i]] = time.Duration(readerLatency) * time.Millisecond
	}
	w := tokendata.NewBackgroundWorker(ctx, tokenDataReaders, numWorkers, 5*time.Second, time.Hour)

	msgs := make([]internal.EVM2EVMOnRampCCIPSendRequestedWithMeta, numMessages)
	for i := range msgs {
		tk := tokens[i%len(tokens)]

		msgs[i] = internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
			EVM2EVMMessage: internal.EVM2EVMMessage{
				SequenceNumber: uint64(i + 1),
				TokenAmounts:   []internal.TokenAmount{{Token: tk}},
			},
		}

		reader := readers[tk]
		reader.On("ReadTokenData", mock.Anything, msgs[i], 0).Run(func(args mock.Arguments) {
			time.Sleep(delays[tk])
		}).Return(tokenData[tk], nil).Maybe()
	}

	w.AddJobsFromMsgs(ctx, msgs)
	// processing of the messages should have started at this point

	tStart := time.Now()
	for _, msg := range msgs {
		b, err := w.GetMsgTokenData(ctx, msg) // fetched from provider
		assert.NoError(t, err)
		assert.Equal(t, tokenData[msg.TokenAmounts[0].Token], b[0])
	}
	assert.True(t, time.Since(tStart) < 600*time.Millisecond)
	assert.True(t, time.Since(tStart) > 200*time.Millisecond)

	tStart = time.Now()
	for _, msg := range msgs {
		b, err := w.GetMsgTokenData(ctx, msg) // fetched from cache
		assert.NoError(t, err)
		assert.Equal(t, tokenData[msg.TokenAmounts[0].Token], b[0])
	}
	assert.True(t, time.Since(tStart) < 200*time.Millisecond)

	w.AddJobsFromMsgs(ctx, msgs) // same messages are added but they should already be in cache
	tStart = time.Now()
	for _, msg := range msgs {
		b, err := w.GetMsgTokenData(ctx, msg)
		assert.NoError(t, err)
		assert.Equal(t, tokenData[msg.TokenAmounts[0].Token], b[0])
	}
	assert.True(t, time.Since(tStart) < 200*time.Millisecond)
}

func TestBackgroundWorker_RetryOnErrors(t *testing.T) {
	ctx := testutils.Context(t)

	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()

	rdr1 := tokendata.NewMockReader(t)
	rdr2 := tokendata.NewMockReader(t)

	w := tokendata.NewBackgroundWorker(ctx, map[common.Address]tokendata.Reader{
		tk1: rdr1,
		tk2: rdr2,
	}, 10, 5*time.Second, time.Hour)

	msgs := []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		{EVM2EVMMessage: internal.EVM2EVMMessage{
			SequenceNumber: uint64(1),
			TokenAmounts:   []internal.TokenAmount{{Token: tk1}},
		}},
		{EVM2EVMMessage: internal.EVM2EVMMessage{
			SequenceNumber: uint64(2),
			TokenAmounts:   []internal.TokenAmount{{Token: tk2}},
		}},
	}

	rdr1.On("ReadTokenData", mock.Anything, msgs[0], 0).
		Return([]byte("some data"), nil).Once()

	// reader2 returns an error
	rdr2.On("ReadTokenData", mock.Anything, msgs[1], 0).
		Return(nil, fmt.Errorf("some err")).Once()

	w.AddJobsFromMsgs(ctx, msgs)
	// processing of the messages should have started at this point

	tokenData, err := w.GetMsgTokenData(ctx, msgs[0])
	assert.NoError(t, err)
	assert.Equal(t, []byte("some data"), tokenData[0])

	_, err = w.GetMsgTokenData(ctx, msgs[1])
	assert.Error(t, err)
	assert.Errorf(t, err, "some error")

	// we make the second reader to return data
	rdr2.On("ReadTokenData", mock.Anything, msgs[1], 0).
		Return([]byte("some other data"), nil).Once()

	// add the jobs again, at this point jobs that previously returned
	// an error are removed from the cache
	w.AddJobsFromMsgs(ctx, msgs)

	// since reader1 returned some data before, we expect to get the cached result
	tokenData, err = w.GetMsgTokenData(ctx, msgs[0])
	assert.NoError(t, err)
	assert.Equal(t, []byte("some data"), tokenData[0])

	// wait some time for msg2 to be re-processed and error overwritten
	time.Sleep(20 * time.Millisecond) // todo: improve the test

	// for reader2 that returned an error before we expect to get data now
	tokenData, err = w.GetMsgTokenData(ctx, msgs[1])
	assert.NoError(t, err)
	assert.Equal(t, []byte("some other data"), tokenData[0])
}

func TestBackgroundWorker_Timeout(t *testing.T) {
	ctx := testutils.Context(t)

	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()

	rdr1 := tokendata.NewMockReader(t)
	rdr2 := tokendata.NewMockReader(t)

	w := tokendata.NewBackgroundWorker(
		ctx, map[common.Address]tokendata.Reader{tk1: rdr1, tk2: rdr2}, 10, 5*time.Second, time.Hour)

	ctx, cf := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cf()

	_, err := w.GetMsgTokenData(ctx, internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		EVM2EVMMessage: internal.EVM2EVMMessage{SequenceNumber: 1}},
	)
	assert.Error(t, err)
}
