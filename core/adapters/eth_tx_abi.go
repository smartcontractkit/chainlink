package adapters

import (
	strpkg "chainlink/core/store"
	"chainlink/core/store/models"
)

type EthTxABI struct {
	EthTxABIEncode
	EthTxRawCalldata
}

func (etx *EthTxABI) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	output := etx.EthTxABIEncode.Perform(input, store)
	if err := output.Error(); err != nil {
		return output
	}

	input2 := models.NewRunInput(input.JobRunID(), output.Data(), models.RunStatusUnstarted)
	return etx.EthTxRawCalldata.Perform(*input2, store)
}
