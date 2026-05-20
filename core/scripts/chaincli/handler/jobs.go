package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

func (k *Keeper) CreateJob(ctx context.Context) {
	k.createJobs(ctx)
}

func (k *Keeper) createJobs(ctx context.Context) {
	lggr, closeLggr := logger.NewLogger()
	logger.Sugared(lggr).ErrorIfFn(closeLggr, "Failed to close logger")

	for i, keeperAddr := range k.cfg.Keepers {
		url := k.cfg.KeeperURLs[i]
		email := k.cfg.KeeperEmails[i]
		if len(email) == 0 {
			email = defaultChainlinkNodeLogin
		}
		pwd := k.cfg.KeeperPasswords[i]
		if len(pwd) == 0 {
			pwd = defaultChainlinkNodePassword
		}

		cl, err := authenticate(ctx, url, email, pwd, lggr)
		if err != nil {
			log.Fatal(err)
		}

		if err = k.createKeeperJob(ctx, cl, k.cfg.RegistryAddress, keeperAddr); err != nil {
			log.Fatal(err)
		}
	}
}

func (k *Keeper) createKeeperJob(ctx context.Context, client cmd.HTTPClient, registryAddr, nodeAddr string) error {
	if !k.cfg.OCR2Keepers {
		return errors.New("legacy keeper jobs are no longer supported; set KEEPER_OCR2=true and configure OCR2 automation")
	}
	if err := k.createOCR2KeeperJob(ctx, client, registryAddr, nodeAddr); err != nil {
		return err
	}
	log.Println("Keeper job has been successfully created in the Chainlink node with address: ", nodeAddr)
	return nil
}

const ocr2keeperJobTemplate = `type = "offchainreporting2"
pluginType = "ocr2automation"
relay = "evm"
name = "ocr2-automation"
forwardingAllowed = false
schemaVersion = 1
contractID = "%s"
contractConfigTrackerPollInterval = "15s"
ocrKeyBundleID = "%s"
transmitterID = "%s"
p2pv2Bootstrappers = [
  "%s"
]

[relayConfig]
chainID = %d

[pluginConfig]
maxServiceWorkers = 100
cacheEvictionInterval = "1s"
contractVersion = "%s"
mercuryCredentialName = "%s"`

func (k *Keeper) createOCR2KeeperJob(ctx context.Context, client cmd.HTTPClient, contractAddr, nodeAddr string) error {
	ocr2KeyConfig, err := getNodeOCR2Config(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to get node OCR2 key bundle ID: %w", err)
	}

	contractVersion := "v2.0"
	if k.cfg.RegistryVersion == config.RegistryVersion2_1 {
		contractVersion = "v2.1"
	}

	request, err := json.Marshal(web.CreateJobRequest{
		TOML: fmt.Sprintf(ocr2keeperJobTemplate,
			contractAddr,
			ocr2KeyConfig.ID,
			nodeAddr,
			k.cfg.BootstrapNodeAddr,
			k.cfg.ChainID,
			contractVersion,
			k.cfg.DataStreamsCredName,
		),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := client.Post(ctx, "/v2/jobs", bytes.NewReader(request))
	if err != nil {
		return fmt.Errorf("failed to create ocr2keeper job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %w", err)
		}

		return fmt.Errorf("unable to create ocr2keeper job: '%v' [%d]", string(body), resp.StatusCode)
	}

	return nil
}
