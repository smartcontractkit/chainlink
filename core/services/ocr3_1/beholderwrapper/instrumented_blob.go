package beholderwrapper

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

type instrumentedBlobBroadcastFetcher struct {
	inner   ocr3_1types.BlobBroadcastFetcher
	metrics *pluginMetrics
	instrumentedBlobFetcher
}

func (i *instrumentedBlobBroadcastFetcher) BroadcastBlob(ctx context.Context, payload []byte, expirationHint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	start := time.Now()
	handle, err := i.inner.BroadcastBlob(ctx, payload, expirationHint)
	i.metrics.recordBlobDuration(ctx, "BroadcastBlob", time.Since(start), err == nil)
	return handle, err
}

type instrumentedBlobFetcher struct {
	inner   ocr3_1types.BlobFetcher
	metrics *pluginMetrics
}

func (i *instrumentedBlobFetcher) FetchBlob(ctx context.Context, handle ocr3_1types.BlobHandle) ([]byte, error) {
	start := time.Now()
	data, err := i.inner.FetchBlob(ctx, handle)
	i.metrics.recordBlobDuration(ctx, "FetchBlob", time.Since(start), err == nil)
	return data, err
}
