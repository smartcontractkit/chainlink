package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	tc "github.com/testcontainers/testcontainers-go"
)

type NodeCoverageHelper struct {
	Nodes            []tc.Container
	GoCoverSrcDir    string   // Path to the source directory on the chainlink image with go coverage data
	NodeCoverageDirs []string // Paths to individual node coverage directories
	CoverageDir      string   // Path to the base directory with all coverage
	MergedDir        string   // Path to the directory where all coverage will be merged
	ChainlinkDir     string   // Path to the root chainlink directory
}

func NewNodeCoverageHelper(ctx context.Context, nodes []tc.Container, chainlinkDir, coverageDir string) (*NodeCoverageHelper, error) {
	coverSrcDir := os.Getenv("GO_COVERAGE_SRC_DIR")
	if coverSrcDir == "" {
		coverSrcDir = "/var/tmp/go-coverage" // Default path
	}

	helper := &NodeCoverageHelper{
		Nodes:         nodes,
		GoCoverSrcDir: coverSrcDir,
		CoverageDir:   coverageDir,
		MergedDir:     filepath.Join(coverageDir, "merged"),
		ChainlinkDir:  chainlinkDir,
	}

	if err := os.MkdirAll(coverageDir, 0755); err != nil {
		return nil, errors.Wrap(err, "failed to create base directory for node coverage")
	}

	// Copy coverage data from nodes
	if err := helper.copyCoverageFromNodes(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to copy coverage from nodes during initialization")
	}

	// Merge the coverage data
	if err := helper.mergeCoverage(); err != nil {
		return nil, errors.Wrap(err, "failed to merge coverage data")
	}

	return helper, nil
}

func (c *NodeCoverageHelper) SaveMergedHTMLReport() (string, error) {
	// Generate the textual coverage report
	txtCommand := exec.Command("go", "tool", "covdata", "textfmt", "-i=.", "-o=cov.txt")
	txtCommand.Dir = c.MergedDir
	if txtOutput, err := txtCommand.CombinedOutput(); err != nil {
		return "", errors.Wrapf(err, "failed to generate textual coverage report: %s", string(txtOutput))
	}

	// Generate the HTML coverage report
	htmlFilePath := filepath.Join(c.CoverageDir, "coverage.html")
	// #nosec G204
	htmlCommand := exec.Command("go", "tool", "cover", "-html="+filepath.Join(c.MergedDir, "cov.txt"), "-o="+htmlFilePath)
	htmlCommand.Dir = c.ChainlinkDir
	if htmlOutput, err := htmlCommand.CombinedOutput(); err != nil {
		return "", errors.Wrapf(err, "failed to generate HTML coverage report: %s", string(htmlOutput))
	}

	return htmlFilePath, nil
}

func (c *NodeCoverageHelper) SaveMergedCoveragePercentage() (string, error) {
	filePath := filepath.Join(c.CoverageDir, "percentage.txt")

	// Calculate coverage percentage from the merged data
	percentCmd := exec.Command("go", "tool", "covdata", "percent", "-i=.")
	percentCmd.Dir = c.MergedDir // Ensure the command runs in the directory with the merged data
	output, err := percentCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get merged coverage percentage report: %w, output: %s", err, string(output))
	}

	// Save the cmd output to a file
	if err := os.WriteFile(filePath, output, 0600); err != nil {
		return "", errors.Wrap(err, "failed to write coverage percentage to file")
	}

	return filePath, nil
}

func (c *NodeCoverageHelper) mergeCoverage() error {
	if err := os.MkdirAll(c.MergedDir, 0755); err != nil {
		return fmt.Errorf("failed to create merged directory: %w", err)
	}

	// Merge the coverage data from all chainlink nodes
	dirInput := strings.Join(c.NodeCoverageDirs, ",")
	// #nosec G204
	mergeCmd := exec.Command("go", "tool", "covdata", "merge", "-o", c.MergedDir, "-i="+dirInput)
	mergeCmd.Dir = filepath.Dir(c.MergedDir)
	output, err := mergeCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing merge command: %w, output: %s", err, string(output))
	}

	// Remove the coverage dirs after merging
	for _, dir := range c.NodeCoverageDirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove directory %s: %w", dir, err)
		}
	}
	c.NodeCoverageDirs = []string{} // Reset the coverage paths after merging

	return nil
}

func (c *NodeCoverageHelper) copyCoverageFromNodes(ctx context.Context) error {
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(c.Nodes))

	for i, node := range c.Nodes {
		wg.Add(1)
		go func(n tc.Container, id int) {
			defer wg.Done()
			finalDestPath := filepath.Join(c.CoverageDir, fmt.Sprintf("node_%d", id))
			if err := os.MkdirAll(finalDestPath, 0755); err != nil {
				errorsChan <- fmt.Errorf("failed to create directory for node %d: %w", id, err)
				return
			}
			err := copyFolderFromContainerUsingDockerCP(ctx, n.GetContainerID(), c.GoCoverSrcDir, finalDestPath)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to copy folder from container for node %d: %w", id, err)
				return
			}
			finalDestPath = filepath.Join(finalDestPath, "go-coverage") // Assuming path structure /var/tmp/go-coverage/TestName/node_X/go-coverage
			c.NodeCoverageDirs = append(c.NodeCoverageDirs, finalDestPath)
		}(node, i)
	}

	wg.Wait()
	close(errorsChan)

	for err := range errorsChan {
		if err != nil {
			return err
		}
	}

	return nil
}
func copyFolderFromContainerUsingDockerCP(ctx context.Context, containerID, srcPath, destPath string) error {
	source := fmt.Sprintf("%s:%s", containerID, srcPath)
	cmd := exec.CommandContext(ctx, "docker", "cp", source, destPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "docker cp command failed: %s, output: %s", cmd, string(output))
	}
	return nil
}
