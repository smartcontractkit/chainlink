package launcher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types"
	mocktypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types/mocks"
)

func Test_ccipDeployment_Close(t *testing.T) {
	type args struct {
		commitBlue           *mocktypes.CCIPOracle
		commitBlueBootstrap  *mocktypes.CCIPOracle
		commitGreen          *mocktypes.CCIPOracle
		commitGreenBootstrap *mocktypes.CCIPOracle
		execBlue             *mocktypes.CCIPOracle
		execBlueBootstrap    *mocktypes.CCIPOracle
		execGreen            *mocktypes.CCIPOracle
		execGreenBootstrap   *mocktypes.CCIPOracle
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args)
		asserts func(t *testing.T, args args)
		wantErr bool
	}{
		{
			name: "no errors, blue only",
			args: args{
				commitBlue:           mocktypes.NewCCIPOracle(t),
				commitGreen:          nil,
				commitGreenBootstrap: nil,
				execBlue:             mocktypes.NewCCIPOracle(t),
				execGreen:            nil,
				execGreenBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.execBlue.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "no errors, blue and green",
			args: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: mocktypes.NewCCIPOracle(t),
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.commitGreen.On("Close").Return(nil).Once()
				args.execBlue.On("Close").Return(nil).Once()
				args.execGreen.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.commitGreen.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
				args.execGreen.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "error on commit blue",
			args: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: nil,
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(errors.New("failed")).Once()
				args.execBlue.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "bootstrap blue also closed",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.commitBlueBootstrap.On("Close").Return(nil).Once()
				args.execBlue.On("Close").Return(nil).Once()
				args.execBlueBootstrap.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.commitBlueBootstrap.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
				args.execBlueBootstrap.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "bootstrap green also closed",
			args: args{
				commitBlue:           mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap:  mocktypes.NewCCIPOracle(t),
				commitGreen:          mocktypes.NewCCIPOracle(t),
				commitGreenBootstrap: mocktypes.NewCCIPOracle(t),
				execBlue:             mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:    mocktypes.NewCCIPOracle(t),
				execGreen:            mocktypes.NewCCIPOracle(t),
				execGreenBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.commitBlueBootstrap.On("Close").Return(nil).Once()
				args.commitGreen.On("Close").Return(nil).Once()
				args.commitGreenBootstrap.On("Close").Return(nil).Once()
				args.execBlue.On("Close").Return(nil).Once()
				args.execBlueBootstrap.On("Close").Return(nil).Once()
				args.execGreen.On("Close").Return(nil).Once()
				args.execGreenBootstrap.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.commitBlueBootstrap.AssertExpectations(t)
				args.commitGreen.AssertExpectations(t)
				args.commitGreenBootstrap.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
				args.execBlueBootstrap.AssertExpectations(t)
				args.execGreen.AssertExpectations(t)
				args.execGreenBootstrap.AssertExpectations(t)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ccipDeployment{
				commit: blueGreenDeployment{
					blue: tt.args.commitBlue,
				},
				exec: blueGreenDeployment{
					blue: tt.args.execBlue,
				},
			}
			if tt.args.commitGreen != nil {
				c.commit.green = tt.args.commitGreen
			}
			if tt.args.commitBlueBootstrap != nil {
				c.commit.bootstrapBlue = tt.args.commitBlueBootstrap
			}
			if tt.args.commitGreenBootstrap != nil {
				c.commit.bootstrapGreen = tt.args.commitGreenBootstrap
			}

			if tt.args.execGreen != nil {
				c.exec.green = tt.args.execGreen
			}
			if tt.args.execBlueBootstrap != nil {
				c.exec.bootstrapBlue = tt.args.execBlueBootstrap
			}
			if tt.args.execGreenBootstrap != nil {
				c.exec.bootstrapGreen = tt.args.execGreenBootstrap
			}

			tt.expect(t, tt.args)
			defer tt.asserts(t, tt.args)
			err := c.Close()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_ccipDeployment_StartBlue(t *testing.T) {
	type args struct {
		commitBlue          *mocktypes.CCIPOracle
		commitBlueBootstrap *mocktypes.CCIPOracle
		execBlue            *mocktypes.CCIPOracle
		execBlueBootstrap   *mocktypes.CCIPOracle
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args)
		asserts func(t *testing.T, args args)
		wantErr bool
	}{
		{
			name: "no errors, no bootstrap",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: nil,
				execBlueBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Start").Return(nil).Once()
				args.execBlue.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "no errors, with bootstrap",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Start").Return(nil).Once()
				args.commitBlueBootstrap.On("Start").Return(nil).Once()
				args.execBlue.On("Start").Return(nil).Once()
				args.execBlueBootstrap.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.commitBlueBootstrap.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
				args.execBlueBootstrap.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "error on commit blue",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: nil,
				execBlueBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Start").Return(errors.New("failed")).Once()
				args.execBlue.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "error on exec blue",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: nil,
				execBlueBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Start").Return(nil).Once()
				args.execBlue.On("Start").Return(errors.New("failed")).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "error on commit blue bootstrap",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Start").Return(nil).Once()
				args.commitBlueBootstrap.On("Start").Return(errors.New("failed")).Once()
				args.execBlue.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.commitBlueBootstrap.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "error on exec blue bootstrap",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: nil,
				execBlueBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Start").Return(nil).Once()
				args.execBlue.On("Start").Return(nil).Once()
				args.execBlueBootstrap.On("Start").Return(errors.New("failed")).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
				args.execBlueBootstrap.AssertExpectations(t)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ccipDeployment{
				commit: blueGreenDeployment{
					blue: tt.args.commitBlue,
				},
				exec: blueGreenDeployment{
					blue: tt.args.execBlue,
				},
			}
			if tt.args.commitBlueBootstrap != nil {
				c.commit.bootstrapBlue = tt.args.commitBlueBootstrap
			}
			if tt.args.execBlueBootstrap != nil {
				c.exec.bootstrapBlue = tt.args.execBlueBootstrap
			}

			tt.expect(t, tt.args)
			defer tt.asserts(t, tt.args)
			err := c.StartBlue()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_ccipDeployment_CloseBlue(t *testing.T) {
	type args struct {
		commitBlue          *mocktypes.CCIPOracle
		commitBlueBootstrap *mocktypes.CCIPOracle
		execBlue            *mocktypes.CCIPOracle
		execBlueBootstrap   *mocktypes.CCIPOracle
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args)
		asserts func(t *testing.T, args args)
		wantErr bool
	}{
		{
			name: "no errors, no bootstrap",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: nil,
				execBlueBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.execBlue.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "no errors, with bootstrap",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.commitBlueBootstrap.On("Close").Return(nil).Once()
				args.execBlue.On("Close").Return(nil).Once()
				args.execBlueBootstrap.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.commitBlueBootstrap.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
				args.execBlueBootstrap.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "error on commit blue",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: nil,
				execBlueBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(errors.New("failed")).Once()
				args.execBlue.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "error on exec blue",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: nil,
				execBlueBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.execBlue.On("Close").Return(errors.New("failed")).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "error on commit blue bootstrap",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:   nil,
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.commitBlueBootstrap.On("Close").Return(errors.New("failed")).Once()
				args.execBlue.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.commitBlueBootstrap.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "error on exec blue bootstrap",
			args: args{
				commitBlue:          mocktypes.NewCCIPOracle(t),
				execBlue:            mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap: nil,
				execBlueBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args) {
				args.commitBlue.On("Close").Return(nil).Once()
				args.execBlue.On("Close").Return(nil).Once()
				args.execBlueBootstrap.On("Close").Return(errors.New("failed")).Once()
			},
			asserts: func(t *testing.T, args args) {
				args.commitBlue.AssertExpectations(t)
				args.execBlue.AssertExpectations(t)
				args.execBlueBootstrap.AssertExpectations(t)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ccipDeployment{
				commit: blueGreenDeployment{
					blue: tt.args.commitBlue,
				},
				exec: blueGreenDeployment{
					blue: tt.args.execBlue,
				},
			}
			if tt.args.commitBlueBootstrap != nil {
				c.commit.bootstrapBlue = tt.args.commitBlueBootstrap
			}
			if tt.args.execBlueBootstrap != nil {
				c.exec.bootstrapBlue = tt.args.execBlueBootstrap
			}

			tt.expect(t, tt.args)
			defer tt.asserts(t, tt.args)
			err := c.CloseBlue()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_ccipDeployment_HandleBlueGreen_PrevDeploymentNil(t *testing.T) {
	require.Error(t, (&ccipDeployment{}).HandleBlueGreen(nil))
}

func Test_ccipDeployment_HandleBlueGreen(t *testing.T) {
	type args struct {
		commitBlue           *mocktypes.CCIPOracle
		commitBlueBootstrap  *mocktypes.CCIPOracle
		commitGreen          *mocktypes.CCIPOracle
		commitGreenBootstrap *mocktypes.CCIPOracle
		execBlue             *mocktypes.CCIPOracle
		execBlueBootstrap    *mocktypes.CCIPOracle
		execGreen            *mocktypes.CCIPOracle
		execGreenBootstrap   *mocktypes.CCIPOracle
	}
	tests := []struct {
		name                 string
		argsPrevDeployment   args
		argsFutureDeployment args
		expect               func(t *testing.T, args args, argsPrevDeployment args)
		asserts              func(t *testing.T, args args, argsPrevDeployment args)
		wantErr              bool
	}{
		{
			name: "promotion blue to green, no bootstrap",
			argsPrevDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: mocktypes.NewCCIPOracle(t),
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			argsFutureDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: nil,
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   nil,
			},
			expect: func(t *testing.T, args args, argsPrevDeployment args) {
				argsPrevDeployment.commitBlue.On("Close").Return(nil).Once()
				argsPrevDeployment.execBlue.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
				argsPrevDeployment.commitBlue.AssertExpectations(t)
				argsPrevDeployment.execBlue.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "promotion blue to green, with bootstrap",
			argsPrevDeployment: args{
				commitBlue:           mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap:  mocktypes.NewCCIPOracle(t),
				commitGreen:          mocktypes.NewCCIPOracle(t),
				commitGreenBootstrap: mocktypes.NewCCIPOracle(t),
				execBlue:             mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:    mocktypes.NewCCIPOracle(t),
				execGreen:            mocktypes.NewCCIPOracle(t),
				execGreenBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			argsFutureDeployment: args{
				commitBlue:           mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap:  mocktypes.NewCCIPOracle(t),
				commitGreen:          nil,
				commitGreenBootstrap: nil,
				execBlue:             mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:    mocktypes.NewCCIPOracle(t),
				execGreen:            nil,
				execGreenBootstrap:   nil,
			},
			expect: func(t *testing.T, args args, argsPrevDeployment args) {
				argsPrevDeployment.commitBlue.On("Close").Return(nil).Once()
				argsPrevDeployment.commitBlueBootstrap.On("Close").Return(nil).Once()
				argsPrevDeployment.execBlue.On("Close").Return(nil).Once()
				argsPrevDeployment.execBlueBootstrap.On("Close").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
				argsPrevDeployment.commitBlue.AssertExpectations(t)
				argsPrevDeployment.commitBlueBootstrap.AssertExpectations(t)
				argsPrevDeployment.execBlue.AssertExpectations(t)
				argsPrevDeployment.execBlueBootstrap.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "new green deployment, no bootstrap",
			argsPrevDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: nil,
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   nil,
			},
			argsFutureDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: mocktypes.NewCCIPOracle(t),
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.On("Start").Return(nil).Once()
				args.execGreen.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.AssertExpectations(t)
				args.execGreen.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "new green deployment, with bootstrap",
			argsPrevDeployment: args{
				commitBlue:           mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap:  mocktypes.NewCCIPOracle(t),
				commitGreen:          nil,
				commitGreenBootstrap: nil,
				execBlue:             mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:    mocktypes.NewCCIPOracle(t),
				execGreen:            nil,
				execGreenBootstrap:   nil,
			},
			argsFutureDeployment: args{
				commitBlue:           mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap:  mocktypes.NewCCIPOracle(t),
				commitGreen:          mocktypes.NewCCIPOracle(t),
				commitGreenBootstrap: mocktypes.NewCCIPOracle(t),
				execBlue:             mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:    mocktypes.NewCCIPOracle(t),
				execGreen:            mocktypes.NewCCIPOracle(t),
				execGreenBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.On("Start").Return(nil).Once()
				args.commitGreenBootstrap.On("Start").Return(nil).Once()
				args.execGreen.On("Start").Return(nil).Once()
				args.execGreenBootstrap.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.AssertExpectations(t)
				args.commitGreenBootstrap.AssertExpectations(t)
				args.execGreen.AssertExpectations(t)
				args.execGreenBootstrap.AssertExpectations(t)
			},
			wantErr: false,
		},
		{
			name: "error on commit green start",
			argsPrevDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: nil,
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   nil,
			},
			argsFutureDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: mocktypes.NewCCIPOracle(t),
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.On("Start").Return(errors.New("failed")).Once()
				args.execGreen.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.AssertExpectations(t)
				args.execGreen.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "error on exec green start",
			argsPrevDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: nil,
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   nil,
			},
			argsFutureDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: mocktypes.NewCCIPOracle(t),
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.On("Start").Return(nil).Once()
				args.execGreen.On("Start").Return(errors.New("failed")).Once()
			},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.AssertExpectations(t)
				args.execGreen.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "error on commit green bootstrap start",
			argsPrevDeployment: args{
				commitBlue:           mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap:  mocktypes.NewCCIPOracle(t),
				commitGreen:          nil,
				commitGreenBootstrap: nil,
				execBlue:             mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:    mocktypes.NewCCIPOracle(t),
				execGreen:            nil,
				execGreenBootstrap:   nil,
			},
			argsFutureDeployment: args{
				commitBlue:           mocktypes.NewCCIPOracle(t),
				commitBlueBootstrap:  mocktypes.NewCCIPOracle(t),
				commitGreen:          mocktypes.NewCCIPOracle(t),
				commitGreenBootstrap: mocktypes.NewCCIPOracle(t),
				execBlue:             mocktypes.NewCCIPOracle(t),
				execBlueBootstrap:    mocktypes.NewCCIPOracle(t),
				execGreen:            mocktypes.NewCCIPOracle(t),
				execGreenBootstrap:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.On("Start").Return(nil).Once()
				args.commitGreenBootstrap.On("Start").Return(errors.New("failed")).Once()
				args.execGreen.On("Start").Return(nil).Once()
				args.execGreenBootstrap.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.AssertExpectations(t)
				args.commitGreenBootstrap.AssertExpectations(t)
				args.execGreen.AssertExpectations(t)
				args.execGreenBootstrap.AssertExpectations(t)
			},
			wantErr: true,
		},
		{
			name: "invalid blue-green deployment transition commit: both prev and future deployment have green",
			argsPrevDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: mocktypes.NewCCIPOracle(t),
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			argsFutureDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: mocktypes.NewCCIPOracle(t),
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			expect:  func(t *testing.T, args args, argsPrevDeployment args) {},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {},
			wantErr: true,
		},
		{
			name: "invalid blue-green deployment transition exec: both prev and future deployment have green",
			argsPrevDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: nil,
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			argsFutureDeployment: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: mocktypes.NewCCIPOracle(t),
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   mocktypes.NewCCIPOracle(t),
			},
			expect: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.On("Start").Return(nil).Once()
			},
			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
				args.commitGreen.AssertExpectations(t)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			futDeployment := &ccipDeployment{
				commit: blueGreenDeployment{
					blue: tt.argsFutureDeployment.commitBlue,
				},
				exec: blueGreenDeployment{
					blue: tt.argsFutureDeployment.execBlue,
				},
			}
			if tt.argsFutureDeployment.commitGreen != nil {
				futDeployment.commit.green = tt.argsFutureDeployment.commitGreen
			}
			if tt.argsFutureDeployment.commitBlueBootstrap != nil {
				futDeployment.commit.bootstrapBlue = tt.argsFutureDeployment.commitBlueBootstrap
			}
			if tt.argsFutureDeployment.commitGreenBootstrap != nil {
				futDeployment.commit.bootstrapGreen = tt.argsFutureDeployment.commitGreenBootstrap
			}
			if tt.argsFutureDeployment.execGreen != nil {
				futDeployment.exec.green = tt.argsFutureDeployment.execGreen
			}
			if tt.argsFutureDeployment.execBlueBootstrap != nil {
				futDeployment.exec.bootstrapBlue = tt.argsFutureDeployment.execBlueBootstrap
			}
			if tt.argsFutureDeployment.execGreenBootstrap != nil {
				futDeployment.exec.bootstrapGreen = tt.argsFutureDeployment.execGreenBootstrap
			}

			prevDeployment := &ccipDeployment{
				commit: blueGreenDeployment{
					blue: tt.argsPrevDeployment.commitBlue,
				},
				exec: blueGreenDeployment{
					blue: tt.argsPrevDeployment.execBlue,
				},
			}
			if tt.argsPrevDeployment.commitGreen != nil {
				prevDeployment.commit.green = tt.argsPrevDeployment.commitGreen
			}
			if tt.argsPrevDeployment.commitBlueBootstrap != nil {
				prevDeployment.commit.bootstrapBlue = tt.argsPrevDeployment.commitBlueBootstrap
			}
			if tt.argsPrevDeployment.commitGreenBootstrap != nil {
				prevDeployment.commit.bootstrapGreen = tt.argsPrevDeployment.commitGreenBootstrap
			}
			if tt.argsPrevDeployment.execGreen != nil {
				prevDeployment.exec.green = tt.argsPrevDeployment.execGreen
			}
			if tt.argsPrevDeployment.execBlueBootstrap != nil {
				prevDeployment.exec.bootstrapBlue = tt.argsPrevDeployment.execBlueBootstrap
			}
			if tt.argsPrevDeployment.execGreenBootstrap != nil {
				prevDeployment.exec.bootstrapGreen = tt.argsPrevDeployment.execGreenBootstrap
			}

			tt.expect(t, tt.argsFutureDeployment, tt.argsPrevDeployment)
			defer tt.asserts(t, tt.argsFutureDeployment, tt.argsPrevDeployment)
			err := futDeployment.HandleBlueGreen(prevDeployment)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_isNewGreenInstance(t *testing.T) {
	type args struct {
		pluginType     cctypes.PluginType
		ocrConfigs     []ccipreaderpkg.OCR3ConfigWithMeta
		prevDeployment ccipDeployment
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"prev deployment only blue",
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
					{}, {},
				},
				prevDeployment: ccipDeployment{
					commit: blueGreenDeployment{
						blue: mocktypes.NewCCIPOracle(t),
					},
				},
			},
			true,
		},
		{
			"green -> blue promotion",
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
					{},
				},
				prevDeployment: ccipDeployment{
					commit: blueGreenDeployment{
						blue:  mocktypes.NewCCIPOracle(t),
						green: mocktypes.NewCCIPOracle(t),
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNewGreenInstance(tt.args.pluginType, tt.args.ocrConfigs, tt.args.prevDeployment)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_isPromotion(t *testing.T) {
	type args struct {
		pluginType     cctypes.PluginType
		ocrConfigs     []ccipreaderpkg.OCR3ConfigWithMeta
		prevDeployment ccipDeployment
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"prev deployment only blue",
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
					{}, {},
				},
				prevDeployment: ccipDeployment{
					commit: blueGreenDeployment{
						blue: mocktypes.NewCCIPOracle(t),
					},
				},
			},
			false,
		},
		{
			"green -> blue promotion",
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
					{},
				},
				prevDeployment: ccipDeployment{
					commit: blueGreenDeployment{
						blue:  mocktypes.NewCCIPOracle(t),
						green: mocktypes.NewCCIPOracle(t),
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPromotion(tt.args.pluginType, tt.args.ocrConfigs, tt.args.prevDeployment); got != tt.want {
				t.Errorf("isPromotion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ccipDeployment_HasGreenInstance(t *testing.T) {
	type fields struct {
		commit blueGreenDeployment
		exec   blueGreenDeployment
	}
	type args struct {
		pluginType cctypes.PluginType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"commit green present",
			fields{
				commit: blueGreenDeployment{
					blue:  mocktypes.NewCCIPOracle(t),
					green: mocktypes.NewCCIPOracle(t),
				},
			},
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
			},
			true,
		},
		{
			"commit green not present",
			fields{
				commit: blueGreenDeployment{
					blue: mocktypes.NewCCIPOracle(t),
				},
			},
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
			},
			false,
		},
		{
			"exec green present",
			fields{
				exec: blueGreenDeployment{
					blue:  mocktypes.NewCCIPOracle(t),
					green: mocktypes.NewCCIPOracle(t),
				},
			},
			args{
				pluginType: cctypes.PluginTypeCCIPExec,
			},
			true,
		},
		{
			"exec green not present",
			fields{
				exec: blueGreenDeployment{
					blue: mocktypes.NewCCIPOracle(t),
				},
			},
			args{
				pluginType: cctypes.PluginTypeCCIPExec,
			},
			false,
		},
		{
			"invalid plugin type",
			fields{},
			args{
				pluginType: cctypes.PluginType(100),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ccipDeployment{}
			if tt.fields.commit.blue != nil {
				c.commit.blue = tt.fields.commit.blue
			}
			if tt.fields.commit.green != nil {
				c.commit.green = tt.fields.commit.green
			}
			if tt.fields.exec.blue != nil {
				c.exec.blue = tt.fields.exec.blue
			}
			if tt.fields.exec.green != nil {
				c.exec.green = tt.fields.exec.green
			}
			got := c.HasGreenInstance(tt.args.pluginType)
			require.Equal(t, tt.want, got)
		})
	}
}
