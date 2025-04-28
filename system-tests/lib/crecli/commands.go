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
}

func CompileWorkflow(creCLICommandPath, workflowFolder, workflowFileName string, configFile *string, settingsFile *os.File) (CompilationResult, error) {
	var outputBuffer bytes.Buffer

	compileArgs := []string{"workflow", "compile", "-S", settingsFile.Name()}
	if configFile != nil {
		compileArgs = append(compileArgs, "-c", *configFile)
	}
	compileArgs = append(compileArgs, workflowFileName)
	compileCmd := exec.Command(creCLICommandPath, compileArgs...) // #nosec G204
	compileCmd.Stdout = &outputBuffer
	compileCmd.Stderr = &outputBuffer
	// the CLI expects the workflow code to be located in the same directory as its `go.mod` file
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

	ansiEscapePattern := `\x1b\[[0-9;]*m`
	re = regexp.MustCompile(ansiEscapePattern)

	result := CompilationResult{}

	expectedGistURLs := 1
	if configFile != nil {
		expectedGistURLs++
	}

	switch len(matches) {
	case 1:
		result.WorkflowURL = re.ReplaceAllString(matches[0][1], "")
	case 2:
		result.WorkflowURL = re.ReplaceAllString(matches[0][1], "")
		result.ConfigURL = re.ReplaceAllString(matches[1][1], "")
	default:
		return CompilationResult{}, errors.New("unsupported number of gist URLs in compile output")
	}

	if len(matches) != expectedGistURLs {
		return CompilationResult{}, fmt.Errorf("unexpected number of gist URLs in compile output: %d, expected %d", len(matches), expectedGistURLs)
	}

	return result, nil
}

// Same command to register a workflow or update an existing one
func DeployWorkflow(creCLICommandPath, workflowName, workflowURL string, configURL, secretsURL *string, settingsFile *os.File) error {
	commandArgs := []string{"workflow", "deploy", workflowName, "-b", workflowURL, "-S", settingsFile.Name(), "-v"}
	if configURL != nil {
		commandArgs = append(commandArgs, "-c", *configURL)
	}
	if secretsURL != nil {
		commandArgs = append(commandArgs, "-s", *secretsURL)
	}

	deployCmd := exec.Command(creCLICommandPath, commandArgs...) // #nosec G204
	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr
	if err := deployCmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start register command")
	}

	return nil
}

func EncryptSecrets(creCLICommandPath, secretsFile string, secrets map[string]string, settingsFile *os.File) (string, error) {
	var outputBuffer bytes.Buffer

	commandArgs := []string{"secrets", "encrypt", "-S", settingsFile.Name(), "-v", "-s", secretsFile}
	encryptCmd := exec.Command(creCLICommandPath, commandArgs...) // #nosec G204
	encryptCmd.Stdout = &outputBuffer
	encryptCmd.Stderr = &outputBuffer

	// Preserve existing environment variables
	encryptCmd.Env = os.Environ()

	// set all secrets as environment variables, so that "encrypt" command can pick them up
	for name, value := range secrets {
		encryptCmd.Env = append(encryptCmd.Env, fmt.Sprintf("%s=%s", name, value))
	}
	if err := encryptCmd.Start(); err != nil {
		return "", errors.Wrap(err, "failed to start encrypt command")
	}

	err := encryptCmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, "failed to wait for encrypt command")
	}

	re := regexp.MustCompile(`Gist URL=([^\s]+)`)
	matches := re.FindAllStringSubmatch(outputBuffer.String(), -1)

	if len(matches) != 1 {
		return "", fmt.Errorf("unexpected number of gist URLs in encrypt output: %d, expected 1", len(matches))
	}

	ansiEscapePattern := `\x1b\[[0-9;]*m`
	re = regexp.MustCompile(ansiEscapePattern)

	return re.ReplaceAllString(matches[0][1], ""), nil
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
