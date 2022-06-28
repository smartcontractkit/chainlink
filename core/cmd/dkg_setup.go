package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/shutdown"
	"github.com/smartcontractkit/chainlink/core/static"
	clipkg "github.com/urfave/cli"
)

type RunNodeDKGPayload struct {
	OnChainPublicKey  string
	OffChainPublicKey string
	ConfigPublicKey   string
	PeerID            string
	Transmitter       string
	DkgEncrypt        string
	DkgSign           string
}

func (cli *Client) RunNodeDKG(c *clipkg.Context) (*RunNodeDKGPayload, error) {
	lggr := cli.Logger.Named("RunNode")

	err := cli.Config.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "config validation failed")
	}

	lggr.Infow(fmt.Sprintf("Configuring Chainlink Node FOR DKG %s at commit %s", static.Version, static.Sha), "Version", static.Version, "SHA", static.Sha)

	ldb := pg.NewLockedDB(cli.Config, lggr)

	// rootCtx will be cancelled when SIGINT|SIGTERM is received
	rootCtx, cancelRootCtx := context.WithCancel(context.Background())

	// cleanExit is used to skip "fail fast" routine
	cleanExit := make(chan struct{})
	var shutdownStartTime time.Time
	defer func() {
		close(cleanExit)
		if !shutdownStartTime.IsZero() {
			log.Printf("Graceful shutdown time: %s", time.Since(shutdownStartTime))
		}
	}()

	go shutdown.HandleShutdown(func(sig string) {
		lggr.Infof("Shutting down due to %s signal received...", sig)

		shutdownStartTime = time.Now()
		cancelRootCtx()

		select {
		case <-cleanExit:
			return
		case <-time.After(cli.Config.ShutdownGracePeriod()):
		}

		lggr.Criticalf("Shutdown grace period of %v exceeded, closing DB and exiting...", cli.Config.ShutdownGracePeriod())
		// LockedDB.Close() will release DB locks and close DB connection
		// Executing this explicitly because defers are not executed in case of os.Exit()
		if err = ldb.Close(); err != nil {
			lggr.Criticalf("Failed to close LockedDB: %v", err)
		}
		if err = cli.CloseLogger(); err != nil {
			log.Printf("Failed to close Logger: %v", err)
		}

		os.Exit(-1)
	})

	if err = ldb.Open(rootCtx); err != nil {
		return nil, cli.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfClosing(ldb, "db")

	app, err := cli.AppFactory.NewApplication(cli.Config, ldb.DB())
	if err != nil {
		return nil, cli.errorOut(errors.Wrap(err, "fatal error instantiating application"))
	}

	keyStore := app.GetKeyStore()
	err = cli.KeyStoreAuthenticator.authenticate(c, keyStore)
	if err != nil {
		return nil, errors.Wrap(err, "error authenticating keystore")
	}

	evmChainSet := app.GetChains().EVM
	if cli.Config.EVMEnabled() {
		if err != nil {
			return nil, errors.Wrap(err, "error migrating keystore")
		}

		for _, ch := range evmChainSet.Chains() {
			err2 := app.GetKeyStore().Eth().EnsureKeys(ch.ID())
			if err2 != nil {
				return nil, errors.Wrap(err2, "failed to ensure keystore keys")
			}
		}
	}

	err2 := app.GetKeyStore().OCR2().EnsureKeys()
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to ensure ocr key")
	}

	err2 = app.GetKeyStore().DKGSign().EnsureKey()
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to ensure ocr key")
	}

	err2 = app.GetKeyStore().DKGEncrypt().EnsureKey()
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to ensure ocr key")
	}

	err2 = app.GetKeyStore().P2P().EnsureKey()
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to ensure p2p key")
	}

	dkgEncrypt, _ := app.GetKeyStore().DKGEncrypt().GetAll()
	dkgSign, _ := app.GetKeyStore().DKGSign().GetAll()
	p2p, _ := app.GetKeyStore().P2P().GetAll()
	ocr2List, _ := app.GetKeyStore().OCR2().GetAll()
	ethKeys, _ := app.GetKeyStore().Eth().GetAll()

	peerID := p2p[0].PeerID().Raw()
	if c.Bool("isBootstrapper") == false {
		peerID = c.String("bootstrapperPeerID")
	}
	ocr2 := ocr2List[0]
	offChainPublicKey := ocr2.OffchainPublicKey()
	configPublicKey := ocr2.ConfigEncryptionPublicKey()
	keyID := c.String("keyID")
	transmitter := ethKeys[0].Address.String()
	dkgEncryptKey := dkgEncrypt[0].PublicKeyString()
	dkgSignKey := dkgSign[0].PublicKeyString()

	if c.Bool("isBootstrapper") {
		sp := fmt.Sprintf(`
		type = "bootstrap"
		schemaVersion = 1
		name = ""
		externalJobID = "4a48cc1b-9091-465b-9f78-2539341009d1"
		id = "3"
		contractID = "%s"
		relay = "evm"
		monitoringEndpoint = "chain.link:4321"
		
		[relayConfig]
		chainID = 4
		`,
			c.String("contractID"),
		)
		var jb job.Job
		err2 = toml.Unmarshal([]byte(sp), &jb)
		if err2 != nil {
			return nil, errors.Wrap(err2, "failed to unmarshal job spec")
		}
		var os job.BootstrapSpec
		err = toml.Unmarshal([]byte(sp), &os)
		if err2 != nil {
			return nil, errors.Wrap(err2, "failed to unmarshal job spec")
		}
		jb.BootstrapSpec = &os

		err2 = app.AddJobV2(rootCtx, &jb)
		if err2 != nil {
			return nil, errors.Wrap(err2, "failed to add job")
		}
		fmt.Println(sp)

		time.Sleep(time.Second)
	}

	sp := fmt.Sprintf(`
			type = "offchainreporting2"
			schemaVersion = 1
			name = "ocr2"
			externalJobID = "6d46d85f-d38c-4f4a-9f00-ac29a25b6330"
			maxTaskDuration = "30s"
			contractID = "%s"
			ocrKeyBundleID = "%s"
			p2pv2Bootstrappers = [
			"%s@127.0.0.1:%s"
			]
			relay = "evm"
			pluginType = "dkg"
			transmitterID = "%s"

			[relayConfig]
			chainID = 4

			[pluginConfig]
			EncryptionPublicKey = "%s"
			KeyID = "%s"
			SigningPublicKey = "%s"
			`,
		c.String("contractID"),
		ocr2.ID(),
		peerID,
		c.String("port"),
		transmitter,
		dkgEncryptKey,
		keyID,
		dkgSignKey,
	)

	var jb job.Job
	err2 = toml.Unmarshal([]byte(sp), &jb)
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to unmarshal job spec")
	}
	var os job.OCR2OracleSpec
	err = toml.Unmarshal([]byte(sp), &os)
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to unmarshal job spec")
	}
	jb.OCR2OracleSpec = &os

	err2 = app.AddJobV2(rootCtx, &jb)
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to add job")
	}
	fmt.Println(sp)

	return &RunNodeDKGPayload{
		OnChainPublicKey:  ocr2.OnChainPublicKey(),
		OffChainPublicKey: hex.EncodeToString(offChainPublicKey[:]),
		ConfigPublicKey:   hex.EncodeToString(configPublicKey[:]),
		PeerID:            peerID,
		Transmitter:       transmitter,
		DkgEncrypt:        dkgEncryptKey,
		DkgSign:           dkgSignKey,
	}, nil
}
