package services

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
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

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Nurse struct {
	utils.StartStopOnce

	cfg Config
	log logger.Logger

	checks   map[string]CheckFunc
	checksMu sync.RWMutex

	chGather chan gatherRequest
	chStop   chan struct{}
	wgDone   sync.WaitGroup
}

type Config interface {
	AutoPprofProfileRoot() string
	AutoPprofPollInterval() models.Duration
	AutoPprofGatherDuration() models.Duration
	AutoPprofGatherTraceDuration() models.Duration
	AutoPprofMaxProfileSize() utils.FileSize
	AutoPprofCPUProfileRate() int
	AutoPprofMemProfileRate() int
	AutoPprofBlockProfileRate() int
	AutoPprofMutexProfileFraction() int
	AutoPprofMemThreshold() utils.FileSize
	AutoPprofGoroutineThreshold() int
}

type CheckFunc func() (unwell bool, meta Meta)

type gatherRequest struct {
	reason string
	meta   Meta
}

type Meta map[string]interface{}

const profilePerms = 0666

func NewNurse(cfg Config, log logger.Logger) *Nurse {
	return &Nurse{
		cfg:      cfg,
		log:      log.Named("nurse"),
		checks:   make(map[string]CheckFunc),
		chGather: make(chan gatherRequest, 1),
		chStop:   make(chan struct{}),
	}
}

func (n *Nurse) Start() error {
	return n.StartOnce("nurse", func() error {
		// This must be set *once*, and it must occur as early as possible
		runtime.MemProfileRate = n.cfg.AutoPprofMemProfileRate()

		runtime.SetCPUProfileRate(n.cfg.AutoPprofCPUProfileRate())

		err := utils.EnsureDirAndMaxPerms(n.cfg.AutoPprofProfileRoot(), 0644)
		if err != nil {
			return err
		}

		n.AddCheck("mem", n.checkMem)
		n.AddCheck("goroutines", n.checkGoroutines)

		n.wgDone.Add(2)

		// Checker
		go func() {
			defer n.wgDone.Done()
			for {
				select {
				case <-n.chStop:
					return
				case <-time.After(n.cfg.AutoPprofPollInterval().Duration()):
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
	return n.StopOnce("nurse", func() error {
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
	unwell := memStats.Alloc >= uint64(n.cfg.AutoPprofMemThreshold())
	if !unwell {
		return false, nil
	}
	return true, Meta{
		"mem_alloc": utils.FileSize(memStats.Alloc),
		"threshold": n.cfg.AutoPprofMemThreshold(),
	}
}

func (n *Nurse) checkGoroutines() (bool, Meta) {
	num := runtime.NumGoroutine()
	unwell := num >= n.cfg.AutoPprofGoroutineThreshold()
	if !unwell {
		return false, nil
	}
	return true, Meta{
		"num_goroutine": num,
		"threshold":     n.cfg.AutoPprofGoroutineThreshold(),
	}
}

func (n *Nurse) gatherVitals(reason string, meta Meta) {
	loggerFields := (logger.Fields{"reason": reason}).Merge(logger.Fields(meta))

	n.log.Debugw("Nurse is gathering vitals", loggerFields.Slice()...)

	size, err := n.totalProfileBytes()
	if err != nil {
		n.log.Errorw("could not fetch total profile bytes", loggerFields.With("error", err).Slice()...)
		return
	} else if size >= uint64(n.cfg.AutoPprofMaxProfileSize()) {
		n.log.Warnw("cannot write pprof profile, total profile size exceeds configured PPROF_MAX_PROFILE_SIZE",
			loggerFields.With("total", size, "max", n.cfg.AutoPprofMaxProfileSize()).Slice()...,
		)
		return
	}

	runtime.SetBlockProfileRate(n.cfg.AutoPprofBlockProfileRate())
	defer runtime.SetBlockProfileRate(0)
	runtime.SetMutexProfileFraction(n.cfg.AutoPprofMutexProfileFraction())
	defer runtime.SetMutexProfileFraction(0)

	now := time.Now()

	var wg sync.WaitGroup
	wg.Add(9)

	err = n.appendLog(now, reason, meta)
	if err != nil {
		n.log.Warnw("cannot write pprof profile", loggerFields.With("error", err).Slice()...)
		return
	}
	go n.gatherCPU(now, &wg)
	go n.gatherTrace(now, &wg)
	go n.gather("allocs", now, &wg)
	go n.gather("block", now, &wg)
	go n.gather("goroutine", now, &wg)
	go n.gather("heap", now, &wg)
	go n.gather("mutex", now, &wg)
	go n.gather("threadcreate", now, &wg)

	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-n.chStop:
	case <-ch:
	}
}

func (n *Nurse) appendLog(now time.Time, reason string, meta Meta) error {
	filename := filepath.Join(n.cfg.AutoPprofProfileRoot(), "nurse.log")
	mode := os.O_APPEND | os.O_CREATE | os.O_WRONLY

	file, err := os.OpenFile(filename, mode, profilePerms)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write([]byte(fmt.Sprintf("==== %v\n", now))); err != nil {
		return err
	}
	if _, err = file.Write([]byte(fmt.Sprintf("reason: %v\n", reason))); err != nil {
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
		if _, err = file.Write([]byte(fmt.Sprintf("- %v: %v\n", k, meta[k]))); err != nil {
			return err
		}
	}
	_, err = file.Write([]byte("\n"))
	return err
}

func (n *Nurse) gatherCPU(now time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := n.openFile(now, "cpu", false)
	if err != nil {
		n.log.Errorw("could not write cpu profile", "error", err)
		return
	}
	defer file.Close()

	err = pprof.StartCPUProfile(file)
	if err != nil {
		n.log.Errorw("could not start cpu profile", "error", err)
		return
	}
	defer pprof.StopCPUProfile()

	select {
	case <-n.chStop:
	case <-time.After(n.cfg.AutoPprofGatherDuration().Duration()):
	}
}

func (n *Nurse) gatherTrace(now time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := n.openFile(now, "trace", true)
	if err != nil {
		n.log.Errorw("could not write trace profile", "error", err)
		return
	}
	defer file.Close()

	err = trace.Start(file)
	if err != nil {
		n.log.Errorw("could not start trace profile", "error", err)
		return
	}
	defer trace.Stop()

	select {
	case <-n.chStop:
	case <-time.After(n.cfg.AutoPprofGatherTraceDuration().Duration()):
	}
}

func (n *Nurse) gather(typ string, now time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	p := pprof.Lookup(typ)
	if p == nil {
		n.log.Errorf("Invariant violation: pprof type '%v' does not exist", typ)
		return
	}

	p0, err := collectProfile(p)
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not collect %v profile", typ), "error", err)
		return
	}

	t := time.NewTimer(n.cfg.AutoPprofGatherDuration().Duration())
	defer t.Stop()

	select {
	case <-n.chStop:
		return
	case <-t.C:
	}

	p1, err := collectProfile(p)
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not collect %v profile", typ), "error", err)
		return
	}
	ts := p1.TimeNanos
	dur := p1.TimeNanos - p0.TimeNanos

	p0.Scale(-1)

	p1, err = profile.Merge([]*profile.Profile{p0, p1})
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not compute delta for %v profile", typ), "error", err)
		return
	}

	p1.TimeNanos = ts // set since we don't know what profile.Merge set for TimeNanos.
	p1.DurationNanos = dur

	file, err := n.openFile(now, typ, false)
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not write %v profile", typ), "error", err)
		return
	}
	defer file.Close()

	err = p1.Write(file)
	if err != nil {
		n.log.Errorw(fmt.Sprintf("could not write %v profile", typ), "error", err)
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

func (n *Nurse) openFile(now time.Time, typ string, shouldGzip bool) (io.WriteCloser, error) {
	filename := fmt.Sprintf("%v.%v.pprof", now, typ)
	if shouldGzip {
		filename += ".gz"
	}
	fullpath := filepath.Join(n.cfg.AutoPprofProfileRoot(), filename)
	mode := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(fullpath, mode, profilePerms)
	if err != nil {
		return nil, err
	}
	if shouldGzip {
		return gzip.NewWriter(file), nil
	}
	return file, nil
}

func (n *Nurse) totalProfileBytes() (uint64, error) {
	entries, err := os.ReadDir(n.cfg.AutoPprofProfileRoot())
	if err != nil {
		return 0, err
	}
	var size uint64
	for _, entry := range entries {
		if entry.IsDir() ||
			(filepath.Ext(entry.Name()) != ".pprof" &&
				entry.Name() != "nurse.log" &&
				!strings.HasSuffix(entry.Name(), ".pprof.gz")) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return 0, err
		}
		size += uint64(info.Size())
	}
	return size, nil
}
