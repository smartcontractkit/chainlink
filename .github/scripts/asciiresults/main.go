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

	jsonFile, err := os.ReadFile(*jsonFileFlag)
	if err != nil {
		panic(fmt.Errorf("error reading JSON file:", err))
	}

	// Parse JSON data
	var jobs []ParsedResult
	if *namedKey == "" {
		err = json.Unmarshal(jsonFile, &jobs)
		if err != nil {
			panic(fmt.Errorf("error unmarshalling JSON array:", err))
		}
	} else {
		var results ResultsMap
		err = json.Unmarshal(jsonFile, &results)
		if err != nil {
			panic(fmt.Errorf("error unmarshalling JSON map:", err))
		}
		if val, ok := results[*namedKey]; ok {
			jobs = val
		} else {
			fmt.Printf("key %s not found in the JSON file\n", *namedKey)
			return
		}
	}

	// Determine column widths
	maxFirstColumnLen := len(*firstColumnHeader)
	maxSecondColumnLen := len(*secondColumnHeader)

	for _, job := range jobs {
		if len(job.Cap) > maxFirstColumnLen {
			maxFirstColumnLen = len(job.Cap)
		}
	}

	// Adjust column widths for section
	if len(*section) > 0 && len(*section) > maxFirstColumnLen {
		maxFirstColumnLen = len(*section)
	}

	// Create table format strings
	firstColumnFormat := fmt.Sprintf("%%-%ds", maxFirstColumnLen)
	secondColumnFormat := fmt.Sprintf("%%-%ds", maxSecondColumnLen)
	rowFormat := fmt.Sprintf("| %s | %s |\n", firstColumnFormat, secondColumnFormat)
	separator := fmt.Sprintf("+-%s-+-%s-+\n", strings.Repeat("-", maxFirstColumnLen), strings.Repeat("-", maxSecondColumnLen))

	// Open or create the output file
	outputFileName := "output.txt"
	if *outputFile != "" {
		outputFileName = *outputFile
	}
	file, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("error opening/creating file:", err))
	}
	defer func() { _ = file.Close() }()

	// Write table header if file is new
	fileInfo, err := file.Stat()
	if err != nil {
		panic(fmt.Errorf("error getting file info:", err))
	}
	if fileInfo.Size() == 0 {
		header := separator
		header += fmt.Sprintf(rowFormat, *firstColumnHeader, *secondColumnHeader)
		header += separator
		_, err = file.WriteString(header)
		if err != nil {
			panic(fmt.Errorf("error writing to file:", err))
		}
	}

	if len(*section) > 0 {
		sectionHeader := fmt.Sprintf("| %s |\n", centerText(*section, maxFirstColumnLen+maxSecondColumnLen+3))
		_, err = file.WriteString(sectionHeader)
		if err != nil {
			panic(fmt.Errorf("error writing to file:", err))
		}
		_, err = file.WriteString(separator)
		if err != nil {
			panic(fmt.Errorf("error writing to file:", err))
		}
	}

	for _, job := range jobs {
		result := "X"
		if job.Conclusion == ":white_check_mark:" {
			result = "✓"
		}
		line := fmt.Sprintf(rowFormat, job.Cap, result)
		_, err = file.WriteString(line)
		if err != nil {
			panic(fmt.Errorf("error writing to file:", err))
		}
	}

	footer := separator
	_, err = file.WriteString(footer)
	if err != nil {
		panic(fmt.Errorf("error writing to file:", err))
	}

	fmt.Println("Table updated successfully in", outputFileName)
}

// centerText centers a string within a given width
func centerText(s string, width int) string {
	spaces := (width - len(s)) / 2
	if (spaces*2+width+len(s))%2 == 0 {
		return strings.Repeat(" ", spaces) + s + strings.Repeat(" ", spaces)
	}
	return strings.Repeat(" ", spaces) + s + strings.Repeat(" ", spaces+1)
}
