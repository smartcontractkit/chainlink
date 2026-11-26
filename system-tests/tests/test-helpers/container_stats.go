package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	dc "github.com/moby/moby/client"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// ContainerUsage is a container's memory and CPU use as `docker stats` reports them. It is read
// from the Docker API because cAdvisor on Docker Desktop does not expose per-container series.
type ContainerUsage struct {
	MemoryMiB  float64 `json:"memory_mib"`  // usage minus inactive page cache, like `docker stats`
	CPUPercent float64 `json:"cpu_percent"` // 100 = one full core, averaged over the sampling window
}

// SampleContainerUsage reads Docker stats for each container twice, window apart, and returns the
// memory of the second read and the CPU used between both. Containers that don't exist are skipped.
func SampleContainerUsage(t *testing.T, names []string, window time.Duration) map[string]ContainerUsage {
	t.Helper()

	cli := newDockerClient(t)
	first := make(map[string]container.StatsResponse, len(names))
	for _, name := range names {
		stats, err := readContainerStats(t.Context(), cli, name)
		if err != nil {
			framework.L.Warn().Err(err).Str("container", name).Msg("skipping container usage")
			continue
		}
		first[name] = stats
	}

	time.Sleep(window)

	out := make(map[string]ContainerUsage, len(first))
	for name, before := range first {
		after, err := readContainerStats(t.Context(), cli, name)
		require.NoError(t, err, "failed to read Docker stats of %s", name)
		out[name] = containerUsageBetween(before, after)
	}
	return out
}

func containerUsageBetween(before, after container.StatsResponse) ContainerUsage {
	usage := ContainerUsage{MemoryMiB: float64(memoryInUse(after.MemoryStats)) / 1024 / 1024}
	elapsed := after.Read.Sub(before.Read)
	if elapsed > 0 && after.CPUStats.CPUUsage.TotalUsage >= before.CPUStats.CPUUsage.TotalUsage {
		cpuNanos := after.CPUStats.CPUUsage.TotalUsage - before.CPUStats.CPUUsage.TotalUsage
		usage.CPUPercent = float64(cpuNanos) / float64(elapsed.Nanoseconds()) * 100
	}
	return usage
}

// memoryInUse mirrors the Docker CLI: page cache that can be reclaimed doesn't count.
func memoryInUse(m container.MemoryStats) uint64 {
	inactive, ok := m.Stats["inactive_file"] // cgroup v2
	if !ok {
		inactive = m.Stats["total_inactive_file"] // cgroup v1
	}
	if inactive > m.Usage {
		return m.Usage
	}
	return m.Usage - inactive
}

func readContainerStats(ctx context.Context, cli *dc.Client, name string) (container.StatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := cli.ContainerStats(ctx, name, dc.ContainerStatsOptions{Stream: false})
	if err != nil {
		return container.StatsResponse{}, fmt.Errorf("docker stats %s: %w", name, err)
	}
	defer res.Body.Close()

	var stats container.StatsResponse
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return container.StatsResponse{}, fmt.Errorf("decode docker stats of %s: %w", name, err)
	}
	return stats, nil
}
