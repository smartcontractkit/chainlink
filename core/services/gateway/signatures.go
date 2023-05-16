package gateway

import (
	"bytes"
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Message signatures are over the following data:
//  1. MessageId aligned to 128 bytes
//  2. Method aligned to 64 bytes
//  3. DonId aligned to 64 bytes
//  4. Payload (before parsing)
func SignMessage(msgBody *MessageBody, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	rawData, err := getRawMessageBody(msgBody)
	if err != nil {
		return nil, err
	}
	return SignData(rawData, privateKey)
}

func ValidateMessageSignature(msg *Message) error {
	rawData, err := getRawMessageBody(&msg.Body)
	if err != nil {
		return err
	}
	signatureBytes, err := utils.TryParseHex(msg.Signature)
	if err != nil {
		return err
	}
	senderBytes, err := utils.TryParseHex(msg.Body.Sender)
	if err != nil {
		return err
	}
	return ValidateSignature(rawData, signatureBytes, senderBytes)
}

func SignData(data [][]byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	hash := crypto.Keccak256Hash(data...)
	return crypto.Sign(hash.Bytes(), privateKey)
}

func ValidateSignature(data [][]byte, signature []byte, signerAddress []byte) error {
	hash := crypto.Keccak256Hash(data...)
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return err
	}
	ecdsaPubKey, _ := crypto.UnmarshalPubkey(sigPublicKey)
	sigAddr := crypto.PubkeyToAddress(*ecdsaPubKey)
	if !bytes.Equal(sigAddr.Bytes(), signerAddress) {
		return errors.New("invalid signer")
	}

	signatureNoRecoverID := signature[:len(signature)-1]
	if !crypto.VerifySignature(sigPublicKey, hash.Bytes(), signatureNoRecoverID) {
		return errors.New("invalid signature")
	}
	return nil
}

func getRawMessageBody(msgBody *MessageBody) ([][]byte, error) {
	if msgBody == nil {
		return nil, errors.New("nil message")
	}
	alignedMessageId := make([]byte, MessageIdMaxLen)
	copy(alignedMessageId, msgBody.MessageId)
	alignedMethod := make([]byte, MessageMethodMaxLen)
	copy(alignedMethod, msgBody.Method)
	alignedDonId := make([]byte, MessageDonIdMaxLen)
	copy(alignedDonId, msgBody.DonId)
	return [][]byte{alignedMessageId, alignedMethod, alignedDonId, msgBody.Payload}, nil
}
