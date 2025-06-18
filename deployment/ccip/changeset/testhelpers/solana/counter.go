package solana

import (
	"fmt"

	ag_binary "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_ccip_receiver"
)

// The generated go bindings reference to BaseState which is marked as an account on it's own
// when the bindings attempt to decode it they check for distriminator, which is wrong since
// the it's embedded into the account as a field.
//
// We work around this using a custom decoder that only decodes the counter value.
type ReceiverCounter struct {
	Value uint8
}

func (obj ReceiverCounter) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(test_ccip_receiver.CounterDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Value` param:
	err = encoder.Encode(obj.Value)
	if err != nil {
		return err
	}
	return nil
}

func (obj *ReceiverCounter) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(test_ccip_receiver.CounterDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %+v, got %s",
				test_ccip_receiver.CounterDiscriminator,
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Value`:
	err = decoder.Decode(&obj.Value)
	if err != nil {
		return err
	}
	return nil
}
