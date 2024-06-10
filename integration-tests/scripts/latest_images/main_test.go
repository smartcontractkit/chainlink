package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func mockFetchImageDetailsSuccess(_ string) ([]byte, error) {
	return []byte(`[
		{
			"imageTags": ["v1.0.0"]
		},
		{
			"imageTags": ["v1.1.0"]
		},
		{
			"imageTags": ["v1.2.0"]
		}
	]`), nil
}

func mockFetchImageDetailsError(_ string) ([]byte, error) {
	return nil, fmt.Errorf("failed to describe images")
}

func TestGetLatestImages(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		images, err := getLatestImages(mockFetchImageDetailsSuccess, "test-repo", "v1.*", 2, "")
		require.NoError(t, err)
		require.Equal(t, "test-repo:v1.2.0,test-repo:v1.1.0", images)
	})

	t.Run("ErrorFetchingDetails", func(t *testing.T) {
		_, err := getLatestImages(mockFetchImageDetailsError, "test-repo", "v1.*", 2, "")
		require.Error(t, err)
		require.Equal(t, "failed to describe images: failed to describe images", err.Error())
	})

	t.Run("ErrorParsingTags", func(t *testing.T) {
		_, err := getLatestImages(mockFetchImageDetailsSuccess, "test-repo", "invalid[regex", 2, "")
		require.Error(t, err)
		require.Equal(t, "failed to parse image tags: failed to compile regex: error parsing regexp: missing closing ]: `[regex`", err.Error())
	})

	t.Run("InsufficientTags", func(t *testing.T) {
		_, err := getLatestImages(mockFetchImageDetailsSuccess, "test-repo", "v1.*", 5, "")
		require.Error(t, err)
		require.Equal(t, "failed to get 5 latest tags for test-repo. found only 3", err.Error())
	})

	t.Run("IgnoredTags", func(t *testing.T) {
		images, err := getLatestImages(mockFetchImageDetailsSuccess, "test-repo", "v1.*", 2, "v1.2.0")
		require.NoError(t, err)
		require.Equal(t, "test-repo:v1.1.0,test-repo:v1.0.0", images)
	})
}

func TestValidateInputs(t *testing.T) {
	t.Run("MissingArguments", func(t *testing.T) {
		os.Args = []string{"main"}
		expectedError := errors.New("usage: <repository_name> <grep_string> <count> [<ignored_tags>]")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("EmptyRepositoryName", func(t *testing.T) {
		os.Args = []string{"main", "", "v1.*", "2"}
		expectedError := errors.New("error: repository_name cannot be empty")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("EmptyGrepString", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "", "2"}
		expectedError := errors.New("error: grep_string cannot be empty")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("InvalidGrepString", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "invalid[regex", "2"}
		expectedError := errors.New("error: grep_string is not a valid regex")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("NonIntegerCount", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "v1.*", "two"}
		expectedError := fmt.Errorf("error: count must be an integer, but %s is not an integer", "two")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("EmptyIgnoredTag", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "v1.*", "2", "v1.0.0,"}
		expectedError := errors.New("error: ignored tag cannot be empty")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})
}
