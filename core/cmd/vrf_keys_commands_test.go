package cmd_test

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVRFKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		compressed   = "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01"
		uncompressed = "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4db44652a69526181101d4aa9a58ecf43b1be972330de99ea5e540f56f4e0a672f"
		hash         = "0x9926c5f19ec3b3ce005e1c183612f05cfc042966fcdd82ec6e78bf128d91695a"
		buffer       = bytes.NewBufferString("")
		r            = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.VRFKeyPresenter{
		VRFKeyResource: presenters.VRFKeyResource{
			Compressed:   compressed,
			Uncompressed: uncompressed,
			Hash:         hash,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, compressed)
	assert.Contains(t, output, uncompressed)
	assert.Contains(t, output, hash)

	// Render many resources
	buffer.Reset()
	ps := cmd.VRFKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, compressed)
	assert.Contains(t, output, uncompressed)
	assert.Contains(t, output, hash)
}

func AssertKeysEqual(t *testing.T, k1, k2 cmd.VRFKeyPresenter) {
	AssertKeysEqualNoTimestamps(t, k1, k2)
}

func AssertKeysEqualNoTimestamps(t *testing.T, k1, k2 cmd.VRFKeyPresenter) {
	assert.Equal(t, k1.Compressed, k2.Compressed)
	assert.Equal(t, k1.Hash, k2.Hash)
	assert.Equal(t, k1.Uncompressed, k2.Uncompressed)
}

func TestClientVRF_CRUD(t *testing.T) {
	t.Parallel()

	// Test application boots with vrf password loaded in memory.
	// i.e. as if a user had booted with --vrfpassword=<vrfPasswordFilePath>
	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	require.NoError(t, client.ListVRFKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	keys := *r.Renders[0].(*cmd.VRFKeyPresenters)
	// No keys yet
	require.Equal(t, 0, len(keys))

	// Create a VRF key
	require.NoError(t, client.CreateVRFKey(cltest.EmptyCLIContext()))
	require.Equal(t, 2, len(r.Renders))
	k1 := *r.Renders[1].(*cmd.VRFKeyPresenter)

	// List the key and ensure it matches
	require.NoError(t, client.ListVRFKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 3, len(r.Renders))
	keys = *r.Renders[2].(*cmd.VRFKeyPresenters)
	AssertKeysEqual(t, k1, keys[0])

	// Create another key
	require.NoError(t, client.CreateVRFKey(cltest.EmptyCLIContext()))
	require.Equal(t, 4, len(r.Renders))
	k2 := *r.Renders[3].(*cmd.VRFKeyPresenter)

	// Ensure the list is valid
	require.NoError(t, client.ListVRFKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 5, len(r.Renders))
	keys = *r.Renders[4].(*cmd.VRFKeyPresenters)
	require.Contains(t, []string{keys[0].ID, keys[1].ID}, k1.ID)
	require.Contains(t, []string{keys[0].ID, keys[1].ID}, k2.ID)

	// Now do a hard delete and ensure its completely removes the key
	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{k2.Compressed})
	set.Bool("hard", true, "")
	set.Bool("yes", true, "")
	c := cli.NewContext(nil, set, nil)
	err := client.DeleteVRFKey(c)
	require.NoError(t, err)
	// Should return the deleted key
	require.Equal(t, 6, len(r.Renders))
	deletedKey := *r.Renders[5].(*cmd.VRFKeyPresenter)
	AssertKeysEqual(t, k2, deletedKey)
	// Should NOT be in the DB as archived
	allKeys, err := app.KeyStore.VRF().GetAll()
	require.NoError(t, err)
	assert.Equal(t, 1, len(allKeys))
}

func TestVRF_ImportExport(t *testing.T) {
	t.Parallel()
	// Test application boots with vrf password loaded in memory.
	// i.e. as if a user had booted with --vrfpassword=<vrfPasswordFilePath>
	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()
	t.Log(client, r)

	// Create a key (encrypted with cltest.VRFPassword)
	require.NoError(t, client.CreateVRFKey(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	k1 := *r.Renders[0].(*cmd.VRFKeyPresenter)
	t.Log(k1.Compressed)

	// Export it, encrypted with cltest.Password instead
	keyName := "vrfkey1"
	set := flag.NewFlagSet("test VRF export", 0)
	set.Parse([]string{k1.Compressed}) // Arguments
	set.String("newpassword", "../internal/fixtures/correct_password.txt", "")
	set.String("output", keyName, "")
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ExportVRFKey(c))
	// File exists
	require.NoError(t, utils.JustError(os.Stat(keyName)))
	t.Cleanup(func() {
		os.Remove(keyName)
	})

	// Should error if we try to import a duplicate key
	importSet := flag.NewFlagSet("test VRF import", 0)
	importSet.Parse([]string{keyName})
	importSet.String("oldpassword", "../internal/fixtures/correct_password.txt", "")
	importCli := cli.NewContext(nil, importSet, nil)
	require.Error(t, client.ImportVRFKey(importCli))

	// Lets delete the key and import it
	set = flag.NewFlagSet("test", 0)
	set.Parse([]string{k1.Compressed})
	set.Bool("hard", true, "")
	set.Bool("yes", true, "")
	require.NoError(t, client.DeleteVRFKey(cli.NewContext(nil, set, nil)))
	// Should succeed
	require.NoError(t, client.ImportVRFKey(importCli))
	require.NoError(t, client.ListVRFKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 4, len(r.Renders))
	keys := *r.Renders[3].(*cmd.VRFKeyPresenters)
	AssertKeysEqualNoTimestamps(t, k1, keys[0])
}
