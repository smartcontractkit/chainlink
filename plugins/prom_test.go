package plugins

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPromServer(t *testing.T) {

	s := NewPromServer(0, logger.TestLogger(t))
	require.NoError(t, s.Start())

	tcpAddr, ok := s.Addr().(*net.TCPAddr)
	require.True(t, ok, "expect tcp listener")

	url := fmt.Sprintf("http://localhost:%d/metrics", tcpAddr.Port)
	resp, err := http.Get(url)
	require.NoError(t, err)
	require.NoError(t, err, "endpoint %s", url)
	require.NotNil(t, resp.Body)
	defer resp.Body.Close()

	require.NoError(t, s.Close())
}
