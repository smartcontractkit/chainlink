package presenters_test

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

type MI = models.Initiator

func TestPresenterInitiatorHasCorrectKeys(t *testing.T) {
	t.Parallel()

	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	now := time.Now()

	tests := []struct {
		i    models.Initiator
		keys []string
	}{
		{MI{Type: models.InitiatorWeb}, []string{"type"}},
		{MI{Type: models.InitiatorCron, Schedule: models.Cron("* * * * *")}, []string{"type", "schedule"}},
		{MI{Type: models.InitiatorRunAt, Time: models.Time{now}}, []string{"type", "time", "ran"}},
		{MI{Type: models.InitiatorEthLog, Address: address}, []string{"type", "address"}},
	}

	for _, test := range tests {
		t.Run(test.i.Type, func(t *testing.T) {
			j, err := json.Marshal(presenters.Initiator{test.i})
			assert.Nil(t, err)

			var value map[string]interface{}
			err = json.Unmarshal(j, &value)
			assert.Nil(t, err)

			keys := utils.GetStringKeys(value)
			sort.Strings(keys)
			sort.Strings(test.keys)
			assert.Equal(t, test.keys, keys)
		})
	}
}
