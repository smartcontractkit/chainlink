package confidentialrelay

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
)

type teeKeyFetcher string

var _ host.EncryptionKeyFetcher = (*teeKeyFetcher)(nil)

func (e teeKeyFetcher) GetEncryptionKeys(_ context.Context) ([]string, error) {
	return []string{string(e)}, nil
}
