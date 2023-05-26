package cmd_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestExternalInitiatorPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		name          = "ExternalInitiator 1"
		url           = cltest.MustWebURL(t, "http://example.com")
		createdAt     = time.Now()
		updatedAt     = time.Now()
		outgoingToken = "anoutgoingtoken"
		accessKey     = "anaccesskey"
		buffer        = bytes.NewBufferString("")
		r             = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.ExternalInitiatorPresenter{
		ExternalInitiatorResource: presenters.ExternalInitiatorResource{
			JAID:          presenters.NewJAID(name),
			Name:          name,
			URL:           url,
			AccessKey:     accessKey,
			OutgoingToken: outgoingToken,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, name)
	assert.Contains(t, output, url.String())
	assert.Contains(t, output, accessKey)
	assert.Contains(t, output, outgoingToken)

	// Render many resources
	buffer.Reset()
	ps := cmd.ExternalInitiatorPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, name)
	assert.Contains(t, output, url.String())
	assert.Contains(t, output, accessKey)
	assert.Contains(t, output, outgoingToken)
}
