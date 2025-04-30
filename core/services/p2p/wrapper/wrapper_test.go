package wrapper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/wrapper"
)

func TestPeerWrapper_CleanStartClose(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	lggr := logger.TestLogger(t)
	port := freeport.GetOne(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.Capabilities.Peering.V2.Enabled = &enabled
		c.Capabilities.Peering.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}
	})
	keystoreP2P := ksmocks.NewP2P(t)
	key, err := p2pkey.NewV2()
	require.NoError(t, err)
	keystoreP2P.On("GetOrFirst", mock.Anything).Return(key, nil)

	wrapper := wrapper.NewExternalPeerWrapper(keystoreP2P, cfg.Capabilities().Peering(), db, lggr)
	require.NotNil(t, wrapper)
	require.NoError(t, wrapper.Start(testutils.Context(t)))
	require.NoError(t, wrapper.Close())
}
