package job

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Transformer interface {
	PipelineStage
	Transform(input interface{}) (interface{}, error)
}

type BaseTransformer struct {
	// ID       uint64   `json:"-" gorm:"primary_key"`
	notifiee Notifiee `json:"-" gorm:"-"`
}

// func (t BaseTransformer) GetID() uint64           { return t.ID }
func (t *BaseTransformer) SetNotifiee(n Notifiee) { t.notifiee = n }

type TransformerType string

const (
	TransformerTypeJSONParse TransformerType = "jsonparse"
	TransformerTypeMultiply  TransformerType = "multiply"
)

type Transformers []Transformer

func (t Transformers) Transform(input interface{}) (interface{}, error) {
	if len(t) == 0 {
		return input, nil
	}

	var err error
	for _, transformer := range t {
		input, err = transformer.Transform(input)
		if err != nil {
			return nil, err
		}
	}
	return input, nil
}

func (t Transformers) SetNotifiee(n Notifiee) {
	for _, transformer := range t {
		transformer.SetNotifiee(n)
	}
}

func (t *Transformers) UnmarshalJSON(bs []byte) (err error) {
	var rawJSON []json.RawMessage
	err = json.Unmarshal(bs, &rawJSON)
	if err != nil {
		return err
	}

	for _, bs := range rawJSON {
		var header struct {
			Type TransformerType `json:"type"`
		}
		err := json.Unmarshal(bs, &header)
		if err != nil {
			return err
		}
		var transformer Transformer
		switch header.Type {
		case TransformerTypeJSONParse:
			jsonTransformer := JSONParseTransformer{}
			err = json.Unmarshal(bs, &jsonTransformer)
			if err != nil {
				return err
			}
			transformer = &jsonTransformer

		case TransformerTypeMultiply:
			multiplyTransformer := MultiplyTransformer{}
			err = json.Unmarshal(bs, &multiplyTransformer)
			if err != nil {
				return err
			}
			transformer = &multiplyTransformer

		default:
			return errors.Errorf("invalid transformer type '%v'", header.Type)
		}
		*t = append(*t, transformer)
	}
	return nil
}
