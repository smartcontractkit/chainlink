package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/static"
	clipkg "github.com/urfave/cli"
)

type SetupDKGNodePayload struct {
	OnChainPublicKey  string
	OffChainPublicKey string
	ConfigPublicKey   string
	PeerID            string
	Transmitter       string
	DkgEncrypt        string
	DkgSign           string
}

type DKGTemplateArgs struct {
	contractID              string
	ocrKeyBundleID          string
	p2pv2BootstrapperPeerID string
	p2pv2BootstrapperPort   string
	transmitterID           string
	chainID                 int64
	EncryptionPublicKey     string
	KeyID                   string
	SigningPublicKey        string
}

const dkgTemplate = `
# DKGSpec
type                 = "offchainreporting2"
schemaVersion        = 1
name                 = "ocr2"
maxTaskDuration      = "30s"
contractID           = "%s"
ocrKeyBundleID       = "%s"
p2pv2Bootstrappers   = ["%s@127.0.0.1:%s"]
relay                = "evm"
pluginType           = "dkg"
transmitterID        = "%s"

[relayConfig]
chainID              = %d

[pluginConfig]
EncryptionPublicKey  = "%s"
KeyID                = "%s"
SigningPublicKey     = "%s"
`

const bootstrapTemplate = `
type                               = "bootstrap"
schemaVersion                      = 1
name                               = ""
id                                 = "1"
contractID                         = "%s"
relay                              = "evm"

[relayConfig]
chainID                            = %d
`

func (cli *Client) ConfigureDKGNode(c *clipkg.Context) (*SetupDKGNodePayload, error) {
	lggr := cli.Logger.Named("SetupDKGJob")
	err := cli.Config.Validate()
	if err != nil {
		return nil, cli.errorOut(errors.Wrap(err, "config validation failed"))
	}
	lggr.Infow(fmt.Sprintf("Configuring Chainlink Node FOR DKG %s at commit %s", static.Version, static.Sha), "Version", static.Version, "SHA", static.Sha)

	ldb := pg.NewLockedDB(cli.Config, lggr)
	rootCtx, _ := context.WithCancel(context.Background())

	if err = ldb.Open(rootCtx); err != nil {
		return nil, cli.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfClosing(ldb, "db")

	app, err := cli.AppFactory.NewApplication(cli.Config, ldb.DB())
	if err != nil {
		return nil, cli.errorOut(errors.Wrap(err, "fatal error instantiating application"))
	}

	// Initialize keystore and generate keys.
	keyStore := app.GetKeyStore()
	err = setupDKGKeystore(cli, c, app, keyStore)
	if err != nil {
		return nil, cli.errorOut(err)
	}

	// Get all configuration parameters.
	keyID := c.String("keyID")
	dkgEncrypt, _ := app.GetKeyStore().DKGEncrypt().GetAll()
	dkgSign, _ := app.GetKeyStore().DKGSign().GetAll()
	dkgEncryptKey := dkgEncrypt[0].PublicKeyString()
	dkgSignKey := dkgSign[0].PublicKeyString()
	p2p, _ := app.GetKeyStore().P2P().GetAll()
	ocr2List, _ := app.GetKeyStore().OCR2().GetAll()
	ethKeys, _ := app.GetKeyStore().Eth().GetAll()
	transmitterID := ethKeys[0].Address.String()
	peerID := p2p[0].PeerID().Raw()
	if c.Bool("isBootstrapper") == false {
		peerID = c.String("bootstrapperPeerID")
	}

	// Find the EVM OCR2 bundle.
	var ocr2 ocr2key.KeyBundle
	for _, ocr2Item := range ocr2List {
		if ocr2Item.ChainType() == chaintype.EVM {
			ocr2 = ocr2Item
		}
	}
	if ocr2 == nil {
		return nil, cli.errorOut(errors.Wrap(job.ErrNoSuchKeyBundle, "evm OCR2 key bundle not found"))
	}
	offChainPublicKey := ocr2.OffchainPublicKey()
	configPublicKey := ocr2.ConfigEncryptionPublicKey()

	// Set up bootstrapper job if bootstrapper.
	if c.Bool("isBootstrapper") {
		err = setupBootstrapperJob(cli, c, app)
		if err != nil {
			return nil, cli.errorOut(err)
		}
	}

	// Set up DKG job.
	dkgTemplateArgs := &DKGTemplateArgs{
		contractID:              c.String("contractID"),
		ocrKeyBundleID:          ocr2.ID(),
		p2pv2BootstrapperPeerID: peerID,
		p2pv2BootstrapperPort:   c.String("bootstrapPort"),
		transmitterID:           transmitterID,
		chainID:                 c.Int64("chainID"),
		EncryptionPublicKey:     dkgEncryptKey,
		KeyID:                   keyID,
		SigningPublicKey:        dkgSignKey,
	}
	err = createDKGJob(cli, c, app, *dkgTemplateArgs)
	if err != nil {
		return nil, cli.errorOut(err)
	}

	return &SetupDKGNodePayload{
		OnChainPublicKey:  ocr2.OnChainPublicKey(),
		OffChainPublicKey: hex.EncodeToString(offChainPublicKey[:]),
		ConfigPublicKey:   hex.EncodeToString(configPublicKey[:]),
		PeerID:            p2p[0].PeerID().Raw(),
		Transmitter:       transmitterID,
		DkgEncrypt:        dkgEncryptKey,
		DkgSign:           dkgSignKey,
	}, nil
}

func setupDKGKeystore(cli *Client, c *clipkg.Context, app chainlink.Application, keyStore keystore.Master) error {
	err := cli.KeyStoreAuthenticator.authenticate(c, keyStore)
	if err != nil {
		return errors.Wrap(err, "error authenticating keystore")
	}

	evmChainSet := app.GetChains().EVM
	if cli.Config.EVMEnabled() {
		if err != nil {
			return errors.Wrap(err, "error migrating keystore")
		}

		for _, ch := range evmChainSet.Chains() {
			err = keyStore.Eth().EnsureKeys(ch.ID())
			if err != nil {
				return errors.Wrap(err, "failed to ensure keystore keys")
			}
		}
	}

	err = keyStore.OCR2().EnsureKeys()
	if err != nil {
		return errors.Wrap(err, "failed to ensure ocr key")
	}

	err = keyStore.DKGSign().EnsureKey()
	if err != nil {
		return errors.Wrap(err, "failed to ensure ocr key")
	}

	err = keyStore.DKGEncrypt().EnsureKey()
	if err != nil {
		return errors.Wrap(err, "failed to ensure ocr key")
	}

	err = keyStore.P2P().EnsureKey()
	if err != nil {
		return errors.Wrap(err, "failed to ensure p2p key")
	}

	return nil
}

func setupBootstrapperJob(cli *Client, c *clipkg.Context, app chainlink.Application) error {
	sp := fmt.Sprintf(bootstrapTemplate,
		c.String("contractID"),
		c.Int64("chainID"),
	)
	var jb job.Job
	err := toml.Unmarshal([]byte(sp), &jb)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal job spec")
	}
	var os job.BootstrapSpec
	err = toml.Unmarshal([]byte(sp), &os)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal job spec")
	}
	jb.BootstrapSpec = &os

	err = app.AddJobV2(context.Background(), &jb)
	if err != nil {
		return errors.Wrap(err, "failed to add job")
	}
	fmt.Println(sp)

	// Give a cooldown
	time.Sleep(time.Second)

	return nil
}

func createDKGJob(cli *Client, c *clipkg.Context, app chainlink.Application, args DKGTemplateArgs) error {
	// Set up DKG job if.
	sp := fmt.Sprintf(dkgTemplate,
		args.contractID,
		args.ocrKeyBundleID,
		args.p2pv2BootstrapperPeerID,
		args.p2pv2BootstrapperPort,
		args.transmitterID,
		args.chainID,
		args.EncryptionPublicKey,
		args.KeyID,
		args.SigningPublicKey,
	)

	var jb job.Job
	err := toml.Unmarshal([]byte(sp), &jb)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal job spec")
	}
	var os job.OCR2OracleSpec
	err = toml.Unmarshal([]byte(sp), &os)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal job spec")
	}
	jb.OCR2OracleSpec = &os

	err = app.AddJobV2(context.Background(), &jb)
	if err != nil {
		return errors.Wrap(err, "failed to add job")
	}
	fmt.Println(sp)

	return nil
}
