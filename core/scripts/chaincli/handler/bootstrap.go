package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

const (
	bootstrapJobSpec = `type = "bootstrap"
schemaVersion = 1
name = "ocr2keeper bootstrap node"
contractID = "%s"
relay = "evm"

[relayConfig]
chainID = %d`

	bootstrapTOML = `[P2P]
[P2P.V2]
ListenAddresses = ["0.0.0.0:%s"]`
)

// StartBootstrapNode starts the ocr2 bootstrap node with the given contract address, returns the tcp address of the node
func (h *baseHandler) StartBootstrapNode(ctx context.Context, addr string, uiPort, p2pv2Port int, force bool) string {
	lggr, closeLggr := logger.NewLogger()
	logger.Sugared(lggr).ErrorIfFn(closeLggr, "Failed to close logger")

	const containerName = "bootstrap"

	urlRaw, _, err := h.launchChainlinkNode(
		ctx,
		uiPort,
		containerName,
		fmt.Sprintf(bootstrapTOML, strconv.Itoa(p2pv2Port)),
		force,
	)
	if err != nil {
		lggr.Fatal("Failed to launch chainlink node, ", err)
	}

	cl, err := authenticate(ctx, urlRaw, defaultChainlinkNodeLogin, defaultChainlinkNodePassword, lggr)
	if err != nil {
		lggr.Fatal("Authentication failed, ", err)
	}

	p2pKeyID, err := getP2PKeyID(ctx, cl)
	if err != nil {
		lggr.Fatal("Failed to get P2P key ID, ", err)
	}

	if err = h.createBootstrapJob(ctx, cl, addr); err != nil {
		lggr.Fatal("Failed to create keeper job: ", err)
	}

	tcpAddr := fmt.Sprintf("%s@%s:%d", p2pKeyID, containerName, p2pv2Port)
	lggr.Info("Bootstrap job has been successfully created in the Chainlink node with address ", urlRaw, ", tcp: ", tcpAddr)

	return tcpAddr
}

// createBootstrapJob creates a bootstrap job in the chainlink node by the given address
func (h *baseHandler) createBootstrapJob(ctx context.Context, client cmd.HTTPClient, contractAddr string) error {
	request, err := json.Marshal(web.CreateJobRequest{
		TOML: fmt.Sprintf(bootstrapJobSpec, contractAddr, h.cfg.ChainID),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := client.Post(ctx, "/v2/jobs", bytes.NewReader(request))
	if err != nil {
		return fmt.Errorf("failed to create bootstrap job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %w", err)
		}

		return fmt.Errorf("unable to create bootstrap job: '%v' [%d]", string(body), resp.StatusCode)
	}
	log.Println("Bootstrap job has been successfully created in the Chainlink node")
	return nil
}
