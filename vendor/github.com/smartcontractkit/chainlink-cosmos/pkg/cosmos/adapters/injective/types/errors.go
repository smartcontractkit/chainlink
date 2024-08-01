package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrStaleReport               = sdkerrors.Register(ModuleName, 1, "stale report")
	ErrIncompleteProposal        = sdkerrors.Register(ModuleName, 2, "incomplete proposal")
	ErrRepeatedAddress           = sdkerrors.Register(ModuleName, 3, "repeated oracle address")
	ErrTooManySigners            = sdkerrors.Register(ModuleName, 4, "too many signers")
	ErrIncorrectConfig           = sdkerrors.Register(ModuleName, 5, "incorrect config")
	ErrConfigDigestNotMatch      = sdkerrors.Register(ModuleName, 6, "config digest doesn't match")
	ErrWrongNumberOfSignatures   = sdkerrors.Register(ModuleName, 7, "wrong number of signatures")
	ErrIncorrectSignature        = sdkerrors.Register(ModuleName, 8, "incorrect signature")
	ErrNoTransmitter             = sdkerrors.Register(ModuleName, 9, "no transmitter specified")
	ErrIncorrectTransmissionData = sdkerrors.Register(ModuleName, 10, "incorrect transmission data")
	ErrNoTransmissionsFound      = sdkerrors.Register(ModuleName, 11, "no transmissions found")
	ErrMedianValueOutOfBounds    = sdkerrors.Register(ModuleName, 12, "median value is out of bounds")
	ErrIncorrectRewardPoolDenom  = sdkerrors.Register(ModuleName, 13, "LINK denom doesn't match")
	ErrNoRewardPool              = sdkerrors.Register(ModuleName, 14, "Reward Pool doesn't exist")
	ErrInvalidPayees             = sdkerrors.Register(ModuleName, 15, "wrong number of payees and transmitters")
	ErrModuleAdminRestricted     = sdkerrors.Register(ModuleName, 16, "action is restricted to the module admin")
	ErrFeedAlreadyExists         = sdkerrors.Register(ModuleName, 17, "feed already exists")
	ErrFeedDoesntExists          = sdkerrors.Register(ModuleName, 19, "feed doesnt exists")
	ErrAdminRestricted           = sdkerrors.Register(ModuleName, 20, "action is admin-restricted")
	ErrInsufficientRewardPool    = sdkerrors.Register(ModuleName, 21, "insufficient reward pool")
	ErrPayeeAlreadySet           = sdkerrors.Register(ModuleName, 22, "payee already set")
	ErrPayeeRestricted           = sdkerrors.Register(ModuleName, 23, "action is payee-restricted")
	ErrFeedConfigNotFound        = sdkerrors.Register(ModuleName, 24, "feed config not found")
)
