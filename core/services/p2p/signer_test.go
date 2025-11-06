package p2p

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

func TestSigner_InitializeAndSign(t *testing.T) {
	keystoreP2P := ksmocks.NewP2P(t)
	key, err := p2pkey.NewV2()
	require.NoError(t, err)
	keystoreP2P.On("GetOrFirst", mock.Anything).Return(key, nil)
	s := NewSigner(keystoreP2P, p2pkey.PeerID{}) // peerID unset gets the default one

	_, err = s.Sign([]byte("msg"))
	require.Error(t, err)
	require.Equal(t, "private key not set", err.Error())

	require.NoError(t, s.Initialize())
	sig, err := s.Sign([]byte("msg"))
	require.NoError(t, err)
	require.NotNil(t, sig)
}
