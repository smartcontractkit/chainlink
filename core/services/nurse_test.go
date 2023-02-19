package services

import (
	"io/fs"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type mockConfig struct {
	//mock.Mock
	t                    *testing.T
	root                 string
	pollInterval         *models.Duration
	gatherDuration       *models.Duration
	traceDuration        *models.Duration
	profileSize          utils.FileSize
	cpuProfileRate       int
	memProfileRate       int
	blockProfileRate     int
	mutexProfileFraction int
	memThreshold         utils.FileSize
	goroutineThreshold   int
}

var (
	testInterval = time.Duration(100 * time.Millisecond)
	testRate     = 1
	testSize     = 1024 * 1024
)

func newMockConfig(t *testing.T) *mockConfig {
	return &mockConfig{
		root:                 t.TempDir(),
		pollInterval:         models.MustNewDuration(testInterval),
		gatherDuration:       models.MustNewDuration(testInterval),
		traceDuration:        models.MustNewDuration(testInterval),
		profileSize:          utils.FileSize(testSize),
		memProfileRate:       testRate,
		blockProfileRate:     testRate,
		mutexProfileFraction: testRate,
		memThreshold:         utils.FileSize(testSize),
		goroutineThreshold:   testRate,
		t:                    t,
	}
}

func (c mockConfig) AutoPprofProfileRoot() string {

	return c.root
}

func (c mockConfig) AutoPprofPollInterval() models.Duration {

	return *c.pollInterval
}

func (c mockConfig) AutoPprofGatherDuration() models.Duration {

	return *c.gatherDuration
}

func (c mockConfig) AutoPprofGatherTraceDuration() models.Duration {

	return *c.traceDuration
}

func (c mockConfig) AutoPprofMaxProfileSize() utils.FileSize {

	return c.profileSize
}

func (c mockConfig) AutoPprofCPUProfileRate() int {

	return c.cpuProfileRate
}

func (c mockConfig) AutoPprofMemProfileRate() int {

	return c.memProfileRate
}

func (c mockConfig) AutoPprofBlockProfileRate() int {

	return c.blockProfileRate
}

func (c mockConfig) AutoPprofMutexProfileFraction() int {

	return c.mutexProfileFraction
}

func (c mockConfig) AutoPprofMemThreshold() utils.FileSize {

	return c.memThreshold
}

func (c mockConfig) AutoPprofGoroutineThreshold() int {

	return c.goroutineThreshold
}

func TestNurse_appendLog(t *testing.T) {

	l := logger.TestLogger(t)
	nrse := NewNurse(newMockConfig(t), l)
	require.NoError(t, nrse.appendLog(time.Now(), "test", Meta{}))
	wc, err := nrse.createFile(time.Now(), "test", false)
	require.NoError(t, err)
	require.NoError(t, wc.Close())

	wc, err = nrse.createFile(time.Now(), "testgz", false)
	require.NoError(t, err)
	require.NoError(t, wc.Close())

	var wg sync.WaitGroup

	wg.Add(1)
	nrse.gatherCPU(time.Now(), &wg)
	wg.Wait()

	profiles, err := nrse.listProfiles()
	require.NoError(t, err)
	assertProfileExists(t, profiles, cpuProf)
	n, err := nrse.totalProfileBytes()
	require.NoError(t, err)
	require.Greater(t, n, 0)
	requireRemoveProfiles(t, profiles)

	wg.Add(1)
	nrse.gatherTrace(time.Now(), &wg)
	wg.Wait()

	profiles, err = nrse.listProfiles()
	require.NoError(t, err)
	assertProfileExists(t, profiles, traceProf)
	requireRemoveProfiles(t, profiles)

	wg.Add(1)
	nrse.gather(cpuProf, time.Now(), &wg)
	wg.Wait()
	assertProfileExists(t, profiles, cpuProf)
	requireRemoveProfiles(t, profiles)

}

func assertProfileExists(t *testing.T, profiles []fs.FileInfo, typ string) {
	var names []string
	for _, p := range profiles {
		names = append(names, p.Name())
		if strings.Contains(p.Name(), typ) {
			return
		}
	}
	assert.Failf(t, "profile doesn't exist", "require profile '%s' does not exist %+v", typ, names)
}

func requireRemoveProfiles(t *testing.T, profiles []fs.FileInfo) {
	for _, p := range profiles {
		require.NoError(t, os.Remove(p.Name()))
	}

}
