package mercury

import (
	"context"

	"github.com/nats-io/nats.go"

	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
)

type natsTransmitter struct {
	feedID mercuryutils.FeedID
	url    string

	conn *nats.Conn
}

func newNATSTransmitter(feedID mercuryutils.FeedID, url string) *natsTransmitter {
	return &natsTransmitter{feedID, url, nil}
}

func (nt *natsTransmitter) Start(context.Context) error {
	// TODO: advanced connect/dialler that uses context
	// user: system
	// password: nOS0OJbtBh4RUA1P5KJ5FQYLU6bl1Vso
	nc, err := nats.Connect(nt.url, nats.UserInfo("system", "nOS0OJbtBh4RUA1P5KJ5FQYLU6bl1Vso"))
	nt.conn = nc
	return err
}
func (nt *natsTransmitter) Close() error {
	nt.conn.Close()
	return nil
}
