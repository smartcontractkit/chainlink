package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// envreplaceCmd represents the environment replace command
var envreplaceCmd = &cobra.Command{
	Use:   "envreplace",
	Short: "Replace environment variable placeholders in a string",
	Long: `Replaces placeholders of the form ${{ env.VAR_NAME }} with the actual environment variable values.
Example usage:
./e2e_tests_tool envreplace --input 'Example with ${{ env.PATH }} and ${{ env.HOME }}.'`,
	Run: func(cmd *cobra.Command, _ []string) {
		input, _ := cmd.Flags().GetString("input")

		output, err := replaceEnvPlaceholders(input)
		if err != nil {
			fmt.Printf("Error replacing placeholders: %v\n", err)
			return
		}

		fmt.Println(output)
	},
}

func replaceEnvPlaceholders(input string) (string, error) {
	// Regular expression to match ${env.VAR_NAME}
	r := regexp.MustCompile(`\$\{\{\s*env\.(\w+)\s*\}\}`)
	var errors []string // Slice to accumulate error messages

	// Replace each match in the input string
	replaced := r.ReplaceAllStringFunc(input, func(m string) string {
		varName := r.FindStringSubmatch(m)[1]
		value := os.Getenv(varName)
		if value == "" {
			// If the environment variable is not set or empty, accumulate an error message
			errors = append(errors, fmt.Sprintf("environment variable '%s' not set or is empty", varName))
			return m // Return the original placeholder if error occurs
		}
		return value
	})

	// Check if there were any errors
	if len(errors) > 0 {
		return replaced, fmt.Errorf("multiple errors occurred: %s", strings.Join(errors, ", "))
	}
	return replaced, nil // No error occurred, return the replaced string
}

func init() {
	envreplaceCmd.Flags().StringP("input", "i", "", "Input string with placeholders")
	err := envreplaceCmd.MarkFlagRequired("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
