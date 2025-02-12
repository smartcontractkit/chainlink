package leafer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Filepath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Log(cwd)

	pkgName := "helloworld"
	suffix := ""

	p := filepath.Join(cwd, "generated", suffix, pkgName, pkgName+".go")
	t.Log(p)
}
