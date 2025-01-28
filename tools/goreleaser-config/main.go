package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	environments := []string{"develop", "production", "devspace"}
	for _, e := range environments {
		cfg := Generate(e)
		data, err := yaml.Marshal(&cfg)
		if err != nil {
			panic(err)
		}
		filename := fmt.Sprintf("../../.goreleaser.%s.yaml", e)
		err = os.WriteFile(filename, data, 0644) //nolint:gosec // G306
		if err != nil {
			panic(err)
		}
		fmt.Printf("Generated %s\n", filename)
	}
}
