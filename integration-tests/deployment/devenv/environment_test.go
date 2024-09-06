package devenv

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestDevEnv(t *testing.T) {
	DeployLocalCluster(t, zerolog.Logger{})
}
