package mercurytransmitter

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/grpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/observation"
)

const (
	// Mercury server error codes
	DuplicateReport  = 2
	commitInterval   = time.Millisecond * 25
	commitBufferSize = 1000
)

var (
	promTransmitSuccessCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "mercurytransmitter",
		Name:      "transmit_success_count",
		Help:      "Number of successful transmissions (duplicates are counted as success)",
	},
		[]string{"donID", "serverURL"},
	)
	promTransmitDuplicateCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "mercurytransmitter",
		Name:      "transmit_duplicate_count",
		Help:      "Number of transmissions where the server told us it was a duplicate",
	},
		[]string{"donID", "serverURL"},
	)
	promTransmitConnectionErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "mercurytransmitter",
		Name:      "transmit_connection_error_count",
		Help:      "Number of errored transmissions that failed due to problem with the connection",
	},
		[]string{"donID", "serverURL"},
	)
)

type Transmission struct {
	ServerURL    string
	ConfigDigest types.ConfigDigest
	SeqNr        uint64
	Report       ocr3types.ReportWithInfo[llotypes.ReportInfo]
	Sigs         []types.AttributedOnchainSignature
}

// Hash takes sha256 hash of all fields
func (t Transmission) Hash() [32]byte {
	h := sha256.New()
	h.Write([]byte(t.ServerURL))
	h.Write(t.ConfigDigest[:])
	if err := binary.Write(h, binary.BigEndian, t.SeqNr); err != nil {
		// This should never happen
		panic(err)
	}
	h.Write(t.Report.Report)
	h.Write([]byte(t.Report.Info.LifeCycleStage))
	if err := binary.Write(h, binary.BigEndian, t.Report.Info.ReportFormat); err != nil {
		// This should never happen
		panic(err)
	}
	for _, sig := range t.Sigs {
		h.Write(sig.Signature)
		if err := binary.Write(h, binary.BigEndian, sig.Signer); err != nil {
			// This should never happen
			panic(err)
		}
	}
	var result [32]byte
	h.Sum(result[:0])
	return result
}

type Transmitter interface {
	llotypes.Transmitter
	services.Service
}

var _ Transmitter = (*transmitter)(nil)

type Config interface {
	Protocol() config.MercuryTransmitterProtocol
	ReaperMaxAge() commonconfig.Duration
	TransmitConcurrency() uint32
	TransmitQueueMaxSize() uint32
	TransmitTimeout() commonconfig.Duration
}

type transmitter struct {
	services.StateMachine
	lggr           logger.SugaredLogger
	verboseLogging bool
	cfg            Config

	orm     ORM
	servers map[string]*server

	donID       uint32
	fromAccount string

	stopCh services.StopChan
	wg     *sync.WaitGroup

	commitCh chan *Transmission
}

type Opts struct {
	Lggr                 logger.Logger
	VerboseLogging       bool
	Cfg                  Config
	Clients              map[string]grpc.Client
	FromAccount          string
	DonID                uint32
	ORM                  ORM
	CapabilitiesRegistry coretypes.CapabilitiesRegistry
}

func New(opts Opts) Transmitter {
	return newTransmitter(opts)
}

func newTransmitter(opts Opts) *transmitter {
	sugared := logger.Sugared(opts.Lggr).Named("LLOMercuryTransmitter")
	servers := make(map[string]*server, len(opts.Clients))
	for serverURL, client := range opts.Clients {
		sLggr := sugared.Named(fmt.Sprintf("%q", serverURL)).With("serverURL", serverURL)
		servers[serverURL] = newServer(sLggr, opts.VerboseLogging, opts.Cfg, client, opts.ORM, serverURL)
	}
	return &transmitter{
		services.StateMachine{},
		sugared.Named("LLOMercuryTransmitter"),
		opts.VerboseLogging,
		opts.Cfg,
		opts.ORM,
		servers,
		opts.DonID,
		opts.FromAccount,
		make(services.StopChan),
		&sync.WaitGroup{},
		make(chan *Transmission, 1000*len(servers)),
	}
}

func (mt *transmitter) Start(ctx context.Context) (err error) {
	return mt.StartOnce("LLOMercuryTransmitter", func() error {
		if mt.verboseLogging {
			mt.lggr.Debugw("Loading transmit requests from database")
		}

		g, startCtx := errgroup.WithContext(ctx)
		// Number of goroutines spawned per server will be
		// TransmitConcurrency+2 (1 for persistence manager, 1 for client)
		//
		// This could potentially be reduced by implementing transmit batching,
		// see: https://smartcontract-it.atlassian.net/browse/MERC-6635
		for _, s := range mt.servers {
			// concurrent start of all servers
			g.Go(func() error {
				// Load DB transmissions and populate server transmit queue
				transmissions, err := s.pm.Load(startCtx)
				if err != nil {
					return err
				}
				s.q.Init(transmissions)

				// Start all associated services
				//
				// client, queue etc should be started before spawning server loops
				//
				// pm must be stopped last to give it a chance to clean up the
				// remaining transmissions
				startClosers := []services.StartClose{s.pm, s.c, s.q}
				if err := (&services.MultiStart{}).Start(startCtx, startClosers...); err != nil {
					return err
				}

				// Spawn transmission loop threads
				s.spawnTransmitLoops(mt.stopCh, mt.wg, mt.donID, int(mt.cfg.TransmitConcurrency()))
				return nil
			})
		}

		mt.spawnCommitLoops()
		return g.Wait()
	})
}

func (mt *transmitter) Close() error {
	return mt.StopOnce("LLOMercuryTransmitter", func() error {
		// Drain all the queues first
		var qs []io.Closer
		for _, s := range mt.servers {
			qs = append(qs, s.q)
		}
		if err := services.CloseAll(qs...); err != nil {
			return err
		}

		close(mt.stopCh)
		mt.wg.Wait()

		// Close all the persistence managers
		// Close all the clients
		var closers []io.Closer
		for _, s := range mt.servers {
			closers = append(closers, s.pm)
			closers = append(closers, s.c)
		}
		return services.CloseAll(closers...)
	})
}

func (mt *transmitter) Name() string { return mt.lggr.Name() }

func (mt *transmitter) HealthReport() map[string]error {
	report := map[string]error{mt.Name(): mt.Healthy()}
	for _, s := range mt.servers {
		services.CopyHealth(report, s.HealthReport())
	}
	return report
}

// Transmit enqueues the report for transmission to the Mercury servers
func (mt *transmitter) Transmit(
	ctx context.Context,
	digest types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[llotypes.ReportInfo],
	sigs []types.AttributedOnchainSignature,
) (err error) {
	ok := mt.IfStarted(func() {
		for serverURL := range mt.servers {
			t := &Transmission{
				ServerURL:    serverURL,
				ConfigDigest: digest,
				SeqNr:        seqNr,
				Report:       report,
				Sigs:         sigs,
			}
			select {
			case mt.commitCh <- t:
			case <-ctx.Done():
				err = fmt.Errorf("failed to add transmission to commit channel: %w", ctx.Err())
			}
		}
		observation.GetCache(digest).SetLastTransmissionSeqNr(seqNr)
	})

	if !ok {
		return errors.New("transmitter is not started")
	}

	return err
}

func (mt *transmitter) transmit(ctx context.Context, transmissions []*Transmission) error {
	// On shutdown appears that libocr can pass us a pre-canceled context;
	// don't even bother trying to insert/transmit in this case
	if ctx.Err() != nil {
		return fmt.Errorf("cannot transmit; context already canceled: %w", ctx.Err())
	}

	// NOTE: This insert on its own can leave orphaned records in the case of
	// shutdown, because:
	// 1. Transmitter is shut down after oracle
	// 2. OCR may pass a pre-canceled context or a context that is canceled mid-transmit
	// 3. Insert can succeed even if the context is canceled, but return error
	//
	// Usually the number of orphaned records will be very small, and they
	// would be transmitted/cleaned up on the next boot anyway.
	//
	// However, there are two ways to avoid this:
	// 1. Use a transaction to rollback the insert on error
	// 2. Allow the insert anyway (it will be transmitted on next boot) and be
	// sure that the persistence manager issues a final cleanup that truncates
	// the table to exactly maxSize records. Since persistenceManager is shut
	// down AFTER the Oracle closes, this should always catch the straggler
	// records.
	//
	// Since this is a hot path, the performance impact of holding a
	// transaction open is too high, hence we choose option 2.
	//
	// In very rare cases if the final delete fails for some reason, we could
	// end up with slightly more than maxSize records persisted to the DB on
	// application exit.
	//
	// Must insert BEFORE pushing to queue since the queue will handle deletion
	// on queue overflow.
	if err := mt.orm.Insert(ctx, transmissions); err != nil {
		return err
	}

	for i := range transmissions {
		t := transmissions[i]
		if mt.verboseLogging {
			mt.lggr.Debugw("Transmit report",
				"digest", t.ConfigDigest.Hex(), "seqNr", t.SeqNr, "reportFormat", t.Report.Info.ReportFormat,
				"reportLifeCycleStage", t.Report.Info.LifeCycleStage,
				"transmissionHash", fmt.Sprintf("%x", t.Hash()))
		}

		// OK to do this synchronously since pushing to queue is just a mutex
		// lock and array append and ought to be extremely fast
		s := mt.servers[t.ServerURL]
		if ok := s.q.Push(t); !ok {
			s.transmitQueuePushErrorCount.Inc()
			// This shouldn't be possible since transmitter is always shut down
			// after oracle
			return errors.New("transmit queue is closed")
		}
	}

	return nil
}

// FromAccount returns the stringified (hex) CSA public key
func (mt *transmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return ocrtypes.Account(mt.fromAccount), nil
}

func (mt *transmitter) spawnCommitLoops() {
	for x := 0; x < len(mt.servers); x++ {
		mt.wg.Add(1)

		go func() {
			defer mt.wg.Done()

			var err error
			ctx, cancel := mt.stopCh.NewCtx()
			defer cancel()

			buff := cap(mt.commitCh) / 10
			transmissions := make([]*Transmission, 0, buff)
			ticker := time.NewTicker(commitInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					if len(transmissions) >= buff {
						closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second)
						defer closeCancel()
						if err = mt.transmit(closeCtx, transmissions); err != nil {
							mt.lggr.Error("Error transmitting records when stopping", "error", err)
						}
					}
					return

				case <-ticker.C:
					if len(transmissions) > 0 {
						err = mt.transmit(ctx, transmissions)
						transmissions = make([]*Transmission, 0, buff)
					}

				case t := <-mt.commitCh:
					transmissions = append(transmissions, t)
					if len(transmissions) >= buff {
						err = mt.transmit(ctx, transmissions)
						transmissions = make([]*Transmission, 0, buff)
					}
				}

				if err != nil {
					mt.lggr.Error("Error transmitting records", "error", err)
				}
			}
		}()
	}
}
