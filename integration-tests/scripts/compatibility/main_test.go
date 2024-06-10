package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMainFunction(t *testing.T) {
	resetEnv := func() {
		os.Args = os.Args[:1]
		if _, err := os.Stat(OutputFile); err == nil {
			_ = os.Remove(OutputFile)
		}
	}

	t.Run("MissingArguments", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "arg1", "arg2", "arg3", "arg4"}
		require.PanicsWithError(t, InsufficientArgsErr, func() { main() })
	})

	t.Run("InvalidDockerImageFormat", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "TestOCR.*", "./smoke/ocr_test.go", "besu", "hyperledger/besu"}
		require.PanicsWithError(t, fmt.Sprintf("docker image format is invalid: %s", "hyperledger/besu"), func() { main() })
	})

	t.Run("FileCreationAndWrite", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "TestOCR.*", "./smoke/ocr_test.go", "besu", "hyperledger/besu:21.0.0,hyperledger/besu:22.0.0"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err := os.ReadFile(OutputFile)
		require.NoError(t, err)

		var output Output
		err = json.Unmarshal(bytes, &output)
		require.NoError(t, err)
		require.Len(t, output.Entries, 2)
		require.Equal(t, "ocr", output.Entries[0].Product)
		require.Equal(t, "TestOCR.*", output.Entries[0].TestRegex)
		require.Equal(t, "./smoke/ocr_test.go", output.Entries[0].File)
		require.Equal(t, "besu", output.Entries[0].EthImplementationName)
		require.Equal(t, "hyperledger/besu:21.0.0", output.Entries[0].DockerImage)
	})

	t.Run("AppendToFile", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "TestOCR.*", "./smoke/ocr_test.go", "besu", "hyperledger/besu:21.0.0,hyperledger/besu:22.0.0"}
		require.NotPanics(t, func() { main() })

		os.Args = []string{"main", "ocr", "TestOCR.*", "./smoke/ocr_test.go", "geth", "ethereum/client-go:1.10.0"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err := os.ReadFile(OutputFile)
		require.NoError(t, err)

		var output Output
		err = json.Unmarshal(bytes, &output)
		require.NoError(t, err)
		require.Len(t, output.Entries, 3)
	})

	t.Run("OverwriteFile", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "TestOCR.*", "./smoke/ocr_test.go", "besu", "hyperledger/besu:21.0.0,hyperledger/besu:22.0.0"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err := os.ReadFile(OutputFile)
		require.NoError(t, err)
		var initialOutput Output
		err = json.Unmarshal(bytes, &initialOutput)
		require.NoError(t, err)
		require.Len(t, initialOutput.Entries, 2)

		os.Args = []string{"main", "ocr", "TestOCR.*", "./smoke/ocr_test.go", "besu", "hyperledger/besu:22.0.0,hyperledger/besu:23.0.0"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err = os.ReadFile(OutputFile)
		require.NoError(t, err)

		var output Output
		err = json.Unmarshal(bytes, &output)
		require.NoError(t, err)
		require.Len(t, output.Entries, 4)
		require.Equal(t, "hyperledger/besu:23.0.0", output.Entries[3].DockerImage)
	})

	t.Run("EmptyProduct", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "", "TestOCR.*", "./smoke/ocr_test.go", "besu", "hyperledger/besu:21.0.0"}
		require.PanicsWithError(t, fmt.Sprintf(EmptyParameterErr, "product"), func() { main() })
	})

	t.Run("EmptyTestRegex", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "", "./smoke/ocr_test.go", "besu", "hyperledger/besu:21.0.0"}
		require.PanicsWithError(t, fmt.Sprintf(EmptyParameterErr, "test_regex"), func() { main() })
	})

	t.Run("InvalidTestRegex", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "[invalid", "./smoke/ocr_test.go", "besu", "hyperledger/besu:21.0.0"}
		require.PanicsWithError(t, fmt.Sprintf("failed to compile regex: %v", "error parsing regexp: missing closing ]: `[invalid`"), func() { main() })
	})

	t.Run("EmptyFile", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "TestOCR.*", "", "besu", "hyperledger/besu:21.0.0"}
		require.PanicsWithError(t, fmt.Sprintf(EmptyParameterErr, "file"), func() { main() })
	})

	t.Run("EmptyEthImplementation", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "TestOCR.*", "./smoke/ocr_test.go", "", "hyperledger/besu:21.0.0"}
		require.PanicsWithError(t, fmt.Sprintf(EmptyParameterErr, "eth_implementation"), func() { main() })
	})

	t.Run("EmptyDockerImages", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "ocr", "TestOCR.*", "./smoke/ocr_test.go", "besu", ""}
		require.PanicsWithError(t, fmt.Sprintf(EmptyParameterErr, "docker_images"), func() { main() })
	})

	defer func() { resetEnv() }()
}
