package components

import (
	"bufio"
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewDockerFakeTester is a small utility to test how docker host resolves in different CI environments
// it just prints logs of curl in CI when calling some URL
// it is temporary and will be removed in the future
func NewDockerFakeTester(url string) error {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:      "curlimages/curl:latest",
		Cmd:        []string{"curl", "-v", url},
		Labels:     framework.DefaultTCLabels(),
		WaitingFor: wait.ForExit(),
	}
	curlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}
	logs, err := curlContainer.Logs(ctx)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(logs)
	fmt.Println("Container logs:")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading container logs: %w", err)
	}
	return nil
}
