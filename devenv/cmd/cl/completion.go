package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

func getCommands() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "", Description: "Choose command, press <space> for more options after selecting command"},
		{Text: "up", Description: "Spin up the development environment"},
		{Text: "down", Description: "Tear down the development environment"},
		{Text: "restart", Description: "Restart the development environment"},
		{Text: "test", Description: "Perform smoke or load/chaos testing"},
		{Text: "bs", Description: "Manage the Blockscout EVM block explorer"},
		{Text: "obs", Description: "Manage the observability stack"},
		{Text: "db", Description: "Inspect Databases"},
		{Text: "exit", Description: "Exit the interactive shell"},
	}
}

func getSubCommands(parent string) []prompt.Suggest {
	switch parent {
	case "test":
		return []prompt.Suggest{
			{Text: "load", Description: "Run OCR2 load test"},
			{Text: "gas", Description: "Run OCR2 load test + simulate gas spikes"},
			{Text: "chaos", Description: "Run OCR2 load test + introduce container kills and latency"},
		}
	case "bs":
		return []prompt.Suggest{
			{Text: "up", Description: "Spin up Blockscout and listen to dst chain (8555)"},
			{Text: "up -u http://host.docker.internal:8545 -c 1337", Description: "Spin up Blockscout and listen to src chain (8545)"},
			{Text: "down", Description: "Remove Blockscout stack"},
			{Text: "restart", Description: "Restart Blockscout and listen to dst chain (8555)"},
			{Text: "restart -u http://host.docker.internal:8545 -c 1337", Description: "Restart Blockscout and listen to src chain (8545)"},
		}
	case "obs":
		return []prompt.Suggest{
			{Text: "up", Description: "Spin up observability stack (Loki/Prometheus/Grafana)"},
			{Text: "up -f", Description: "Spin up full observability stack (Pyroscope, cadvisor, postgres exporter)"},
			{Text: "down", Description: "Spin down observability stack"},
			{Text: "restart", Description: "Restart observability stack"},
			{Text: "restart -f", Description: "Restart full observability stack"},
		}
	case "u":
		fallthrough
	case "up":
		fallthrough
	case "r":
		fallthrough
	case "restart":
		return []prompt.Suggest{
			{Text: "env.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes"},
			{Text: "env.toml,env-cl-rebuild.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes (custom build)"},
			{Text: "env.toml,env-geth.toml", Description: "Spin up Geth <> Geth local chains (clique), all services, 4 CL nodes"},
			{Text: "env.toml,env-fuji-fantom.toml", Description: "Spin up testnets: Fuji <> Fantom, all services, 4 CL nodes"},
		}
	default:
		return []prompt.Suggest{}
	}
}

func executor(in string) {
	checkDockerIsRunning()
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}
	if in == "exit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	args := strings.Fields(in)
	os.Args = append([]string{"cl"}, args...)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// completer provides autocomplete suggestions for multi-word commands.
func completer(in prompt.Document) []prompt.Suggest {
	text := in.TextBeforeCursor()
	words := strings.Fields(text)
	lastCharIsSpace := len(text) > 0 && text[len(text)-1] == ' '

	switch {
	case len(words) == 0:
		return getCommands()
	case len(words) == 1:
		if lastCharIsSpace {
			return getSubCommands(words[0])
		}
		return prompt.FilterHasPrefix(getCommands(), words[0], true)

	case len(words) >= 2:
		if lastCharIsSpace {
			return []prompt.Suggest{}
		}
		parent := words[0]
		currentWord := words[len(words)-1]
		return prompt.FilterHasPrefix(getSubCommands(parent), currentWord, true)
	default:
		return []prompt.Suggest{}
	}
}

// resetTerm resets terminal settings to Unix defaults.
func resetTerm() {
	cmd := exec.CommandContext(context.Background(), "stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func StartShell() {
	defer resetTerm()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("cl> "),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionTitle("Chainlink Environment Interactive Shell"),
		prompt.OptionMaxSuggestion(15),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionCompletionWordSeparator(" "),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionSuggestionTextColor(prompt.Green),
		prompt.OptionScrollbarThumbColor(prompt.DarkGray),
		prompt.OptionScrollbarBGColor(prompt.Black),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(buf *prompt.Buffer) {
				fmt.Println("Interrupted, exiting...")
				resetTerm()
				os.Exit(0)
			},
		}),
	)
	p.Run()
}
