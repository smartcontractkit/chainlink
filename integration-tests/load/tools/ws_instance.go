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
	srv *client.MercuryServer
}

func NewWSInstance(srv *client.MercuryServer) *WSInstance {
	return &WSInstance{
		srv: srv,
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
		defer l.ResponsesWaitGroup.Done()
		for {
			select {
			case <-l.ResponsesCtx.Done():
				//nolint
				c.Close(websocket.StatusNormalClosure, "")
				return
			default:
				startedAt := time.Now()
				v := map[string]string{}
				err = wsjson.Read(context.Background(), c, &v)
				if err != nil {
					l.Log.Error().Err(err).Msg("failed read ws msg from instance")
					l.ResponsesChan <- loadgen.CallResult{StartedAt: &startedAt, Failed: true, Error: "ws read error"}
				}
				log.Debug().Interface("Results", v).Msg("Report results")
				if v["report"] == "" {
					log.Error().Msg("report is empty")
					continue
				}
				reportElements := map[string]interface{}{}
				if err := mercury.ValidateReport([]byte(v["report"])); err != nil {
					l.ResponsesChan <- loadgen.CallResult{Error: "report validation error"}
					continue
				}
				log.Debug().Interface("Report", reportElements).Msg("Decoded report")
				l.ResponsesChan <- loadgen.CallResult{StartedAt: &startedAt}
			}
		}
	}()
}
