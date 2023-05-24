package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	gh "github.com/cli/go-gh/v2"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
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
	)
	if err != nil {
		fmt.Printf("Error running gh workflow run: %v\n", err)
		fmt.Println(stdErr.String())
		return
	}
	fmt.Println(stdOut.String())
	err = waitForWorkflowRun(branch, ghUser)
	if err != nil {
		fmt.Printf("Error waiting for workflow to start: %v\n", err)
		return
	}
}

func waitForWorkflowRun(branch, ghUser string) error {
	fmt.Println("Waiting for workflow to start")
	startTime := time.Now()
	checkWorkflow, timeout := time.NewTicker(time.Second), time.After(time.Second*15)
	defer checkWorkflow.Stop()
	for {
		select {
		case <-checkWorkflow.C:
			workflowId, err := checkWorkflowRun(startTime, branch, ghUser)
			if err != nil {
				return err
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
				return err
			}
			return nil
		case <-timeout:
			return fmt.Errorf("timed out waiting for workflow run to start")
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
	return names
}
