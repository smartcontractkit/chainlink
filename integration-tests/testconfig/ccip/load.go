package ccip

import (
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/stretchr/testify/require"
)

const (
	TokenOnlyTransfer    string = "Token"
	DataOnlyTransfer     string = "Data"
	DataAndTokenTransfer string = "DataWithToken"
)

type LoadConfig struct {
	LoadDuration         *string
	MessageDetails       *[]MsgDetails
	RequestFrequency     *string
	CribEnvDirectory     *string
	NumDestinationChains *int
	TimeoutDuration      *string
	TestLabel            *string
}

func (l *LoadConfig) Validate(t *testing.T, e *deployment.Environment) {
	_, err := time.ParseDuration(*l.LoadDuration)
	require.NoError(t, err, "LoadDuration must be a valid duration")

	_, err = time.ParseDuration(*l.TimeoutDuration)
	require.NoError(t, err, "TimeoutDuration must be a valid duration")

	agg := 0
	for _, md := range *l.MessageDetails {
		require.NoError(t, md.Validate())
		agg += int(*md.Ratio)
	}
	require.Equal(t, 100, agg, "Sum of MessageDetails Ratios must be 100")

	require.GreaterOrEqual(t, *l.NumDestinationChains, 1, "NumDestinationChains must be greater than or equal to 1")
	require.GreaterOrEqual(t, len(e.Chains), *l.NumDestinationChains, "NumDestinationChains must be less than or equal to the number of chains in the environment")

}

func (l *LoadConfig) GetLoadDuration() time.Duration {
	ld, _ := time.ParseDuration(*l.LoadDuration)
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
	return pointer.GetString(m.MsgType) == TokenOnlyTransfer || pointer.GetString(m.MsgType) == DataAndTokenTransfer
}

func (m *MsgDetails) IsDataTransfer() bool {
	return pointer.GetString(m.MsgType) == DataOnlyTransfer || pointer.GetString(m.MsgType) == DataAndTokenTransfer
}

func (m *MsgDetails) Validate() error {
	if m == nil {
		return errors.New("msg details should be set")
	}
	if m.MsgType == nil {
		return errors.New("msg type should be set")
	}
	if pointer.GetString(m.MsgType) != DataOnlyTransfer &&
		pointer.GetString(m.MsgType) != TokenOnlyTransfer &&
		pointer.GetString(m.MsgType) != DataAndTokenTransfer {
		return errors.New(fmt.Sprintf("msg type should be - %s/%s/%s", DataOnlyTransfer, TokenOnlyTransfer, DataAndTokenTransfer))
	}

	if m.DestGasLimit == nil {
		return errors.New("dest gas limit should be set")
	}
	if *m.DestGasLimit < 0 {
		return errors.New("dest gas limit should be greater than 0")
	}

	if m.Ratio == nil {
		return errors.New("ratio should be set")
	}
	if *m.Ratio < 0 || *m.Ratio > 100 {
		return errors.New("ratio should be between 0 and 100")
	}

	if pointer.GetString(m.MsgType) == DataAndTokenTransfer {
		if m.DataLengthBytes == nil {
			return errors.New("data length should be set for data and token transfer")
		}
		if *m.DataLengthBytes < 0 {
			return errors.New("data length should be greater than 0")
		}
	}
	if pointer.GetString(m.MsgType) == DataOnlyTransfer {
		if m.DataLengthBytes == nil {
			return errors.New("data length should be set for data transfer")
		}
		if *m.DataLengthBytes < 0 {
			return errors.New("data length should be greater than 0")
		}
	}

	return nil
}
