package logger

import (
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// diskSpaceAvailableFn is used for testing to replace the default diskSpaceAvailable.
type diskSpaceAvailableFn func(path string) (utils.FileSize, error)

// diskSpaceAvailable returns the available/free disk space in the requested `path`. Returns an error if it fails to find the path.
func diskSpaceAvailable(path string) (utils.FileSize, error) {
	free, err := availableDiskSpace(path)
	if err != nil {
		return 0, err
	}
	return utils.FileSize(free), nil
}
