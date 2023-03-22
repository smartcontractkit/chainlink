package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type AbiEvent struct {
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Inputs []AbiInputs `json:"inputs"`
}

type AbiInputs struct {
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Indexed      bool            `json:"indexed"`
	InternalType string          `json:"internalType,omitempty"`
	Components   []AbiComponents `json:"components,omitempty"`
}

type AbiComponents struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	InternalType string `json:"internalType"`
}

type SolEvent struct {
	Name   string
	Params []SolEventParam
}

type SolEventParam struct {
	Type         string
	InternalType string
	Name         string
	Indexed      bool
}

type SolStruct struct {
	Name            string
	SolStructParams []SolStructParam
}

type SolStructParam struct {
	Type string
	Name string
}

type SolFunction struct {
	Name      string
	EventName string
	Params    []SolFunctionParam
}

type SolFunctionParam struct {
	Type   string
	Memory bool
	Name   string
}

type TemplateData struct {
	Events    []SolEvent
	Structs   []SolStruct
	Functions []SolFunction
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

func extractEventsAndStructs(abiJSON []byte, fileName string) ([]SolEvent, []SolStruct, error) {
	var parsedABI []AbiEvent
	err := json.Unmarshal(abiJSON, &parsedABI)
	if err != nil {
		return nil, nil, err
	}

	var solEvents []SolEvent
	var solStructs []SolStruct

	for _, item := range parsedABI {
		if item.Type == "event" {
			eventName := fmt.Sprintf("%s_%s", fileName, item.Name)
			var eventParams []SolEventParam

			for i, input := range item.Inputs {
				if input.Name == "" {
					input.Name = "param" + fmt.Sprintf("%d", i+1)
				}
				if input.Type == "tuple" && strings.Contains(input.InternalType, "struct") {
					internalType := strings.TrimPrefix(input.InternalType, "struct ")
					if strings.Contains(internalType, ".") {
						internalType = strings.Split(internalType, ".")[1]
					}
					structName := fmt.Sprintf("%s_%s", fileName, internalType)

					var solStructParams []SolStructParam
					for _, component := range input.Components {
						solStructParam := SolStructParam{
							Type: component.Type,
							Name: component.Name,
						}
						solStructParams = append(solStructParams, solStructParam)
					}

					solStruct := SolStruct{
						Name:            structName,
						SolStructParams: solStructParams,
					}
					solStructs = append(solStructs, solStruct)

					eventParams = append(eventParams, SolEventParam{
						Type:         structName,
						Name:         input.Name,
						Indexed:      input.Indexed,
						InternalType: "struct",
					})
				} else {
					eventParams = append(eventParams, SolEventParam{
						Type:    input.Type,
						Name:    input.Name,
						Indexed: input.Indexed,
					})
				}
			}

			solEvents = append(solEvents, SolEvent{
				Name:   eventName,
				Params: eventParams,
			})
		}
	}

	return solEvents, solStructs, nil
}

func generateFunctions(solEvents []SolEvent) ([]SolFunction, error) {
	var solFunctions []SolFunction

	for _, event := range solEvents {
		var solParams []SolFunctionParam

		for _, eventParam := range event.Params {
			memory := false
			if eventParam.Type == "bytes" || eventParam.Type == "string" || strings.HasSuffix(eventParam.Type, "[]") ||
				eventParam.InternalType == "struct" {
				memory = true
			}
			funcParam := SolFunctionParam{
				Type:   eventParam.Type,
				Memory: memory,
				Name:   eventParam.Name,
			}
			solParams = append(solParams, funcParam)
		}

		solFunctions = append(solFunctions, SolFunction{
			Name:      "emit" + event.Name,
			EventName: event.Name,
			Params:    solParams,
		})
	}

	return solFunctions, nil
}

func generateContract(solEvents []SolEvent, solStructs []SolStruct, solFunctions []SolFunction) (string, error) {
	const templateStr = `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

contract EventsMock {
    {{- range .Structs }}
    struct {{ .Name }} { {{- range .SolStructParams }}{{ .Type }} {{ .Name }}; {{- end }} }
    {{- end }}

    {{- range $event := .Events }}
    event {{ $event.Name }}({{- range $paramIndex, $param := $event.Params }}{{ $param.Type }}{{ if $param.Indexed }} indexed{{ end }} {{ $param.Name }}{{- if lt $paramIndex (sub1 (len $event.Params)) }},{{ end }}{{ end }});
    {{- end }}

	{{- range $function := .Functions }}
    function {{ $function.Name }}({{- range $paramIndex, $param := $function.Params }}{{ $param.Type }}{{ if $param.Memory }} memory{{ end }} {{ $param.Name }}{{- if lt $paramIndex (sub1 (len $function.Params)) }},{{- end }}{{ end }}) public {
        emit {{ $function.EventName }}({{- range $paramIndex, $param := $function.Params }}{{ $param.Name }}{{- if lt $paramIndex (sub1 (len $function.Params)) }},{{ end }}{{ end }});
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
	}

	tmpl, err := template.New("mockContract").Funcs(funcMap).Funcs(sprig.TxtFuncMap()).Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := TemplateData{
		Events:    solEvents,
		Structs:   solStructs,
		Functions: solFunctions,
	}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func main() {
	abiPath := os.Args[1]
	solPath := os.Args[2]

	abiFiles, err := getABIFiles(abiPath)
	if err != nil {
		fmt.Println("Error finding ABI files:", err)
		os.Exit(1)
	}

	events := []SolEvent{}
	structs := []SolStruct{}
	functions := []SolFunction{}

	for _, abiFile := range abiFiles {
		abiJSON, err := ioutil.ReadFile(abiFile)
		if err != nil {
			fmt.Println("Error reading ABI file:", err)
			os.Exit(1)
		}

		fileName := strings.TrimSuffix(filepath.Base(abiFile), ".abi")
		fileEvents, fileStructs, err := extractEventsAndStructs(abiJSON, fileName)
		if err != nil {
			fmt.Println("Error parsing events:", err)
			os.Exit(1)
		}
		fileFunctions, err := generateFunctions(fileEvents)
		if err != nil {
			fmt.Println("Error generating functions:", err)
			os.Exit(1)
		}

		events = append(events, fileEvents...)
		structs = append(structs, fileStructs...)
		functions = append(functions, fileFunctions...)
	}

	// Generate the contract
	contract, err := generateContract(events, structs, functions)
	if err != nil {
		fmt.Println("Error generating mock contract:", err)
		os.Exit(1)
	}

	// Save the mock contract to a file
	err = os.WriteFile(solPath, []byte(contract), 0644)
	if err != nil {
		fmt.Println("Error writing mock contract to a file:", err)
		os.Exit(1)
	}

	fmt.Println("Generated EventsMock.sol mock contract!")
}
