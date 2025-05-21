package crecli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

type CompilationResult struct {
	WorkflowURL string
	ConfigURL   string
}

func CompileWorkflow(creCLICommandPath, workflowFolder, workflowFileName string, configFile *string, workflowSettingsFile, settingsFile *os.File) (CompilationResult, error) {
	var outputBuffer bytes.Buffer

	// the CLI expects the workflow code to be located in the same directory as its `go.mod`` file. That's why we assume that the file, which
	// the CLI also expects `cre.yaml` settings file to be present either in the present directory or any of its parent tree directories.

	cliFile, err := os.Create(filepath.Join(workflowFolder, CRECLISettingsFileName))
	if err != nil {
		return CompilationResult{}, err
	}

	settingsFileBytes, err := os.ReadFile(settingsFile.Name())
	if err != nil {
		return CompilationResult{}, err
	}

	_, err = cliFile.Write(settingsFileBytes)
	if err != nil {
		return CompilationResult{}, err
	}

	compileArgs := []string{"workflow", "compile", "-S", workflowSettingsFile.Name()}
	if configFile != nil {
		compileArgs = append(compileArgs, "-c", *configFile)
	}
	compileArgs = append(compileArgs, workflowFileName)
	compileCmd := exec.Command(creCLICommandPath, compileArgs...) // #nosec G204
	compileCmd.Stdout = &outputBuffer
	compileCmd.Stderr = &outputBuffer
	// the CLI expects the workflow code to be located in the same directory as its `go.mod` file
	compileCmd.Dir = workflowFolder
	err = compileCmd.Start()
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
func DeployWorkflow(creCLICommandPath, workflowURL string, configURL, secretsURL *string, settingsFile *os.File) error {
	commandArgs := []string{"workflow", "deploy", "-b", workflowURL, "-S", settingsFile.Name(), "-v"}
	if configURL != nil {
		commandArgs = append(commandArgs, "-c", *configURL)
	}
	if secretsURL != nil {
		commandArgs = append(commandArgs, "-s", *secretsURL)
	}

	deployCmd := exec.Command(creCLICommandPath, commandArgs...) // #nosec G204
	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr
	if startErr := deployCmd.Start(); startErr != nil {
		return errors.Wrap(startErr, "failed to start deploy command")
	}

	waitErr := deployCmd.Wait()
	if waitErr != nil {
		return errors.Wrap(waitErr, "failed to wait for deploy command")
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

func PauseWorkflow(creCLICommandPath string, settingsFile *os.File) error {
	commandArgs := []string{"workflow", "pause", "-S", settingsFile.Name(), "-v"}

	pauseCmd := exec.Command(creCLICommandPath, commandArgs...) // #nosec G204
	pauseCmd.Stdout = os.Stdout
	pauseCmd.Stderr = os.Stderr
	if startErr := pauseCmd.Start(); startErr != nil {
		return errors.Wrap(startErr, "failed to start pause command")
	}

	waitErr := pauseCmd.Wait()
	if waitErr != nil {
		return errors.Wrap(waitErr, "failed to wait for pause command")
	}

	return nil
}

func ActivateWorkflow(creCLICommandPath string, settingsFile *os.File) error {
	commandArgs := []string{"workflow", "activate", "-S", settingsFile.Name(), "-v"}

	activateCmd := exec.Command(creCLICommandPath, commandArgs...) // #nosec G204
	activateCmd.Stdout = os.Stdout
	activateCmd.Stderr = os.Stderr
	if startErr := activateCmd.Start(); startErr != nil {
		return errors.Wrap(startErr, "failed to start activate command")
	}

	waitErr := activateCmd.Wait()
	if waitErr != nil {
		return errors.Wrap(waitErr, "failed to wait for activate command")
	}

	return nil
}
