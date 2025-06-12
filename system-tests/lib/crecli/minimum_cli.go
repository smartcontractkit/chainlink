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
	Compile(workflowFilePath string, workflowConfigFilePath string, outputFileName string) error
	Upload(uploadType UploadType, wasmFilePath string, workflowConfigPath string) (*UploadOutput, error)
	Deploy(wasmURL string, configURL string) error
}

type CLI struct {
	CreCLICommandPath string
}

func NewCreCli(creCLIPath string) *CLI {
	return &CLI{
		CreCLICommandPath: creCLIPath,
	}
}

func (c CLI) Compile(workflowFileName string, workflowConfigFilePath string, outputFileName string) error {
	workflowFolder := filepath.Dir(workflowFileName)

	_, err := os.Create(filepath.Join(workflowFolder, CRECLISettingsFileName))
	if err != nil {
		return err
	}

	compileArgs := []string{
		"workflow",
		"compile",
		"-g=false",
		"-c",
		workflowConfigFilePath,
		"-o",
		outputFileName,
	}
	compileArgs = append(compileArgs, workflowFileName)

	fmt.Printf("%s %#v\n", c.CreCLICommandPath, compileArgs)

	compileCmd := exec.Command(c.CreCLICommandPath, compileArgs...) // #nosec G204

	var outputBuffer bytes.Buffer
	compileCmd.Stdout = &outputBuffer
	compileCmd.Stderr = &outputBuffer
	compileCmd.Dir = workflowFolder

	err = compileCmd.Start()
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

func (c CLI) Upload(uploadType UploadType, wasmfilePath string, configPath string) (*UploadOutput, error) {
	workflowFolder := filepath.Dir(wasmfilePath)

	switch uploadType {
	case MINIO:
		args := []string{
			"upload",
			"minio",
			"batch",
			"-v",
			"-b",
			s3provider.DefaultBucket,
			"-f",
			wasmfilePath,
			"-f",
			configPath,
		}

		fmt.Printf("%s %#v\n", c.CreCLICommandPath, args)

		uploadCmd := exec.Command(c.CreCLICommandPath, args...) // #nosec G204

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

		if len(matches) == 2 {
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
