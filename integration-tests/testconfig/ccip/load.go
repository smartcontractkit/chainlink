package ccip

import (
	"errors"
	"fmt"

	"testing"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/stretchr/testify/require"
)

const (
	TokenTransfer             string = "TokenTransfer"
	MessagingTransfer         string = "Messaging"
	ProgrammableTokenTransfer string = "ProgrammableTokenTransfer"
)

type LoadConfig struct {
	LoadDuration         *string
	ChaosMode            int
	RPCLatency           *string
	RPCJitter            *string
	MessageDetails       *[]MsgDetails
	RequestFrequency     *string
	CribEnvDirectory     *string
	NumDestinationChains *int
	TimeoutDuration      *string
	TestLabel            *string
	GasLimit             *uint64
	OOOExecution         *bool
}

const (
	ChaosModeNone = iota
	ChaosModeTypeRPCLatency
	ChaosModeTypeFull
)

func (l *LoadConfig) Validate(t *testing.T, e *deployment.Environment) {
	_, err := time.ParseDuration(*l.LoadDuration)
	require.NoError(t, err, "LoadDuration must be a valid duration")

	_, err = time.ParseDuration(*l.TimeoutDuration)
	require.NoError(t, err, "TimeoutDuration must be a valid duration")

	agg := 0
	for _, md := range *l.MessageDetails {
		require.NoError(t, md.Validate())
		agg += *md.Ratio
	}
	require.Equal(t, 100, agg, "Sum of MessageDetails Ratios must be 100")

	require.GreaterOrEqual(t, *l.NumDestinationChains, 1, "NumDestinationChains must be greater than or equal to 1")
	require.GreaterOrEqual(t, len(e.Chains), *l.NumDestinationChains, "NumDestinationChains must be less than or equal to the number of chains in the environment")
}

func (l *LoadConfig) GetLoadDuration() time.Duration {
	ld, _ := time.ParseDuration(*l.LoadDuration)
	return ld
}

func (l *LoadConfig) GetRPCLatency() time.Duration {
	ld, _ := time.ParseDuration(*l.RPCLatency)
	return ld
}

func (l *LoadConfig) GetRPCJitter() time.Duration {
	ld, _ := time.ParseDuration(*l.RPCJitter)
	return ld
}

func (l *LoadConfig) GetTimeoutDuration() time.Duration {
	ld, _ := time.ParseDuration(*l.TimeoutDuration)
	if ld == 0 {
		return 30 * time.Minute
	}
	return ld
}

type MsgDetails struct {
	MsgType         *string `toml:",omitempty"`
	DestGasLimit    *int64  `toml:",omitempty"`
	DataLengthBytes *int    `toml:",omitempty"`
	Ratio           *int    `toml:",omitempty"` // Percentage ratio of this message type (0-100)
}

func (m *MsgDetails) IsTokenTransfer() bool {
	return pointer.GetString(m.MsgType) == TokenTransfer || pointer.GetString(m.MsgType) == ProgrammableTokenTransfer
}

func (m *MsgDetails) IsTokenOnlyTransfer() bool {
	return pointer.GetString(m.MsgType) == TokenTransfer
}

func (m *MsgDetails) IsDataTransfer() bool {
	return pointer.GetString(m.MsgType) == MessagingTransfer || pointer.GetString(m.MsgType) == ProgrammableTokenTransfer
}

func (m *MsgDetails) Validate() error {
	if m == nil {
		return errors.New("msg details should be set")
	}
	if m.MsgType == nil {
		return errors.New("msg type should be set")
	}
	msgType := pointer.GetString(m.MsgType)
	if msgType != MessagingTransfer &&
		msgType != TokenTransfer &&
		msgType != ProgrammableTokenTransfer {
		return fmt.Errorf("msg type should be one of %s/%s/%s. Got %s",
			MessagingTransfer, TokenTransfer, ProgrammableTokenTransfer, msgType)
	}

	if m.Ratio == nil {
		return errors.New("ratio should be set")
	}
	if *m.Ratio < 0 || *m.Ratio > 100 {
		return errors.New("ratio should be between 0 and 100")
	}

	if msgType == ProgrammableTokenTransfer {
		if m.DataLengthBytes == nil {
			return errors.New("data length should be set for data and token transfer")
		}
		if *m.DataLengthBytes < 0 {
			return errors.New("data length should be greater than 0")
		}
	}
	if msgType == MessagingTransfer {
		if m.DataLengthBytes == nil {
			return errors.New("data length should be set for data transfer")
		}
		if *m.DataLengthBytes < 0 {
			return errors.New("data length should be greater than 0")
		}
	}

	return nil
}
