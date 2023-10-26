package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	defaultPath := "./override.toml"
	defaultOutputPath := "./.env"
	overrideconfigPath := flag.String("path", defaultPath, "config for overriding for default test config")
	outputPath := flag.String("output", defaultOutputPath, "output path for env file")

	flag.Parse()
	if *overrideconfigPath == "" {
		overrideconfigPath = &defaultPath
	}
	if *outputPath == "" {
		outputPath = &defaultOutputPath
	}

	cData, err := os.ReadFile(*overrideconfigPath)
	if err != nil {
		log.Println("unable to read the toml at ", *overrideconfigPath, "error - ", err)
		os.Exit(1)
	}

	// convert the data to Base64 encoded string
	encoded := base64.StdEncoding.EncodeToString(cData)
	// set the env var
	if os.Setenv("BASE64_TEST_CONFIG_OVERRIDE", encoded) != nil {
		os.Exit(1)
	}
	// create an env file for the env var
	envFile, err := os.OpenFile(*outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("unable to create .env file - ", err)
		os.Exit(1)
	}
	defer envFile.Close()
	envVars := fmt.Sprintf(`
export BASE64_TEST_CONFIG_OVERRIDE=%s
export TEST_BASE64_TEST_CONFIG_OVERRIDE=%s
`, encoded, encoded)
	_, err = envFile.WriteString(envVars)
	if err != nil {
		log.Println("unable to write to .env file - ", err)
		os.Exit(1)
	}
	fmt.Println("Successfully set the env var BASE64_TEST_CONFIG_OVERRIDE with the contents of ", *overrideconfigPath, "as Base64 encoded string")
}
