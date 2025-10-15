package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/root"
)

// CompletionNode represents a node in the command tree with completions and navigation
type CompletionNode struct {
	Suggestions     []prompt.Suggest
	Children        map[string]*CompletionNode
	DynamicProvider func([]string) []prompt.Suggest
	Flags           []prompt.Suggest
}

var commandTree *CompletionNode

func init() {
	commandTree = buildCommandTree()
}

// buildCommandTree constructs the complete command tree with all subcommands, flags, and dynamic providers
func buildCommandTree() *CompletionNode {
	root := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "env", Description: "Manage local CRE environments"},
			{Text: "bs", Description: "Manage the Blockscout EVM block explorer"},
			{Text: "obs", Description: "Manage the observability stack"},
			{Text: "exit", Description: "Exit the interactive shell"},
		},
		Children: make(map[string]*CompletionNode),
	}

	// ENV command tree
	envNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "start", Description: "Spin up the development environment"},
			{Text: "stop", Description: "Tear down the development environment"},
			{Text: "restart", Description: "Restart the development environment"},
			{Text: "workflow", Description: "Workflow management commands"},
			{Text: "beholder", Description: "Beholder stack management commands"},
			{Text: "swap", Description: "Swap capabilities or nodes in running environment"},
			{Text: "state", Description: "Manage and view environment state"},
			{Text: "billing", Description: "Billing Platform Service management commands"},
		},
		Children: make(map[string]*CompletionNode),
	}

	// ENV START - dynamic TOML file completion + flags
	envStartNode := &CompletionNode{
		DynamicProvider: getWorkflowTomlFiles,
		Flags: []prompt.Suggest{
			{Text: "--wait-on-error-timeout", Description: "Time to wait before removing Docker containers if environment fails to start (e.g. 10s, 1m, 1h) (default: 15s)"},
			{Text: "--extra-allowed-gateway-ports", Description: "Extra allowed ports for outgoing connections from the Gateway Connector (e.g. 8080,8081)"},
			{Text: "--with-example", Description: "Deploys and registers example workflow (default: false)"},
			{Text: "--example-workflow-timeout", Description: "Time to wait until example workflow succeeds (e.g. 10s, 1m, 1h) (default: 5m)"},
			{Text: "--with-plugins-docker-image", Description: "Docker image to use (must have all capabilities included)"},
			{Text: "--example-workflow-trigger", Description: "Trigger for example workflow to deploy (web-trigger or cron) (default: web-trigger)"},
			{Text: "--with-beholder", Description: "Deploys Beholder (Chip Ingress + Red Panda) (default: false)"},
			{Text: "--with-dashboards", Description: "Deploys Observability Stack and Grafana Dashboards (default: false)"},
			{Text: "--with-billing", Description: "Deploys Billing Platform Service (default: false)"},
			{Text: "--with-proto-configs", Description: "Paths to protobuf config files for Beholder, comma separated (default: ./proto-configs/default.toml)"},
			{Text: "--auto-setup", Description: "Runs setup before starting the environment (default: false)"},
			{Text: "--with-contracts-version", Description: "Version of workflow and capabilities registry contracts to use (v1 or v2) (default: v1)"},
			{Text: "--setup-config", Description: "Path to the TOML configuration file for the setup command"},
		},
	}

	// ENV STOP - flags
	envStopNode := &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--all", Description: "Remove also all extra services (beholder, billing) (default: false)"},
		},
	}

	// ENV RESTART - same as start
	envRestartNode := &CompletionNode{
		DynamicProvider: getWorkflowTomlFiles,
		Flags:           envStartNode.Flags, // Same flags as start
	}

	envNode.Children["start"] = envStartNode
	envNode.Children["stop"] = envStopNode
	envNode.Children["restart"] = envRestartNode

	// ENV WORKFLOW - workflow management
	workflowNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "deploy-and-verify-example", Description: "Deploy and verify example workflow"},
			{Text: "delete", Description: "Delete a specific workflow"},
			{Text: "delete-all", Description: "Delete all workflows"},
			{Text: "compile", Description: "Compile a workflow specification"},
			{Text: "deploy", Description: "Deploy a compiled workflow"},
		},
		Children: make(map[string]*CompletionNode),
	}

	// Workflow subcommand flags
	workflowNode.Children["compile"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--workflow-file-path", Description: "⚠️  Path to the workflow main Go file (required)"},
			{Text: "--workflow-name", Description: "Workflow name (default: exampleworkflow)"},
		},
	}

	workflowNode.Children["deploy"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--workflow-file-path", Description: "⚠️  Path to a base64-encoded workflow WASM file or to a Go file that contains the workflow (if --compile flag is used) (required)"},
			{Text: "--config-file-path", Description: "⚠️  Path to the workflow config file (required)"},
			{Text: "--secrets-file-path", Description: "⚠️  Path to the secrets file with env var to secret name mappings (not the encrypted one) (required)"},
			{Text: "--container-target-dir", Description: "Path to the target directory in the Docker container (default: /home/chainlink)"},
			{Text: "--container-name-pattern", Description: "Pattern to match Docker containers workflow DON containers (e.g. 'workflow-node')"},
			{Text: "--rpc-url", Description: "RPC URL (default: http://localhost:8545)"},
			{Text: "--workflow-owner-address", Description: "Workflow owner address (default: 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266)"},
			{Text: "--workflow-registry-address", Description: "Workflow registry address (default: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512)"},
			{Text: "--capabilities-registry-address", Description: "Capabilities registry address (default: 0x5FbDB2315678afecb367f032d93F642f64180aa3)"},
			{Text: "--don-id", Description: "donID used in the workflow registry contract (integer starting with 1) (default: 1)"},
			{Text: "--name", Description: "⚠️  Workflow name (required)"},
			{Text: "--delete-workflow-file", Description: "Deletes the workflow file after deployment (default: false)"},
			{Text: "--compile", Description: "Compiles the workflow before deploying it (default: false)"},
			{Text: "--with-contracts-version", Description: "Version of workflow and capabilities registry contracts to use (v1 or v2) (default: v1)"},
		},
	}

	workflowNode.Children["delete"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--rpc-url", Description: "RPC URL (default: http://localhost:8545)"},
			{Text: "--owner-address", Description: "Workflow owner address (default: 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266)"},
			{Text: "--workflow-registry-address", Description: "Workflow registry address (default: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512)"},
			{Text: "--name", Description: "⚠️  Workflow name (required)"},
			{Text: "--with-contracts-version", Description: "Version of workflow and capabilities registry contracts to use (v1 or v2) (default: v1)"},
		},
	}

	envNode.Children["workflow"] = workflowNode

	// ENV BEHOLDER - beholder management
	beholderNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "start", Description: "Start the Beholder stack"},
			{Text: "stop", Description: "Stop the Beholder stack"},
			{Text: "create-kafka-topics", Description: "Create Kafka topics for Beholder"},
			{Text: "fetch-and-register-protos", Description: "Fetch and register protobuf definitions"},
		},
		Children: make(map[string]*CompletionNode),
	}

	beholderNode.Children["start"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--with-proto-configs", Description: "Paths to protobuf config files for Beholder, comma separated (default: ./proto-configs/default.toml)"},
			{Text: "--wait-on-error-timeout", Description: "Time to wait before removing Docker containers if environment fails to start (e.g. 10s, 1m, 1h) (default: 15s)"},
		},
	}

	beholderNode.Children["create-kafka-topics"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--red-panda-kafka-url", Description: "⚠️  Red Panda Kafka URL (required)"},
			{Text: "--topics", Description: "⚠️  Kafka topics to create (e.g. 'topic1,topic2') (required)"},
			{Text: "--purge-topics", Description: "Remove existing Kafka topics (default: false)"},
		},
	}

	beholderNode.Children["fetch-and-register-protos"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--red-panda-schema-registry-url", Description: "Red Panda Schema Registry URL (default: http://localhost:8081)"},
			{Text: "--with-proto-configs", Description: "Paths to protobuf config files for Beholder, comma separated (default: ./proto-configs/default.toml)"},
		},
	}

	envNode.Children["beholder"] = beholderNode

	// ENV SWAP - swap capabilities or nodes
	swapNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "capability", Description: "Swap capability binary in running environment"},
			{Text: "nodes", Description: "Swap Chainlink node binary in running environment"},
		},
		Children: make(map[string]*CompletionNode),
	}

	swapNode.Children["capability"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--name", Description: "Name of the capability to swap (need to mach the value of capability flag used in the environment TOML config) (default: \"\")"},
			{Text: "--binary", Description: "⚠️  Location of the binary to swap on the host machine (required)"},
			{Text: "--force", Description: "Force removal of Docker containers. Set to false to enable graceful shutdown of the containers (be mindful that it will take longer to remove the them) (default: true)"},
		},
	}

	swapNode.Children["nodes"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--force", Description: "Force removal of Docker containers. Set to false to enable graceful shutdown of the containers (be mindful that it will take longer to remove the them) (default: true)"},
			{Text: "--wait-time", Description: "Time to wait for the containers to be removed (default: 2m)"},
		},
	}

	envNode.Children["swap"] = swapNode

	// ENV STATE - state management
	stateNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "list", Description: "List all state files in the environment"},
			{Text: "purge", Description: "Purge all state and cache files"},
		},
		Children: make(map[string]*CompletionNode),
	}
	envNode.Children["state"] = stateNode

	// ENV BILLING - billing service management
	billingNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "start", Description: "Start the Billing Platform Service"},
			{Text: "stop", Description: "Stop the Billing Platform Service"},
		},
		Children: make(map[string]*CompletionNode),
	}
	envNode.Children["billing"] = billingNode

	root.Children["env"] = envNode
	root.Children["e"] = envNode // Alias

	// BS command tree
	bsNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "up", Description: "Spin up Blockscout EVM block explorer"},
			{Text: "down", Description: "Spin down Blockscout EVM block explorer"},
			{Text: "restart", Description: "Restart the Blockscout EVM block explorer"},
		},
		Children: make(map[string]*CompletionNode),
		Flags: []prompt.Suggest{
			{Text: "--url", Description: "EVM RPC node URL (default: http://host.docker.internal:8555)"},
			{Text: "--chain-id", Description: "RPC's Chain ID (default: 2337)"},
		},
	}

	// BS subcommands inherit parent flags
	bsNode.Children["up"] = &CompletionNode{
		Flags: bsNode.Flags,
	}

	bsNode.Children["down"] = &CompletionNode{
		Flags: bsNode.Flags,
	}

	bsNode.Children["restart"] = &CompletionNode{
		Flags: bsNode.Flags,
	}

	root.Children["bs"] = bsNode

	// OBS command tree
	obsNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "up", Description: "Spin up the observability stack"},
			{Text: "down", Description: "Spin down the observability stack"},
			{Text: "restart", Description: "Restart the observability stack (data wipe)"},
		},
		Children: make(map[string]*CompletionNode),
		Flags: []prompt.Suggest{
			{Text: "--full", Description: "Enable full observability stack with additional components (default: false)"},
		},
	}

	// OBS subcommands inherit parent flags
	obsNode.Children["up"] = &CompletionNode{
		Flags: obsNode.Flags,
	}

	obsNode.Children["down"] = &CompletionNode{
		Flags: obsNode.Flags,
	}

	obsNode.Children["restart"] = &CompletionNode{
		Flags: obsNode.Flags,
	}

	root.Children["obs"] = obsNode

	// EXAMPLES command tree
	examplesNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "contracts", Description: "Deploy example contracts"},
			{Text: "deploy-fake-price-provider", Description: "Deploy a fake price provider service locally"},
		},
		Children: make(map[string]*CompletionNode),
	}

	// EXAMPLES CONTRACTS - contract deployment
	contractsNode := &CompletionNode{
		Suggestions: []prompt.Suggest{
			{Text: "deploy-permissionless-feeds-consumer", Description: "Deploy a Permissionless Feeds Consumer contract"},
			{Text: "deploy-balance-reader", Description: "Deploy a Balance Reader contract"},
		},
		Children: make(map[string]*CompletionNode),
	}

	contractsNode.Children["deploy-permissionless-feeds-consumer"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--rpc-url", Description: "RPC URL (default: http://localhost:8545)"},
		},
	}

	contractsNode.Children["deploy-balance-reader"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--rpc-url", Description: "RPC URL (default: http://localhost:8545)"},
		},
	}

	examplesNode.Children["contracts"] = contractsNode

	// EXAMPLES DEPLOY-FAKE-PRICE-PROVIDER
	examplesNode.Children["deploy-fake-price-provider"] = &CompletionNode{
		Flags: []prompt.Suggest{
			{Text: "--auth-key", Description: "Authentication key for the price provider (default: Bearer test-auth-key)"},
			{Text: "--port", Description: "Port to run the fake price provider on (default: 80)"},
			{Text: "--feed-ids", Description: "Feed IDs to provide prices for (default: 0x1234567890123456789012345678901234567890123456789012345678901234)"},
		},
	}

	root.Children["examples"] = examplesNode

	return root
}

// getWorkflowTomlFiles returns curated workflow TOML files with meaningful descriptions
// Also dynamically discovers any new files not in the curated list
func getWorkflowTomlFiles(commandPath []string) []prompt.Suggest {
	// Curated list with meaningful descriptions
	curatedFiles := []prompt.Suggest{
		{Text: "workflow-don.toml", Description: "Basic DON with 5 Chainlink nodes"},
		{Text: "workflow-don-solana.toml", Description: "Workflow DON with Solana chain support"},
		{Text: "workflow-don-tron.toml", Description: "Workflow DON with Tron chain support"},
		{Text: "workflow-gateway-don.toml", Description: "Workflow DON with Gateway connector in a separate node"},
		{Text: "workflow-gateway-capabilities-don.toml", Description: "Workflow DON and Capabilities DON with Gateway connector in a separate node"},
	}

	// Try to find additional files not in the curated list
	possiblePaths := []string{
		"./configs/", // if running from core/scripts/cre/environment
		"./core/scripts/cre/environment/configs/", // if running from repo root
		"../environment/configs/",                 // if running from parent directory
	}

	var files []os.DirEntry
	var err error

	for _, path := range possiblePaths {
		files, err = os.ReadDir(path)
		if err == nil {
			break
		}
	}

	// If we found the directory, add any workflow files not in the curated list
	if err == nil {
		curatedNames := make(map[string]bool)
		for _, s := range curatedFiles {
			curatedNames[s.Text] = true
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), "workflow") && strings.HasSuffix(file.Name(), ".toml") {
				if !curatedNames[file.Name()] {
					curatedFiles = append(curatedFiles, prompt.Suggest{
						Text:        file.Name(),
						Description: "Workflow configuration (auto-discovered)",
					})
				}
			}
		}
	}

	return curatedFiles
}

func executor(in string) {
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}
	if in == "exit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}
	args := strings.Fields(in)
	os.Args = append([]string{"local_cre"}, args...)
	if err := root.RootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// completer provides autocomplete suggestions using the command tree
func completer(in prompt.Document) []prompt.Suggest {
	text := in.TextBeforeCursor()
	words := strings.Fields(text)
	lastCharIsSpace := len(text) > 0 && text[len(text)-1] == ' '

	// Special case: if text ends with a dash or partial flag, manually capture it
	// strings.Fields won't capture a trailing "-" or partial flags properly
	if len(text) > 0 && !lastCharIsSpace {
		// Get the last "word" by splitting on space manually
		parts := strings.Split(text, " ")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			// If it starts with - but wasn't captured in words, or if words is empty
			if strings.HasPrefix(lastPart, "-") {
				// Make sure it's in the words slice
				if len(words) == 0 || !strings.HasPrefix(words[len(words)-1], "-") {
					words = append(words, lastPart)
				}
			}
		}
	}

	return traverseTree(commandTree, words, lastCharIsSpace, []string{})
}

// traverseTree navigates the command tree and returns appropriate completions
func traverseTree(node *CompletionNode, words []string, lastCharIsSpace bool, commandPath []string) []prompt.Suggest {
	// Filter out flags from navigation
	commandWords := filterFlags(words)

	// Determine if we're currently typing/completing a flag
	// Only consider it a flag if we're actively typing it (no trailing space)
	isTypingFlag := !lastCharIsSpace && len(words) > 0 && strings.HasPrefix(words[len(words)-1], "-")

	// No command words: return current level suggestions
	if len(commandWords) == 0 {
		// If typing a flag, show only flags
		if isTypingFlag {
			if len(node.Flags) == 0 {
				return []prompt.Suggest{}
			}
			return prompt.FilterHasPrefix(node.Flags, words[len(words)-1], true)
		}

		suggestions := node.Suggestions
		if node.DynamicProvider != nil {
			suggestions = node.DynamicProvider(commandPath)
		}

		// Filter by current word if not at space
		if !lastCharIsSpace && len(words) > 0 {
			return prompt.FilterHasPrefix(suggestions, words[len(words)-1], true)
		}

		return suggestions
	}

	currentWord := commandWords[0]

	// Last command word without space: filter current level
	if len(commandWords) == 1 && !lastCharIsSpace && !isTypingFlag {
		suggestions := node.Suggestions
		if node.DynamicProvider != nil {
			suggestions = node.DynamicProvider(commandPath)
		}
		return prompt.FilterHasPrefix(suggestions, currentWord, true)
	}

	// Navigate to child node
	child, exists := node.Children[currentWord]
	if !exists {
		return []prompt.Suggest{}
	}

	newPath := append(commandPath, currentWord)

	// If we have more command words, recurse deeper first
	if len(commandWords) > 1 {
		// Find where the second command word starts in the original words slice
		// We need to pass all remaining words (including flags) after the first command word
		remainingWords := findWordsAfterCommand(words, currentWord)
		return traverseTree(child, remainingWords, lastCharIsSpace, newPath)
	}

	// At this point, len(commandWords) == 1, we're at the deepest command level

	// Command word with space: show next level (files, subcommands, etc.)
	if lastCharIsSpace && !isTypingFlag {
		suggestions := child.Suggestions
		if child.DynamicProvider != nil {
			suggestions = child.DynamicProvider(newPath)
		}
		return suggestions
	}

	// Last command word without space and not typing flag: filter current level
	if !lastCharIsSpace && !isTypingFlag {
		suggestions := child.Suggestions
		if child.DynamicProvider != nil {
			suggestions = child.DynamicProvider(newPath)
		}
		return prompt.FilterHasPrefix(suggestions, currentWord, true)
	}

	// If typing a flag, show flags from child node
	if isTypingFlag {
		if len(child.Flags) == 0 {
			return []prompt.Suggest{}
		}
		return prompt.FilterHasPrefix(child.Flags, words[len(words)-1], true)
	}

	// If we reach here, we should show next level completions (shouldn't normally happen)
	suggestions := child.Suggestions
	if child.DynamicProvider != nil {
		suggestions = child.DynamicProvider(newPath)
	}
	return suggestions
}

// filterFlags removes flag arguments from words for tree navigation
func filterFlags(words []string) []string {
	var result []string
	skipNext := false

	for i, w := range words {
		if skipNext {
			skipNext = false
			continue
		}

		if strings.HasPrefix(w, "-") {
			// Check if next word is flag value (not starting with -)
			if i+1 < len(words) && !strings.HasPrefix(words[i+1], "-") {
				skipNext = true
			}
			continue
		}

		result = append(result, w)
	}

	return result
}

// findWordsAfterCommand returns all words after the first occurrence of the command word
func findWordsAfterCommand(words []string, command string) []string {
	for i, w := range words {
		if w == command && i+1 < len(words) {
			return words[i+1:]
		}
	}
	return []string{}
}

// resetTerm resets terminal settings to Unix defaults.
func resetTerm() {
	cmd := exec.Command("stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func StartShell() {
	defer resetTerm()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("local_cre> "),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionTitle("Local CRE Interactive Shell"),
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
