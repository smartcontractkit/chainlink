package main

import (
	"flag"
	"fmt"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	evm2 "github.com/smartcontractkit/smart-contract-spec/internal/gen/evm"
	utils2 "github.com/smartcontractkit/smart-contract-spec/internal/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	contracts           = flag.String("contracts", "contracts", "comma-separated list of directories containing EVM smart contracts source and ABI files; must be set")
	output              = flag.String("output", "generated/evm/bindings", "output folder for the generated code")
	clean               = flag.Bool("clean", false, "output folder for the generated code")
	silentIfNoContracts = flag.Bool("silent-if-no-contracts", false, "do not fails if there are not contracts to be processed")
	verbose             = flag.Bool("verbose", false, "generates debugging output")
)

type ExitStatus int

const (
	Success                    ExitStatus = iota // 0
	DisplayHelp                                  // 1
	ContractOrOutputEmpty                        // 2
	FailedToListContractFolder                   // 3
	FailedLoadingAbiFile
	FailedCleaningOrCreatingOutputDir
	FailedGeneratingCRConfig
	FailedGeneratingContractBindingCode
	FailedGeneratingContractBindingFile
	FailedGeneratingCRCWConfigCode
	FailedFormatingFile
	NoInputContracts
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of evm-bindings:\n")
	fmt.Fprintf(os.Stderr, "\tevm-bindings[flags] -contracts T [directory]\n")
	flag.PrintDefaults()
}

func debugLogging(text string, args ...interface{}) {
	if *verbose {
		fmt.Printf(text, args)
	}
}

func logError(err error, exitStatus ExitStatus) {
	fmt.Printf("Error: ", err)
	os.Exit(int(exitStatus))
}

func main() {
	parseArguments()
	validateInputs()
	contractFolders, goModRelatedOutputDir := processInputs()
	contracts := processContracts(contractFolders)
	failIfNoContractsAndFlagSet(contracts)
	err := prepareOutputDirectory(contracts)
	chainReaderConfig, chainWriterConfig := processChainReaderChainWriterConfig(err, contracts)
	outputDirs := strings.Split(*output, "/")
	packageName := outputDirs[len(outputDirs)-1]
	generatedContracts := map[string][]byte{}
	writeGoCode(contracts, packageName, generatedContracts, err, goModRelatedOutputDir, chainReaderConfig, chainWriterConfig)
	formatCode(goModRelatedOutputDir)
}

func parseArguments() {
	log.SetFlags(0)
	log.SetPrefix("evm-bindings: ")
	flag.Usage = Usage
	flag.Parse()
}

func formatCode(goModRelatedOutputDir string) {
	// Format the file using gofmt
	cmd := exec.Command("gofmt", "-w", goModRelatedOutputDir)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error formatting file with gofmt -w %s\n", goModRelatedOutputDir)
		logError(err, FailedFormatingFile)
	}
}

func writeGoCode(contracts map[string]evm2.CodeDetails, packageName string, generatedContracts map[string][]byte, err error, goModRelatedOutputDir string, chainReaderConfig types.ChainReaderConfig, chainWriterConfig types.ChainWriterConfig) {
	for contractName, contractDetails := range contracts {
		generatedContractContent, err := evm2.GenerateContractBinding(packageName, contractDetails)
		if err != nil {
			fmt.Printf("Failed generating contract binding code for %s", contractName)
			logError(err, FailedGeneratingContractBindingCode)
		}
		generatedContracts[contractName] = generatedContractContent
	}

	for contractName, generatedContractContent := range generatedContracts {
		fmt.Printf("Generating go file for: %s\n", contractName)
		contractFileName := utils2.CamelToSnake(contractName)
		err = utils2.GenerateFile(goModRelatedOutputDir, contractFileName, "go", generatedContractContent)
		if err != nil {
			fmt.Printf("Error generating contract biding file %s: %s\n", contractFileName)
			logError(err, FailedGeneratingContractBindingFile)
		}
	}

	configFactoryContent, err := evm2.GenerateChainReaderAndWriterFactory(packageName, chainReaderConfig, chainWriterConfig)
	if err != nil {
		fmt.Printf("Failed generating chain reader and writer config\n")
		logError(err, FailedGeneratingCRCWConfigCode)
	}

	fmt.Printf("Generating ChainReader and ChainWriter config in %s\n", goModRelatedOutputDir)
	utils2.GenerateFile(goModRelatedOutputDir, "chain_config_factory", "go", configFactoryContent)
}

func processChainReaderChainWriterConfig(err error, contracts map[string]evm2.CodeDetails) (types.ChainReaderConfig, types.ChainWriterConfig) {
	chainReaderConfig, chainWriterConfig, err := evm2.GenerateChainReaderChainWriterConfig(contracts)
	if err != nil {
		fmt.Printf("Failed generating chain reader config: %s\n", err)
		logError(err, FailedGeneratingCRConfig)
	}
	return chainReaderConfig, chainWriterConfig
}

func prepareOutputDirectory(contracts map[string]evm2.CodeDetails) error {
	err := utils2.CreateDirectories([]string{*output}, *clean)
	if err != nil {
		fmt.Printf("Failed creating or cleaning directories %s\n", contracts)
		logError(err, FailedCleaningOrCreatingOutputDir)
	}
	return err
}

func validateInputs() {
	if len(*contracts) == 0 || len(*output) == 0 {
		flag.Usage()
		os.Exit(int(ContractOrOutputEmpty))
	}
}

func processContracts(contractFolders []string) map[string]evm2.CodeDetails {
	contracts := map[string]evm2.CodeDetails{}
	for _, contractFolder := range contractFolders {
		files, err := utils2.ListFiles(contractFolder, "abi")
		if err != nil {
			fmt.Printf("Failed listing directories %s\n", contractFolder)
			logError(err, FailedToListContractFolder)
		}

		for _, file := range files {
			fmt.Printf("Processing contract in: %s\n", file)
			contractName, contractABI, err := utils2.LoadFile(file)
			if err != nil {
				debugLogging("Failed loading ABI file  %s", file)
				logError(err, FailedLoadingAbiFile)
			}
			contractDetails, err := evm2.ConvertABIToCodeDetails(contractName, contractABI)
			contracts[contractName] = contractDetails
		}
	}
	return contracts
}

func failIfNoContractsAndFlagSet(contracts map[string]evm2.CodeDetails) {
	if len(contracts) == 0 && !*silentIfNoContracts {
		fmt.Printf("Failing since silent-if-no-contracts is not set and there are no contracts in the input folder")
		os.Exit(int(NoInputContracts))
	} else if len(contracts) == 0 {
		fmt.Printf("No input contracts to process")
		os.Exit(int(Success))
	}
}

func processInputs() ([]string, string) {
	debugLogging("Option values:")
	debugLogging("contracts: %s\n", *contracts)
	debugLogging("output %s\n", *output)
	debugLogging("clean %s\n", *clean)
	debugLogging("verbose %s\n", *verbose)
	debugLogging("fail-if-no-contracts %s\n", *silentIfNoContracts)
	contractFolders := toAbsolutePathsBasedOnGoModLocation(strings.Split(*contracts, ","))
	fmt.Printf("Procesing contract folders: %s\n", contractFolders)
	goModRelatedOutputDir := toAbsolutePath(*output)
	fmt.Printf("Output folder for generated code: %s\n", goModRelatedOutputDir)
	return contractFolders, goModRelatedOutputDir
}

func toAbsolutePathsBasedOnGoModLocation(paths []string) []string {
	// Create a new slice to store the results with the same length as the input slice
	results := make([]string, len(paths))

	// Iterate over the input slice
	for i, path := range paths {
		results[i] = toAbsolutePath(path)
	}

	return results
}

func toAbsolutePath(path string) string {
	workingDir, err := os.Getwd()

	for !fileExists(workingDir + "/" + "go.mod") {
		workingDir = filepath.Dir(workingDir)
	}
	if err != nil {
		log.Fatal(err)
		os.Exit(9)
	}
	return workingDir + "/" + path
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
