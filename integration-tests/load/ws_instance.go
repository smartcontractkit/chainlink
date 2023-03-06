package load

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
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
	c, _, err := websocket.Dial(context.Background(), fmt.Sprintf("%s/ws", m.srv.URL), &websocket.DialOptions{
		HTTPHeader: http.Header{"Authorization": []string{"Basic Y2xpZW50OmNsaWVudHBhc3M="}},
	})
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
				if err = ReportTypes.UnpackIntoMap(reportElements, []byte(v["report"])); err != nil {
					l.Log.Error().Err(err).Msg("failed to unpack report")
					l.ResponsesChan <- loadgen.CallResult{Error: "blob unpacking error"}
					continue
				}
				if err := testsetups.ValidateReport(reportElements); err != nil {
					l.ResponsesChan <- loadgen.CallResult{Error: "report validation error"}
					continue
				}
				log.Debug().Interface("Report", reportElements).Msg("Decoded report")
				l.ResponsesChan <- loadgen.CallResult{StartedAt: &startedAt}
			}
		}
	}()
}
