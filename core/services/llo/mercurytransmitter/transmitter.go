package mercurytransmitter

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
)

const (
	// Mercury server error codes
	DuplicateReport = 2
)

var (
	transmitSuccessCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_mercury_transmit_success_count",
		Help: "Number of successful transmissions (duplicates are counted as success)",
	},
		[]string{"donID", "serverURL"},
	)
	transmitDuplicateCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_mercury_transmit_duplicate_count",
		Help: "Number of transmissions where the server told us it was a duplicate",
	},
		[]string{"donID", "serverURL"},
	)
	transmitConnectionErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_mercury_transmit_connection_error_count",
		Help: "Number of errored transmissions that failed due to problem with the connection",
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
	TransmitQueueMaxSize() uint32
	TransmitTimeout() commonconfig.Duration
}

type transmitter struct {
	services.StateMachine
	lggr logger.SugaredLogger
	cfg  Config

	orm     ORM
	servers map[string]*server

	donID       uint32
	fromAccount string

	stopCh services.StopChan
	wg     *sync.WaitGroup
}

type Opts struct {
	Lggr        logger.Logger
	Cfg         Config
	Clients     map[string]wsrpc.Client
	FromAccount ed25519.PublicKey
	DonID       uint32
	ORM         ORM
}

func New(opts Opts) Transmitter {
	return newTransmitter(opts)
}

func newTransmitter(opts Opts) *transmitter {
	sugared := logger.Sugared(opts.Lggr).Named("LLOMercuryTransmitter")
	servers := make(map[string]*server, len(opts.Clients))
	for serverURL, client := range opts.Clients {
		sLggr := sugared.Named(serverURL).With("serverURL", serverURL)
		servers[serverURL] = newServer(sLggr, opts.Cfg, client, opts.ORM, serverURL)
	}
	return &transmitter{
		services.StateMachine{},
		sugared.Named("LLOMercuryTransmitter").With("donID", opts.ORM.DonID()),
		opts.Cfg,
		opts.ORM,
		servers,
		opts.DonID,
		fmt.Sprintf("%x", opts.FromAccount),
		make(services.StopChan),
		&sync.WaitGroup{},
	}
}

func (mt *transmitter) Start(ctx context.Context) (err error) {
	return mt.StartOnce("LLOMercuryTransmitter", func() error {
		mt.lggr.Debugw("Loading transmit requests from database")

		{
			var startClosers []services.StartClose
			for _, s := range mt.servers {
				transmissions, err := s.pm.Load(ctx)
				if err != nil {
					return err
				}
				s.q.Init(transmissions)
				// starting pm after loading from it is fine because it simply spawns some garbage collection/prune goroutines
				startClosers = append(startClosers, s.c, s.q, s.pm)

				mt.wg.Add(2)
				go s.runDeleteQueueLoop(mt.stopCh, mt.wg)
				go s.runQueueLoop(mt.stopCh, mt.wg, fmt.Sprintf("%d", mt.donID))
			}
			if err := (&services.MultiStart{}).Start(ctx, startClosers...); err != nil {
				return err
			}
		}

		return nil
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
	sigs []types.AttributedOnchainSignature) error {
	transmissions := make([]*Transmission, 0, len(mt.servers))
	for serverURL := range mt.servers {
		transmissions = append(transmissions, &Transmission{
			ServerURL:    serverURL,
			ConfigDigest: digest,
			SeqNr:        seqNr,
			Report:       report,
			Sigs:         sigs,
		})
	}
	if err := mt.orm.Insert(ctx, transmissions); err != nil {
		return err
	}

	g := new(errgroup.Group)
	for i := range transmissions {
		t := transmissions[i]
		mt.lggr.Debugw("LLOMercuryTransmit", "digest", digest.Hex(), "seqNr", seqNr, "reportFormat", report.Info.ReportFormat, "reportLifeCycleStage", report.Info.LifeCycleStage, "transmissionHash", t.Hash())
		g.Go(func() error {
			s := mt.servers[t.ServerURL]
			if ok := s.q.Push(t); !ok {
				s.transmitQueuePushErrorCount.Inc()
				return errors.New("transmit queue is closed")
			}
			return nil
		})
	}

	return g.Wait()
}

// FromAccount returns the stringified (hex) CSA public key
func (mt *transmitter) FromAccount() (ocrtypes.Account, error) {
	return ocrtypes.Account(mt.fromAccount), nil
}
