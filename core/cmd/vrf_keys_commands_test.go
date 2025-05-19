package cmd_test

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
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

func TestShellVRF_CRUD(t *testing.T) {
	t.Parallel()

	// Test application boots with vrf password loaded in memory.
	// i.e. as if a user had booted with --vrfpassword=<vrfPasswordFilePath>
	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()

	require.NoError(t, client.ListVRFKeys(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 1)
	keys := *r.Renders[0].(*cmd.VRFKeyPresenters)
	// No keys yet
	require.Empty(t, keys)

	// Create a VRF key
	require.NoError(t, client.CreateVRFKey(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 2)
	k1 := *r.Renders[1].(*cmd.VRFKeyPresenter)

	// List the key and ensure it matches
	require.NoError(t, client.ListVRFKeys(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 3)
	keys = *r.Renders[2].(*cmd.VRFKeyPresenters)
	AssertKeysEqual(t, k1, keys[0])

	// Create another key
	require.NoError(t, client.CreateVRFKey(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 4)
	k2 := *r.Renders[3].(*cmd.VRFKeyPresenter)

	// Ensure the list is valid
	require.NoError(t, client.ListVRFKeys(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 5)
	keys = *r.Renders[4].(*cmd.VRFKeyPresenters)
	require.Contains(t, []string{keys[0].ID, keys[1].ID}, k1.ID)
	require.Contains(t, []string{keys[0].ID, keys[1].ID}, k2.ID)

	// Now do a hard delete and ensure its completely removes the key
	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.DeleteVRFKey, set, "")

	require.NoError(t, set.Parse([]string{k2.Compressed}))
	require.NoError(t, set.Set("hard", "true"))
	require.NoError(t, set.Set("yes", "true"))

	c := cli.NewContext(nil, set, nil)
	err := client.DeleteVRFKey(c)
	require.NoError(t, err)
	// Should return the deleted key
	require.Len(t, r.Renders, 6)
	deletedKey := *r.Renders[5].(*cmd.VRFKeyPresenter)
	AssertKeysEqual(t, k2, deletedKey)
	// Should NOT be in the DB as archived
	allKeys, err := app.KeyStore.VRF().GetAll()
	require.NoError(t, err)
	assert.Len(t, allKeys, 1)
}

func TestVRF_ImportExport(t *testing.T) {
	t.Parallel()
	// Test application boots with vrf password loaded in memory.
	// i.e. as if a user had booted with --vrfpassword=<vrfPasswordFilePath>
	app := startNewApplicationV2(t, nil)
	client, r := app.NewShellAndRenderer()
	t.Log(client, r)

	// Create a key (encrypted with cltest.VRFPassword)
	require.NoError(t, client.CreateVRFKey(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 1)
	k1 := *r.Renders[0].(*cmd.VRFKeyPresenter)
	t.Log(k1.Compressed)

	// Export it, encrypted with cltest.Password instead
	keyName := "vrfkey1"
	set := flag.NewFlagSet("test VRF export", 0)
	flagSetApplyFromAction(client.ExportVRFKey, set, "")

	require.NoError(t, set.Parse([]string{k1.Compressed})) // Arguments
	require.NoError(t, set.Set("new-password", "../internal/fixtures/correct_password.txt"))
	require.NoError(t, set.Set("output", keyName))

	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ExportVRFKey(c))
	// File exists
	require.NoError(t, utils.JustError(os.Stat(keyName)))
	t.Cleanup(func() {
		os.Remove(keyName)
	})

	// Should error if we try to import a duplicate key
	importSet := flag.NewFlagSet("test VRF import", 0)
	flagSetApplyFromAction(client.ImportVRFKey, importSet, "")

	require.NoError(t, importSet.Parse([]string{keyName}))
	require.NoError(t, importSet.Set("old-password", "../internal/fixtures/correct_password.txt"))

	importCli := cli.NewContext(nil, importSet, nil)
	require.Error(t, client.ImportVRFKey(importCli))

	// Lets delete the key and import it
	set = flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.DeleteVRFKey, set, "")

	require.NoError(t, set.Parse([]string{k1.Compressed}))
	require.NoError(t, set.Set("hard", "true"))
	require.NoError(t, set.Set("yes", "true"))

	require.NoError(t, client.DeleteVRFKey(cli.NewContext(nil, set, nil)))
	// Should succeed
	require.NoError(t, client.ImportVRFKey(importCli))
	require.NoError(t, client.ListVRFKeys(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 4)
	keys := *r.Renders[3].(*cmd.VRFKeyPresenters)
	AssertKeysEqualNoTimestamps(t, k1, keys[0])
}
