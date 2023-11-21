package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/cli/go-gh/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/networks"
)

func waitForWorkflowRun(branch, ghUser string) (string, error) {
	fmt.Println("Waiting for workflow to start")
	startTime := time.Now()
	checkWorkflow, timeout := time.NewTicker(time.Second), time.After(time.Second*15)
	defer checkWorkflow.Stop()
	for {
		select {
		case <-checkWorkflow.C:
			workflowId, err := checkWorkflowRun(startTime, branch, ghUser)
			if err != nil {
				return "", err
			}
			if workflowId == "" {
				fmt.Println("Checking...")
				continue
			}
			fmt.Printf("Triggered Workflow with ID: %s\n", workflowId)
			fmt.Println("Opening run in browser...")
			_, stdErr, err := gh.Exec( // Opens the run in browser
				"run", "view", workflowId, "-w",
			)
			if err != nil {
				fmt.Println(stdErr.String())
				return "", err
			}
			return workflowId, nil
		case <-timeout:
			return "", fmt.Errorf("timed out waiting for workflow run to start")
		}
	}
}

func checkWorkflowRun(startTime time.Time, branch, ghUser string) (string, error) {
	stdOut, stdErr, err := gh.Exec( // Retrieves the runId of the workflow we just started
		"run", "list", "-b", branch, "-w", workflowFile, "-u", ghUser,
		"--json", "startedAt,databaseId", "-q", ".[0]",
	)
	if err != nil {
		fmt.Println(stdErr.String())
		return "", err
	}
	if stdOut.String() == "" {
		return "", nil
	}
	workflowRun := struct {
		DatabaseId int       `json:"databaseId"`
		StartedAt  time.Time `json:"startedAt"`
	}{}
	err = json.Unmarshal(stdOut.Bytes(), &workflowRun)
	if err != nil {
		return "", err
	}
	if workflowRun.StartedAt.Before(startTime) { // Make sure the workflow run started after we started waiting
		return "", nil
	}
	return fmt.Sprint(workflowRun.DatabaseId), nil
}

// getUser retrieves the current GitHub user's username
func getUser() (string, error) {
	stdOut, stdErr, err := gh.Exec(
		"api", "user", "-q", ".login",
	)
	if err != nil {
		fmt.Println(stdErr.String())
		return "", err
	}
	return stdOut.String(), nil
}

// getTestBranch prompts the user to select a test branch
func getTestBranch(options []string) (string, error) {
	fmt.Println("Ensure your branch has had its latest work pushed to GitHub before running a test.")
	testBranchPrompt := promptui.Select{
		Label: "Test Branch or Tag",
		Items: options,
		Searcher: func(input string, index int) bool {
			return strings.Contains(options[index], input)
		},
		StartInSearchMode: true,
	}
	_, branch, err := testBranchPrompt.Run()
	if err != nil {
		return "", err
	}
	return branch, nil
}

// getAllBranchesAndTags uses the github API to retrieve all branches and tags for the chainlink repo
// this call can take a while, so start it at the beginning asynchronously
func collectBranchesAndTags(results chan []string, errChan chan error) {
	defer close(errChan)
	defer close(results)

	branchChan, tagChan := make(chan []string, 1), make(chan []string, 1)
	defer close(branchChan)
	defer close(tagChan)

	// branches
	go func() {
		stdOut, stdErr, err := gh.Exec("api", fmt.Sprintf("repos/%s/branches", chainlinkRepo), "-q", ".[][\"name\"]", "--paginate")
		if err != nil {
			errChan <- fmt.Errorf("%w: %s", err, stdErr.String())
		}
		branches := strings.Split(stdOut.String(), "\n")
		cleanBranches := []string{}
		for _, branch := range branches {
			trimmed := strings.TrimSpace(branch)
			if branch != "" {
				cleanBranches = append(cleanBranches, trimmed)
			}
		}
		branchChan <- cleanBranches
	}()

	// tags
	go func() {
		stdOut, stdErr, err := gh.Exec("api", fmt.Sprintf("repos/%s/tags", chainlinkRepo), "-q", ".[][\"name\"]", "--paginate")
		if err != nil {
			errChan <- fmt.Errorf("%w: %s", err, stdErr.String())
		}
		tags := strings.Split(stdOut.String(), "\n")
		cleanTags := []string{}
		for _, tag := range tags {
			trimmed := strings.TrimSpace(tag)
			if tag != "" {
				cleanTags = append(cleanTags, trimmed)
			}
		}
		tagChan <- cleanTags
	}()

	// combine results
	branches, tags := <-branchChan, <-tagChan
	combined := append(branches, tags...)
	sort.Slice(combined, func(i, j int) bool {
		if combined[i] == "develop" {
			return true
		} else if combined[j] == "develop" {
			return false
		}
		return strings.Compare(combined[i], combined[j]) < 0
	})
	results <- combined
}

const helpDirectoryText = `Smoke tests are designed to be quick checks on basic functionality. 

Soak tests are designed to run for a long time and test the stability of the system under minimal or regular load.

Performance tests are designed to test the system under heavy load and measure performance metrics.

Chaos tests are designed to break the system in various ways and ensure it recovers gracefully.

Reorg tests are designed to test the system's ability to handle reorgs on the blockchain.

Benchmark tests are designed to check how far the system can go before running into issues.`

// getTestDirectory prompts the user to select a test directory
func getTestDirectory() (string, error) {
	testDirectoryPrompt := promptui.Select{
		Label: "Test Type",
		Items: testDirectories,
		Size:  10,
		Searcher: func(input string, index int) bool {
			return strings.Contains(testDirectories[index], input)
		},
		StartInSearchMode: true,
	}
	_, dir, err := testDirectoryPrompt.Run()
	if err != nil {
		return "", err
	}
	if dir == helpText {
		fmt.Println(helpDirectoryText)
		return getTestDirectory()
	}
	return dir, nil
}

// getTest searches the chosen test directory for valid tests to run
func getTest(dir string) (string, error) {
	items := testNames(dir)
	testPrompt := promptui.Select{
		Label: "Test Name",
		Items: items,
		Size:  15,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(items[index]), strings.ToLower(input))
		},
		StartInSearchMode: true,
	}
	_, test, err := testPrompt.Run()
	if err != nil {
		return "", err
	}
	return test, nil
}

// testNames returns a list of test names in the given directory
func testNames(directory string) []string {
	// Regular expression pattern to search for
	pattern := "func Test(\\w+?)\\(t \\*testing.T\\)"

	names := []string{}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() { // Skip directories
			return nil
		}
		if !strings.HasSuffix(info.Name(), "_test.go") { // Skip non-test files
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		regex := regexp.MustCompile(pattern)
		// Iterate over each line in the file
		for scanner.Scan() {
			line := scanner.Text()
			submatches := regex.FindStringSubmatch(line)
			if len(submatches) > 0 {
				names = append(names, submatches[1])
			}
		}

		if scanner.Err() != nil {
			log.Error().Str("File", info.Name()).Msg("Error scanning file")
		}
		return scanner.Err()
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Error looking for tests")
	}
	sort.Strings(names)
	return names
}

// getNetwork prompts the user for a network to run the test on, including urls and keys if necessary
func getNetwork() (networkName, networkWs, networkHTTP, fundingKey string, err error) {
	validNetworks, i := make([]string, len(networks.MappedNetworks)), 0
	for network := range networks.MappedNetworks {
		validNetworks[i] = network
		i++
	}
	sort.Slice(validNetworks, func(i, j int) bool { // Get in (mostly) alphabetical order
		if validNetworks[i] == "SIMULATED" {
			return true
		} else if validNetworks[j] == "SIMULATED" {
			return false
		}
		return strings.Compare(validNetworks[i], validNetworks[j]) < 0
	})

	networkPrompt := promptui.Select{
		Label: "Network",
		Items: validNetworks,
		Size:  10,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(validNetworks[index]), strings.ToLower(input))
		},
		StartInSearchMode: true,
	}
	_, network, err := networkPrompt.Run()
	if err != nil {
		return "", "", "", "", err
	}
	if strings.Contains(network, "SIMULATED") { // We take care of simulated network URLs
		return network, "", "", "", nil
	}

	networkWsPrompt := promptui.Prompt{
		Label: "Network WS URL",
		Validate: func(s string) error {
			if s == "" {
				return errors.New("URL cannot be empty")
			}
			if !strings.HasPrefix(s, "ws") {
				return errors.New("URL must start with ws")
			}
			return nil
		},
	}
	networkWs, err = networkWsPrompt.Run()
	if err != nil {
		return "", "", "", "", err
	}

	networkHTTPPrompt := promptui.Prompt{
		Label: "Network HTTP URL",
		Validate: func(s string) error {
			if s == "" {
				return errors.New("URL cannot be empty")
			}
			if !strings.HasPrefix(s, "http") {
				return errors.New("URL must start with http")
			}
			return nil
		},
	}
	networkHTTP, err = networkHTTPPrompt.Run()
	if err != nil {
		return "", "", "", "", err
	}

	networkFundingKeyPrompt := promptui.Prompt{
		Label: "Network Funding Key",
		Validate: func(s string) error {
			if s == "" {
				return errors.New("funding key cannot be empty for a non-simulated network")
			}
			_, err := crypto.HexToECDSA(s)
			if err != nil {
				return fmt.Errorf("funding key must be a valid hex string: %w", err)
			}
			return nil
		},
	}
	fundingKey, err = networkFundingKeyPrompt.Run()
	if err != nil {
		return "", "", "", "", err
	}

	return network, networkWs, networkHTTP, fundingKey, nil
}
