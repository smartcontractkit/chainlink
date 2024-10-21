package components

import (
	"bufio"
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewMockTester starts a curl container to make a GET request and print logs line by line.
func NewMockTester(url string) error {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "curlimages/curl:latest",
		Cmd:          []string{"curl", "-v", url},
		WaitingFor:   wait.ForExit(),
		ExposedPorts: []string{"9111/tcp"},
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
