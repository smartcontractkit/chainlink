package pb

import (
	"encoding/base64"

	"github.com/shopspring/decimal"
)

func NewBoolValue(b bool) *Value {
	return &Value{
		Value: &Value_BoolValue{
			BoolValue: b,
		},
	}
}

func NewBytesValue(b []byte) *Value {
	bs := base64.StdEncoding.EncodeToString(b)
	return &Value{
		Value: &Value_BytesValue{
			BytesValue: bs,
		},
	}
}

func NewDecimalValue(d decimal.Decimal) *Value {
	return &Value{
		Value: &Value_DecimalValue{
			DecimalValue: d.String(),
		},
	}
}

func NewStringValue(s string) *Value {
	return &Value{
		Value: &Value_StringValue{
			StringValue: s,
		},
	}
}

func NewMapValue(m map[string]*Value) *Value {
	return &Value{
		Value: &Value_MapValue{
			MapValue: &Map{
				Fields: m,
			},
		},
	}
}

func NewListValue(m []*Value) *Value {
	return &Value{
		Value: &Value_ListValue{
			ListValue: &List{
				Fields: m,
			},
		},
	}
}

func NewInt64Value(i int64) *Value {
	return &Value{
		Value: &Value_Int64Value{
			Int64Value: i,
		},
	}
}
