package client

import (
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/web/sdk/internal/generated"
)

func TestDecodeInput(t *testing.T) {
	type args struct {
		in  any
		out any
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		errMessage string
	}{
		{
			name: "success",
			args: args{&JobDistributorInput{
				Name:      "name",
				Uri:       "uri",
				PublicKey: "publicKey",
			}, &generated.CreateFeedsManagerInput{}},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "non-pointer",
			args: args{&JobDistributorInput{
				Name:      "name",
				Uri:       "uri",
				PublicKey: "publicKey",
			}, generated.CreateFeedsManagerInput{}},
			wantErr:    true,
			errMessage: "out type must be a non-nil pointer",
		},
		{
			name: "incorrect type",
			args: args{&JobDistributorInput{
				Name:      "name",
				Uri:       "uri",
				PublicKey: "publicKey",
			}, generated.CreateFeedsManagerChainConfigInput{}},
			wantErr:    true,
			errMessage: "json: cannot unmarshal object into Go value of type generated.CreateFeedsManagerChainConfigInput",
		},
		{
			name: "success",
			args: args{&JobDistributorInput{
				Name:      "name",
				Uri:       "uri",
				PublicKey: "publicKey",
			}, &generated.UpdateFeedsManagerInput{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecodeInput(tt.args.in, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("DecodeInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
