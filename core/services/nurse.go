package services

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/pprof/profile"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Nurse struct {
	services.StateMachine

	cfg Config
	log logger.Logger

	checks   map[string]CheckFunc
	checksMu sync.RWMutex

	chGather chan gatherRequest
	chStop   chan struct{}
	wgDone   sync.WaitGroup
}

type Config interface {
	BlockProfileRate() int
	CPUProfileRate() int
	GatherDuration() commonconfig.Duration
	GatherTraceDuration() commonconfig.Duration
	GoroutineThreshold() int
	MaxProfileSize() utils.FileSize
	MemProfileRate() int
	MemThreshold() utils.FileSize
	MutexProfileFraction() int
	PollInterval() commonconfig.Duration
	ProfileRoot() string
}

type CheckFunc func() (unwell bool, meta Meta)

type gatherRequest struct {
	reason string
	meta   Meta
}

type Meta map[string]interface{}

const (
	cpuProfName   = "cpu"
	traceProfName = "trace"
)

func NewNurse(cfg Config, log logger.Logger) *Nurse {
	return &Nurse{
		cfg:      cfg,
		log:      log.Named("Nurse"),
		checks:   make(map[string]CheckFunc),
		chGather: make(chan gatherRequest, 1),
		chStop:   make(chan struct{}),
	}
}

func (n *Nurse) Start() error {
	return n.StartOnce("Nurse", func() error {
		// This must be set *once*, and it must occur as early as possible
		if n.cfg.MemProfileRate() != runtime.MemProfileRate {
			runtime.MemProfileRate = n.cfg.BlockProfileRate()
		}

		n.log.Debugf("Starting nurse with config %+v", n.cfg)
		runtime.SetCPUProfileRate(n.cfg.CPUProfileRate())
		runtime.SetBlockProfileRate(n.cfg.BlockProfileRate())
		runtime.SetMutexProfileFraction(n.cfg.MutexProfileFraction())

		err := utils.EnsureDirAndMaxPerms(n.cfg.ProfileRoot(), 0744)
		if err != nil {
			return err
		}

		n.AddCheck("mem", n.checkMem)
		n.AddCheck("goroutines", n.checkGoroutines)

		n.wgDone.Add(1)
		// Checker
		go func() {
			defer n.wgDone.Done()
			for {
				select {
				case <-n.chStop:
					return
				case <-time.After(n.cfg.PollInterval().Duration()):
				}

				func() {
					n.checksMu.RLock()
					defer n.checksMu.RUnlock()
					for reason, checkFunc := range n.checks {
						if unwell, meta := checkFunc(); unwell {
							n.GatherVitals(reason, meta)
							break
						}
					}
				}()
			}
		}()

		n.wgDone.Add(1)
		// Responder
		go func() {
			defer n.wgDone.Done()
			for {
				select {
				case <-n.chStop:
					return
				case req := <-n.chGather:
					n.gatherVitals(req.reason, req.meta)
				}
			}
		}()

		return nil
	})
}

func (n *Nurse) Close() error {
	return n.StopOnce("Nurse", func() error {
		n.log.Debug("Nurse closing...")
		defer n.log.Debug("Nurse closed")
		close(n.chStop)
		n.wgDone.Wait()
		return nil
	})
}

func (n *Nurse) AddCheck(reason string, checkFunc CheckFunc) {
	n.checksMu.Lock()
	defer n.checksMu.Unlock()
	n.checks[reason] = checkFunc
}

func (n *Nurse) GatherVitals(reason string, meta Meta) {
	select {
	case <-n.chStop:
	case n.chGather <- gatherRequest{reason, meta}:
	default:
	}
}

func (n *Nurse) checkMem() (bool, Meta) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	unwell := memStats.Alloc >= uint64(n.cfg.MemThreshold())
	if !unwell {
		return false, nil
	}
	return true, Meta{
		"mem_alloc": utils.FileSize(memStats.Alloc),
		"threshold": n.cfg.MemThreshold(),
	}
}

func (n *Nurse) checkGoroutines() (bool, Meta) {
	num := runtime.NumGoroutine()
	unwell := num >= n.cfg.GoroutineThreshold()
	if !unwell {
		return false, nil
	}
	return true, Meta{
		"num_goroutine": num,
		"threshold":     n.cfg.GoroutineThreshold(),
	}
}

func (n *Nurse) gatherVitals(reason string, meta Meta) {
	loggerFields := (logger.Fields{"reason": reason}).Merge(logger.Fields(meta))

	n.log.Debugw("Nurse is gathering vitals", loggerFields.Slice()...)

	size, err := n.totalProfileBytes()
	if err != nil {
		n.log.Errorw("could not fetch total profile bytes", loggerFields.With("err", err).Slice()...)
		return
	} else if size >= uint64(n.cfg.MaxProfileSize()) {
		n.log.Warnw("cannot write pprof profile, total profile size exceeds configured PPROF_MAX_PROFILE_SIZE",
			loggerFields.With("total", size, "max", n.cfg.MaxProfileSize()).Slice()...,
		)
		return
	}

	now := time.Now()

	err = n.appendLog(now, reason, meta)
	if err != nil {
		n.log.Warnw("cannot write pprof profile", loggerFields.With("err", err).Slice()...)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go n.gatherCPU(now, &wg)
	wg.Add(1)
	go n.gatherTrace(now, &wg)
	wg.Add(1)
	go n.gather("allocs", now, &wg)
	wg.Add(1)
	go n.gather("block", now, &wg)
	wg.Add(1)
	go n.gather("goroutine", now, &wg)

	// pprof docs state memory profile is not
	// created if the MemProfileRate is zero
	if runtime.MemProfileRate != 0 {
		wg.Add(1)
		go n.gather("heap", now, &wg)
	} else {
		n.log.Info("skipping heap collection because runtime.MemProfileRate = 0")
	}

	wg.Add(1)
	go n.gather("mutex", now, &wg)
	wg.Add(1)
	go n.gather("threadcreate", now, &wg)

	ch := make(chan struct{})
	n.wgDone.Add(1)
	go func() {
		defer n.wgDone.Done()
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-n.chStop:
	case <-ch:
	}
}

func (n *Nurse) appendLog(now time.Time, reason string, meta Meta) error {
	filename := filepath.Join(n.cfg.ProfileRoot(), "nurse.log")

	n.log.Debugf("creating nurse log %s", filename)
	file, err := os.Create(filename)

	if err != nil {
		return err
	}
	wc := utils.NewDeferableWriteCloser(file)
	defer wc.Close()

	if _, err = wc.Write([]byte(fmt.Sprintf("==== %v\n", now))); err != nil {
		return err
	}
	if _, err = wc.Write([]byte(fmt.Sprintf("reason: %v\n", reason))); err != nil {
		return err
	}
	ks := make([]string, len(meta))
	var i int
	for k := range meta {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	for _, k := range ks {
		if _, err = wc.Write([]byte(fmt.Sprintf("- %v: %v\n", k, meta[k]))); err != nil {
			return err
		}
	}
	_, err = wc.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return wc.Close()
}

func (n *Nurse) gatherCPU(now time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	n.log.Debugf("gather cpu %d ...", now.UnixMicro())
	defer n.log.Debugf("gather cpu %d done", now.UnixMicro())
	wc, err := n.createFile(now, cpuProfName, false)
	if err != nil {
		n.log.Errorw("could not write cpu profile", "err", err)
		return
	}
	defer wc.Close()

	err = pprof.StartCPUProfile(wc)
	if err != nil {
		n.log.Errorw("could not start cpu profile", "err", err)
		return
	}

	select {
	case <-n.chStop:
		n.log.Debug("gather cpu received stop")

	case <-time.After(n.cfg.GatherDuration().Duration()):
		n.log.Debugf("gather cpu duration elapsed %s. stoping profiling.", n.cfg.GatherDuration().Duration().String())
	}

	pprof.StopCPUProfile()

	err = wc.Close()
	if err != nil {
		n.log.Errorw("could not close cpu profile", "err", err)
		return
	}
}

func (n *Nurse) gatherTrace(now time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	n.log.Debugf("gather trace %d ...", now.UnixMicro())
	defer n.log.Debugf("gather trace %d done", now.UnixMicro())
	wc, err := n.createFile(now, traceProfName, true)
	if err != nil {
		n.log.Errorw("could not write trace profile", "err", err)
		return
	}
	defer wc.Close()

	err = trace.Start(wc)
	if err != nil {
		n.log.Errorw("could not start trace profile", "err", err)
		return
	}

	select {
	case <-n.chStop:
	case <-time.After(n.cfg.GatherTraceDuration().Duration()):
	}

	trace.Stop()

	err = wc.Close()
	if err != nil {
		n.log.Errorw("could not close trace profile", "err", err)
		return
	}
}

func (n *Nurse) gather(typ string, now time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	n.log.Debugf("gather %s %d ...", typ, now.UnixMicro())
	n.log.Debugf("gather %s %d done", typ, now.UnixMicro())

	p := pprof.Lookup(typ)
	if p == nil {
		n.log.Errorf("Invariant violation: pprof type '%v' does not exist", typ)
		return
	}

	p0, err := collectProfile(p)
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not collect %v profile", typ), "err", err)
		return
	}

	t := time.NewTimer(n.cfg.GatherDuration().Duration())
	defer t.Stop()

	select {
	case <-n.chStop:
		return
	case <-t.C:
	}

	p1, err := collectProfile(p)
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not collect %v profile", typ), "err", err)
		return
	}
	ts := p1.TimeNanos
	dur := p1.TimeNanos - p0.TimeNanos

	p0.Scale(-1)

	p1, err = profile.Merge([]*profile.Profile{p0, p1})
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not compute delta for %v profile", typ), "err", err)
		return
	}

	p1.TimeNanos = ts // set since we don't know what profile.Merge set for TimeNanos.
	p1.DurationNanos = dur

	wc, err := n.createFile(now, typ, false)
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not write %v profile", typ), "err", err)
		return
	}
	defer wc.Close()

	err = p1.Write(wc)
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not write %v profile", typ), "err", err)
		return
	}
	err = wc.Close()
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not close file for %v profile", typ), "err", err)
		return
	}
}

func collectProfile(p *pprof.Profile) (*profile.Profile, error) {
	var buf bytes.Buffer
	if err := p.WriteTo(&buf, 0); err != nil {
		return nil, err
	}
	ts := time.Now().UnixNano()
	p0, err := profile.Parse(&buf)
	if err != nil {
		return nil, err
	}
	p0.TimeNanos = ts
	return p0, nil
}

func (n *Nurse) createFile(now time.Time, typ string, shouldGzip bool) (*utils.DeferableWriteCloser, error) {
	filename := fmt.Sprintf("%v.%v.pprof", now.UnixMicro(), typ)
	if shouldGzip {
		filename += ".gz"
	}
	fullpath := filepath.Join(n.cfg.ProfileRoot(), filename)
	n.log.Debugf("creating file %s", fullpath)

	file, err := os.Create(fullpath)
	if err != nil {
		return nil, err
	}
	if shouldGzip {
		gw := gzip.NewWriter(file)
		return utils.NewDeferableWriteCloser(gw), nil
	}

	return utils.NewDeferableWriteCloser(file), nil
}

func (n *Nurse) totalProfileBytes() (uint64, error) {
	profiles, err := n.listProfiles()
	if err != nil {
		return 0, err
	}
	var size uint64
	for _, p := range profiles {
		size += uint64(p.Size())
	}
	return size, nil
}

func (n *Nurse) listProfiles() ([]fs.FileInfo, error) {
	out := make([]fs.FileInfo, 0)
	entries, err := os.ReadDir(n.cfg.ProfileRoot())

	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() ||
			(filepath.Ext(entry.Name()) != ".pprof" &&
				entry.Name() != "nurse.log" &&
				!strings.HasSuffix(entry.Name(), ".pprof.gz")) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		out = append(out, info)
	}
	return out, nil
}
