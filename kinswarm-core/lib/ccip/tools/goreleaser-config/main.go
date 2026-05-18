package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	environments := []string{"develop", "devspace"}
	for _, e := range environments {
		cfg := Generate(e)
		data, err := yaml.Marshal(&cfg)
		if err != nil {
			panic(err)
		}
		filename := fmt.Sprintf("../../.goreleaser.%s.yaml", e)
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Generated %s\n", filename)
	}
}
