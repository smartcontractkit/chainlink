package gethwrappers

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// VersionHash is the hash used to detect changes in the underlying contract
func VersionHash(abiPath string, binPath string) (hash string) {
	abi, err := os.ReadFile(abiPath)
	if err != nil {
		Exit("Could not read abi path to create version hash", err)
	}
	bin := []byte("")
	if binPath != "-" {
		bin, err = os.ReadFile(binPath)
		if err != nil {
			Exit("Could not read abi path to create version hash", err)
		}
	}
	hashMsg := string(abi) + string(bin) + "\n"
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hashMsg)))
}

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
		if _, err := os.Stat(filepath.Join(root, "go.mod")); !os.IsNotExist(err) {
			return root
		}
		root = filepath.Dir(root)
	}
	Exit("could not find project root", nil)
	panic("can't get here")
}

func TempDir(dirPrefix string) (string, func()) {
	tmpDir, err := os.MkdirTemp("", dirPrefix+"-contractWrapper")
	if err != nil {
		Exit("failed to create temporary working directory", err)
	}
	return tmpDir, func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Println("failure while cleaning up temporary working directory:", err)
		}
	}
}

// BoxOutput formats its arguments as fmt.Printf, and encloses them in a box of
// arrows pointing at their content, in order to better highlight it. See
// ExampleBoxOutput
func BoxOutput(errorMsgTemplate string, errorMsgValues ...interface{}) string {
	errorMsgTemplate = fmt.Sprintf(errorMsgTemplate, errorMsgValues...)
	lines := strings.Split(errorMsgTemplate, "\n")
	maxlen := 0
	for _, line := range lines {
		if len(line) > maxlen {
			maxlen = len(line)
		}
	}
	internalLength := maxlen + 4
	output := "↘" + strings.Repeat("↓", internalLength) + "↙\n" // top line
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	readme := strings.Repeat("README ", maxlen/7)
	output += "→  " + readme + strings.Repeat(" ", maxlen-len(readme)) + "  ←\n"
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	for _, line := range lines {
		output += "→  " + line + strings.Repeat(" ", maxlen-len(line)) + "  ←\n"
	}
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	output += "→  " + readme + strings.Repeat(" ", maxlen-len(readme)) + "  ←\n"
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	return "\n" + output + "↗" + strings.Repeat("↑", internalLength) + "↖" + // bottom line
		"\n\n"
}
