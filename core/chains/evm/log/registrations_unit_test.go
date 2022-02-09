package log

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ Listener = testListener{}

type testListener struct {
	jobID int32
}

func (tl testListener) JobID() int32        { return tl.jobID }
func (tl testListener) HandleLog(Broadcast) { panic("not implemented") }

func newTestListener(t *testing.T, jobID int32) testListener {
	return testListener{jobID}
}

func newTestRegistrations(t *testing.T) *registrations {
	return newRegistrations(logger.TestLogger(t), *testutils.FixtureChainID)
}

func TestRegistrationsUnit_InvariantViolations(t *testing.T) {
	l := newTestListener(t, 1)
	r := newTestRegistrations(t)

	contractAddr := testutils.NewAddress()
	opts := ListenerOpts{Contract: contractAddr, MinIncomingConfirmations: 1}
	sub := &subscriber{l, opts}

	r.addSubscriber(sub)

	// Different subscriber same job ID different contract address is ok
	subB := &subscriber{l, ListenerOpts{Contract: testutils.NewAddress(), MinIncomingConfirmations: 1}}
	r.addSubscriber(subB)

	// Different subscriber same jobID/contract address is not ok
	assert.Panics(t, func() {
		opts := ListenerOpts{Contract: contractAddr, MinIncomingConfirmations: 1}
		subError := &subscriber{l, opts}

		r.addSubscriber(subError)
	})

	l2 := newTestListener(t, 2)
	sub2 := &subscriber{l2, opts}

	// Different subscriber different job ID same contract address is ok
	r.addSubscriber(sub2)

	// Adding same subscriber twice is not ok
	assert.Panics(t, func() {
		r.addSubscriber(sub2)
	}, "expected adding same subscription twice to panic")

	r.removeSubscriber(sub)

	// Removing subscriber twice also panics
	assert.Panics(t, func() {
		r.removeSubscriber(sub)
	}, "expected removing a subscriber twice to panic")

	// Now we can add it again
	r.addSubscriber(sub)
}

func newTopic() Topic {
	return Topic(utils.NewHash())
}

func TestRegistrationsUnit_addSubscriber(t *testing.T) {
	contractAddr := testutils.NewAddress()
	r := newTestRegistrations(t)

	l := newTestListener(t, 1)
	topic1 := utils.NewHash()
	topicValueFilters1 := [][]Topic{[]Topic{newTopic(), newTopic()}, []Topic{newTopic()}, []Topic{}}
	topic2 := utils.NewHash()
	topicValueFilters2 := [][]Topic{[]Topic{newTopic()}}
	topic3 := utils.NewHash()
	topicValueFilters3 := [][]Topic{}
	logsWithTopics := make(map[common.Hash][][]Topic)
	logsWithTopics[topic1] = topicValueFilters1
	logsWithTopics[topic2] = topicValueFilters2
	logsWithTopics[topic3] = topicValueFilters3
	opts := ListenerOpts{Contract: contractAddr, LogsWithTopics: logsWithTopics, MinIncomingConfirmations: 1}
	sub := &subscriber{l, opts}
	r.addSubscriber(sub)

	// same contract, same topics
	l2 := newTestListener(t, 2)
	opts2 := opts
	sub2 := &subscriber{l2, opts2}
	r.addSubscriber(sub2)

	// same contract, different topics
	l3 := newTestListener(t, 3)
	topic4 := utils.NewHash()
	topicValueFilters4 := [][]Topic{[]Topic{newTopic()}}
	logsWithTopics3 := make(map[common.Hash][][]Topic)
	logsWithTopics3[topic4] = topicValueFilters4
	opts3 := opts
	opts3.LogsWithTopics = logsWithTopics3
	sub3 := &subscriber{l3, opts3}
	r.addSubscriber(sub3)

	assert.Equal(t, 1, int(r.highestNumConfirmations))

	// same contract, same topics, different MinIncomingConfirmations
	l4 := newTestListener(t, 4)
	opts4 := opts3
	opts4.MinIncomingConfirmations = 42
	sub4 := &subscriber{l4, opts4}
	r.addSubscriber(sub4)

	assert.Equal(t, 42, int(r.highestNumConfirmations))

	assert.Len(t, r.registeredSubs, 4)
	assert.Contains(t, r.registeredSubs, sub)
	assert.Contains(t, r.registeredSubs, sub2)
	assert.Contains(t, r.registeredSubs, sub3)
	assert.Contains(t, r.registeredSubs, sub4)

	assert.Len(t, r.handlersByConfs, 2)
	require.Contains(t, r.handlersByConfs, uint32(1))
	require.Contains(t, r.handlersByConfs, uint32(42))

	// contractAddress => logTopic => Listener
	handlers1 := r.handlersByConfs[1].lookupSubs
	assert.Len(t, handlers1, 1)
	assert.Contains(t, handlers1, contractAddr)
	h1 := handlers1[contractAddr]
	// 4 topics on this contract addr
	assert.Len(t, h1, 4)
	assert.Contains(t, h1, topic1)
	assert.Contains(t, h1, topic2)
	assert.Contains(t, h1, topic3)
	assert.Contains(t, h1, topic4)
	// topics map to their subscribers
	assert.Len(t, h1[topic1], 2) // listeners 1 and 2
	assert.Contains(t, h1[topic1], sub)
	assert.Contains(t, h1[topic1], sub2)
	assert.Len(t, h1[topic2], 2) // listeners 1 and 2
	assert.Contains(t, h1[topic2], sub)
	assert.Contains(t, h1[topic2], sub2)
	assert.Len(t, h1[topic3], 2) // listeners 1 and 2
	assert.Contains(t, h1[topic3], sub)
	assert.Contains(t, h1[topic3], sub2)
	assert.Len(t, h1[topic4], 1) // listener 3
	assert.Contains(t, h1[topic4], sub3)

	handlers42 := r.handlersByConfs[42].lookupSubs
	assert.Len(t, handlers42, 1)
	assert.Contains(t, handlers1, contractAddr)
	h42 := handlers42[contractAddr]
	// 1 topic on this contract addr
	assert.Len(t, h42, 1)
	assert.Contains(t, h1, topic4)
	// topic maps to its subscriber
	assert.Len(t, h42[topic4], 1) // listener 4
	assert.Contains(t, h42[topic4], sub4)
}
