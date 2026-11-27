package utils

import (
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func SetupBeholder(lggr logger.Logger) error {
	writer := &lggrWriter{lggr: lggr}

	client, err := beholder.NewWriterClient(writer)
	if err != nil {
		return err
	}

	beholder.SetClient(client)

	return nil
}

func cleanupBeholder() error {
	client := beholder.GetClient()
	if client != nil {
		return client.Close()
	}

	return nil
}

type lggrWriter struct {
	lggr logger.Logger
}

func (w lggrWriter) Write(bts []byte) (int, error) {
	w.lggr.Info(string(bts))

	return len(bts), nil
}
