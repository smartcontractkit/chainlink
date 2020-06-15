package gethwrappers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Exit(msg string, err error) {
	if err != nil {
		fmt.Println(msg+":", err)
	} else {
		fmt.Println(msg)
	}
	os.Exit(1)
}

// GetProjectRoot returns the root of the chainlink project
func GetProjectRoot() (rootPath string) {
	root, err := os.Getwd()
	if err != nil {
		Exit("could not get current working directory while seeking project root",
			err)
	}
	for root != "/" { // Walk up path to find dir containing go.mod
		if _, err := os.Stat(filepath.Join(root, "go.mod")); os.IsNotExist(err) {
			root = filepath.Dir(root)
		} else {
			return root
		}
	}
	Exit("could not find project root", nil)
	panic("can't get here")
}

func TempDir(dirPrefix string) (string, func()) {
	tmpDir, err := ioutil.TempDir("", dirPrefix+"-contractWrapper")
	if err != nil {
		Exit("failed to create temporary working directory", err)
	}
	return tmpDir, func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Println("failure while cleaning up temporary working directory:", err)
		}
	}
}
