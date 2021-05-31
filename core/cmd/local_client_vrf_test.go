package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

const (
	vrfPasswordFilePath = "../testdata/secrets/vrf_password.txt"
	vrfKeyFilePath      = "../testdata/secrets/vrf_key.json"
	// This is the public key found in the vrf key file
	vrfPublicKey = "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01"
)

func TestVRFKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		compressed   = "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01"
		uncompressed = "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4db44652a69526181101d4aa9a58ecf43b1be972330de99ea5e540f56f4e0a672f"
		hash         = "0x9926c5f19ec3b3ce005e1c183612f05cfc042966fcdd82ec6e78bf128d91695a"
		createdAt    = time.Now()
		updatedAt    = time.Now().Add(time.Second)
		deletedAt    = time.Now().Add(2 * time.Second)
		buffer       = bytes.NewBufferString("")
		r            = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.VRFKeyPresenter{
		Compressed:   compressed,
		Uncompressed: uncompressed,
		Hash:         hash,
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
		DeletedAt:    &deletedAt,
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, compressed)
	assert.Contains(t, output, uncompressed)
	assert.Contains(t, output, hash)
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, deletedAt.String())

	// Render many resources
	buffer.Reset()
	ps := cmd.VRFKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, compressed)
	assert.Contains(t, output, uncompressed)
	assert.Contains(t, output, hash)
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, deletedAt.String())
}

func TestLocalClientVRF_ListVRFKeys(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)

	r := &cltest.RendererMock{}
	client := cmd.Client{
		Renderer:   r,
		Config:     store.Config,
		AppFactory: cltest.InstanceAppFactory{App: app},
	}

	// Import a key
	set := flag.NewFlagSet("test", 0)
	set.String("password", vrfPasswordFilePath, "")
	set.String("file", vrfKeyFilePath, "")
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportVRFKey(c))

	assert.Nil(t, client.ListVRFKeys(cltest.EmptyCLIContext()))

	require.Equal(t, 1, len(r.Renders))
	p := *r.Renders[0].(*cmd.VRFKeyPresenters)
	fmt.Printf("%+v", p)
	assert.Equal(t, "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01", p[0].Compressed)
}

func TestLocalClientVRF_CreateVRFKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)

	client := cmd.Client{
		Config:     store.Config,
		AppFactory: cltest.InstanceAppFactory{App: app},
	}

	// Must supply password
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, "must specify password file", client.CreateVRFKey(c).Error())

	set = flag.NewFlagSet("test", 0)
	set.String("password", vrfPasswordFilePath, "")
	c = cli.NewContext(nil, set, nil)

	requireVRFKeysCount(t, store, 0)

	require.NoError(t, client.CreateVRFKey(c))

	requireVRFKeysCount(t, store, 1)
}

func TestLocalClientVRF_ImportVRFKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)

	client := cmd.Client{
		Config:     store.Config,
		AppFactory: cltest.InstanceAppFactory{App: app},
	}

	// Must supply password
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, "must specify password file", client.ImportVRFKey(c).Error())

	// Must supply file
	set = flag.NewFlagSet("test", 0)
	set.String("password", vrfPasswordFilePath, "")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, "must specify key file", client.ImportVRFKey(c).Error())

	// Failed to read file
	set = flag.NewFlagSet("test", 0)
	set.String("password", vrfPasswordFilePath, "")
	set.String("file", "./testdata/does_not_exist.json", "")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, "failed to read file ./testdata/does_not_exist.json: open ./testdata/does_not_exist.json: no such file or directory", client.ImportVRFKey(c).Error())

	// Success
	set = flag.NewFlagSet("test", 0)
	set.String("password", vrfPasswordFilePath, "")
	set.String("file", vrfKeyFilePath, "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportVRFKey(c))

	keys := []vrf.EncryptedVRFKey{}
	app.GetStore().DB.Find(&keys)
	assert.Len(t, keys, 1)

	pubKey, err := secp256k1.NewPublicKeyFromHex(vrfPublicKey)
	require.NoError(t, err)
	assert.Equal(t, pubKey, keys[0].PublicKey)

	// Already exists
	require.Equal(t,
		"while attempting to import key from CL: key with matching public key already stored in DB",
		client.ImportVRFKey(c).Error(),
	)
}

func TestLocalClientVRF_ExportVRFKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)

	client := cmd.Client{
		Config:     store.Config,
		AppFactory: cltest.InstanceAppFactory{App: app},
	}

	// Must supply public key
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, "must specify public key", client.ExportVRFKey(c).Error())

	// No key found
	set = flag.NewFlagSet("test", 0)
	set.String("publicKey", vrfPublicKey, "")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, "while retrieving keys with matching public key 0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01: could not find any keys with public key 0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01", client.ExportVRFKey(c).Error())

	// Import the file
	set = flag.NewFlagSet("test", 0)
	set.String("password", vrfPasswordFilePath, "")
	set.String("file", vrfKeyFilePath, "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportVRFKey(c))

	// Must supply file
	set = flag.NewFlagSet("test", 0)
	set.String("publicKey", vrfPublicKey, "")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, "must specify file to export to", client.ExportVRFKey(c).Error())

	// Success
	testdir := filepath.Join(os.TempDir(), t.Name())
	err := os.MkdirAll(testdir, 0700|os.ModeDir)
	assert.NoError(t, err)
	defer os.RemoveAll(testdir)

	keyfilepath := filepath.Join(testdir, "key")
	defer deleteKeyExportFile(t)

	set = flag.NewFlagSet("test", 0)
	set.String("publicKey", vrfPublicKey, "")
	set.String("file", keyfilepath, "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ExportVRFKey(c))

	// Will not overwrite existing file
	require.Contains(t, client.ExportVRFKey(c).Error(), "refusing to overwrite existing file")
}

func TestLocalClientVRF_DeleteVRFKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)

	client := cmd.Client{
		Config:     store.Config,
		AppFactory: cltest.InstanceAppFactory{App: app},
	}

	// Import a key
	set := flag.NewFlagSet("test", 0)
	set.String("password", vrfPasswordFilePath, "")
	set.String("file", vrfKeyFilePath, "")
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportVRFKey(c))

	// Must supply public key
	set = flag.NewFlagSet("test", 0)
	set.Bool("yes", true, "")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, "must specify public key", client.DeleteVRFKey(c).Error())

	// Delete the key
	set = flag.NewFlagSet("test", 0)
	set.String("publicKey", vrfPublicKey, "")
	set.Bool("yes", true, "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.DeleteVRFKey(c))

	keys := []vrf.EncryptedVRFKey{}
	app.GetStore().DB.Find(&keys)
	assert.Len(t, keys, 0)
}

func requireVRFKeysCount(t *testing.T, store *store.Store, length int) []*secp256k1.PublicKey {
	keys, err := store.VRFKeyStore.ListKeys()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
