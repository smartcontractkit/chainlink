package envcommon

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ParseJSONFile(path string, v any) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	b, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}
