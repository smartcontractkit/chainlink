package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

type Gauntlet struct {
	path string
	yarn string
	arg  []string
}

type FlowReport []struct {
	Name string `json:"name"`
	Txs  []struct {
		Contract string `json:"contract"`
		Hash     string `json:"hash"`
		Success  bool   `json:"success"`
	}
	Data   map[string]string `json:"data"`
	StepId int               `json:"stepId"`
}

type Report struct {
	Responses []struct {
		Tx struct {
			Hash    string `json:"hash"`
			Address string `json:"address"`
		}
		Contract string `json:"contract"`
	} `json:"responses"`
	Data map[string]string `json:"data"`
}

func NewGauntlet(path string) (Gauntlet, error) {
	yarn, err := exec.LookPath("yarn")
	if err != nil {
		return Gauntlet{}, errors.New("'yarn' is not installed (required by Gauntlet)")
	}

	// Change path to root directory
	cwd, _ := os.Getwd()
	if err = os.Chdir(path); err != nil {
		return Gauntlet{}, errors.Wrap(err, "error in changing to root directory")
	}

	fmt.Println("Installing dependencies")
	if _, err = exec.Command(yarn).Output(); err != nil {
		return Gauntlet{}, errors.New("error installing dependencies")
	}
	// Move back into ops folder
	if err = os.Chdir(cwd); err != nil {
		return Gauntlet{}, errors.Wrap(err, "error in changing to ops folder")
	}

	arg := []string{"--cwd", path, "gauntlet"}
	fmt.Printf("Runing gauntlet via yarn with args: %s\n", arg)

	err = exec.Command(yarn, arg...).Run()
	if err != nil {
		return Gauntlet{}, err
	}

	return Gauntlet{
		path: path,
		yarn: yarn,
		arg:  arg,
	}, nil
}

func (g Gauntlet) Flag(flag string, value string) string {
	return fmt.Sprintf("--%s=%s", flag, value)
}

func (g Gauntlet) ExecCommand(arg ...string) error {
	// Collect all the args
	a := []string{}
	a = append(a, g.arg...)
	a = append(a, arg...)
	// Execute
	cmd := exec.Command(g.yarn, a...)
	output, err := cmd.CombinedOutput()
	// For diagnostics we print full combined output on error
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	return nil
}

func (g Gauntlet) ReadCommandReport() (Report, error) {
	path := filepath.Join(g.path, "report.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		return Report{}, err
	}

	var report Report
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err = json.Unmarshal(byteValue, &report); err != nil {
		return Report{}, errors.Wrap(err, "error in unmarshalling report")
	}

	return report, nil
}

func (g Gauntlet) ReadCommandFlowReport() (FlowReport, error) {
	path := filepath.Join(g.path, "flow-report.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		return FlowReport{}, err
	}

	var report FlowReport
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &report)
	if err != nil {
		return FlowReport{}, err
	}

	return report, nil
}
