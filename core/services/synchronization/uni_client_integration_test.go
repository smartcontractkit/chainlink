package synchronization

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/smartcontractkit/wsrpc"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/synchronization/telem"
)

func TestUniClient(t *testing.T) {
	t.Skip()
	privKey, err := hex.DecodeString("TODO")
	require.NoError(t, err)
	pubKey, err := hex.DecodeString("TODO")
	require.NoError(t, err)
	t.Log(len(privKey), len(pubKey))
	lggr := logger.TestLogger(t)
	c, err := wsrpc.DialUniWithContext(context.Background(),
		lggr,
		"TODO",
		privKey,
		pubKey)
	require.NoError(t, err)
	t.Log(c)
	client := telem.NewTelemClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	resp, err := client.Telem(ctx, &telem.TelemRequest{
		Telemetry: []byte(`hello world`),
		Address:   "myaddress",
	})
	cancel()
	t.Log(resp, err)
	require.NoError(t, c.Close())
}
