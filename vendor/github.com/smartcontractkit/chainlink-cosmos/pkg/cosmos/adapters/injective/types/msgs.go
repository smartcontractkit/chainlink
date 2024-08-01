package types

import (
	errors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateFeed             = "createFeed"
	TypeMsgUpdateFeed             = "updateFeed"
	TypeMsgTransmit               = "transmit"
	TypeMsgFundFeedRewardPool     = "fundFeedRewardPool"
	TypeMsgWithdrawFeedRewardPool = "withdrawFeedRewardPool"
	TypeMsgSetPayees              = "setPayees"
	TypeMsgTransferPayeeship      = "transferPayeeship"
	TypeMsgAcceptPayeeship        = "acceptPayeeship"
)

var (
	_ sdk.Msg = &MsgCreateFeed{}
	_ sdk.Msg = &MsgUpdateFeed{}
	_ sdk.Msg = &MsgTransmit{}
	_ sdk.Msg = &MsgFundFeedRewardPool{}
	_ sdk.Msg = &MsgWithdrawFeedRewardPool{}
	_ sdk.Msg = &MsgSetPayees{}
	_ sdk.Msg = &MsgTransferPayeeship{}
	_ sdk.Msg = &MsgAcceptPayeeship{}
)

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgCreateFeed) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgCreateFeed) Type() string { return TypeMsgCreateFeed }

// ValidateBasic implements the sdk.Msg interface for MsgCreateFeed.
func (msg MsgCreateFeed) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	return msg.Config.ValidateBasic()
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgCreateFeed) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgCreateFeed) GetSigners() []sdk.AccAddress {
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgUpdateFeed) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgUpdateFeed) Type() string { return TypeMsgUpdateFeed }

// ValidateBasic implements the sdk.Msg interface for MsgUpdateFeed.
func (msg MsgUpdateFeed) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if msg.FeedId == "" || len(msg.FeedId) > FeedIDMaxLength {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "feedId not valid")
	}

	seenTransmitters := make(map[string]struct{}, len(msg.Transmitters))
	for _, transmitter := range msg.Transmitters {
		addr, err := sdk.AccAddressFromBech32(transmitter)
		if err != nil {
			return err
		}

		if _, ok := seenTransmitters[addr.String()]; ok {
			return ErrRepeatedAddress
		}
		seenTransmitters[addr.String()] = struct{}{}
	}

	seenSigners := make(map[string]struct{}, len(msg.Signers))
	for _, signer := range msg.Signers {
		addr, err := sdk.AccAddressFromBech32(signer)
		if err != nil {
			return err
		}

		if _, ok := seenSigners[addr.String()]; ok {
			return ErrRepeatedAddress
		}
		seenSigners[addr.String()] = struct{}{}
	}

	if msg.LinkPerTransmission != nil {
		if msg.LinkPerTransmission.IsNil() || !msg.LinkPerTransmission.IsPositive() {
			return errors.Wrap(ErrIncorrectConfig, "LinkPerTransmission must be positive")
		}
	}

	if msg.LinkPerObservation != nil {
		if msg.LinkPerObservation.IsNil() || !msg.LinkPerObservation.IsPositive() {
			return errors.Wrap(ErrIncorrectConfig, "LinkPerObservation must be positive")
		}
	}

	if msg.LinkDenom == "" {
		return sdkerrors.ErrInvalidCoins
	}

	if msg.FeedAdmin != "" {
		if _, err := sdk.AccAddressFromBech32(msg.FeedAdmin); err != nil {
			return err
		}
	}

	if msg.BillingAdmin != "" {
		if _, err := sdk.AccAddressFromBech32(msg.BillingAdmin); err != nil {
			return err
		}
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgUpdateFeed) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgUpdateFeed) GetSigners() []sdk.AccAddress {
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgTransmit) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgTransmit) Type() string { return TypeMsgTransmit }

// ValidateBasic implements the sdk.Msg interface for MsgTransmit.
func (msg MsgTransmit) ValidateBasic() error {
	if len(msg.Transmitter) == 0 {
		return ErrNoTransmitter
	}

	if len(msg.ConfigDigest) == 0 {
		return errors.Wrap(ErrIncorrectTransmissionData, "missing config digest")
	} else if len(msg.FeedId) == 0 {
		return errors.Wrap(ErrIncorrectTransmissionData, "missing feed hash")
	} else if msg.Report == nil {
		return errors.Wrap(ErrIncorrectTransmissionData, "missing report")
	}

	if len(msg.Report.Observers) > MaxNumOracles {
		return errors.Wrap(ErrIncorrectTransmissionData, "too many observers")
	} else if len(msg.Report.Observations) != len(msg.Report.Observers) {
		return errors.Wrap(ErrIncorrectTransmissionData, "wrong observations count")
	}

	if len(msg.Report.Observations) > MaxNumOracles {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "num observations out of bounds")
	}

	for i := 0; i < len(msg.Report.Observations)-1; i++ {
		inOrder := msg.Report.Observations[i].LTE(msg.Report.Observations[i+1])
		if !inOrder {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "observations not sorted")
		}
	}
	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgTransmit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgTransmit) GetSigners() []sdk.AccAddress {
	transmitter := sdk.MustAccAddressFromBech32(msg.Transmitter)
	return []sdk.AccAddress{transmitter}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgFundFeedRewardPool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgFundFeedRewardPool) Type() string { return TypeMsgFundFeedRewardPool }

// ValidateBasic implements the sdk.Msg interface for MsgFundFeedRewardPool.
func (msg MsgFundFeedRewardPool) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if msg.FeedId == "" || len(msg.FeedId) > FeedIDMaxLength {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "feedId not valid")
	}

	if !msg.Amount.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgFundFeedRewardPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgFundFeedRewardPool) GetSigners() []sdk.AccAddress {
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgWithdrawFeedRewardPool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgWithdrawFeedRewardPool) Type() string { return TypeMsgWithdrawFeedRewardPool }

// ValidateBasic implements the sdk.Msg interface for MsgWithdrawFeedRewardPool.
func (msg MsgWithdrawFeedRewardPool) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if msg.FeedId == "" || len(msg.FeedId) > FeedIDMaxLength {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "feedId not valid")
	}

	if !msg.Amount.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgWithdrawFeedRewardPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgWithdrawFeedRewardPool) GetSigners() []sdk.AccAddress {
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgSetPayees) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgSetPayees) Type() string { return TypeMsgSetPayees }

// ValidateBasic implements the sdk.Msg interface for MsgSetPayees.
func (msg MsgSetPayees) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if msg.FeedId == "" || len(msg.FeedId) > FeedIDMaxLength {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "feedId not valid")
	}

	if len(msg.Transmitters) != len(msg.Payees) || len(msg.Payees) == 0 {
		return ErrInvalidPayees
	}

	seenTransmitters := make(map[string]struct{}, len(msg.Transmitters))
	for _, transmitter := range msg.Transmitters {
		addr, err := sdk.AccAddressFromBech32(transmitter)
		if err != nil {
			return err
		}

		if _, ok := seenTransmitters[addr.String()]; ok {
			return ErrRepeatedAddress
		}
		seenTransmitters[addr.String()] = struct{}{}
	}

	seenPayees := make(map[string]struct{}, len(msg.Payees))
	for _, payee := range msg.Payees {
		addr, err := sdk.AccAddressFromBech32(payee)
		if err != nil {
			return err
		}

		if _, ok := seenPayees[addr.String()]; ok {
			return ErrRepeatedAddress
		}
		seenPayees[addr.String()] = struct{}{}
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgSetPayees) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgSetPayees) GetSigners() []sdk.AccAddress {
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgTransferPayeeship) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgTransferPayeeship) Type() string { return TypeMsgTransferPayeeship }

// ValidateBasic implements the sdk.Msg interface for MsgTransferPayeeship.
func (msg MsgTransferPayeeship) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if msg.FeedId == "" || len(msg.FeedId) > FeedIDMaxLength {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "feedId not valid")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Transmitter); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Transmitter)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Proposed); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Proposed)
	}

	if msg.Transmitter == msg.Proposed {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "proposed cannot be the same as transmitter")
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgTransferPayeeship) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgTransferPayeeship) GetSigners() []sdk.AccAddress {
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgAcceptPayeeship) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgAcceptPayeeship) Type() string { return TypeMsgAcceptPayeeship }

// ValidateBasic implements the sdk.Msg interface for MsgAcceptPayeeship.
func (msg MsgAcceptPayeeship) ValidateBasic() error {
	if msg.Payee == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Payee)
	}

	if msg.FeedId == "" || len(msg.FeedId) > FeedIDMaxLength {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "feedId not valid")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Transmitter); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Transmitter)
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgAcceptPayeeship) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgAcceptPayeeship) GetSigners() []sdk.AccAddress {
	sender := sdk.MustAccAddressFromBech32(msg.Payee)
	return []sdk.AccAddress{sender}
}
