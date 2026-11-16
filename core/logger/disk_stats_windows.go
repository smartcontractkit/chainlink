//go:build windows

package logger

import "golang.org/x/sys/windows"

func availableDiskSpace(path string) (uint64, error) {
	directory, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	var freeAvailable, total, totalFree uint64
	if err := windows.GetDiskFreeSpaceEx(directory, &freeAvailable, &total, &totalFree); err != nil {
		return 0, err
	}
	return totalFree, nil
}
