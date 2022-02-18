package utils

import "github.com/shirou/gopsutil/v3/disk"

//go:generate mockery --name DiskStatsProvider --output ./mocks --case=underscore

type DiskStatsProvider interface {
	AvailableSpace(path string) (FileSize, error)
}

type provider struct{}

func NewDiskStatsProvider() DiskStatsProvider {
	return &provider{}
}

func (p provider) AvailableSpace(path string) (FileSize, error) {
	diskUsage, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}

	return FileSize(diskUsage.Free), nil
}
