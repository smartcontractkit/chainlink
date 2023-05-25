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

	gh "github.com/cli/go-gh/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink/integration-tests/networks"
)

const (
	art string = `
-------------------------------------------------------------------------------------------------
 _____ _           _       _ _       _      _____         _    ______                            
/  __ \ |         (_)     | (_)     | |    |_   _|       | |   | ___ \                           
| /  \/ |__   __ _ _ _ __ | |_ _ __ | | __   | | ___  ___| |_  | |_/ /   _ _ __  _ __   ___ _ __ 
| |   | '_ \ / _, | | '_ \| | | '_ \| |/ /   | |/ _ \/ __| __| |    / | | | '_ \| '_ \ / _ \ '__|
| \__/\ | | | (_| | | | | | | | | | |   <    | |  __/\__ \ |_  | |\ \ |_| | | | | | | |  __/ |   
 \____/_| |_|\__,_|_|_| |_|_|_|_| |_|_|\_\   \_/\___||___/\__| \_| \_\__,_|_| |_|_| |_|\___|_|
-------------------------------------------------------------------------------------------------

Follow the prompts to run an E2E test. Use arrow keys to scroll and Enter to select an option.

Make sure you have the GitHub CLI (https://cli.github.com/) downloaded and authorized.
`
	helpText      string = "What do these mean?"
	chainlinkRepo string = "smartcontractkit/chainlink"
	workflowFile  string = "generic-test-runner.yml"
)

var (
	testDirectories = []string{helpText, "smoke", "soak", "performance", "reorg", "chaos", "benchmark"}
)

func main() {
	fmt.Print(art)

	ghUser, err := getUser()
	if err != nil {
		fmt.Printf("error getting GitHub user, make sure you're signed in to the GitHub CLI: %v\n", err)
		return
	}
	fmt.Printf("Running as %s\n", ghUser)

	network, wsURL, httpURL, fundingKey, err := getNetwork()
	if err != nil {
		fmt.Printf("error getting network: %v\n", err)
		return
	}

	branch, err := getTestBranch()
	if err != nil {
		fmt.Printf("error getting test branch: %v\n", err)
		return
	}

	dir, err := getTestDirectory()
	if err != nil {
		fmt.Printf("error getting test directory: %v\n", err)
		return
	}

	test, err := getTest(dir)
	if err != nil {
		fmt.Printf("error getting test: %v\n", err)
		return
	}

	stdOut, stdErr, err := gh.Exec( // Triggers the workflow with specified test
		"workflow", "run", workflowFile,
		"--repo", chainlinkRepo,
		"--ref", branch,
		"-f", fmt.Sprintf("directory=%s", dir),
		"-f", fmt.Sprintf("test=Test%s", test),
		"-f", fmt.Sprintf("network=%s", network),
		"-f", fmt.Sprintf("wsURL=%s", wsURL),
		"-f", fmt.Sprintf("httpURL=%s", httpURL),
		"-f", fmt.Sprintf("fundingKey=%s", fundingKey),
	)
	if err != nil {
		fmt.Printf("Error running gh workflow run: %v\n", err)
		fmt.Println(stdErr.String())
		return
	}
	fmt.Println(stdOut.String())

	_, err = waitForWorkflowRun(branch, ghUser)
	if err != nil {
		fmt.Printf("Error waiting for workflow to start: %v\n", err)
		return
	}
}

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
func getTestBranch() (string, error) {
	fmt.Println("Ensure your branch has had its latest work pushed to GitHub before running a test.")
	testBranchPrompt := promptui.Prompt{
		Label:     "Test Branch or Tag",
		AllowEdit: false,
		Default:   "develop",
	}
	branch, err := testBranchPrompt.Run()
	if err != nil {
		return "", err
	}
	return branch, nil
}

// getTestDirectory prompts the user to select a test directory
func getTestDirectory() (string, error) {
	testDirectoryPrompt := promptui.Select{
		Label: "Test Type",
		Items: testDirectories,
		Size:  10,
	}
	_, dir, err := testDirectoryPrompt.Run()
	if err != nil {
		return "", err
	}
	if dir == helpText { // TODO: Write help text
		fmt.Println("Smoke tests are designed to be quick ")
		return getTestDirectory()
	}
	return dir, nil
}

// getTest searches the chosen test directory for valid tests to run
func getTest(dir string) (string, error) {
	testPrompt := promptui.Select{
		Label: "Test Name",
		Items: testNames(dir),
		Size:  15,
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

		// Skip directories
		if info.IsDir() {
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
	sort.Strings(validNetworks) // Get in alphabetical order

	networkPrompt := promptui.Select{
		Label: "Network",
		Items: validNetworks,
		Size:  10,
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
