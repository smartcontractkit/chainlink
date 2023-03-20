package tools

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WSConfig struct {
	TargetURl string
}

type WSInstance struct {
	srv  *client.MercuryServer
	stop chan struct{}
}

func (m *WSInstance) Stop(l *loadgen.Generator) {
	m.stop <- struct{}{}
}

func (m *WSInstance) Clone(l *loadgen.Generator) loadgen.Instance {
	return &WSInstance{
		srv:  m.srv,
		stop: make(chan struct{}, 1),
	}
}

func NewWSInstance(srv *client.MercuryServer) *WSInstance {
	return &WSInstance{
		srv:  srv,
		stop: make(chan struct{}, 1),
	}
}

// Run create an instance firing read requests against mock ws server
func (m *WSInstance) Run(l *loadgen.Generator) {
	c, _, err := m.srv.DialWS()
	if err != nil {
		l.Log.Error().Err(err).Msg("failed to connect from instance")
		//nolint
		c.Close(websocket.StatusInternalError, "")
	}
	l.ResponsesWaitGroup.Add(1)
	go func() {
		//nolint
		defer l.ResponsesWaitGroup.Done()
		defer c.Close(websocket.StatusNormalClosure, "")
		for {
			select {
			case <-l.ResponsesCtx.Done():
				return
			case <-m.stop:
				return
			default:
				startedAt := time.Now()
				v := map[string]string{}
				err = wsjson.Read(context.Background(), c, &v)
				if err != nil {
					l.Log.Error().Err(err).Msg("failed read ws msg from instance")
					l.ResponsesChan <- loadgen.CallResult{StartedAt: &startedAt, Failed: true, Error: "ws read error"}
				}
				report, err := mercury.DecodeReport([]byte(v["report"]))
				if err != nil {
					l.ResponsesChan <- loadgen.CallResult{Error: "report validation error", Failed: true}
					continue
				}
				log.Info().Interface("Report", report).Msg("Decoded report")
				l.ResponsesChan <- loadgen.CallResult{StartedAt: &startedAt}
			}
		}
	}()
}
