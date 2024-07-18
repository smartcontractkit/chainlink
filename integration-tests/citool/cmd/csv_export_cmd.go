package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var csvExportCmd = &cobra.Command{
	Use:   "csvexport",
	Short: "Export tests to CSV format",
	Run: func(cmd *cobra.Command, _ []string) {
		configFile, _ := cmd.Flags().GetString("file")
		if err := exportConfigToCSV(configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	csvExportCmd.Flags().StringP("file", "f", "", "Path to YML file")
	err := csvExportCmd.MarkFlagRequired("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func exportConfigToCSV(configFile string) error {
	// Read the YAML file
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	// Unmarshal the YAML into the Config struct
	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return err
	}

	// Create a CSV file
	file, err := os.Create("output.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV headers
	headers := []string{"ID", "Test Path", "Test Env Type", "Runs On", "Test Cmd", "Test Config Override Required", "Test Secrets Required", "Remote Runner Memory", "Pyroscope Env", "Workflows", "Test Inputs"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Iterate over Tests and write data to CSV
	for _, test := range config.Tests {
		workflows := strings.Join(test.Workflows, ", ") // Combine workflows into a single CSV field
		// Serialize TestInputs
		testInputs := serializeMap(test.TestInputs)

		record := []string{
			test.ID,
			test.Path,
			test.TestEnvType,
			test.RunsOn,
			test.TestCmd,
			fmt.Sprintf("%t", test.TestConfigOverrideRequired),
			fmt.Sprintf("%t", test.TestSecretsRequired),
			test.RemoteRunnerMemory,
			test.PyroscopeEnv,
			workflows,
			testInputs,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func serializeMap(inputs map[string]string) string {
	pairs := make([]string, 0, len(inputs))
	for key, value := range inputs {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(pairs, ", ")
}
