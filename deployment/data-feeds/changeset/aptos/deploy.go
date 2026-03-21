package aptos

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"

	modulefeeds "github.com/smartcontractkit/chainlink-aptos/bindings/data_feeds"
	moduleplatform "github.com/smartcontractkit/chainlink-aptos/bindings/platform"
	moduleplatform_secondary "github.com/smartcontractkit/chainlink-aptos/bindings/platform_secondary"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

const aptosCLIPathEnvVar = "APTOS_CLI_PATH"

func ensureHostAptosCLI() error {
	cliPath, err := firstWorkingAptosCLI(aptosCLICandidates())
	if err != nil {
		return err
	}

	cliDir := filepath.Dir(cliPath)
	pathDirs := filepath.SplitList(os.Getenv("PATH"))
	if len(pathDirs) > 0 && filepath.Clean(pathDirs[0]) == cliDir {
		return nil
	}

	if currentPath := os.Getenv("PATH"); currentPath != "" {
		return os.Setenv("PATH", cliDir+string(os.PathListSeparator)+currentPath)
	}
	return os.Setenv("PATH", cliDir)
}

func aptosCLICandidates() []string {
	var candidates []string
	seen := make(map[string]struct{})
	add := func(path string) {
		path = strings.TrimSpace(path)
		if path == "" {
			return
		}
		path = filepath.Clean(path)
		if _, ok := seen[path]; ok {
			return
		}
		seen[path] = struct{}{}
		candidates = append(candidates, path)
	}

	add(os.Getenv(aptosCLIPathEnvVar))
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		add(filepath.Join(dir, "aptos"))
	}

	add("/tmp/aptos-cli-7.8.0/aptos")
	add("/opt/homebrew/bin/aptos")
	add("/usr/local/bin/aptos")
	return candidates
}

func firstWorkingAptosCLI(candidates []string) (string, error) {
	var problems []string
	for _, candidate := range candidates {
		err := validateAptosCLI(candidate)
		if err == nil {
			return candidate, nil
		}

		problems = append(problems, fmt.Sprintf("%s (%v)", candidate, err))
	}

	return "", fmt.Errorf("failed to find a working Aptos CLI; set %s or install a valid aptos binary. Checked: %s", aptosCLIPathEnvVar, strings.Join(problems, ", "))
}

func validateAptosCLI(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("is a directory")
	}
	if info.Mode()&0o111 == 0 {
		return errors.New("is not executable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, path, "--version").CombinedOutput()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("version check failed: %s", strings.TrimSpace(string(out)))
	}

	versionOutput := strings.ToLower(strings.TrimSpace(string(out)))
	if !strings.Contains(versionOutput, "aptos") {
		return fmt.Errorf("unexpected version output: %q", strings.TrimSpace(string(out)))
	}
	return nil
}

func DeployDataFeeds(chain cldf_aptos.Chain, owner aptos.AccountAddress, platform aptos.AccountAddress, secondaryPlatform aptos.AccountAddress, labels []string) (*types.DeployDataFeedsResponse, error) {
	if err := ensureHostAptosCLI(); err != nil {
		return nil, err
	}
	address, pendingTX, feedsModule, err := modulefeeds.DeployToObject(chain.DeployerSigner, chain.Client, owner, platform, owner, secondaryPlatform)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChainlinkDataFeeds: %w", err)
	}

	tx, err := chain.Client.WaitForTransaction(pendingTX.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkDataFeeds: %s, %w", pendingTX.Hash, err)
	}

	if !tx.Success {
		return nil, errors.New("ChainlinkDataFeeds deployment transaction failed: " + tx.VmStatus)
	}

	// ChainlinkDataFeeds package contracts don't implement typeAndVersion interface, so we have to set it manually
	tvStr := "ChainlinkDataFeeds 1.0.0"
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployDataFeedsResponse{
		Address:  address,
		Tx:       tx.Hash,
		Tv:       tv,
		Contract: &feedsModule,
	}
	return resp, nil
}

func DeployPlatform(chain cldf_aptos.Chain, owner aptos.AccountAddress, labels []string) (*types.DeployPlatformResponse, error) {
	if err := ensureHostAptosCLI(); err != nil {
		return nil, err
	}
	if owner == (aptos.AccountAddress{}) {
		owner = chain.DeployerSigner.AccountAddress()
	}
	address, pendingTX, platformModule, err := moduleplatform.DeployToObject(chain.DeployerSigner, chain.Client, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChainlinkPlatform: %w", err)
	}

	tx, err := chain.Client.WaitForTransaction(pendingTX.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkPlatform: %s, %w", pendingTX.Hash, err)
	}

	if !tx.Success {
		return nil, errors.New("ChainlinkPlatform deployment transaction failed: " + tx.Hash)
	}
	// ChainlinkPlatform package contracts don't implement typeAndVersion interface, so we have to set it manually
	tvStr := "ChainlinkPlatform 1.0.0"
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployPlatformResponse{
		Address:  address,
		Tx:       tx.Hash,
		Tv:       tv,
		Contract: &platformModule,
	}
	return resp, nil
}

func DeployPlatformSecondary(chain cldf_aptos.Chain, owner aptos.AccountAddress, labels []string) (*types.DeployPlatformSecondaryResponse, error) {
	if err := ensureHostAptosCLI(); err != nil {
		return nil, err
	}
	if owner == (aptos.AccountAddress{}) {
		owner = chain.DeployerSigner.AccountAddress()
	}
	address, pendingTX, platformModule, err := moduleplatform_secondary.DeployToObject(chain.DeployerSigner, chain.Client, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChainlinkPlatformSecondary: %w", err)
	}

	tx, err := chain.Client.WaitForTransaction(pendingTX.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkPlatformSecondary: %s, %w", pendingTX.Hash, err)
	}

	if !tx.Success {
		return nil, errors.New("ChainlinkPlatformSecondary deployment transaction failed: " + tx.Hash)
	}
	// ChainlinkPlatformSecondary package contracts don't implement typeAndVersion interface, so we have to set it manually
	tvStr := "ChainlinkPlatformSecondary 1.0.0"
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployPlatformSecondaryResponse{
		Address:  address,
		Tx:       tx.Hash,
		Tv:       tv,
		Contract: &platformModule,
	}
	return resp, nil
}
