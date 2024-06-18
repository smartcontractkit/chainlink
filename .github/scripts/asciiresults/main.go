package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ParsedResult struct {
	Conclusion string `json:"conclusion"`
	Cap        string `json:"cap"`
	URL        string `json:"html_url"`
}

type ResultsMap map[string][]ParsedResult

func readJSONFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func parseJSON(jsonData []byte, namedKey string) ([]ParsedResult, error) {
	var jobs []ParsedResult
	if namedKey == "" {
		err := json.Unmarshal(jsonData, &jobs)
		return jobs, err
	}

	var results ResultsMap
	err := json.Unmarshal(jsonData, &results)
	if err != nil {
		return nil, err
	}

	if val, ok := results[namedKey]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("key %s not found in the JSON file", namedKey)
}

func calculateColumnWidths(firstColumnHeader, secondColumnHeader string, sections map[string][]ParsedResult) (int, int) {
	maxFirstColumnLen := len(firstColumnHeader)
	maxSecondColumnLen := len(secondColumnHeader)

	for section, jobs := range sections {
		if len(section) > maxFirstColumnLen {
			maxFirstColumnLen = len(section)
		}
		for _, job := range jobs {
			if len(job.Cap) > maxFirstColumnLen {
				maxFirstColumnLen = len(job.Cap)
			}
		}
	}

	return maxFirstColumnLen, maxSecondColumnLen
}

func writeResultsToFile(fileName string, firstColumnHeader, secondColumnHeader, currentSection string, jobs []ParsedResult) error {
	// Read existing data
	sections := make(map[string][]ParsedResult)
	if _, err := os.Stat(fileName); err == nil {
		data, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		var sectionName string
		for _, line := range lines {
			if strings.HasPrefix(line, "|") {
				parts := strings.Split(line, "|")
				if len(parts) == 3 { // It's a section header
					sectionName = strings.TrimSpace(parts[1])
				} else if len(parts) == 4 { // It might be a job entry, but can also be the header
					if strings.TrimSpace(parts[1]) == firstColumnHeader || strings.TrimSpace(parts[2]) == secondColumnHeader {
						continue
					}
					if sectionName != "" {
						sections[sectionName] = append(sections[sectionName], ParsedResult{
							Cap:        strings.TrimSpace(parts[1]),
							Conclusion: strings.TrimSpace(parts[2]),
						})
					} else {
						sections[""] = append(sections[""], ParsedResult{
							Cap:        strings.TrimSpace(parts[1]),
							Conclusion: strings.TrimSpace(parts[2]),
						})
					}
				}
			}
		}
	}

	// Add new jobs to the current section or default section
	if currentSection != "" {
		sections[currentSection] = append(sections[currentSection], jobs...)
	} else {
		sections[""] = append(sections[""], jobs...)
	}

	// Calculate column widths
	maxFirstColumnLen, maxSecondColumnLen := calculateColumnWidths(firstColumnHeader, secondColumnHeader, sections)

	firstColumnFormat := fmt.Sprintf("%%-%ds", maxFirstColumnLen)
	secondColumnFormat := fmt.Sprintf("%%-%ds", maxSecondColumnLen)
	rowFormat := fmt.Sprintf("| %s | %s |\n", firstColumnFormat, secondColumnFormat)
	separator := fmt.Sprintf("+-%s-+-%s-+\n", strings.Repeat("-", maxFirstColumnLen), strings.Repeat("-", maxSecondColumnLen))

	// Open file for writing
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// Write header
	header := separator
	header += fmt.Sprintf(rowFormat, firstColumnHeader, secondColumnHeader)
	header += separator
	_, err = file.WriteString(header)
	if err != nil {
		return err
	}

	// Write sections and jobs
	for section, jobs := range sections {
		if section != "" {
			sectionHeader := fmt.Sprintf("| %s |\n", centerText(section, maxFirstColumnLen+maxSecondColumnLen+3))
			_, err = file.WriteString(sectionHeader)
			if err != nil {
				return err
			}
			_, err = file.WriteString(separator)
			if err != nil {
				return err
			}
		}

		for _, job := range jobs {
			result := "X"
			if job.Conclusion == ":white_check_mark:" || job.Conclusion == "√" {
				result = "√"
			}
			line := fmt.Sprintf(rowFormat, job.Cap, result)
			_, err = file.WriteString(line)
			if err != nil {
				return err
			}
		}

		_, err = file.WriteString(separator)
		if err != nil {
			return err
		}
	}

	return nil
}

func centerText(s string, width int) string {
	spaces := (width - len(s)) / 2
	if (spaces*2+width+len(s))%2 == 0 {
		return strings.Repeat(" ", spaces) + s + strings.Repeat(" ", spaces)
	}
	return strings.Repeat(" ", spaces) + s + strings.Repeat(" ", spaces+1)
}

func main() {
	firstColumnHeader := flag.String("firstColumn", "Value", "Header for the first column")
	secondColumnHeader := flag.String("secondColumn", "Result", "Header for the second column")
	jsonFileFlag := flag.String("jsonfile", "", "Path to JSON input file")
	section := flag.String("section", "", "Optional section name")
	namedKey := flag.String("namedKey", "", "Optional named key to look for in the JSON input")
	outputFile := flag.String("outputFile", "", "Optional output file to save results (default: output.txt)")

	flag.Parse()

	if *jsonFileFlag == "" {
		panic(fmt.Errorf("please provide a path to the JSON input file using --jsonfile flag"))
	}

	jsonFile, err := readJSONFile(*jsonFileFlag)
	if err != nil {
		panic(fmt.Errorf("error reading JSON file: %v", err))
	}

	jobs, err := parseJSON(jsonFile, *namedKey)
	if err != nil {
		panic(fmt.Errorf("error parsing JSON file: %v", err))
	}

	if len(jobs) == 0 {
		fmt.Println("No results found in the JSON file")
		return
	}

	outputFileName := "output.txt"
	if *outputFile != "" {
		outputFileName = *outputFile
	}

	err = writeResultsToFile(outputFileName, *firstColumnHeader, *secondColumnHeader, *section, jobs)
	if err != nil {
		panic(fmt.Errorf("error writing to file: %v", err))
	}

	msg := fmt.Sprintf("Found results for '%s'. Updating file %s", *namedKey, outputFileName)
	if *namedKey == "" {
		msg = fmt.Sprintf("Results updated successfully in %s", outputFileName)
	}

	fmt.Println(msg)
}
