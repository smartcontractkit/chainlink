package main

import (
	"fmt"

	gh "github.com/cli/go-gh/v2"
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

Make sure you have the GitHub CLI and it's authorized. Find it at https://cli.github.com/

Follow the prompts to run an E2E test. Type to search, use arrow keys to scroll, and Enter to select an option.
`
	helpText      string = "What do these mean?"
	chainlinkRepo string = "smartcontractkit/chainlink"
	workflowFile  string = "generic-test-runner.yml"
)

var (
	testDirectories = []string{helpText, "smoke", "soak", "performance", "reorg", "chaos", "benchmark"}
)

func main() {
	// This can take a while to retrieve, start it at the beginning asynchronously
	branchesAndTags, branchesAndTagsErr := make(chan []string, 1), make(chan error, 1)
	go collectBranchesAndTags(branchesAndTags, branchesAndTagsErr)

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

	dir, err := getTestDirectory()
	if err != nil {
		fmt.Printf("error getting test directory: %v\n", err)
		return
	}

	err = <-branchesAndTagsErr
	if err != nil {
		fmt.Printf("error getting branches and tags: %v\n", err)
		return
	}
	branch, err := getTestBranch(<-branchesAndTags)
	if err != nil {
		fmt.Printf("error selecting test branch: %v\n", err)
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
