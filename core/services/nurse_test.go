package services

import (
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type mockConfig struct {
	t                    *testing.T
	root                 string
	pollInterval         *commonconfig.Duration
	gatherDuration       *commonconfig.Duration
	traceDuration        *commonconfig.Duration
	profileSize          utils.FileSize
	cpuProfileRate       int
	memProfileRate       int
	blockProfileRate     int
	mutexProfileFraction int
	memThreshold         utils.FileSize
	goroutineThreshold   int
}

var (
	testInterval = 50 * time.Millisecond
	testDuration = 20 * time.Millisecond
	testRate     = 100
	testSize     = 16 * 1024 * 1024
)

func newMockConfig(t *testing.T) *mockConfig {
	return &mockConfig{
		root:                 t.TempDir(),
		pollInterval:         commonconfig.MustNewDuration(testInterval),
		gatherDuration:       commonconfig.MustNewDuration(testDuration),
		traceDuration:        commonconfig.MustNewDuration(testDuration),
		profileSize:          utils.FileSize(testSize),
		memProfileRate:       runtime.MemProfileRate,
		blockProfileRate:     testRate,
		mutexProfileFraction: testRate,
		memThreshold:         utils.FileSize(testSize),
		goroutineThreshold:   testRate,
		t:                    t,
	}
}

func (c mockConfig) ProfileRoot() string {
	return c.root
}

func (c mockConfig) PollInterval() commonconfig.Duration {
	return *c.pollInterval
}

func (c mockConfig) GatherDuration() commonconfig.Duration {
	return *c.gatherDuration
}

func (c mockConfig) GatherTraceDuration() commonconfig.Duration {
	return *c.traceDuration
}

func (c mockConfig) MaxProfileSize() utils.FileSize {
	return c.profileSize
}

func (c mockConfig) CPUProfileRate() int {
	return c.cpuProfileRate
}

func (c mockConfig) MemProfileRate() int {
	return c.memProfileRate
}

func (c mockConfig) BlockProfileRate() int {
	return c.blockProfileRate
}

func (c mockConfig) MutexProfileFraction() int {
	return c.mutexProfileFraction
}

func (c mockConfig) MemThreshold() utils.FileSize {
	return c.memThreshold
}

func (c mockConfig) GoroutineThreshold() int {
	return c.goroutineThreshold
}

func TestNurse(t *testing.T) {
	l := logger.TestLogger(t)
	nrse := NewNurse(newMockConfig(t), l)
	nrse.AddCheck("test", func() (bool, Meta) { return true, Meta{} })

	require.NoError(t, nrse.Start(t.Context()))
	defer func() { require.NoError(t, nrse.Close()) }()

	require.NoError(t, nrse.appendLog(time.Now(), "test", Meta{}))

	wc, err := nrse.createFile(time.Now(), "test", false)
	require.NoError(t, err)
	n, err := wc.Write([]byte("junk"))
	require.NoError(t, err)
	require.Positive(t, n)
	require.NoError(t, wc.Close())

	wc, err = nrse.createFile(time.Now(), "testgz", false)
	require.NoError(t, err)
	require.NoError(t, wc.Close())

	// check both of the files exist. synchronous, check immediately
	assert.True(t, profileExists(t, nrse, "test"))
	assert.True(t, profileExists(t, nrse, "testgz"))

	testutils.AssertEventually(t, func() bool { return profileExists(t, nrse, cpuProfName) })
	testutils.AssertEventually(t, func() bool { return profileExists(t, nrse, traceProfName) })
	n2, err := nrse.totalProfileBytes()
	require.NoError(t, err)
	require.Positive(t, n2)
}

func profileExists(t *testing.T, nrse *Nurse, typ string) bool {
	profiles, err := nrse.listProfiles()
	require.NoError(t, err)
	for _, p := range profiles {
		if strings.Contains(p.Name(), typ) {
			return true
		}
	}
	return false
}
