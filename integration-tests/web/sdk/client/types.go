package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type CreateFeedsManagerInput struct {
	Name      string `json:"name"`
	Uri       string `json:"uri"`
	PublicKey string `json:"publicKey"`
}

func DecodeInput(in, out any) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr || reflect.ValueOf(out).IsNil() {
		return fmt.Errorf("out type must be a non-nil pointer")
	}
	jsonBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		return err
	}

	return nil
}
