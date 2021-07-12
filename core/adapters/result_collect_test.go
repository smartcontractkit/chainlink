package adapters_test

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestResultCollect_Perform(t *testing.T) {
	var tt = []struct {
		name                     string
		json                     string
		expectedResultCollection []gjson.Result
		expectedResult           gjson.Result
	}{
		{
			name:                     "empty add bool",
			json:                     `{"result":false}`,
			expectedResultCollection: gjson.ParseBytes([]byte(`[false]`)).Array(),
			expectedResult:           gjson.ParseBytes([]byte(`false`)),
		},
		{
			name:                     "exists add bool",
			json:                     fmt.Sprintf(`{"result":false,"%s":[false,true]}`, models.ResultCollectionKey),
			expectedResultCollection: gjson.ParseBytes([]byte(`[false,true,false]`)).Array(),
			expectedResult:           gjson.ParseBytes([]byte(`false`)),
		},
		{
			name:                     "exists add int",
			json:                     fmt.Sprintf(`{"result":20,"%s":[false]}`, models.ResultCollectionKey),
			expectedResultCollection: gjson.ParseBytes([]byte(`[false,20]`)).Array(),
			expectedResult:           gjson.ParseBytes([]byte(`20`)),
		},
		{
			name:                     "exists add float",
			json:                     fmt.Sprintf(`{"result":20.214,"%s":[false]}`, models.ResultCollectionKey),
			expectedResultCollection: gjson.ParseBytes([]byte(`[false,20.214]`)).Array(),
			expectedResult:           gjson.ParseBytes([]byte(`20.214`)),
		},
		{
			name:                     "exists non-array",
			json:                     fmt.Sprintf(`{"result":20.214,"%s":false}`, models.ResultCollectionKey),
			expectedResultCollection: gjson.ParseBytes([]byte(`[false,20.214]`)).Array(),
			expectedResult:           gjson.ParseBytes([]byte(`20.214`)),
		},
	}

	for _, tc := range tt {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			past := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.ResultCollect{}
			result := adapter.Perform(past, nil, nil)
			assert.NoError(t, result.Error())
			require.Len(t, result.ResultCollection().Array(), len(test.expectedResultCollection))
			for i, r := range result.ResultCollection().Array() {
				assert.Equal(t, test.expectedResultCollection[i].Value(), r.Value())
			}
			assert.Equal(t, test.expectedResult.Value(), result.Result().Value())
		})
	}
}
