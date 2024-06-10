package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

func getLatestImages(repositoryName, grepString string, count int, ignoredTags string) (string, error) {
	// Run the AWS CLI command to get the image details
	cmd := exec.Command("aws", "ecr", "describe-images", "--repository-name", repositoryName, "--region", os.Getenv("AWS_REGION"), "--output", "json", "--query", "imageDetails[?imageTags!=`null` && imageTags!=`[]`]")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to describe images: %w", err)
	}

	// Parse the JSON output using gojq
	var imageDetails []interface{}
	if err := json.Unmarshal(output, &imageDetails); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	query, err := gojq.Parse(".[] | .imageTags[0]")
	if err != nil {
		return "", fmt.Errorf("failed to parse gojq query: %w", err)
	}

	var tags []string
	iter := query.Run(imageDetails)
	for {
		tag, ok := iter.Next()
		if !ok {
			break
		}
		if tagStr, ok := tag.(string); ok {
			tags = append(tags, tagStr)
		} else if err, ok := tag.(error); ok {
			return "", fmt.Errorf("failed to run gojq query: %w", err)
		}
	}

	// Sort tags in reverse order
	sort.Sort(sort.Reverse(sort.StringSlice(tags)))

	// Filter tags based on the regex
	re, err := regexp.Compile(grepString)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}

	var filteredTags []string
	for _, tag := range tags {
		if re.MatchString(tag) {
			filteredTags = append(filteredTags, tag)
		}
		if len(filteredTags) == count {
			break
		}
	}

	if len(filteredTags) < count {
		return "", fmt.Errorf("failed to get %d latest tags for %s", count, repositoryName)
	}

	// Parse ignored tags
	ignoredTagsArray := strings.Split(ignoredTags, ",")

	var imagesArr []string
	for _, tag := range filteredTags {
		ignore := false
		for _, ignoredTag := range ignoredTagsArray {
			if tag == ignoredTag {
				ignore = true
				break
			}
		}
		if !ignore {
			imagesArr = append(imagesArr, fmt.Sprintf("%s:%s", repositoryName, tag))
		}
	}

	images := strings.Join(imagesArr, ",")
	return images, nil
}

func main() {
	// Read command line arguments
	if len(os.Args) < 4 {
		log.Fatalf("Usage: %s <repository_name> <grep_string> <count> <ignored_tags>", os.Args[0])
	}

	repositoryName := os.Args[1]
	grepString := os.Args[2]
	count, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Error: count must be an integer")
	}

	var ignoredTags string
	if len(os.Args) == 5 {
		ignoredTags = os.Args[4]
	}

	images, err := getLatestImages(repositoryName, grepString, count, ignoredTags)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(images)
}
