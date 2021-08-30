package cmd_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSAKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.CSAKeyPresenter{
		JAID: cmd.JAID{ID: pubKey},
		CSAKeyResource: presenters.CSAKeyResource{
			JAID:   presenters.NewJAID(pubKey),
			PubKey: pubKey,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, pubKey)

	// Render many resources
	buffer.Reset()
	ps := cmd.CSAKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, pubKey)
}

func TestClient_ListCSAKeys(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	key, err := app.GetKeyStore().CSA().Create()
	require.NoError(t, err)

	requireCSAKeyCount(t, app, 1)

	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListCSAKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	keys := *r.Renders[0].(*cmd.CSAKeyPresenters)
	assert.Equal(t, key.PublicKeyString(), keys[0].PubKey)
}

func TestClient_CreateCSAKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	requireCSAKeyCount(t, app, 0)

	require.NoError(t, client.CreateCSAKey(nilContext))

	requireCSAKeyCount(t, app, 1)
}

func requireCSAKeyCount(t *testing.T, app chainlink.Application, length int) {
	t.Helper()

	keys, err := app.GetKeyStore().CSA().GetAll()
	require.NoError(t, err)
	require.Equal(t, length, len(keys))
}
