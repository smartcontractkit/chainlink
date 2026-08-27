package ocr2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestNewServices_NilPeerWrapper(t *testing.T) {
	t.Parallel()

	d := &Delegate{
		lggr: logger.TestLogger(t),
	}

	_, err := d.NewServices(context.Background(), "dontime@1.0.0", 1, commontypes.DonTimePlugin, "", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "libp2p peer was missing or not started")
}
