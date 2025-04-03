package crecli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type CompilationResult struct {
	WorkflowURL string
	ConfigURL   string
	SecretsURL  string
}

func SetFeedAdmin(creCLICommandPath string, chainID int, adminAddress common.Address, settingsFile *os.File) error {
	setFeedAdminCmd := exec.Command(creCLICommandPath, "-S", settingsFile.Name(), "df", "set-feed-admin", "--chain-id", strconv.Itoa(chainID), "--feed-admin", adminAddress.Hex()) // #nosec G204
	var outputBuffer bytes.Buffer
	setFeedAdminCmd.Stdout = &outputBuffer
	setFeedAdminCmd.Stderr = &outputBuffer
	if err := setFeedAdminCmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start DF set feed admin command")
	}

	waitErr := setFeedAdminCmd.Wait()
	fmt.Println("Set Feed Admin output:\n", outputBuffer.String())
	if waitErr != nil {
		return errors.Wrap(waitErr, "failed to wait for compile command")
	}

	return nil
}

func SetFeedConfig(creCLICommandPath, feedID, feedDecimals, feedDescription string, chainID int, allowedSenders, allowedWorkflowOwners []common.Address, allowedWorkflowNames []string, settingsFile *os.File) error {
	allowedSendersHex := make([]string, len(allowedSenders))
	for i, addr := range allowedSenders {
		allowedSendersHex[i] = addr.Hex()
	}
	allowedSendersStr := strings.Join(allowedSendersHex, ",")

	allowedWorkflowOwnersHex := make([]string, len(allowedWorkflowOwners))
	for i, addr := range allowedWorkflowOwners {
		allowedWorkflowOwnersHex[i] = addr.Hex()
	}
	allowedWorkflowOwnersStr := strings.Join(allowedWorkflowOwnersHex, ",")

	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	feedLength := len(cleanFeedID)

	if feedLength < 32 {
		return errors.Errorf("feed ID must be at least 32 characters long, but was %d", feedLength)
	}

	if feedLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	setFeedConfigCmd := exec.Command(creCLICommandPath,
		"-S", settingsFile.Name(),
		"df",
		"set-feed-config",
		"--chain-id", strconv.Itoa(chainID),
		"--allowed-senders", allowedSendersStr,
		"--allowed-workflow-owners", allowedWorkflowOwnersStr,
		"--allowed-workflow-names", strings.Join(allowedWorkflowNames, ","),
		"--data-id", cleanFeedID,
		"--decimals-arr", fmt.Sprintf("[%s]", feedDecimals),
		"--description", feedDescription,
	) // #nosec G204

	var outputBuffer bytes.Buffer
	setFeedConfigCmd.Stdout = &outputBuffer
	setFeedConfigCmd.Stderr = &outputBuffer
	if err := setFeedConfigCmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start DF set feed config command")
	}

	waitErr := setFeedConfigCmd.Wait()
	fmt.Println("Set Feed Config output:\n", outputBuffer.String())
	if waitErr != nil {
		return errors.Wrap(waitErr, "failed to wait for compile command")
	}

	return nil
}

func CompileWorkflow(creCLICommandPath, workflowFolder string, configFile, settingsFile *os.File) (CompilationResult, error) {
	var outputBuffer bytes.Buffer

	// the CLI expects the workflow code to be located in the same directory as its `go.mod`` file. That's why we assume that the file, which
	// contains the entrypoint method is always named `main.go`. This is a limitation of the CLI, which we can't change.
	compileCmd := exec.Command(creCLICommandPath, "workflow", "compile", "-S", settingsFile.Name(), "-c", configFile.Name(), "main.go") // #nosec G204
	compileCmd.Stdout = &outputBuffer
	compileCmd.Stderr = &outputBuffer
	compileCmd.Dir = workflowFolder
	err := compileCmd.Start()
	if err != nil {
		return CompilationResult{}, errors.Wrap(err, "failed to start compile command")
	}

	err = compileCmd.Wait()
	fmt.Println("Compile output:\n", outputBuffer.String())
	if err != nil {
		return CompilationResult{}, errors.Wrap(err, "failed to wait for compile command")
	}

	re := regexp.MustCompile(`Gist URL=([^\s]+)`)
	matches := re.FindAllStringSubmatch(outputBuffer.String(), -1)
	if len(matches) < 2 {
		return CompilationResult{}, errors.New("failed to find gist URLs in compile output")
	}

	ansiEscapePattern := `\x1b\[[0-9;]*m`
	re = regexp.MustCompile(ansiEscapePattern)

	workflowGistURL := re.ReplaceAllString(matches[0][1], "")
	workflowConfigURL := re.ReplaceAllString(matches[1][1], "")

	if workflowGistURL == "" || workflowConfigURL == "" {
		return CompilationResult{}, errors.New("failed to find gist URLs in compile output")
	}

	return CompilationResult{
		WorkflowURL: workflowGistURL,
		ConfigURL:   workflowConfigURL,
	}, nil
}

// Same command to register a workflow or update an existing one
func DeployWorkflow(creCLICommandPath, workflowName, workflowURL, configURL string, settingsFile *os.File) error {
	deployCmd := exec.Command(creCLICommandPath, "workflow", "deploy", workflowName, "-b", workflowURL, "-c", configURL, "-S", settingsFile.Name(), "-v") // #nosec G204
	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr
	if err := deployCmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start register command")
	}

	return nil
}
