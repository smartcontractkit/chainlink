package ocr2key

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	chainKey, _ := newTerraKeyBundle()
	b, _ := chainKey.Marshal()

	var genKey keyBundleRawData
	json.Unmarshal(b, &genKey)

	fmt.Printf("%+v", genKey)
}
