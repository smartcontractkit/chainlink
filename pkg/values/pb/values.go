package pb

import (
	"encoding/base64"
	"encoding/json"

	"github.com/shopspring/decimal"
)

func NewBoolValue(b bool) (*Value, error) {
	return &Value{
		Value: &Value_BoolValue{
			BoolValue: b,
		},
	}, nil
}

func NewBytesValue(b []byte) (*Value, error) {
	bs := base64.StdEncoding.EncodeToString(b)
	return &Value{
		Value: &Value_BytesValue{
			BytesValue: bs,
		},
	}, nil
}

func NewDecimalValue(d decimal.Decimal) (*Value, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return &Value{
		Value: &Value_DecimalValue{
			DecimalValue: string(b),
		},
	}, nil
}

func NewStringValue(s string) (*Value, error) {
	return &Value{
		Value: &Value_StringValue{
			StringValue: s,
		},
	}, nil
}

func NewMapValue(m map[string]*Value) (*Value, error) {
	return &Value{
		Value: &Value_MapValue{
			MapValue: &Map{
				Fields: m,
			},
		},
	}, nil
}

func NewListValue(m []*Value) (*Value, error) {
	return &Value{
		Value: &Value_ListValue{
			ListValue: &List{
				Fields: m,
			},
		},
	}, nil
}

func NewInt64Value(i int64) (*Value, error) {
	return &Value{
		Value: &Value_Int64Value{
			Int64Value: i,
		},
	}, nil
}

func NewNilValue() (*Value, error) {
	return &Value{
		Value: &Value_NilValue{
			NilValue: &Nil{},
		},
	}, nil
}
