package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
)

const (
	art string = `
------------------------------------------------------------------------------------------------------
 _____ _           _       _ _       _      _____         _    ______                            
/  __ \ |         (_)     | (_)     | |    |_   _|       | |   | ___ \                           
| /  \/ |__   __ _ _ _ __ | |_ _ __ | | __   | | ___  ___| |_  | |_/ /   _ _ __  _ __   ___ _ __ 
| |   | '_ \ / _, | | '_ \| | | '_ \| |/ /   | |/ _ \/ __| __| |    / | | | '_ \| '_ \ / _ \ '__|
| \__/\ | | | (_| | | | | | | | | | |   <    | |  __/\__ \ |_  | |\ \ |_| | | | | | | |  __/ |   
 \____/_| |_|\__,_|_|_| |_|_|_|_| |_|_|\_\   \_/\___||___/\__| \_| \_\__,_|_| |_|_| |_|\___|_|
------------------------------------------------------------------------------------------------------

Follow the prompts to run an E2E test. Use arrow keys to scroll and Enter to select an option.
`
	helpText string = "What do these mean?"
	allTests string = "All"
)

var (
	testDirectories = []string{helpText, "smoke", "soak", "performance", "reorg", "chaos", "benchmark"}
)

func main() {
	fmt.Print(art)
	testDirectoryPrompt := promptui.Select{
		Label: "Test Type",
		Items: testDirectories,
	}

	_, dir, err := testDirectoryPrompt.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Prompt failed")
	}
	if dir == helpText { // TODO: Write help text
		fmt.Println("Smoke tests are designed to be quick ")
		return
	}

	testChoiches := testNames(dir)
	if dir == "smoke" { // You can run all smoke if you want
		testChoiches = append([]string{allTests}, testChoiches...)
	}
	testPrompt := promptui.Select{
		Label: "Test",
		Items: testChoiches,
	}

	_, test, err := testPrompt.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Prompt failed")
	}
	fmt.Println(test)

}

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

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a scanner to read the file line by line
		scanner := bufio.NewScanner(file)

		// Compile the regex pattern
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
