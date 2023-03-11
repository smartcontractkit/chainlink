package evm

import (
	"math/big"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_decodeValue(t *testing.T) {
	type args struct {
		nodeKey string
		input   []byte
	}
	upkeepID := big.NewInt(100)
	upkeepIdBytes := upkeepID.Bytes()
	apiKey := "THIS_IS_API_KEY"
	plaintext := [2]interface{}{
		upkeepIdBytes,
		apiKey,
	}
	marshal, err := cbor.Marshal(plaintext)
	assert.Nil(t, err, t.Name())
	// TODO test encryption when we have that sorted out
	tests := []struct {
		name    string
		args    args
		want    DecryptedValue
		wantErr error
	}{
		{
			name: "success",
			args: args{
				nodeKey: "dumb-key",
				input:   marshal,
			},
			want: DecryptedValue{
				UpkeepID: upkeepID,
				Value:    apiKey,
			},
		},
		{
			name: "error",
			args: args{
				nodeKey: "dumb-key",
				input:   []byte("blahblah"),
			},
			wantErr: errors.New("cbor: cannot unmarshal UTF-8 text string into Go value of type [2]interface {}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := decodeValue(tt.args.nodeKey, tt.args.input)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), gotErr.Error(), tt.name)
				assert.NotNil(t, gotErr, tt.name)
			}
			assert.Equalf(t, tt.want, got, "decodeValue(%v, %v)", tt.args.nodeKey, tt.args.input)
		})
	}
}

func Test_getAPIKeys(t *testing.T) {
	upkeepID := big.NewInt(100)
	upkeepIdBytes := upkeepID.Bytes()
	apiKey := "THIS_IS_API_KEY"
	plaintext := [2]interface{}{
		upkeepIdBytes,
		apiKey,
	}
	valMarshal, err := cbor.Marshal(plaintext)
	require.Nil(t, err, "plaintext needs to marshal")

	offchainConfig := OffchainAPIKeys{Keys: []Key{
		{
			Name:  "Authorization",
			Type:  "HeAdEr",
			Value: valMarshal,
		},
	}}
	marshal, err := cbor.Marshal(offchainConfig)
	require.Nil(t, err, "offchainConfig needs to marshal")
	offchainConfigBadDecrpt := OffchainAPIKeys{Keys: []Key{
		{
			Name:  "Authorization",
			Type:  "HeAdEr",
			Value: []byte{1, 1, 1, 1, 1},
		},
	}}
	marshalDecryptErr, err := cbor.Marshal(offchainConfigBadDecrpt)
	require.Nil(t, err, "offchainConfig needs to marshal")
	result := OffchainAPIKeys{Keys: []Key{
		{
			Name:       "Authorization",
			Type:       "HeAdEr",
			Value:      valMarshal,
			DecryptVal: apiKey,
		},
	}}
	type args struct {
		upkeepID       *big.Int
		offchainConfig []byte
	}
	tests := []struct {
		name    string
		args    args
		want    OffchainAPIKeys
		wantErr error
	}{
		{
			name: "success",
			args: args{
				upkeepID:       upkeepID,
				offchainConfig: marshal,
			},
			want: result,
		},
		{
			name: "fail - not right upkeepID",
			args: args{
				upkeepID:       big.NewInt(1),
				offchainConfig: marshal,
			},
			want: offchainConfig,
		},
		{
			name: "fail - decrypt error",
			args: args{
				upkeepID:       upkeepID,
				offchainConfig: marshalDecryptErr,
			},
			want: offchainConfigBadDecrpt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := getAPIKeys(tt.args.upkeepID, tt.args.offchainConfig)
			if tt.wantErr != nil {
				assert.NotNil(t, gotErr, "expected error")
				assert.Equal(t, tt.wantErr.Error(), gotErr.Error(), tt.name)
			} else {
				assert.Nil(t, err, "should be no errors")
			}
			assert.Equalf(t, tt.want, got, "getAPIKeys(%v, %v)", tt.args.upkeepID, tt.args.offchainConfig)
		})
	}
}
