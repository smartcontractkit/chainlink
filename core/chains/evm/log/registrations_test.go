package log

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
)

var _ Listener = testListener{}

type testListener struct {
	jobID int32
}

func (tl testListener) JobID() int32                         { return tl.jobID }
func (tl testListener) HandleLog(context.Context, Broadcast) { panic("not implemented") }

func newTestListener(t *testing.T, jobID int32) testListener {
	return testListener{jobID}
}

func newTestRegistrations(t *testing.T) *registrations {
	return newRegistrations(logger.Test(t), *testutils.FixtureChainID)
}

func newTopic() Topic {
	return Topic(utils.NewHash())
}

func TestUnit_Registrations_InvariantViolations(t *testing.T) {
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
		subError := &subscriber{l, ListenerOpts{Contract: contractAddr, MinIncomingConfirmations: 1}}

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

func TestUnit_Registrations_addSubscriber_removeSubscriber(t *testing.T) {
	contractAddr := testutils.NewAddress()
	r := newTestRegistrations(t)

	l := newTestListener(t, 1)
	topic1 := utils.NewHash()
	topicValueFilters1 := [][]Topic{{newTopic(), newTopic()}, {newTopic()}, {}}
	topic2 := utils.NewHash()
	topicValueFilters2 := [][]Topic{{newTopic()}}
	topic3 := utils.NewHash()
	topicValueFilters3 := [][]Topic{}
	logsWithTopics := make(map[common.Hash][][]Topic)
	logsWithTopics[topic1] = topicValueFilters1
	logsWithTopics[topic2] = topicValueFilters2
	logsWithTopics[topic3] = topicValueFilters3
	opts := ListenerOpts{Contract: contractAddr, LogsWithTopics: logsWithTopics, MinIncomingConfirmations: 1}
	sub := &subscriber{l, opts}

	// same contract, same topics
	l2 := newTestListener(t, 2)
	opts2 := opts
	sub2 := &subscriber{l2, opts2}

	// same contract, different topics
	l3 := newTestListener(t, 3)
	topic4 := utils.NewHash()
	topicValueFilters4 := [][]Topic{{newTopic()}}
	logsWithTopics3 := make(map[common.Hash][][]Topic)
	logsWithTopics3[topic4] = topicValueFilters4
	opts3 := opts
	opts3.LogsWithTopics = logsWithTopics3
	sub3 := &subscriber{l3, opts3}

	// same contract, same topics, greater MinIncomingConfirmations
	l4 := newTestListener(t, 4)
	opts4 := opts3
	opts4.MinIncomingConfirmations = 42
	sub4 := &subscriber{l4, opts4}

	// same contract, same topics, midrange MinIncomingConfirmations
	l5 := newTestListener(t, 5)
	opts5 := opts3
	opts5.MinIncomingConfirmations = 21
	sub5 := &subscriber{l5, opts5}

	t.Run("addSubscriber", func(t *testing.T) {
		needsResub := r.addSubscriber(sub)
		assert.True(t, needsResub)

		// same contract, same topics
		needsResub = r.addSubscriber(sub2)
		assert.False(t, needsResub)

		// same contract, different topics
		needsResub = r.addSubscriber(sub3)
		assert.True(t, needsResub)

		assert.Equal(t, 1, int(r.highestNumConfirmations))

		// same contract, same topics, different MinIncomingConfirmations
		needsResub = r.addSubscriber(sub4)
		// resub required because confirmations went higher
		assert.True(t, needsResub)
		assert.Equal(t, 42, int(r.highestNumConfirmations))

		// same contract, same topics, midrange MinIncomingConfirmations
		needsResub = r.addSubscriber(sub5)
		// resub NOT required because confirmations is lower than the highest
		assert.False(t, needsResub)
		assert.Equal(t, 42, int(r.highestNumConfirmations))

		assert.Len(t, r.registeredSubs, 5)
		assert.Contains(t, r.registeredSubs, sub)
		assert.Contains(t, r.registeredSubs, sub2)
		assert.Contains(t, r.registeredSubs, sub3)
		assert.Contains(t, r.registeredSubs, sub4)
		assert.Contains(t, r.registeredSubs, sub5)

		assert.Len(t, r.handlersByConfs, 3)
		require.Contains(t, r.handlersByConfs, uint32(1))
		require.Contains(t, r.handlersByConfs, uint32(21))
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

		handlers21 := r.handlersByConfs[21].lookupSubs
		assert.Len(t, handlers21, 1)
		assert.Contains(t, handlers1, contractAddr)
		h21 := handlers21[contractAddr]
		// 1 topic on this contract addr
		assert.Len(t, h21, 1)
		assert.Contains(t, h1, topic4)
		// topic maps to its subscriber
		assert.Len(t, h21[topic4], 1) // listener 5
		assert.Contains(t, h21[topic4], sub5)
	})

	t.Run("removeSubscriber", func(t *testing.T) {
		needsResub := r.removeSubscriber(sub)
		// No resub necessary: sub2 also needs all these topics
		assert.False(t, needsResub)

		assert.Len(t, r.registeredSubs, 4)
		assert.NotContains(t, r.registeredSubs, sub)
		assert.Contains(t, r.registeredSubs, sub2)
		assert.Contains(t, r.registeredSubs, sub3)
		assert.Contains(t, r.registeredSubs, sub4)
		assert.Contains(t, r.registeredSubs, sub5)

		needsResub = r.removeSubscriber(sub2)
		// sub2 has topics in it that other subs don't cover
		assert.True(t, needsResub)
		assert.Len(t, r.registeredSubs, 3)
		assert.NotContains(t, r.registeredSubs, sub2)
		assert.Contains(t, r.registeredSubs, sub3)
		assert.Contains(t, r.registeredSubs, sub4)
		assert.Contains(t, r.registeredSubs, sub4)

		needsResub = r.removeSubscriber(sub3)
		// sub5 and sub4 cover everything that sub3 does already, resub not necessary
		assert.False(t, needsResub)
		assert.Len(t, r.registeredSubs, 2)
		assert.NotContains(t, r.registeredSubs, sub3)
		assert.Contains(t, r.registeredSubs, sub4)
		assert.Contains(t, r.registeredSubs, sub4)

		needsResub = r.removeSubscriber(sub4)
		// sub5 covers everything that sub4 does already, resub not necessary
		assert.False(t, needsResub)
		assert.Len(t, r.registeredSubs, 1)
		assert.NotContains(t, r.registeredSubs, sub4)
		assert.Contains(t, r.registeredSubs, sub5)

		needsResub = r.removeSubscriber(sub5)
		// Nothing left, need to refresh subscriptions
		assert.True(t, needsResub)
		assert.Len(t, r.registeredSubs, 0)
	})
}
