package plugin

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

// noopBlobBroadcastFetcher implements BlobBroadcastFetcher without network I/O.
// The consensus plugin does not use blob broadcast; the fetcher is only required
// by the OCR 3.1 ReportingPluginFactory signature.
type noopBlobBroadcastFetcher struct{}

func (noopBlobBroadcastFetcher) BroadcastBlob(_ context.Context, _ []byte, _ ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	return ocr3_1types.BlobHandle{}, fmt.Errorf("blob broadcast not used by blobconsensus plugin")
}

func (noopBlobBroadcastFetcher) FetchBlob(_ context.Context, _ ocr3_1types.BlobHandle) ([]byte, error) {
	return nil, fmt.Errorf("blob fetch not used by blobconsensus plugin")
}

var defaultNoopBlobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher = noopBlobBroadcastFetcher{}
