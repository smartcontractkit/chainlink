package utils

import (
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// SetupBeholder wires beholder's OTLP output into lggr and returns a func that
// closes the client.
func SetupBeholder(lggr logger.Logger) (func(), error) {
	writer := &lggrWriter{lggr: lggr}

	client, err := beholder.NewWriterClient(writer)
	if err != nil {
		return nil, err
	}

	beholder.SetClient(client)

	return func() { _ = client.Close() }, nil
}

type lggrWriter struct {
	lggr logger.Logger
}

func (w lggrWriter) Write(bts []byte) (int, error) {
	w.lggr.Info(string(bts))

	return len(bts), nil
}
