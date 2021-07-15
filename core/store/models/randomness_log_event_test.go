package models_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestValidate_RandomnessLogEvent(t *testing.T) {
	t.Parallel()

	j := cltest.NewJobWithRandomnessLog()

	tests := []struct {
		description string
		emitter     common.Address
		want        bool
	}{
		{"valid", j.Initiators[0].Address, true},
		{"incorrect address", cltest.NewAddress(), false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.description, func(t *testing.T) {
			fee := assets.Link(*big.NewInt(1))
			randLog := models.RandomnessRequestLog{
				KeyHash:   utils.NewHash(),
				Seed:      big.NewInt(1),
				JobID:     utils.NewHash(),
				Sender:    cltest.NewAddress(),
				Fee:       &fee,
				RequestID: utils.NewHash(),
				Raw:       models.RawRandomnessRequestLog{},
			}
			log := cltest.NewRandomnessRequestLog(t, randLog, test.emitter, 1)
			le := models.RandomnessLogEvent{models.InitiatorLogEvent{Initiator: j.Initiators[0], Log: log}}
			result := le.Validate()
			assert.Equal(t, test.want, result)
		})
	}

}
