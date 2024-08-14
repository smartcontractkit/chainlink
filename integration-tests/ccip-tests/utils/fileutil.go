package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// FilesWithRegex returns all filepaths under root folder matching with regex pattern
func FilesWithRegex(root, pattern string) ([]string, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	var filenames []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file matches the regex pattern
		if !info.IsDir() && r.MatchString(info.Name()) {
			filenames = append(filenames, path)
		}

		return nil
	})
	return filenames, err
}

func FileNameFromPath(path string) string {
	if !strings.Contains(path, "/") {
		return path
	}
	return strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
}

// FirstFileFromMatchingPath formats the given filepathWithPattern with actual file path
// if filepathWithPattern is provided with a regex expression it returns the first filepath
// matching with the regex.
// if there is no regex provided in filepathWithPattern it just returns the provided filepath
func FirstFileFromMatchingPath(filepathWithPattern string) (string, error) {
	filename := FileNameFromPath(filepathWithPattern)
	if strings.Contains(filepathWithPattern, "/") {
		rootFolder := strings.Split(filepathWithPattern, filename)[0]
		allFiles, err := FilesWithRegex(rootFolder, filename)
		if err != nil {
			return "", fmt.Errorf("error trying to find file %s:%w", filepathWithPattern, err)
		}
		if len(allFiles) == 0 {
			return "", fmt.Errorf("error trying to find file %s", filepathWithPattern)
		}
		if len(allFiles) > 1 {
			log.Warn().Str("path", filepathWithPattern).Msg("more than one contract config files found in location, using the first one")
		}
		return allFiles[0], nil
	}
	return filepathWithPattern, nil
}
