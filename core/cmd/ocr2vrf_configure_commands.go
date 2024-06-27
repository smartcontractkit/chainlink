package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type SetupOCR2VRFNodePayload struct {
	OnChainPublicKey  string
	OffChainPublicKey string
	ConfigPublicKey   string
	PeerID            string
	Transmitter       string
	DkgEncrypt        string
	DkgSign           string
	SendingKeys       []string
}

type dkgTemplateArgs struct {
	contractID              string
	ocrKeyBundleID          string
	p2pv2BootstrapperPeerID string
	p2pv2BootstrapperPort   string
	transmitterID           string
	useForwarder            bool
	chainID                 int64
	encryptionPublicKey     string
	keyID                   string
	signingPublicKey        string
}

type ocr2vrfTemplateArgs struct {
	dkgTemplateArgs
	vrfBeaconAddress      string
	vrfCoordinatorAddress string
	linkEthFeedAddress    string
	sendingKeys           []string
}

const DKGTemplate = `
# DKGSpec
type                 = "offchainreporting2"
schemaVersion        = 1
name                 = "ocr2"
maxTaskDuration      = "30s"
contractID           = "%s"
ocrKeyBundleID       = "%s"
relay                = "evm"
pluginType           = "dkg"
transmitterID        = "%s"
forwardingAllowed    = %t
%s

[relayConfig]
chainID              = %d

[pluginConfig]
EncryptionPublicKey  = "%s"
KeyID                = "%s"
SigningPublicKey     = "%s"
`

const OCR2VRFTemplate = `
type                 = "offchainreporting2"
schemaVersion        = 1
name                 = "ocr2vrf-chainID-%d"
maxTaskDuration      = "30s"
contractID           = "%s"
ocrKeyBundleID       = "%s"
relay                = "evm"
pluginType           = "ocr2vrf"
transmitterID        = "%s"
forwardingAllowed    = %t
%s

[relayConfig]
chainID              = %d
sendingKeys          = [%s]

[pluginConfig]
dkgEncryptionPublicKey = "%s"
dkgSigningPublicKey    = "%s"
dkgKeyID               = "%s"
dkgContractAddress     = "%s"

vrfCoordinatorAddress  = "%s"
linkEthFeedAddress     = "%s"
`

const BootstrapTemplate = `
type                               = "bootstrap"
schemaVersion                      = 1
name                               = "bootstrap-chainID-%d"
id                                 = "1"
contractID                         = "%s"
relay                              = "evm"

[relayConfig]
chainID                            = %d
`

const forwarderAdditionalEOACount = 4

func (s *Shell) ConfigureOCR2VRFNode(c *cli.Context, owner *bind.TransactOpts, ec *ethclient.Client) (*SetupOCR2VRFNodePayload, error) {
	ctx := s.ctx()
	lggr := logger.Sugared(s.Logger.Named("ConfigureOCR2VRFNode"))
	lggr.Infow(
		fmt.Sprintf("Configuring Chainlink Node for job type %s %s at commit %s", c.String("job-type"), static.Version, static.Sha),
		"Version", static.Version, "SHA", static.Sha)

	var pwd, vrfpwd *string
	if passwordFile := c.String("password"); passwordFile != "" {
		p, err := utils.PasswordFromFile(passwordFile)
		if err != nil {
			return nil, errors.Wrap(err, "error reading password from file")
		}
		pwd = &p
	}
	if vrfPasswordFile := c.String("vrfpassword"); len(vrfPasswordFile) != 0 {
		p, err := utils.PasswordFromFile(vrfPasswordFile)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading VRF password from vrfpassword file \"%s\"", vrfPasswordFile)
		}
		vrfpwd = &p
	}

	s.Config.SetPasswords(pwd, vrfpwd)

	err := s.Config.Validate()
	if err != nil {
		return nil, s.errorOut(errors.Wrap(err, "config validation failed"))
	}

	cfg := s.Config
	ldb := pg.NewLockedDB(cfg.AppID(), cfg.Database(), cfg.Database().Lock(), lggr)

	if err = ldb.Open(ctx); err != nil {
		return nil, s.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfFn(ldb.Close, "Error closing db")

	app, err := s.AppFactory.NewApplication(ctx, s.Config, lggr, ldb.DB())
	if err != nil {
		return nil, s.errorOut(errors.Wrap(err, "fatal error instantiating application"))
	}

	chainID := c.Int64("chainID")

	// Initialize keystore and generate keys.
	keyStore := app.GetKeyStore()
	err = setupKeystore(ctx, s, app, keyStore)
	if err != nil {
		return nil, s.errorOut(err)
	}

	// Start application.
	err = app.Start(ctx)
	if err != nil {
		return nil, s.errorOut(err)
	}

	// Close application.
	defer lggr.ErrorIfFn(app.Stop, "Failed to Stop application")

	// Initialize transmitter settings.
	var sendingKeys []string
	var sendingKeysAddresses []common.Address
	useForwarder := c.Bool("use-forwarder")
	ethKeys, err := app.GetKeyStore().Eth().EnabledKeysForChain(ctx, big.NewInt(chainID))
	if err != nil {
		return nil, s.errorOut(err)
	}
	transmitterID := ethKeys[0].Address.String()

	// Populate sendingKeys with current ETH keys.
	for _, k := range ethKeys {
		sendingKeys = append(sendingKeys, k.Address.String())
		sendingKeysAddresses = append(sendingKeysAddresses, k.Address)
	}

	if useForwarder {
		// Add extra sending keys if using a forwarder.
		sendingKeys, sendingKeysAddresses, err = s.appendForwarders(ctx, chainID, app.GetKeyStore().Eth(), sendingKeys, sendingKeysAddresses)
		if err != nil {
			return nil, err
		}
		err = s.authorizeForwarder(c, ldb.DB(), chainID, ec, owner, sendingKeysAddresses)
		if err != nil {
			return nil, err
		}
	}

	// Get all configuration parameters.
	keyID := c.String("keyID")
	dkgEncrypt, _ := app.GetKeyStore().DKGEncrypt().GetAll()
	dkgSign, _ := app.GetKeyStore().DKGSign().GetAll()
	dkgEncryptKey := dkgEncrypt[0].PublicKeyString()
	dkgSignKey := dkgSign[0].PublicKeyString()
	p2p, _ := app.GetKeyStore().P2P().GetAll()
	ocr2List, _ := app.GetKeyStore().OCR2().GetAll()
	peerID := p2p[0].PeerID().Raw()
	if !c.Bool("isBootstrapper") {
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
		return nil, s.errorOut(errors.Wrap(job.ErrNoSuchKeyBundle, "evm OCR2 key bundle not found"))
	}
	offChainPublicKey := ocr2.OffchainPublicKey()
	configPublicKey := ocr2.ConfigEncryptionPublicKey()

	if c.Bool("isBootstrapper") {
		// Set up bootstrapper job if bootstrapper.
		err = createBootstrapperJob(ctx, lggr, c, app)
	} else if c.String("job-type") == "DKG" {
		// Set up DKG job.
		err = createDKGJob(ctx, lggr, app, dkgTemplateArgs{
			contractID:              c.String("contractID"),
			ocrKeyBundleID:          ocr2.ID(),
			p2pv2BootstrapperPeerID: peerID,
			p2pv2BootstrapperPort:   c.String("bootstrapPort"),
			transmitterID:           transmitterID,
			useForwarder:            useForwarder,
			chainID:                 chainID,
			encryptionPublicKey:     dkgEncryptKey,
			keyID:                   keyID,
			signingPublicKey:        dkgSignKey,
		})
	} else if c.String("job-type") == "OCR2VRF" {
		// Set up OCR2VRF job.
		err = createOCR2VRFJob(ctx, lggr, app, ocr2vrfTemplateArgs{
			dkgTemplateArgs: dkgTemplateArgs{
				contractID:              c.String("dkg-address"),
				ocrKeyBundleID:          ocr2.ID(),
				p2pv2BootstrapperPeerID: peerID,
				p2pv2BootstrapperPort:   c.String("bootstrapPort"),
				transmitterID:           transmitterID,
				useForwarder:            useForwarder,
				chainID:                 chainID,
				encryptionPublicKey:     dkgEncryptKey,
				keyID:                   keyID,
				signingPublicKey:        dkgSignKey,
			},
			vrfBeaconAddress:      c.String("vrf-beacon-address"),
			vrfCoordinatorAddress: c.String("vrf-coordinator-address"),
			linkEthFeedAddress:    c.String("link-eth-feed-address"),
			sendingKeys:           sendingKeys,
		})
	} else {
		err = fmt.Errorf("unknown job type: %s", c.String("job-type"))
	}

	if err != nil {
		return nil, err
	}

	return &SetupOCR2VRFNodePayload{
		OnChainPublicKey:  ocr2.OnChainPublicKey(),
		OffChainPublicKey: hex.EncodeToString(offChainPublicKey[:]),
		ConfigPublicKey:   hex.EncodeToString(configPublicKey[:]),
		PeerID:            p2p[0].PeerID().Raw(),
		Transmitter:       transmitterID,
		DkgEncrypt:        dkgEncryptKey,
		DkgSign:           dkgSignKey,
		SendingKeys:       sendingKeys,
	}, nil
}

func (s *Shell) appendForwarders(ctx context.Context, chainID int64, ks keystore.Eth, sendingKeys []string, sendingKeysAddresses []common.Address) ([]string, []common.Address, error) {
	for i := 0; i < forwarderAdditionalEOACount; i++ {
		// Create the sending key in the keystore.
		k, err := ks.Create(ctx)
		if err != nil {
			return nil, nil, err
		}

		// Enable the sending key for the current chain.
		err = ks.Enable(ctx, k.Address, big.NewInt(chainID))
		if err != nil {
			return nil, nil, err
		}

		sendingKeys = append(sendingKeys, k.Address.String())
		sendingKeysAddresses = append(sendingKeysAddresses, k.Address)
	}

	return sendingKeys, sendingKeysAddresses, nil
}

func (s *Shell) authorizeForwarder(c *cli.Context, db *sqlx.DB, chainID int64, ec *ethclient.Client, owner *bind.TransactOpts, sendingKeysAddresses []common.Address) error {
	ctx := s.ctx()
	// Replace the transmitter ID with the forwarder address.
	forwarderAddress := c.String("forwarder-address")

	// We have to set the authorized senders on-chain here, otherwise the job spawner will fail as the
	// forwarder will not be recognized.
	ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()
	f, err := authorized_forwarder.NewAuthorizedForwarder(common.HexToAddress(forwarderAddress), ec)
	if err != nil {
		return err
	}
	tx, err := f.SetAuthorizedSenders(owner, sendingKeysAddresses)
	if err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, ec, tx)
	if err != nil {
		return err
	}

	// Create forwarder for management in forwarder_manager.go.
	orm := forwarders.NewORM(db)
	_, err = orm.CreateForwarder(ctx, common.HexToAddress(forwarderAddress), *ubig.NewI(chainID))
	if err != nil {
		return err
	}

	return nil
}

func setupKeystore(ctx context.Context, cli *Shell, app chainlink.Application, keyStore keystore.Master) error {
	if err := cli.KeyStoreAuthenticator.authenticate(ctx, keyStore, cli.Config.Password()); err != nil {
		return errors.Wrap(err, "error authenticating keystore")
	}

	if cli.Config.EVMEnabled() {
		chains, err := app.GetRelayers().LegacyEVMChains().List()
		if err != nil {
			return fmt.Errorf("failed to get legacy evm chains")
		}
		for _, ch := range chains {
			if err = keyStore.Eth().EnsureKeys(ctx, ch.ID()); err != nil {
				return errors.Wrap(err, "failed to ensure keystore keys")
			}
		}
	}

	var enabledChains []chaintype.ChainType
	if cli.Config.EVMEnabled() {
		enabledChains = append(enabledChains, chaintype.EVM)
	}
	if cli.Config.CosmosEnabled() {
		enabledChains = append(enabledChains, chaintype.Cosmos)
	}
	if cli.Config.SolanaEnabled() {
		enabledChains = append(enabledChains, chaintype.Solana)
	}
	if cli.Config.StarkNetEnabled() {
		enabledChains = append(enabledChains, chaintype.StarkNet)
	}

	if err := keyStore.OCR2().EnsureKeys(ctx, enabledChains...); err != nil {
		return errors.Wrap(err, "failed to ensure ocr key")
	}

	if err := keyStore.DKGSign().EnsureKey(ctx); err != nil {
		return errors.Wrap(err, "failed to ensure dkgsign key")
	}

	if err := keyStore.DKGEncrypt().EnsureKey(ctx); err != nil {
		return errors.Wrap(err, "failed to ensure dkgencrypt key")
	}

	if err := keyStore.P2P().EnsureKey(ctx); err != nil {
		return errors.Wrap(err, "failed to ensure p2p key")
	}

	return nil
}

func createBootstrapperJob(ctx context.Context, lggr logger.Logger, c *cli.Context, app chainlink.Application) error {
	sp := fmt.Sprintf(BootstrapTemplate,
		c.Int64("chainID"),
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

	err = app.AddJobV2(ctx, &jb)
	if err != nil {
		return errors.Wrap(err, "failed to add job")
	}
	lggr.Info("bootstrap spec:", sp)

	// Give a cooldown
	time.Sleep(time.Second)

	return nil
}

func createDKGJob(ctx context.Context, lggr logger.Logger, app chainlink.Application, args dkgTemplateArgs) error {
	sp := fmt.Sprintf(DKGTemplate,
		args.contractID,
		args.ocrKeyBundleID,
		args.transmitterID,
		args.useForwarder,
		fmt.Sprintf(`p2pv2Bootstrappers   = ["%s@127.0.0.1:%s"]`, args.p2pv2BootstrapperPeerID, args.p2pv2BootstrapperPort),
		args.chainID,
		args.encryptionPublicKey,
		args.keyID,
		args.signingPublicKey,
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

	err = app.AddJobV2(ctx, &jb)
	if err != nil {
		return errors.Wrap(err, "failed to add job")
	}
	lggr.Info("dkg spec:", sp)

	return nil
}

func createOCR2VRFJob(ctx context.Context, lggr logger.Logger, app chainlink.Application, args ocr2vrfTemplateArgs) error {
	var sendingKeysString = fmt.Sprintf(`"%s"`, args.sendingKeys[0])
	for x := 1; x < len(args.sendingKeys); x++ {
		sendingKeysString = fmt.Sprintf(`%s,"%s"`, sendingKeysString, args.sendingKeys[x])
	}
	sp := fmt.Sprintf(OCR2VRFTemplate,
		args.chainID,
		args.vrfBeaconAddress,
		args.ocrKeyBundleID,
		args.transmitterID,
		args.useForwarder,
		fmt.Sprintf(`p2pv2Bootstrappers   = ["%s@127.0.0.1:%s"]`, args.p2pv2BootstrapperPeerID, args.p2pv2BootstrapperPort),
		args.chainID,
		sendingKeysString,
		args.encryptionPublicKey,
		args.signingPublicKey,
		args.keyID,
		args.contractID,
		args.vrfCoordinatorAddress,
		args.linkEthFeedAddress,
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

	err = app.AddJobV2(ctx, &jb)
	if err != nil {
		return errors.Wrap(err, "failed to add job")
	}
	lggr.Info("ocr2vrf spec:", sp)

	return nil
}
