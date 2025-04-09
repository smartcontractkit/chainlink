package crecli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/pkg/errors"
)

type CompilationResult struct {
	WorkflowURL string
	ConfigURL   string
}

func CompileWorkflow(creCLICommandPath, workflowFolder string, configFile *string, settingsFile *os.File) (CompilationResult, error) {
	var outputBuffer bytes.Buffer

	// the CLI expects the workflow code to be located in the same directory as its `go.mod`` file. That's why we assume that the file, which
	// contains the entrypoint method is always named `main.go`. This is a limitation of the CLI, which we can't change.

	compileArgs := []string{"workflow", "compile", "-S", settingsFile.Name()}
	if configFile != nil {
		compileArgs = append(compileArgs, "-c", *configFile)
	}
	compileArgs = append(compileArgs, "main.go")
	compileCmd := exec.Command(creCLICommandPath, compileArgs...) // #nosec G204
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

func EncryptSecrets(creCLICommandPath, secretsFile string, settingsFile *os.File) (string, error) {
	return "", errors.New("not implemented")

	// TODO finish this in the scope of https://smartcontract-it.atlassian.net/browse/DX-81
	// commandArgs := []string{"workflow", "secrets", "encrypt", "-S", settingsFile.Name(), "-v", "-s", "secretsFile"}
	// encryptCmd := exec.Command(creCLICommandPath, commandArgs...) // #nosec G204
	// encryptCmd.Stdout = os.Stdout
	// encryptCmd.Stderr = os.Stderr
	// if err := encryptCmd.Start(); err != nil {
	// 	return "", errors.Wrap(err, "failed to start encrypt command")
	// }

	// return "", nil
}
