package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
)

type CustomInternalType struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Components []CustomArg `json:"components,omitempty"`
}

type CustomArg struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Indexed      bool              `json:"indexed,omitempty"`
	Components   []CustomArg       `json:"components,omitempty"`
	InternalType string `json:"internalType,omitempty"`
}


type CustomEvent struct {
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	Parameters []CustomArg `json:"inputs"`
}

type ABIJson []CustomEvent

func parseEvents(abiJSON []byte, eventMap map[string]CustomEvent, fileName string, structMap map[string]CustomInternalType) (map[string]CustomEvent, error) {
	var parsedABI []CustomEvent
	err := json.Unmarshal(abiJSON, &parsedABI)
	if err != nil {
		return nil, err
	}

	eventCount := make(map[string]int)

	for _, item := range parsedABI {
		if item.Type == "event" {
			eventCount[item.Name]++
			eventKey := fmt.Sprintf("%s_%s_%d", fileName, item.Name, eventCount[item.Name])

			// The following loop replaces the call to flattenTupleComponents
			for i, param := range item.Parameters {
				if param.Type == "tuple" && param.InternalType != "" {
					structName := strings.TrimPrefix(param.InternalType, "struct ")
					if _, ok := structMap[structName]; ok {
						item.Parameters[i].Type = structName
					} else {
						return nil, fmt.Errorf("struct %s not found in structMap", structName)
					}
				}
			}

			item.Name = eventKey
			eventMap[eventKey] = item
		} else if item.Type == "struct" {
			structMap[item.Name] = CustomInternalType{
				Name:       item.Name,
				Components: item.Parameters,
			}
		}
	}

	return eventMap, nil
}


func flattenTupleComponents(params []CustomArg) ([]CustomArg, error) {
	flattenedParams := []CustomArg{}

	for _, param := range params {
		if param.Name == "" {
			param.Name = generateRandomParamName()
		}

		if param.Type == "tuple" {
			if param.Components == nil {
				return nil, fmt.Errorf("tuple parameter missing components")
			}
			for _, component := range param.Components {
				flattenedParam := CustomArg{
					Name:    fmt.Sprintf("%s_%s", param.Name, component.Name),
					Type:    component.Type,
					Indexed: param.Indexed,
				}
				flattenedParams = append(flattenedParams, flattenedParam)
			}
		} else {
			flattenedParams = append(flattenedParams, param)
		}
	}

	return flattenedParams, nil
}

func generateRandomParamName() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())

	var paramName string
	for i := 0; i < 6; i++ {
		paramName += string(alphabet[rand.Intn(26)])
	}
	return paramName
}

func generateMockContract(events []CustomEvent) (string, error) {
	const templateStr = `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

contract EventsMock {
	{{- range $eventIndex, $event := . }}
	event {{ $event.Name }}({{- range $paramIndex, $param := $event.Parameters }} {{ $param.Type }} {{ if $param.Indexed }} indexed{{ end }}  {{ $param.Name }}{{- if lt $paramIndex (sub1 (len $event.Parameters)) }},{{ end }}{{ end }});
	{{- end }}

	{{- range $eventIndex, $event := . }}
	function {{ printf "emit%s" $event.Name }}({{- range $paramIndex, $param := $event.Parameters }}{{ $param.Type }}{{ if needsMemoryKeyword $param.Type }} memory{{ end }} {{ $param.Name }}{{- if lt $paramIndex (sub1 (len $event.Parameters)) }},{{- end }}{{ end }}) public {
		emit {{ $event.Name }}({{- range $paramIndex, $param := $event.Parameters }}{{ $param.Name }}{{- if lt $paramIndex (sub1 (len $event.Parameters)) }},{{ end }}{{ end }});
	}
	{{- end }}
}
`
	funcMap := template.FuncMap{
		"add1": func(x int) int {
			return x + 1
		},
		"sub1": func(x int) int {
			return x - 1
		},
		"needsMemoryKeyword": func(paramType string) bool {
			return strings.HasSuffix(paramType, "[]") || paramType == "string" || paramType == "bytes"
		},
	}

	tmpl, err := template.New("mockContract").Funcs(funcMap).Funcs(sprig.TxtFuncMap()).Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, events)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getABIFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".abi") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func main() {
	root := "/Users/gheorghestrimtu/Documents/chainlink/git/chainlink/contracts/solc/v0.8.6" // Change this to the path containing your ABI files
	abiFiles, err := getABIFiles(root)
	if err != nil {
		fmt.Println("Error finding ABI files:", err)
		os.Exit(1)
	}

	eventMap := make(map[string]CustomEvent)

	for _, abiFile := range abiFiles {
		abiJSON, err := ioutil.ReadFile(abiFile)
		if err != nil {
			fmt.Println("Error reading ABI file:", err)
			os.Exit(1)
		}

		fileName := strings.TrimSuffix(filepath.Base(abiFile), ".abi")
		eventMap, err = parseEvents(abiJSON, eventMap, fileName)
		if err != nil {
			fmt.Println("Error parsing events:", err)
			os.Exit(1)
		}
	}

	events := make([]CustomEvent, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, event)
	}

	// Generate the mock contract
	mockContract, err := generateMockContract(events)
	if err != nil {
		fmt.Println("Error generating mock contract:", err)
		os.Exit(1)
	}

	// Save the mock contract to a file
	err = os.WriteFile("/Users/gheorghestrimtu/Documents/chainlink/git/chainlink/contracts/src/v0.8/mocks/EventsMock.sol", []byte(mockContract), 0644)
	if err != nil {
		fmt.Println("Error writing mock contract to a file:", err)
		os.Exit(1)
	}

	fmt.Println("Generated EventsMock.sol mock contract!")
}
