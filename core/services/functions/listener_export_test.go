package functions

import "github.com/ethereum/go-ethereum/common"

func (l *FunctionsListener) VerifyRequestSignature(requestID RequestID, subscriptionOwner common.Address, requestData *RequestData) error {
	return l.verifyRequestSignature(requestID, subscriptionOwner, requestData)
}

func (l *FunctionsListener) ParseCBOR(requestId RequestID, cborData []byte, maxSizeBytes uint32) (*RequestData, error) {
	return l.parseCBOR(requestId, cborData, maxSizeBytes)
}
