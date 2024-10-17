package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func (k *Keeper) PrintLogs(ctx context.Context, pattern string, grep, vgrep []string) {
	k.streamLogs(ctx, pattern, grep, vgrep)
}

func (k *Keeper) streamLogs(ctx context.Context, pattern string, grep, vgrep []string) {
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return
	}

	// Make sure everything works well
	if _, err = dockerClient.Ping(ctx); err != nil {
		return
	}

	allContainers, err := dockerClient.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		panic(err.Error())
	}

	re := regexp.MustCompile(pattern)

	var containerNames []string
	for _, container := range allContainers {
		for _, name := range container.Names {
			if re.MatchString(name) {
				containerNames = append(containerNames, name)
			}
		}
	}

	if len(containerNames) == 0 {
		panic(fmt.Sprintf("no container names matching regex: %s", pattern))
	}

	containerChannels := make([]chan string, len(containerNames))
	for i := range containerChannels {
		containerChannels[i] = make(chan string)
	}

	for i, containerName := range containerNames {
		go k.containerLogs(ctx, dockerClient, containerName, containerChannels[i], grep, vgrep)
	}

	mergedChannel := k.mergeChannels(containerChannels)

	go func() {
		for logLine := range mergedChannel {
			fmt.Println(logLine)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
	log.Println("Stopping...")
}

func (k *Keeper) containerLogs(ctx context.Context, cli *client.Client, containerID string, logsChan chan<- string, grep, vgrep []string) {
	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			panic(err)
		}
	}()

	reader, writer := io.Pipe()
	go func() {
		if _, err := stdcopy.StdCopy(writer, writer, out); err != nil {
			panic(err)
		}
		if err := writer.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(reader)
Scan:
	for scanner.Scan() {
		rawLogLine := scanner.Text()
		for _, exclude := range vgrep {
			if strings.Contains(rawLogLine, exclude) {
				continue Scan
			}
		}
		for _, include := range grep {
			if !strings.Contains(rawLogLine, include) {
				continue Scan
			}
		}
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(rawLogLine), &m); err != nil {
			continue
		}
		m["containerID"] = containerID
		decoratedLogLine, err := json.Marshal(m)
		if err != nil {
			continue
		}
		select {
		case logsChan <- string(decoratedLogLine):
		case <-ctx.Done():
			return
		}
	}
}

func (k *Keeper) mergeChannels(containerChannels []chan string) <-chan string {
	mergeChannel := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(containerChannels))

	for _, containerChannel := range containerChannels {
		go func(containerCh <-chan string) {
			defer wg.Done()
			for containerLogLine := range containerCh {
				mergeChannel <- containerLogLine
			}
		}(containerChannel)
	}

	go func() {
		wg.Wait()
		close(mergeChannel)
	}()

	return mergeChannel
}
