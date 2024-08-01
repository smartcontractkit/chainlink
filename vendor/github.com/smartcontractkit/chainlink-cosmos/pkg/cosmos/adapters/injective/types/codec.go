package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterLegacyAminoCodec registers the necessary modules/ocr interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateFeed{}, "ocr/MsgCreateFeed", nil)
	cdc.RegisterConcrete(&MsgUpdateFeed{}, "ocr/MsgUpdateFeed", nil)
	cdc.RegisterConcrete(&MsgTransmit{}, "ocr/MsgTransmit", nil)
	cdc.RegisterConcrete(&MsgFundFeedRewardPool{}, "ocr/MsgFundFeedRewardPool", nil)
	cdc.RegisterConcrete(&MsgWithdrawFeedRewardPool{}, "ocr/MsgWithdrawFeedRewardPool", nil)
	cdc.RegisterConcrete(&MsgSetPayees{}, "ocr/MsgSetPayees", nil)
	cdc.RegisterConcrete(&MsgTransferPayeeship{}, "ocr/MsgTransferPayeeship", nil)
	cdc.RegisterConcrete(&MsgAcceptPayeeship{}, "ocr/MsgAcceptPayeeship", nil)

	cdc.RegisterConcrete(&SetConfigProposal{}, "ocr/SetConfigProposal", nil)
	cdc.RegisterConcrete(&SetBatchConfigProposal{}, "ocr/SetBatchConfigProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateFeed{},
		&MsgUpdateFeed{},
		&MsgTransmit{},
		&MsgFundFeedRewardPool{},
		&MsgWithdrawFeedRewardPool{},
		&MsgSetPayees{},
		&MsgTransferPayeeship{},
		&MsgAcceptPayeeship{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&SetConfigProposal{},
		&SetBatchConfigProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global modules/ocr module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/insurance and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
