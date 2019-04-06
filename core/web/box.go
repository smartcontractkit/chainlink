package web

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MatchWildcardBoxPath returns the box path when there is a wildcard match
// and an empty string otherwise
func MatchWildcardBoxPath(boxList []string, path string, file string) (matchedPath string) {
	idPattern := regexp.MustCompile(`(_[a-zA-Z0-9]+_)`)
	pathSeparator := (string)(os.PathSeparator)
	escapedPathSeparator := pathSeparator
	if pathSeparator == `\` {
		escapedPathSeparator = `\\`
	}
	pathAndFile := filepath.Clean(strings.Join(
		[]string{path, file},
		pathSeparator,
	))
	normalizedPathAndFile := strings.Replace(
		strings.TrimPrefix(pathAndFile, pathSeparator),
		`/`,
		pathSeparator,
		-1,
	)

	for i := 0; i < len(boxList) && matchedPath == ""; i++ {
		boxPathWithIDPattern := idPattern.ReplaceAllString(boxList[i], `[a-zA-Z0-9]+`)
		pathPattern := fmt.Sprintf(
			`^%s$`,
			strings.Replace(boxPathWithIDPattern, `\`, escapedPathSeparator, -1),
		)
		match, _ := regexp.MatchString(pathPattern, normalizedPathAndFile)

		if match {
			matchedPath = boxList[i]
		}
	}

	return matchedPath
}

// MatchExactBoxPath returns the box path when there is an exact match for the
// resource and an empty string otherwise
func MatchExactBoxPath(boxList []string, path string) (matchedPath string) {
	pathSeparator := (string)(os.PathSeparator)
	pathWithoutPrefix := strings.TrimPrefix(path, "/")
	normalizedPathAndFile := strings.Replace(
		strings.TrimPrefix(pathWithoutPrefix, pathSeparator),
		`/`,
		pathSeparator,
		-1,
	)

	for i := 0; i < len(boxList) && matchedPath == ""; i++ {
		boxPath := boxList[i]
		if boxPath == normalizedPathAndFile {
			matchedPath = boxPath
		}
	}

	return matchedPath
}
