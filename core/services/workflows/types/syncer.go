package types

import (
	"context"

	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
)

// FetcherFunc is an abstraction for fetching the contents stored at a URL.
type FetcherFunc func(ctx context.Context, messageID string, req ghcapabilities.Request) ([]byte, error)

// LocationRetrieverFunc is an abstraction for getting a URL that can be used to retrieve an artifact.
type LocationRetrieverFunc func(ctx context.Context, req *storage_service.DownloadArtifactRequest) (string, error)
