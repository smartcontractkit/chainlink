package utils

import "github.com/shirou/gopsutil/v3/disk"

//go:generate mockery --name DiskStatsProvider --output ./mocks --case=underscore

// DiskStatsProvider describes the abstraction to the `shirou/gopsutil/v3` for mocking purposes
type DiskStatsProvider interface {
	AvailableSpace(path string) (FileSize, error)
}

type provider struct{}

// NewDiskStatsProvider returns a new `DiskStatsProvider` instance
func NewDiskStatsProvider() DiskStatsProvider {
	return &provider{}
}

// AvailableSpace returns the available/free disk space in the requested `path`. Returns an error if it fails to find the path.
func (p provider) AvailableSpace(path string) (FileSize, error) {
	diskUsage, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}

	return FileSize(diskUsage.Free), nil
}
