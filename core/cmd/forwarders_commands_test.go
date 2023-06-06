package cmd_test

import (
	"bytes"
	"flag"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestEVMForwarderPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id         = "1"
		address    = common.HexToAddress("0x5431F5F973781809D18643b87B44921b11355d81")
		evmChainID = utils.NewBigI(4)
		createdAt  = time.Now()
		updatedAt  = time.Now().Add(time.Second)
		buffer     = bytes.NewBufferString("")
		r          = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.EVMForwarderPresenter{
		EVMForwarderResource: presenters.EVMForwarderResource{
			JAID:       presenters.NewJAID(id),
			Address:    address,
			EVMChainID: *evmChainID,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, address.String())
	assert.Contains(t, output, evmChainID.ToInt().String())
	assert.Contains(t, output, createdAt.Format(time.RFC3339))

	// Render many resources
	buffer.Reset()
	ps := cmd.EVMForwarderPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, address.String())
	assert.Contains(t, output, evmChainID.ToInt().String())
	assert.Contains(t, output, createdAt.Format(time.RFC3339))
}

func TestClient_TrackEVMForwarder(t *testing.T) {
	t.Parallel()

	id := newRandChainID()
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].ChainID = id
		c.EVM[0].Enabled = ptr(true)
	})
	client, r := app.NewClientAndRenderer()

	// Create the fwdr
	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.TrackForwarder, set, "")

	require.NoError(t, set.Set("address", "0x5431F5F973781809D18643b87B44921b11355d81"))
	require.NoError(t, set.Set("evmChainID", id.String()))

	err := client.TrackForwarder(cli.NewContext(nil, set, nil))
	require.NoError(t, err)
	require.Len(t, r.Renders, 1)
	createOutput, ok := r.Renders[0].(*cmd.EVMForwarderPresenter)
	require.True(t, ok, "Expected Renders[0] to be *cmd.EVMForwarderPresenter, got %T", r.Renders[0])

	// Assert fwdr is listed
	require.Nil(t, client.ListForwarders(cltest.EmptyCLIContext()))
	fwds := *r.Renders[1].(*cmd.EVMForwarderPresenters)
	require.Equal(t, 1, len(fwds))
	assert.Equal(t, createOutput.ID, fwds[0].ID)

	// Delete fwdr
	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteForwarder, set, "")

	require.NoError(t, set.Parse([]string{createOutput.ID}))

	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.DeleteForwarder(c))

	// Assert fwdr is not listed
	require.Nil(t, client.ListForwarders(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 3)
	fwds = *r.Renders[2].(*cmd.EVMForwarderPresenters)
	require.Equal(t, 0, len(fwds))
}

func TestClient_TrackEVMForwarder_BadAddress(t *testing.T) {
	t.Parallel()

	id := newRandChainID()
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].ChainID = id
		c.EVM[0].Enabled = ptr(true)
	})
	client, _ := app.NewClientAndRenderer()

	// Create the fwdr
	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.TrackForwarder, set, "")

	require.NoError(t, set.Set("address", "0xWrongFormatAddress"))
	require.NoError(t, set.Set("evmChainID", id.String()))

	err := client.TrackForwarder(cli.NewContext(nil, set, nil))
	require.Contains(t, err.Error(), "could not decode address: invalid hex string")
}

func TestClient_DeleteEVMForwarders_MissingFwdId(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
	})
	client, _ := app.NewClientAndRenderer()

	// Delete fwdr without id
	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteForwarder, set, "")

	c := cli.NewContext(nil, set, nil)
	require.Equal(t, "must pass the forwarder id to be archived", client.DeleteForwarder(c).Error())
}
