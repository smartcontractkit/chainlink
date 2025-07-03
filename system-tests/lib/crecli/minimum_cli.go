package crecli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"
)

type UploadType int

const (
	GITS UploadType = iota
	MINIO
	S3
)

var uploadType = map[UploadType]string{
	GITS: "GITS",
	S3:   "S3",
}

func (t UploadType) String() string {
	return uploadType[t]
}

type MinimumCLI interface {
	Compile(
		workflowFilePath string,
		workflowConfigFilePath string,
		workflowSettingsFilePath string,
		outputFileName string,
	) error
	Upload(
		uploadType UploadType,
		wasmFilePath string,
		workflowConfigFilePath string,
		workflowSettingsFilePath string,
	) (*UploadOutput, error)
	Deploy(
		wasmURL string,
		configURL string,
	) error
	Pause (
		workflowID string,
	) error
}

type CLI struct {
	CreCLICommandPath string
}

func NewCreCli(creCLIPath string) *CLI {
	return &CLI{
		CreCLICommandPath: creCLIPath,
	}
}

func (c CLI) Compile(
	workflowFilePath string,
	workflowConfigFilePath string,
	workflowSettingsFilePath string,
	outputFileName string,
) error {
	workflowFolder := filepath.Dir(workflowFilePath)

	//err := createWorkflowSettingsFile(
	//	workflowFolder,
	//	workflowSettingsFilePath,
	//)
	//if err != nil {
	//	return errors.Wrap(err, "failed to create workflow settings file")
	//}

	compileArgs := []string{
		"workflow",
		"compile",
		"-g=false",
		"-o",
		filepath.Base(outputFileName),
	}
	if workflowConfigFilePath != "" {
		compileArgs = append(compileArgs, "-c", workflowConfigFilePath)
	}
	if workflowSettingsFilePath != "" {
		compileArgs = append(compileArgs, "-S", workflowSettingsFilePath)
	}
	compileArgs = append(compileArgs, filepath.Base(workflowFilePath))

	fmt.Printf("%s %#v\n", c.CreCLICommandPath, compileArgs)

	compileCmd := exec.Command(c.CreCLICommandPath, compileArgs...) // #nosec G204

	var outputBuffer bytes.Buffer
	compileCmd.Stdout = &outputBuffer
	compileCmd.Stderr = &outputBuffer
	compileCmd.Dir = workflowFolder

	err := compileCmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start compile command")
	}

	err = compileCmd.Wait()

	fmt.Println("Compile output:\n", outputBuffer.String())

	if err != nil {
		return errors.Wrap(err, "failed to wait for compile command")
	}

	return nil
}

type UploadOutput struct {
	BinaryURL string
	ConfigURL string
}

func createWorkflowSettingsFile(workflowFolder string, workflowSettingsFilePath string) error {
	cliFilePath := filepath.Join(workflowFolder, CRECLISettingsFileName)
	cliFile, err := os.Create(cliFilePath)
	if err != nil {
		return err
	}
	defer func(cliFile *os.File) {
		err := cliFile.Close()
		if err != nil {
			fmt.Printf("error closing CLI (%s) file: %v\n", CRECLISettingsFileName, err)
		}
	}(cliFile)

	settingsFile, err := os.OpenFile(workflowSettingsFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to open workflow settings file")
	}
	defer func(settingsFile *os.File) {
		err := settingsFile.Close()
		if err != nil {
			fmt.Printf("error closing settings file: %v\n", err)
		}
	}(settingsFile)

	settingsFileBytes, err := os.ReadFile(settingsFile.Name())
	if err != nil {
		return errors.Wrap(err, "failed to read workflow settings file")
	}

	_, err = cliFile.Write(settingsFileBytes)
	if err != nil {
		return errors.Wrap(err, "failed to write workflow settings to CLI file")
	}

	fmt.Printf("created CLI settings file: %s\n", cliFilePath)

	return nil
}

func (c CLI) Upload(
	uploadType UploadType,
	wasmFilePath string,
	configPath string,
	workflowSettingsFilePath string,
) (*UploadOutput, error) {
	workflowFolder := filepath.Dir(wasmFilePath)

	err := createWorkflowSettingsFile(
		workflowFolder,
		workflowSettingsFilePath,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create workflow settings file")
	}

	switch uploadType {
	case MINIO:
		uploadArgs := []string{
			"upload",
			"minio",
			"batch",
			"-v",
			"-b",
			s3provider.DefaultBucket,
			"-f",
			filepath.Base(wasmFilePath),
			"-f",
			configPath,
		}

		//if workflowConfigFilePath != "" {
		//	uploadArgs = append(uploadArgs, "-c", workflowConfigFilePath)
		//}

		if workflowSettingsFilePath != "" {
			uploadArgs = append(uploadArgs, "-S", workflowSettingsFilePath)
		}

		fmt.Printf("%s %#v\n", c.CreCLICommandPath, uploadArgs)

		uploadCmd := exec.Command(c.CreCLICommandPath, uploadArgs...) // #nosec G204

		var outputBuffer bytes.Buffer
		uploadCmd.Stdout = &outputBuffer
		uploadCmd.Stderr = &outputBuffer
		uploadCmd.Dir = workflowFolder

		err := uploadCmd.Start()

		if err != nil {
			return nil, errors.Wrap(err, "failed to start upload command")
		}

		err = uploadCmd.Wait()

		fmt.Println("Compile output:\n", outputBuffer.String())

		if err != nil {
			return nil, errors.Wrap(err, "failed to wait for upload command")
		}

		output := &UploadOutput{}

		re := regexp.MustCompile(`URL=([^\s]+)`)
		matches := re.FindAllStringSubmatch(outputBuffer.String(), -1)

		ansiEscapePattern := `\x1b\[[0-9;]*m`
		cleanRe := regexp.MustCompile(ansiEscapePattern)

		fmt.Printf("%#v\n", matches)

		const urlsFound = 2
		if len(matches) == urlsFound {
			output.BinaryURL = cleanRe.ReplaceAllString(matches[0][1], "")
			output.ConfigURL = cleanRe.ReplaceAllString(matches[1][1], "")
		}

		return output, nil

	default:
		return nil, errors.New("unsupported by this interface, yet")
	}
}

func (c CLI) Deploy(wasmURL string, configURL string) error {
	// TODO implement me
	panic("implement me")
}

func Deref[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
