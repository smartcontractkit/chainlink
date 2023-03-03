package contract

import (
	"go.dedis.ch/kyber/v3"
)

type EncryptionSecretKey kyber.Scalar

type EncryptionPublicKeys []kyber.Point

type SigningSecretKey kyber.Scalar

type SigningPublicKeys []kyber.Point
