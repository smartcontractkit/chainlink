package src

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	defaultArtefactsDir   = "artefacts"
	defaultNodeSetsPath   = ".cache/node_sets.json"
	deployedContractsJSON = "deployed_contracts.json"
)

func mustReadJSON[T any](fileName string) (output T) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Sprintf("failed to open file at %s: %v", fileName, err))
	}
	defer jsonFile.Close()
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read file at %s: %v", fileName, err))
	}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal data: %v", err))
	}
	return
}

func mustWriteJSON[T any](fileName string, data T) {
	jsonFile, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Sprintf("failed to create file at %s: %v", fileName, err))
	}
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", " ")
	err = encoder.Encode(data)
	if err != nil {
		panic(fmt.Sprintf("failed to encode data: %v", err))
	}
}

func ensureArtefactsDir(artefactsDir string) {
	_, err := os.Stat(artefactsDir)
	if err != nil {
		fmt.Println("Creating artefacts directory" + artefactsDir)
		err = os.MkdirAll(artefactsDir, 0700)
		PanicErr(err)
	}
}
