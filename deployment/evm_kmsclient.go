package deployment

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var (
	secp256k1N     = crypto.S256().Params().N
	secp256k1HalfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

// See https://docs.aws.amazon.com/kms/latest/APIReference/API_GetPublicKey.html#API_GetPublicKey_ResponseSyntax
// and https://datatracker.ietf.org/doc/html/rfc5280 for why we need to unpack the KMS public key.
type asn1SubjectPublicKeyInfo struct {
	AlgorithmIdentifier asn1AlgorithmIdentifier
	SubjectPublicKey    asn1.BitString
}

type asn1AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.ObjectIdentifier
}

// See https://aws.amazon.com/blogs/database/part2-use-aws-kms-to-securely-manage-ethereum-accounts/ for why we
// need to manually prep the signature for Ethereum.
type asn1ECDSASig struct {
	R asn1.RawValue
	S asn1.RawValue
}

// TODO: Mockery gen then test with a regular eth key behind the interface.
type KMSClient interface {
	GetPublicKey(input *kms.GetPublicKeyInput) (*kms.GetPublicKeyOutput, error)
	Sign(input *kms.SignInput) (*kms.SignOutput, error)
}

type KMS struct {
	KmsDeployerKeyId     string
	KmsDeployerKeyRegion string
	AwsProfileName       string
}

func NewKMSClient(config KMS) (KMSClient, error) {
	if config.KmsDeployerKeyId == "" {
		return nil, errors.New("KMS key ID is required")
	}
	if config.KmsDeployerKeyRegion == "" {
		return nil, errors.New("KMS key region is required")
	}
	var awsSessionFn AwsSessionFn
	if config.AwsProfileName != "" {
		awsSessionFn = awsSessionFromProfileFn
	} else {
		awsSessionFn = awsSessionFromEnvVarsFn
	}
	return kms.New(awsSessionFn(config)), nil
}

type EVMKMSClient struct {
	Client KMSClient
	KeyID  string
}

func NewEVMKMSClient(client KMSClient, keyID string) *EVMKMSClient {
	return &EVMKMSClient{
		Client: client,
		KeyID:  keyID,
	}
}

func (c *EVMKMSClient) GetKMSTransactOpts(ctx context.Context, chainID *big.Int) (*bind.TransactOpts, error) {
	ecdsaPublicKey, err := c.GetECDSAPublicKey()
	if err != nil {
		return nil, err
	}

	pubKeyBytes := secp256k1.S256().Marshal(ecdsaPublicKey.X, ecdsaPublicKey.Y)
	keyAddr := crypto.PubkeyToAddress(*ecdsaPublicKey)
	if chainID == nil {
		return nil, errors.New("chainID is required")
	}
	signer := types.LatestSignerForChainID(chainID)

	signerFn := func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if address != keyAddr {
			return nil, bind.ErrNotAuthorized
		}

		txHashBytes := signer.Hash(tx).Bytes()

		mType := kms.MessageTypeDigest
		algo := kms.SigningAlgorithmSpecEcdsaSha256
		signOutput, err := c.Client.Sign(
			&kms.SignInput{
				KeyId:            &c.KeyID,
				SigningAlgorithm: &algo,
				MessageType:      &mType,
				Message:          txHashBytes,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to call kms.Sign() on transaction: %w", err)
		}

		ethSig, err := kmsToEthSig(signOutput.Signature, pubKeyBytes, txHashBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to convert KMS signature to Ethereum signature: %w", err)
		}

		return tx.WithSignature(signer, ethSig)
	}

	return &bind.TransactOpts{
		From:    keyAddr,
		Signer:  signerFn,
		Context: ctx,
	}, nil
}

// GetECDSAPublicKey retrieves the public key from KMS and converts it to its ECDSA representation.
func (c *EVMKMSClient) GetECDSAPublicKey() (*ecdsa.PublicKey, error) {
	getPubKeyOutput, err := c.Client.GetPublicKey(&kms.GetPublicKeyInput{
		KeyId: aws.String(c.KeyID),
	})
	if err != nil {
		return nil, fmt.Errorf("can not get public key from KMS for KeyId=%s: %w", c.KeyID, err)
	}

	var asn1pubKeyInfo asn1SubjectPublicKeyInfo
	_, err = asn1.Unmarshal(getPubKeyOutput.PublicKey, &asn1pubKeyInfo)
	if err != nil {
		return nil, fmt.Errorf("can not parse asn1 public key for KeyId=%s: %w", c.KeyID, err)
	}

	pubKey, err := crypto.UnmarshalPubkey(asn1pubKeyInfo.SubjectPublicKey.Bytes)
	if err != nil {
		return nil, fmt.Errorf("can not unmarshal public key bytes: %w", err)
	}
	return pubKey, nil
}

func kmsToEthSig(kmsSig, ecdsaPubKeyBytes, hash []byte) ([]byte, error) {
	var asn1Sig asn1ECDSASig
	_, err := asn1.Unmarshal(kmsSig, &asn1Sig)
	if err != nil {
		return nil, err
	}

	rBytes := asn1Sig.R.Bytes
	sBytes := asn1Sig.S.Bytes

	// Adjust S value from signature to match Eth standard.
	//   See: https://aws.amazon.com/blogs/database/part2-use-aws-kms-to-securely-manage-ethereum-accounts/
	// "After we extract r and s successfully, we have to test if the value of s is greater than secp256k1n/2 as
	// specified in EIP-2 and flip it if required."
	sBigInt := new(big.Int).SetBytes(sBytes)
	if sBigInt.Cmp(secp256k1HalfN) > 0 {
		sBytes = new(big.Int).Sub(secp256k1N, sBigInt).Bytes()
	}

	return recoverEthSignature(ecdsaPubKeyBytes, hash, rBytes, sBytes)
}

// See: https://aws.amazon.com/blogs/database/part2-use-aws-kms-to-securely-manage-ethereum-accounts/
func recoverEthSignature(expectedPublicKeyBytes, txHash, r, s []byte) ([]byte, error) {
	rsSig := append(padTo32Bytes(r), padTo32Bytes(s)...)
	ethSig := append(rsSig, []byte{0}...)

	recoveredPublicKeyBytes, err := crypto.Ecrecover(txHash, ethSig)
	if err != nil {
		return nil, fmt.Errorf("failing to call Ecrecover: %w", err)
	}

	if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(expectedPublicKeyBytes) {
		ethSig = append(rsSig, []byte{1}...)
		recoveredPublicKeyBytes, err = crypto.Ecrecover(txHash, ethSig)
		if err != nil {
			return nil, fmt.Errorf("failing to call Ecrecover: %w", err)
		}

		if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(expectedPublicKeyBytes) {
			return nil, errors.New("can not reconstruct public key from sig")
		}
	}

	return ethSig, nil
}

func padTo32Bytes(buffer []byte) []byte {
	buffer = bytes.TrimLeft(buffer, "\x00")
	for len(buffer) < 32 {
		zeroBuf := []byte{0}
		buffer = append(zeroBuf, buffer...)
	}
	return buffer
}

type AwsSessionFn func(config KMS) *session.Session

var awsSessionFromEnvVarsFn = func(config KMS) *session.Session {
	return session.Must(
		session.NewSession(&aws.Config{
			Region:                        aws.String(config.KmsDeployerKeyRegion),
			CredentialsChainVerboseErrors: aws.Bool(true),
		}))
}

var awsSessionFromProfileFn = func(config KMS) *session.Session {
	return session.Must(
		session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Profile:           config.AwsProfileName,
			Config: aws.Config{
				Region:                        aws.String(config.KmsDeployerKeyRegion),
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		}))
}

func KMSConfigFromEnvVars() (KMS, error) {
	var config KMS
	var exists bool
	config.KmsDeployerKeyId, exists = os.LookupEnv("KMS_DEPLOYER_KEY_ID")
	if !exists {
		return config, errors.New("KMS_DEPLOYER_KEY_ID is required")
	}
	config.KmsDeployerKeyRegion, exists = os.LookupEnv("KMS_DEPLOYER_KEY_REGION")
	if !exists {
		return config, errors.New("KMS_DEPLOYER_KEY_REGION is required")
	}
	config.AwsProfileName, exists = os.LookupEnv("AWS_PROFILE")
	if !exists {
		return config, errors.New("AWS_PROFILE is required")
	}
	return config, nil
}
