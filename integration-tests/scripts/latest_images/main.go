package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

func fetchImageDetails(repositoryName string) ([]byte, error) {
	cmd := exec.Command("aws", "ecr", "describe-images", "--repository-name", repositoryName, "--region", os.Getenv("AWS_REGION"), "--output", "json", "--query", "imageDetails[?imageTags!=`null` && imageTags!=`[]`]")
	return cmd.Output()
}

func parseImageTags(output []byte, grepString string, ignoredTags []string) ([]string, error) {
	var imageDetails []interface{}
	if err := json.Unmarshal(output, &imageDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	query, err := gojq.Parse(".[] | .imageTags[0]")
	if err != nil {
		return nil, fmt.Errorf("failed to parse gojq query: %w", err)
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
			return nil, fmt.Errorf("failed to run gojq query: %w", err)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(tags)))

	re, err := regexp.Compile(grepString)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	var filteredTags []string
	for _, tag := range tags {
		if re.MatchString(tag) {
			ignore := false
			for _, ignoredTag := range ignoredTags {
				if tag == ignoredTag {
					ignore = true
					break
				}
			}
			if !ignore {
				filteredTags = append(filteredTags, tag)
			}
		}
	}

	return filteredTags, nil
}

func getLatestImages(fetchFunc func(string) ([]byte, error), repositoryName, grepString string, count int, ignoredTags string) (string, error) {
	output, err := fetchFunc(repositoryName)
	if err != nil {
		return "", fmt.Errorf("failed to describe images: %w", err)
	}

	ignoredTagsArray := strings.Split(ignoredTags, ",")
	tags, err := parseImageTags(output, grepString, ignoredTagsArray)
	if err != nil {
		return "", fmt.Errorf("failed to parse image tags: %w", err)
	}

	if len(tags) < count {
		return "", fmt.Errorf("failed to get %d latest tags for %s. found only %d", count, repositoryName, len(tags))
	}

	var imagesArr []string
	for i := 0; i < count; i++ {
		imagesArr = append(imagesArr, fmt.Sprintf("%s:%s", repositoryName, tags[i]))
	}

	images := strings.Join(imagesArr, ",")
	return images, nil
}

func main() {
	if err := validateInputs(); err != nil {
		panic(err)
	}

	repositoryName := os.Args[1]
	grepString := os.Args[2]
	count, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(fmt.Errorf("error: count must be an integer, but %s is not an integer", os.Args[3]))
	}

	var ignoredTags string
	if len(os.Args) == 5 {
		ignoredTags = os.Args[4]
	}

	images, err := getLatestImages(fetchImageDetails, repositoryName, grepString, count, ignoredTags)
	if err != nil {
		panic(fmt.Errorf("error getting latest images: %v", err))
	}

	fmt.Println(images)
}

func validateInputs() error {
	if len(os.Args) < 4 {
		return errors.New("usage: <repository_name> <grep_string> <count> [<ignored_tags>]")
	}

	if os.Args[1] == "" {
		return errors.New("error: repository_name cannot be empty")
	}

	if os.Args[2] == "" {
		return errors.New("error: grep_string cannot be empty")
	}

	if _, err := regexp.Compile(os.Args[2]); err != nil {
		return errors.New("error: grep_string is not a valid regex")
	}

	if _, err := strconv.Atoi(os.Args[3]); err != nil {
		return fmt.Errorf("error: count must be an integer, but %s is not an integer", os.Args[3])
	}

	if len(os.Args) == 5 && os.Args[4] != "" {
		for _, ignoredTag := range strings.Split(os.Args[4], ",") {
			if ignoredTag == "" {
				return errors.New("error: ignored tag cannot be empty")
			}
		}
	}

	return nil
}
