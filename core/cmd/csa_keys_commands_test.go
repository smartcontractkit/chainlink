package cmd_test

import (
	"bytes"
	"testing"
	"time"

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
		id        = "1"
		pubKey    = "somepubkey"
		createdAt = time.Now()
		updatedAt = time.Now().Add(time.Second)
		buffer    = bytes.NewBufferString("")
		r         = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.CSAKeyPresenter{
		JAID: cmd.JAID{ID: id},
		CSAKeyResource: presenters.CSAKeyResource{
			JAID:      presenters.NewJAID(id),
			PubKey:    pubKey,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())

	// Render many resources
	buffer.Reset()
	ps := cmd.CSAKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
}

func TestClient_ListCSAKeys(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	key, err := app.GetKeyStore().CSA().CreateCSAKey()
	require.NoError(t, err)

	requireCSAKeyCount(t, app, 1)

	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListCSAKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	keys := *r.Renders[0].(*cmd.CSAKeyPresenters)
	assert.Equal(t, key.PublicKey.String(), keys[0].PubKey)
}

func TestClient_CreateCSAKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	requireCSAKeyCount(t, app, 0)

	require.NoError(t, client.CreateCSAKey(nilContext))

	requireCSAKeyCount(t, app, 1)
}

func requireCSAKeyCount(t *testing.T, app chainlink.Application, length int64) {
	t.Helper()

	count, err := app.GetKeyStore().CSA().CountCSAKeys()
	require.NoError(t, err)
	require.Equal(t, length, count)
}
