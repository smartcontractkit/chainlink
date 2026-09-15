package image

import (
	"errors"
	"fmt"
	"strings"
)

// ResolveOptions contains parameters for resolving an image URI.
type ResolveOptions struct {
	ECRType          string
	RepositoryPath   string
	ImageTag         string
	AWSAccountNumber string
	AWSRegion        string
}

// Resolve resolves the Docker image URI based on ECR type and configuration.
func Resolve(opts ResolveOptions) (string, error) {
	ecrType := strings.ToLower(strings.TrimSpace(opts.ECRType))
	repoPath := strings.ToLower(strings.TrimSpace(opts.RepositoryPath))
	imageTag := strings.TrimSpace(opts.ImageTag)
	awsAccount := strings.ToLower(strings.TrimSpace(opts.AWSAccountNumber))
	awsRegion := strings.ToLower(strings.TrimSpace(opts.AWSRegion))

	if ecrType == "" {
		return "", errors.New("'ECR_TYPE' must be set and non-empty. Allowed values: 'sdlc' or 'public'")
	}
	if ecrType != "sdlc" && ecrType != "public" {
		return "", fmt.Errorf("invalid 'ECR_TYPE': '%s'. Allowed values: 'sdlc' or 'public'", ecrType)
	}
	if repoPath == "" {
		return "", errors.New("'CHAINLINK_IMAGE_REPO_PATH' must be set and non-empty")
	}
	if imageTag == "" {
		return "", errors.New("'CHAINLINK_IMAGE_TAG' must be set and non-empty")
	}

	if ecrType == "public" {
		return fmt.Sprintf("public.ecr.aws/%s:%s", repoPath, imageTag), nil
	}

	// sdlc
	if awsAccount == "" || awsRegion == "" {
		return "", errors.New("for 'ECR_TYPE=sdlc', both 'AWS_ACCOUNT_NUMBER' and 'AWS_REGION' must be set and non-empty")
	}

	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s", awsAccount, awsRegion, repoPath, imageTag), nil
}
