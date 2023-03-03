package networking

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/subprocesses"
)

var _ commontypes.BinaryNetworkEndpoint = &ocrEndpointV1V2{}

const (
	v2InactivityTimeout     = 1 * time.Minute
	v2Headstart             = v2InactivityTimeout / 2
	lastHeardReportInterval = 5 * time.Minute
)

type ocrEndpointV1V2State int

const (
	_ ocrEndpointV1V2State = iota
	ocrEndpointV1V2Unstarted
	ocrEndpointV1V2Started
	ocrEndpointV1V2Closed
)

type ocrEndpointV1V2 struct {
	stateMu   sync.RWMutex
	state     ocrEndpointV1V2State
	v2Started bool // we operate even if v2 won't properly start

	logger    commontypes.Logger
	peerIDs   []string
	v1        commontypes.BinaryNetworkEndpoint
	v2        commontypes.BinaryNetworkEndpoint
	chRecv    chan commontypes.BinaryMessageWithSender
	processes subprocesses.Subprocesses
	chClose   chan struct{}
}

func (o *ocrEndpointV1V2) SendTo(payload []byte, to commontypes.OracleID) {
	o.stateMu.RLock()
	defer o.stateMu.RUnlock()
	if o.state != ocrEndpointV1V2Started {
		o.logger.Error("OCREndpointV1V2: Asked to SentTo while not in started state", commontypes.LogFields{"state": o.state})
		return
	}
	o.v1.SendTo(payload, to)
	if o.v2Started {
		o.v2.SendTo(payload, to)
	}
}

func (o *ocrEndpointV1V2) Broadcast(payload []byte) {
	o.stateMu.RLock()
	defer o.stateMu.RUnlock()
	if o.state != ocrEndpointV1V2Started {
		o.logger.Error("OCREndpointV1V2: Asked to Broadcast while not in started state", commontypes.LogFields{"state": o.state})
		return
	}
	o.v1.Broadcast(payload)
	if o.v2Started {
		o.v2.Broadcast(payload)
	}
}

func (o *ocrEndpointV1V2) Receive() <-chan commontypes.BinaryMessageWithSender {
	return o.chRecv
}

func (o *ocrEndpointV1V2) mergeRecvs() {
	const V1, V2 = 0, 1
	chRecvs := make([]<-chan commontypes.BinaryMessageWithSender, 2)
	o.stateMu.RLock()
	if o.state != ocrEndpointV1V2Started {
		o.stateMu.RUnlock()
		return
	}
	chRecvs[V1] = o.v1.Receive()
	if o.v2Started {
		chRecvs[V2] = o.v2.Receive()
	}
	o.stateMu.RUnlock()
	lastHeardV2 := make([]time.Time, len(o.peerIDs))
	messagesSinceLastReportV1, messagesSinceLastReportV2 := make([]int, len(o.peerIDs)), make([]int, len(o.peerIDs))
	messagesSinceStartupV1, messagesSinceStartupV2 := make([]int, len(o.peerIDs)), make([]int, len(o.peerIDs))
	lastMessageWasV1OrV2 := make([]string, len(o.peerIDs))
	switchesSinceLastReport, switchesSinceStartup := make([]int, len(o.peerIDs)), make([]int, len(o.peerIDs))
	for i := 0; i < len(o.peerIDs); i++ {
		lastHeardV2[i] = time.Now().Add(-v2InactivityTimeout + v2Headstart)
		lastMessageWasV1OrV2[i] = "none"
	}
	ticker := time.NewTicker(lastHeardReportInterval)
	defer ticker.Stop()
	for {
		select {
		case msg := <-chRecvs[V1]:
			if time.Since(lastHeardV2[msg.Sender]) > v2InactivityTimeout {
				select {
				case o.chRecv <- msg:
				case <-o.chClose:
					return
				}
				if lastMessageWasV1OrV2[msg.Sender] != "V1" {
					switchesSinceLastReport[msg.Sender]++
					switchesSinceStartup[msg.Sender]++
				}
				lastMessageWasV1OrV2[msg.Sender] = "V1"
			}

			messagesSinceLastReportV1[msg.Sender]++
			messagesSinceStartupV1[msg.Sender]++
		case msg := <-chRecvs[V2]:
			lastHeardV2[msg.Sender] = time.Now()
			select {
			case o.chRecv <- msg:
			case <-o.chClose:
				return
			}
			if lastMessageWasV1OrV2[msg.Sender] != "V2" {
				switchesSinceLastReport[msg.Sender]++
				switchesSinceStartup[msg.Sender]++
			}
			lastMessageWasV1OrV2[msg.Sender] = "V2"

			messagesSinceLastReportV2[msg.Sender]++
			messagesSinceStartupV2[msg.Sender]++
		case <-ticker.C:
			durationSinceLastHeardV2 := make([]time.Duration, len(lastHeardV2))
			now := time.Now()
			for i, lastTime := range lastHeardV2 {
				durationSinceLastHeardV2[i] = now.Sub(lastTime)
			}
			o.logger.Info("OCREndpointV1V2: Status report", commontypes.LogFields{
				"peerIDs":                   o.peerIDs,
				"durationSinceLastHeardV2":  durationSinceLastHeardV2,
				"messagesSinceLastReportV2": messagesSinceLastReportV2,
				"messagesSinceStartupV2":    messagesSinceStartupV2,
				"messagesSinceLastReportV1": messagesSinceLastReportV1,
				"messagesSinceStartupV1":    messagesSinceStartupV1,
				"switchesSinceLastReport":   switchesSinceLastReport,
				"switchesSinceStartup":      switchesSinceStartup,
				"lastMessageWasV1OrV2":      lastMessageWasV1OrV2,
			})
			for i := 0; i < len(o.peerIDs); i++ {
				messagesSinceLastReportV1[i] = 0
				messagesSinceLastReportV2[i] = 0
				switchesSinceLastReport[i] = 0
			}
		case <-o.chClose:
			return
		}
	}
}

// Start starts the underlying v1 and v2 OCR endpoints. In case the v2 endpoint
// fails, we log an error but do not return it. ocrEndpointV1V2 is designed to
// be resilient against v2 failing to Start and will operate using only v1 if
// needed.
func (o *ocrEndpointV1V2) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			o.logger.Warn("OCREndpointV1V2: Start: errored, auto-closing", nil)
			o.Close()
		}
	}()

	o.stateMu.Lock()
	defer o.stateMu.Unlock()
	if o.state != ocrEndpointV1V2Unstarted {
		return fmt.Errorf("cannot Start ocrEndpointV1V2 that is in state %v", o.state)
	}
	o.state = ocrEndpointV1V2Started

	if err := o.v1.Start(); err != nil {
		o.logger.Warn("OCREndpointV1V2: Start: Failed to start v1", commontypes.LogFields{"err": err})
		return err
	}
	if err := o.v2.Start(); err != nil {
		o.logger.Critical("OCREndpointV1V2: Start: Failed to start v2 OCR endpoint as part of v1v2, operating only with v1", commontypes.LogFields{"error": err})
	} else {
		o.v2Started = true
	}
	o.processes.Go(o.mergeRecvs)
	succeeded = true
	return nil
}

func (o *ocrEndpointV1V2) Close() error {
	err := func() error {
		o.stateMu.Lock()
		defer o.stateMu.Unlock()
		if o.state != ocrEndpointV1V2Started {
			return fmt.Errorf("cannot Close ocrEndpointV1V2 that is in state %v", o.state)
		}
		o.state = ocrEndpointV1V2Closed
		return nil
	}()
	if err != nil {
		return err
	}

	close(o.chClose)
	o.processes.Wait()

	var allErrors error
	allErrors = multierr.Append(allErrors, o.v1.Close())
	allErrors = multierr.Append(allErrors, o.v2.Close())
	return allErrors
}

func newOCREndpointV1V2(
	logger loghelper.LoggerWithContext,
	peerIDs []string,
	ocrV1 commontypes.BinaryNetworkEndpoint,
	ocrV2 commontypes.BinaryNetworkEndpoint,
) (*ocrEndpointV1V2, error) {
	if ocrV1 == nil || ocrV2 == nil {
		return nil, errors.New("cannot accept nil ocr endpoints")
	}
	return &ocrEndpointV1V2{
		sync.RWMutex{},
		ocrEndpointV1V2Unstarted,
		false,
		logger,
		peerIDs,
		ocrV1,
		ocrV2,
		make(chan commontypes.BinaryMessageWithSender),
		subprocesses.Subprocesses{},
		make(chan struct{}),
	}, nil
}
