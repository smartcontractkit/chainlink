//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

func main() {
	// variations := []string{ "develop", "devspace", "ci", "prod"}
	variations := []string{ "develop"}
	for _, v := range variations {
		cfg := Generate(v)
		data, err := yaml.Marshal(&cfg)
		if err != nil {
			panic(err)
		}
		filename := fmt.Sprintf("goreleaser.%s.yaml", v)
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Generated %s\n", filename)
	}
}
