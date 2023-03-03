package pyroscope

import (
	"bytes"
	"github.com/pyroscope-io/client/internal/alignedticker"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/pyroscope-io/client/internal/flameql"
	"github.com/pyroscope-io/client/upstream"
)

type Session struct {
	// configuration, doesn't change
	upstream               upstream.Upstream
	sampleRate             uint32
	profileTypes           []ProfileType
	uploadRate             time.Duration
	disableGCRuns          bool
	DisableAutomaticResets bool

	logger    Logger
	stopOnce  sync.Once
	stopCh    chan struct{}
	flushCh   chan *flush
	trieMutex sync.Mutex

	// these things do change:
	cpuBuf         *bytes.Buffer
	memBuf         *bytes.Buffer
	memPrevBytes   []byte
	goroutinesBuf  *bytes.Buffer
	mutexBuf       *bytes.Buffer
	mutexPrevBytes []byte
	blockBuf       *bytes.Buffer
	blockPrevBytes []byte

	lastGCGeneration uint32
	appName          string
	startTime        time.Time
}

type SessionConfig struct {
	Upstream               upstream.Upstream
	Logger                 Logger
	AppName                string
	Tags                   map[string]string
	ProfilingTypes         []ProfileType
	DisableGCRuns          bool
	DisableAutomaticResets bool
	SampleRate             uint32
	UploadRate             time.Duration
}

type flush struct {
	wg   sync.WaitGroup
	wait bool
}

func NewSession(c SessionConfig) (*Session, error) {
	appName, err := mergeTagsWithAppName(c.AppName, c.Tags)
	if err != nil {
		return nil, err
	}

	ps := &Session{
		upstream:               c.Upstream,
		appName:                appName,
		profileTypes:           c.ProfilingTypes,
		disableGCRuns:          c.DisableGCRuns,
		DisableAutomaticResets: c.DisableAutomaticResets,
		sampleRate:             c.SampleRate,
		uploadRate:             c.UploadRate,
		stopCh:                 make(chan struct{}),
		flushCh:                make(chan *flush),
		logger:                 c.Logger,
		cpuBuf:                 &bytes.Buffer{},
		memBuf:                 &bytes.Buffer{},
		goroutinesBuf:          &bytes.Buffer{},
		mutexBuf:               &bytes.Buffer{},
		blockBuf:               &bytes.Buffer{},
	}

	return ps, nil
}

// mergeTagsWithAppName validates user input and merges explicitly specified
// tags with tags from app name.
//
// App name may be in the full form including tags (app.name{foo=bar,baz=qux}).
// Returned application name is always short, any tags that were included are
// moved to tags map. When merged with explicitly provided tags (config/CLI),
// last take precedence.
//
// App name may be an empty string. Tags must not contain reserved keys,
// the map is modified in place.
func mergeTagsWithAppName(appName string, tags map[string]string) (string, error) {
	k, err := flameql.ParseKey(appName)
	if err != nil {
		return "", err
	}
	for tagKey, tagValue := range tags {
		if flameql.IsTagKeyReserved(tagKey) {
			continue
		}
		if err = flameql.ValidateTagKey(tagKey); err != nil {
			return "", err
		}
		k.Add(tagKey, tagValue)
	}
	return k.Normalized(), nil
}

// revive:disable-next-line:cognitive-complexity complexity is fine
func (ps *Session) takeSnapshots() {
	var automaticResetTicker <-chan time.Time
	if ps.DisableAutomaticResets {
		automaticResetTicker = make(chan time.Time)
	} else {
		t := alignedticker.NewAlignedTicker(ps.uploadRate)
		automaticResetTicker = t.C
		defer t.Stop()
	}
	for {
		select {
		case endTime := <-automaticResetTicker:
			ps.reset(ps.startTime, endTime)
		case f := <-ps.flushCh:
			ps.reset(ps.startTime, ps.truncatedTime())
			ps.upstream.Flush()
			f.wg.Done()
			break
		case <-ps.stopCh:
			return
		}
	}
}

func copyBuf(b []byte) []byte {
	r := make([]byte, len(b))
	copy(r, b)
	return r
}

func (ps *Session) start() error {
	t := ps.truncatedTime()
	ps.reset(t, t)

	go ps.takeSnapshots()
	return nil
}

func (ps *Session) isCPUEnabled() bool {
	for _, t := range ps.profileTypes {
		if t == ProfileCPU {
			return true
		}
	}
	return false
}

func (ps *Session) isMemEnabled() bool {
	for _, t := range ps.profileTypes {
		if t == ProfileInuseObjects || t == ProfileAllocObjects || t == ProfileInuseSpace || t == ProfileAllocSpace {
			return true
		}
	}
	return false
}

func (ps *Session) isBlockEnabled() bool {
	for _, t := range ps.profileTypes {
		if t == ProfileBlockCount || t == ProfileBlockDuration {
			return true
		}
	}
	return false
}

func (ps *Session) isMutexEnabled() bool {
	for _, t := range ps.profileTypes {
		if t == ProfileMutexCount || t == ProfileMutexDuration {
			return true
		}
	}
	return false
}

func (ps *Session) isGoroutinesEnabled() bool {
	for _, t := range ps.profileTypes {
		if t == ProfileGoroutines {
			return true
		}
	}
	return false
}

func (ps *Session) reset(startTime, endTime time.Time) {

	ps.logger.Debugf("profiling session reset %s", startTime.String())

	// first reset should not result in an upload
	if !ps.startTime.IsZero() {
		ps.uploadData(startTime, endTime)
	} else {
		if ps.isCPUEnabled() {
			pprof.StartCPUProfile(ps.cpuBuf)
		}
	}

	ps.startTime = endTime
}

func (ps *Session) uploadData(startTime, endTime time.Time) {
	if ps.isCPUEnabled() {
		pprof.StopCPUProfile()
		defer func() {
			pprof.StartCPUProfile(ps.cpuBuf)
		}()
		ps.upstream.Upload(&upstream.UploadJob{
			Name:            ps.appName,
			StartTime:       startTime,
			EndTime:         endTime,
			SpyName:         "gospy",
			SampleRate:      100,
			Units:           "samples",
			AggregationType: "sum",
			Format:          upstream.FormatPprof,
			Profile:         copyBuf(ps.cpuBuf.Bytes()),
		})
		ps.cpuBuf.Reset()
	}

	if ps.isGoroutinesEnabled() {
		p := pprof.Lookup("goroutine")
		if p != nil {
			p.WriteTo(ps.goroutinesBuf, 0)
			ps.upstream.Upload(&upstream.UploadJob{
				Name:            ps.appName,
				StartTime:       startTime,
				EndTime:         endTime,
				SpyName:         "gospy",
				Units:           "goroutines",
				AggregationType: "average",
				Format:          upstream.FormatPprof,
				Profile:         copyBuf(ps.goroutinesBuf.Bytes()),
				SampleTypeConfig: map[string]*upstream.SampleType{
					"goroutine": {
						DisplayName: "goroutines",
						Units:       "goroutines",
						Aggregation: "average",
					},
				},
			})
			ps.goroutinesBuf.Reset()
		}
	}

	if ps.isBlockEnabled() {
		p := pprof.Lookup("block")
		if p != nil {
			p.WriteTo(ps.blockBuf, 0)
			curBlockBuf := copyBuf(ps.blockBuf.Bytes())
			ps.blockBuf.Reset()
			if ps.blockPrevBytes != nil {
				ps.upstream.Upload(&upstream.UploadJob{
					Name:        ps.appName,
					StartTime:   startTime,
					EndTime:     endTime,
					SpyName:     "gospy",
					Format:      upstream.FormatPprof,
					Profile:     curBlockBuf,
					PrevProfile: ps.blockPrevBytes,
					SampleTypeConfig: map[string]*upstream.SampleType{
						"contentions": {
							DisplayName: "block_count",
							Units:       "lock_samples",
							Cumulative:  true,
						},
						"delay": {
							DisplayName: "block_duration",
							Units:       "lock_nanoseconds",
							Cumulative:  true,
						},
					},
				})
			}
			ps.blockPrevBytes = curBlockBuf
		}
	}
	if ps.isMutexEnabled() {
		p := pprof.Lookup("mutex")
		if p != nil {
			p.WriteTo(ps.mutexBuf, 0)
			curMutexBuf := copyBuf(ps.mutexBuf.Bytes())
			ps.mutexBuf.Reset()
			if ps.mutexPrevBytes != nil {
				ps.upstream.Upload(&upstream.UploadJob{
					Name:        ps.appName,
					StartTime:   startTime,
					EndTime:     endTime,
					SpyName:     "gospy",
					Format:      upstream.FormatPprof,
					Profile:     curMutexBuf,
					PrevProfile: ps.mutexPrevBytes,
					SampleTypeConfig: map[string]*upstream.SampleType{
						"contentions": {
							DisplayName: "mutex_count",
							Units:       "lock_samples",
							Cumulative:  true,
						},
						"delay": {
							DisplayName: "mutex_duration",
							Units:       "lock_nanoseconds",
							Cumulative:  true,
						},
					},
				})
			}
			ps.mutexPrevBytes = curMutexBuf
		}
	}

	if ps.isMemEnabled() {
		currentGCGeneration := numGC()
		// sometimes GC doesn't run within 10 seconds
		//   in such cases we force a GC run
		//   users can disable it with disableGCRuns option
		if currentGCGeneration == ps.lastGCGeneration && !ps.disableGCRuns {
			runtime.GC()
			currentGCGeneration = numGC()
		}
		if currentGCGeneration != ps.lastGCGeneration {
			pprof.WriteHeapProfile(ps.memBuf)
			curMemBytes := copyBuf(ps.memBuf.Bytes())
			ps.memBuf.Reset()
			if ps.memPrevBytes != nil {
				ps.upstream.Upload(&upstream.UploadJob{
					Name:        ps.appName,
					StartTime:   startTime,
					EndTime:     endTime,
					SpyName:     "gospy",
					SampleRate:  100,
					Format:      upstream.FormatPprof,
					Profile:     curMemBytes,
					PrevProfile: ps.memPrevBytes,
				})
			}
			ps.memPrevBytes = curMemBytes
			ps.lastGCGeneration = currentGCGeneration
		}
	}
}

func (ps *Session) Stop() {
	ps.trieMutex.Lock()
	defer ps.trieMutex.Unlock()

	ps.stopOnce.Do(func() {
		// TODO: wait for stopCh consumer to finish!
		close(ps.stopCh)
		// before stopping, upload the tries
		ps.uploadLastBitOfData(time.Now())
	})
}

func (ps *Session) uploadLastBitOfData(now time.Time) {
	if ps.isCPUEnabled() {
		pprof.StopCPUProfile()
		ps.upstream.Upload(&upstream.UploadJob{
			Name:            ps.appName,
			StartTime:       ps.startTime,
			EndTime:         now,
			SpyName:         "gospy",
			SampleRate:      100,
			Units:           "samples",
			AggregationType: "sum",
			Format:          upstream.FormatPprof,
			Profile:         copyBuf(ps.cpuBuf.Bytes()),
		})
	}
}

func (ps *Session) flush(wait bool) {
	f := &flush{
		wg:   sync.WaitGroup{},
		wait: wait,
	}
	f.wg.Add(1)
	ps.flushCh <- f
	if wait {
		f.wg.Wait()
	}
}

func (ps *Session) truncatedTime() time.Time {
	return time.Now().Truncate(ps.uploadRate)
}

func numGC() uint32 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.NumGC
}
