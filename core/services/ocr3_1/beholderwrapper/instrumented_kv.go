package beholderwrapper

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

type instrumentedKVStateReader struct {
	inner   ocr3_1types.KeyValueStateReader
	ctx     context.Context //nolint:containedctx // libocr 3.1's API doesn't support passing in ctx via the Read/Write method.
	metrics *pluginMetrics
}

func (i *instrumentedKVStateReader) Read(key []byte) ([]byte, error) {
	start := time.Now()
	data, err := i.inner.Read(key)
	i.metrics.recordKVDuration(i.ctx, "Read", time.Since(start), err == nil)
	return data, err
}

type instrumentedKVStateReadWriter struct {
	instrumentedKVStateReader
	writer ocr3_1types.KeyValueStateReadWriter
}

func (i *instrumentedKVStateReadWriter) Write(key []byte, value []byte) error {
	start := time.Now()
	err := i.writer.Write(key, value)
	i.metrics.recordKVDuration(i.ctx, "Write", time.Since(start), err == nil)
	return err
}

func (i *instrumentedKVStateReadWriter) Delete(key []byte) error {
	start := time.Now()
	err := i.writer.Delete(key)
	i.metrics.recordKVDuration(i.ctx, "Delete", time.Since(start), err == nil)
	return err
}
