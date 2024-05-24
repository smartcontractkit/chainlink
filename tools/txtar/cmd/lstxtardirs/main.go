package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/txtar"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	testDir := filepath.Join(wd, "./testdata/scripts")

	dirPtr := flag.String(
		"dir",
		testDir,
		"the directory to run the tests in; defaults to running all the tests in testdata/scripts",
	)

	recursePtr := flag.Bool(
		"recurse",
		false,
		"whether to recurse or not",
	)

	flag.Parse()

	dirs := []string{}
	visitor := txtar.NewDirVisitor(*dirPtr, txtar.RecurseOpt(*recursePtr), func(path string) error {
		dirs = append(dirs, path)
		return nil
	})
	err = visitor.Walk()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(strings.Join(dirs, "\n"))
}
