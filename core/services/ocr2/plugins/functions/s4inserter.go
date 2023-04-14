package functions

import (
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4storage"
)

type S4Inserter struct {
	s4Api *s4storage.S4APIService
	lggr  logger.Logger
}

var _ gateway.Handler = &S4Inserter{}

func NewS4Inserter(s4Api *s4storage.S4APIService, lggr logger.Logger) *S4Inserter {
	return &S4Inserter{s4Api: s4Api, lggr: lggr}
}

func (s *S4Inserter) HandleUserMessage(msg *gateway.Message, cb gateway.Callback) {
	if msg.Method == "secrets_upload" {
		s.lggr.Info("GatewayConnector uploading user secrets for ", msg.SenderAddress)
		payload := msg.Payload["secrets"].(string)
		s.s4Api.Put(msg.SenderAddress, 0, []byte(payload), 0)
		msg.Method = "secrets_upload_ack" // BAD
		cb.SendResponse(msg)
	} else {
		s.lggr.Error("GatewayConnector received unsupported method")
	}
}

func (s *S4Inserter) HandleNodeMessage(msg *gateway.Message, nodeAddr string) {
	//no-op
}

func (s *S4Inserter) Init(connMgr gateway.ConnectionManager, donConfig *gateway.GatewayDONConfig) {
	//no-op
}
