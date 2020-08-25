package job

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type JSONParseTransformer struct {
	BaseTransformer
	Path []string `json:"path"`
}

var (
	_ Transformer   = JSONParseTransformer{}
	_ PipelineStage = JSONParseTransformer{}
)

func (t JSONParseTransformer) Transform(input interface{}) (out interface{}, err error) {
	defer func() { t.notifiee.OnEndStage(t, out, err) }()
	t.notifiee.OnBeginStage(t, input)

	var bs []byte
	switch v := input.(type) {
	case []byte:
		bs = v
	case string:
		bs = []byte(v)
	default:
		return nil, errors.Errorf("JSONParseTransformer does not accept inputs of type %T", input)
	}

	var decoded interface{}
	err = json.Unmarshal(bs, &decoded)
	if err != nil {
		return nil, err
	}

	for i, part := range t.Path {
		switch d := decoded.(type) {
		case map[string]interface{}:
			var exists bool
			decoded, exists = d[part]
			if !exists && i == len(t.Path)-1 {
				return nil, nil
			} else if !exists {
				return nil, errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))
			}

		case []interface{}:
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			if index < 0 {
				index = len(d) + index
			}

			exists := index >= 0 && index < len(d)
			if !exists && i == len(t.Path)-1 {
				return nil, nil
			} else if !exists {
				return nil, errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))
			}
			decoded = d[index]

		default:
			return nil, errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))
		}
	}
	return decoded, nil
}

func (t JSONParseTransformer) MarshalJSON() ([]byte, error) {
	type preventInfiniteRecursion JSONParseTransformer
	type transformerWithType struct {
		Type TransformerType `json:"type"`
		preventInfiniteRecursion
	}
	return json.Marshal(transformerWithType{
		TransformerTypeJSONParse,
		preventInfiniteRecursion(t),
	})
}
