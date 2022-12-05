package cmd_test

import (
	"bytes"
	"io"
	"regexp"
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/chainlink" //nolint:typecheck
	"github.com/smartcontractkit/chainlink/core/web"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRendererJSON_RenderVRFKeys(t *testing.T) {
	t.Parallel()

	r := cmd.RendererJSON{Writer: io.Discard}
	keys := []cmd.VRFKeyPresenter{
		{
			VRFKeyResource: webpresenters.VRFKeyResource{
				Compressed:   "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01",
				Uncompressed: "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4db44652a69526181101d4aa9a58ecf43b1be972330de99ea5e540f56f4e0a672f",
				Hash:         "0x9926c5f19ec3b3ce005e1c183612f05cfc042966fcdd82ec6e78bf128d91695a",
			},
		},
	}
	assert.NoError(t, r.Render(&keys))
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestRendererTable_RenderConfiguration(t *testing.T) {
	t.Parallel()

	app := cltest.NewLegacyApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	resp, cleanup := client.Get("/v2/config")
	defer cleanup()
	var cp config.ConfigPrinter
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &cp))

	r := cmd.RendererTable{Writer: io.Discard}
	assert.NoError(t, r.Render(&cp))
}

func TestRendererTable_RenderConfigurationV2(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	cfg, ok := app.Config.(chainlink.ConfigV2)
	require.True(t, ok)
	wantUser, wantEffective := cfg.ConfigTOML()
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	t.Run("effective", func(t *testing.T) {
		resp, cleanup := client.Get("/v2/config/v2")
		t.Cleanup(cleanup)
		var effective web.ConfigV2Resource
		require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &effective))

		assert.Equal(t, wantEffective, effective.Config)
	})

	t.Run("user", func(t *testing.T) {
		resp, cleanup := client.Get("/v2/config/v2?userOnly=true")
		t.Cleanup(cleanup)
		var user web.ConfigV2Resource
		require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &user))

		assert.Equal(t, wantUser, user.Config)
	})
}

type testWriter struct {
	expected string
	t        testing.TB
	found    bool
}

func (w *testWriter) Write(actual []byte) (int, error) {
	if bytes.Contains(actual, []byte(w.expected)) {
		w.found = true
	}
	return len(actual), nil
}

func TestRendererTable_RenderExternalInitiatorAuthentication(t *testing.T) {
	t.Parallel()

	eia := webpresenters.ExternalInitiatorAuthentication{
		Name:           "bitcoin",
		URL:            cltest.WebURL(t, "http://localhost:8888"),
		AccessKey:      "accesskey",
		Secret:         "secret",
		OutgoingToken:  "outgoingToken",
		OutgoingSecret: "outgoingSecret",
	}
	tests := []struct {
		name, content string
	}{
		{"Name", eia.Name},
		{"URL", eia.URL.String()},
		{"AccessKey", eia.AccessKey},
		{"Secret", eia.Secret},
		{"OutgoingToken", eia.OutgoingToken},
		{"OutgoingSecret", eia.OutgoingSecret},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tw := &testWriter{test.content, t, false}
			r := cmd.RendererTable{Writer: tw}

			assert.NoError(t, r.Render(&eia))
			assert.True(t, tw.found)
		})
	}
}

func TestRendererTable_PatchResponse(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}

	patchResponse := web.ConfigPatchResponse{
		EvmGasPriceDefault: web.Change{
			From: "98721",
			To:   "53276",
		},
	}

	assert.NoError(t, r.Render(&patchResponse))
	output := buffer.String()
	assert.Regexp(t, regexp.MustCompile("98721"), output)
	assert.Regexp(t, regexp.MustCompile("53276"), output)
}

func TestRendererTable_RenderUnknown(t *testing.T) {
	t.Parallel()
	r := cmd.RendererTable{Writer: io.Discard}
	anon := struct{ Name string }{"Romeo"}
	assert.Error(t, r.Render(&anon))
}
