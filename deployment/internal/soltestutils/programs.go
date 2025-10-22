package soltestutils

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

// MCMSProgramIDs is a map of predeployed MCMS Solana program IDs used in tests.
var MCMSProgramIDs = map[string]string{
	"mcm":               "5vNJx78mz7KVMjhuipyr9jKBKcMrKYGdjGkgE4LUmjKk",
	"timelock":          "DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA",
	"access_controller": "6KsN58MTnRQ8FfPaXHiFPPFGDRioikj9CdPvPxZJdCjb",
}

// MCMSPrograms downloads the MCMS program artifacts and returns the path to the cached artifacts
// and the map of program IDs to paths.
//
// This can be used to preload the MCMS program artifacts into a test environment as arguments to
// the WithSolanaContainer function.
//
// TODO: Remove the dependency on the memory package by extracting the download logic into a
// separate solutils package.
func ProgramsForMCMS(t *testing.T) (string, map[string]string) {
	targetDir := t.TempDir()

	// Download the MCMS program artifacts
	memory.DownloadSolanaProgramArtifactsForTest(t)

	// Copy the specific artifacts to the path provided
	for name := range MCMSProgramIDs {
		src := filepath.Join(memory.ProgramsPath, name+".so")
		dst := filepath.Join(targetDir, name+".so")

		// Copy the cached artifacts to the target directory
		srcFile, err := os.Open(src)
		require.NoError(t, err)

		dstFile, err := os.Create(dst)
		require.NoError(t, err)

		_, err = io.Copy(dstFile, srcFile)
		require.NoError(t, err)

		srcFile.Close()
		dstFile.Close()
	}

	// Return the path to the cached artifacts and the map of program IDs to paths
	return targetDir, MCMSProgramIDs
}
