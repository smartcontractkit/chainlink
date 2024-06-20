package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestSnoozedRoots(t *testing.T) {
	c := NewCommitRootsCache(logger.TestLogger(t), 1*time.Minute, 1*time.Minute)

	k1 := [32]byte{1}
	k2 := [32]byte{2}

	// return false for non existing element
	snoozed := c.IsSkipped(k1)
	assert.False(t, snoozed)

	// after an element is marked as executed it should be snoozed
	c.MarkAsExecuted(k1)
	snoozed = c.IsSkipped(k1)
	assert.True(t, snoozed)

	// after snoozing an element it should be snoozed
	c.Snooze(k2)
	snoozed = c.IsSkipped(k2)
	assert.True(t, snoozed)
}

func TestEvictingElements(t *testing.T) {
	c := newCommitRootsCache(logger.TestLogger(t), 1*time.Hour, 1*time.Millisecond, 1*time.Millisecond, 1*time.Millisecond)

	k1 := [32]byte{1}
	c.Snooze(k1)

	time.Sleep(10 * time.Millisecond)

	assert.False(t, c.IsSkipped(k1))
}

func Test_UnexecutedRoots(t *testing.T) {
	type rootWithTs struct {
		root [32]byte
		ts   time.Time
	}

	r1 := [32]byte{1}
	r2 := [32]byte{2}
	r3 := [32]byte{3}

	t1 := time.Now().Add(-4 * time.Hour)
	t2 := time.Now().Add(-3 * time.Hour)
	t3 := time.Now().Add(-2 * time.Hour)

	tests := []struct {
		name                    string
		roots                   []rootWithTs
		executedRoots           [][32]byte
		permissionLessThreshold time.Duration
		expectedTimestamp       time.Time
	}{
		{
			name:                    "empty",
			roots:                   []rootWithTs{},
			permissionLessThreshold: 1 * time.Hour,
		},
		{
			name: "returns first root when all are not executed",
			roots: []rootWithTs{
				{r1, t1},
				{r2, t2},
				{r3, t3},
			},
			permissionLessThreshold: 10 * time.Hour,
			expectedTimestamp:       t1,
		},
		{
			name: "returns first root when tail of queue is executed",
			roots: []rootWithTs{
				{r1, t1},
				{r2, t2},
				{r3, t3},
			},
			executedRoots:           [][32]byte{r2, r3},
			permissionLessThreshold: 10 * time.Hour,
			expectedTimestamp:       t1,
		},
		{
			name: "returns first not executed root",
			roots: []rootWithTs{
				{r1, t1},
				{r2, t2},
				{r3, t3},
			},
			executedRoots:           [][32]byte{r1, r2},
			permissionLessThreshold: 10 * time.Hour,
			expectedTimestamp:       t3,
		},
		{
			name: "returns r2 timestamp when r1 and r3 are executed",
			roots: []rootWithTs{
				{r1, t1},
				{r2, t2},
				{r3, t3},
			},
			executedRoots:           [][32]byte{r1, r3},
			permissionLessThreshold: 10 * time.Hour,
			expectedTimestamp:       t2,
		},
		{
			name: "returns oldest root even when all are executed",
			roots: []rootWithTs{
				{r1, t1},
				{r2, t2},
				{r3, t3},
			},
			executedRoots:           [][32]byte{r1, r2, r3},
			permissionLessThreshold: 10 * time.Hour,
			expectedTimestamp:       t3,
		},
		{
			name: "returns permissionLessThreshold when all roots ale older that threshold",
			roots: []rootWithTs{
				{r1, t1},
				{r2, t2},
				{r3, t3},
			},
			permissionLessThreshold: 1 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCommitRootsCache(logger.TestLogger(t), tt.permissionLessThreshold, 1*time.Hour, 1*time.Millisecond, 1*time.Millisecond)

			for _, r := range tt.roots {
				c.AppendUnexecutedRoot(r.root, r.ts)
			}

			for _, r := range tt.executedRoots {
				c.MarkAsExecuted(r)
			}

			commitTs := c.OldestRootTimestamp()
			if tt.expectedTimestamp.IsZero() {
				assert.True(t, commitTs.Before(time.Now().Add(-tt.permissionLessThreshold)))
			} else {
				assert.Equal(t, tt.expectedTimestamp.Add(-time.Second), commitTs)
			}
		})
	}
}

func Test_UnexecutedRootsScenario(t *testing.T) {
	permissionLessThreshold := 10 * time.Hour
	c := newCommitRootsCache(logger.TestLogger(t), permissionLessThreshold, 1*time.Hour, 1*time.Millisecond, 1*time.Millisecond)

	k1 := [32]byte{1}
	k2 := [32]byte{2}
	k3 := [32]byte{3}
	k4 := [32]byte{4}

	t1 := time.Now().Add(-4 * time.Hour)
	t2 := time.Now().Add(-3 * time.Hour)
	t3 := time.Now().Add(-2 * time.Hour)
	t4 := time.Now().Add(-1 * time.Hour)

	// First check should return permissionLessThreshold window
	commitTs := c.OldestRootTimestamp()
	assert.True(t, commitTs.Before(time.Now().Add(-permissionLessThreshold)))

	c.AppendUnexecutedRoot(k1, t1)
	c.AppendUnexecutedRoot(k2, t2)
	c.AppendUnexecutedRoot(k3, t3)

	// After loading roots it should return the first one
	commitTs = c.OldestRootTimestamp()
	assert.Equal(t, t1.Add(-time.Second), commitTs)

	// Marking root in the middle as executed shouldn't change the commitTs
	c.MarkAsExecuted(k2)
	commitTs = c.OldestRootTimestamp()
	assert.Equal(t, t1.Add(-time.Second), commitTs)

	// Marking k1 as executed when k2 is already executed should return timestamp of k3
	c.MarkAsExecuted(k1)
	commitTs = c.OldestRootTimestamp()
	assert.Equal(t, t3.Add(-time.Second), commitTs)

	// Marking all as executed should return timestamp of the latest
	c.MarkAsExecuted(k3)
	commitTs = c.OldestRootTimestamp()
	assert.Equal(t, t3.Add(-time.Second), commitTs)

	// Adding k4 should return timestamp of k4
	c.AppendUnexecutedRoot(k4, t4)
	commitTs = c.OldestRootTimestamp()
	assert.Equal(t, t4.Add(-time.Second), commitTs)

	c.MarkAsExecuted(k4)
	commitTs = c.OldestRootTimestamp()
	assert.Equal(t, t4.Add(-time.Second), commitTs)

	// Appending already executed roots should be ignored
	c.AppendUnexecutedRoot(k1, t1)
	c.AppendUnexecutedRoot(k2, t2)
	commitTs = c.OldestRootTimestamp()
	assert.Equal(t, t4.Add(-time.Second), commitTs)
}

func Test_UnexecutedRootsStaleQueue(t *testing.T) {
	permissionLessThreshold := 5 * time.Hour
	c := newCommitRootsCache(logger.TestLogger(t), permissionLessThreshold, 1*time.Hour, 1*time.Millisecond, 1*time.Millisecond)

	k1 := [32]byte{1}
	k2 := [32]byte{2}
	k3 := [32]byte{3}

	t1 := time.Now().Add(-4 * time.Hour)
	t2 := time.Now().Add(-3 * time.Hour)
	t3 := time.Now().Add(-2 * time.Hour)

	c.AppendUnexecutedRoot(k1, t1)
	c.AppendUnexecutedRoot(k2, t2)
	c.AppendUnexecutedRoot(k3, t3)

	// First check should return permissionLessThreshold window
	commitTs := c.OldestRootTimestamp()
	assert.Equal(t, t1.Add(-time.Second), commitTs)

	// Reducing permissionLessExecutionThreshold works as speeding the clock
	c.messageVisibilityInterval = 1 * time.Hour

	commitTs = c.OldestRootTimestamp()
	assert.True(t, commitTs.Before(time.Now().Add(-1*time.Hour)))
	assert.True(t, commitTs.After(t1))
	assert.True(t, commitTs.After(t2))
	assert.True(t, commitTs.After(t3))
}
