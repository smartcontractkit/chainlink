package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestBridgePresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		name          = "Bridge 1"
		url           = "http://example.com"
		createdAt     = time.Now()
		outgoingToken = "anoutgoingtoken"
		buffer        = bytes.NewBufferString("")
		r             = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.BridgePresenter{
		BridgeResource: presenters.BridgeResource{
			JAID:          presenters.NewJAID(name),
			Name:          name,
			URL:           url,
			Confirmations: 10,
			OutgoingToken: outgoingToken,
			CreatedAt:     createdAt,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, name)
	assert.Contains(t, output, url)
	assert.Contains(t, output, "10")
	assert.Contains(t, output, outgoingToken)

	// Render many resources
	buffer.Reset()
	ps := cmd.BridgePresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, name)
	assert.Contains(t, output, url)
	assert.Contains(t, output, "10")
	assert.NotContains(t, output, outgoingToken)
}

func TestShell_IndexBridges(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	bt1 := &bridges.BridgeType{
		Name:          bridges.MustParseBridgeName("cliindexbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err := app.BridgeORM().CreateBridgeType(ctx, bt1)
	require.NoError(t, err)

	bt2 := &bridges.BridgeType{
		Name:          bridges.MustParseBridgeName("cliindexbridges2"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err = app.BridgeORM().CreateBridgeType(ctx, bt2)
	require.NoError(t, err)

	require.NoError(t, client.IndexBridges(cltest.EmptyCLIContext()))
	bridges := *r.Renders[0].(*cmd.BridgePresenters)
	require.Len(t, bridges, 2)
	p := bridges[0]
	assert.Equal(t, bt1.Name.String(), p.Name)
	assert.Equal(t, bt1.URL.String(), p.URL)
	assert.Equal(t, bt1.Confirmations, p.Confirmations)

	p = bridges[1]
	assert.Equal(t, bt2.Name.String(), p.Name)
	assert.Equal(t, bt2.URL.String(), p.URL)
	assert.Equal(t, bt2.Confirmations, p.Confirmations)
}

func TestShell_ShowBridge(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	bt := &bridges.BridgeType{
		Name:          bridges.MustParseBridgeName(testutils.RandomizeName("showbridge")),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	require.NoError(t, app.BridgeORM().CreateBridgeType(testutils.Context(t), bt))

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.ShowBridge, set, "")

	require.NoError(t, set.Parse([]string{bt.Name.String()}))

	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.ShowBridge(c))
	require.Len(t, r.Renders, 1)
	p := r.Renders[0].(*cmd.BridgePresenter)
	assert.Equal(t, bt.Name.String(), p.Name)
	assert.Equal(t, bt.URL.String(), p.URL)
	assert.Equal(t, bt.Confirmations, p.Confirmations)
}

func TestShell_CreateBridge(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()

	tests := []struct {
		name    string
		param   string
		errored bool
	}{
		{"EmptyString", "", true},
		{"ValidString", `{ "name": "TestBridge", "url": "http://localhost:3000/randomNumber" }`, false},
		{"InvalidString", `{ "noname": "", "nourl": "" }`, true},
		{"InvalidChar", `{ "badname": "path/bridge", "nourl": "" }`, true},
		{"ValidPath", "../testdata/apiresponses/create_random_number_bridge_type.json", false},
		{"InvalidPath", "bad/filepath/", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			set := flag.NewFlagSet("bridge", 0)
			flagSetApplyFromAction(client.CreateBridge, set, "")

			require.NoError(t, set.Parse([]string{test.param}))

			c := cli.NewContext(nil, set, nil)
			if test.errored {
				assert.Error(t, client.CreateBridge(c))
			} else {
				assert.NoError(t, client.CreateBridge(c))
			}
		})
	}
}

func TestShell_RemoveBridge(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	bt := &bridges.BridgeType{
		Name:          bridges.MustParseBridgeName(testutils.RandomizeName("removebridge")),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err := app.BridgeORM().CreateBridgeType(testutils.Context(t), bt)
	require.NoError(t, err)

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.RemoveBridge, set, "")

	require.NoError(t, set.Parse([]string{bt.Name.String()}))

	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.RemoveBridge(c))

	require.Len(t, r.Renders, 1)
	p := r.Renders[0].(*cmd.BridgePresenter)
	assert.Equal(t, bt.Name.String(), p.Name)
	assert.Equal(t, bt.URL.String(), p.URL)
	assert.Equal(t, bt.Confirmations, p.Confirmations)
}
func TestShell_UpdateBridge(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()
	name := testutils.RandomizeName("updatebridge")

	bt := &bridges.BridgeType{
		Name:          bridges.MustParseBridgeName(name),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	require.NoError(t, app.BridgeORM().CreateBridgeType(testutils.Context(t), bt))
	tests := []struct {
		name    string
		args    []string
		errored bool
	}{
		{"NoArgs", []string{}, true},
		{"OnlyName", []string{name}, true},
		{"ValidUpdate", []string{name, fmt.Sprintf(`{ "name": "%s", "url": "http://localhost:3000/updated" }`, name)}, false},
		{"InvalidJSON", []string{name, `{ "url": "http://localhost:3000/updated"`}, true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			set := flag.NewFlagSet("bridge", 0)
			flagSetApplyFromAction(client.UpdateBridge, set, "")

			require.NoError(t, set.Parse(test.args))

			c := cli.NewContext(nil, set, nil)
			if test.errored {
				assert.Error(t, client.UpdateBridge(c))
			} else {
				assert.NoError(t, client.UpdateBridge(c))
			}
		})
	}
}
