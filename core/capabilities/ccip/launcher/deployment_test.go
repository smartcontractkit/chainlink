package launcher

import (
	"errors"
	"testing"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	mocktypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types/mocks"

	"github.com/stretchr/testify/require"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
)

func Test_ccipDeployment_Close(t *testing.T) {
	type args struct {
		commitBlue  *mocktypes.CCIPOracle
		commitGreen *mocktypes.CCIPOracle
		execBlue    *mocktypes.CCIPOracle
		execGreen   *mocktypes.CCIPOracle
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args)
		asserts func(t *testing.T, args args)
		wantErr bool
	}{
		{
			name: "no errors, active only",
			args: args{
				commitBlue:  mocktypes.NewCCIPOracle(t),
				commitGreen: nil,
				execBlue:    mocktypes.NewCCIPOracle(t),
				execGreen:   nil,
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
			name: "no errors, active and candidate",
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
			name: "error on commit active",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ccipDeployment{
				commit: activeCandidateDeployment{
					active: tt.args.commitBlue,
				},
				exec: activeCandidateDeployment{
					active: tt.args.execBlue,
				},
			}
			if tt.args.commitGreen != nil {
				c.commit.candidate = tt.args.commitGreen
			}

			if tt.args.execGreen != nil {
				c.exec.candidate = tt.args.execGreen
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
		commitBlue *mocktypes.CCIPOracle
		execBlue   *mocktypes.CCIPOracle
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args)
		asserts func(t *testing.T, args args)
		wantErr bool
	}{
		{
			name: "no errors",
			args: args{
				commitBlue: mocktypes.NewCCIPOracle(t),
				execBlue:   mocktypes.NewCCIPOracle(t),
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
			name: "error on commit active",
			args: args{
				commitBlue: mocktypes.NewCCIPOracle(t),
				execBlue:   mocktypes.NewCCIPOracle(t),
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
			name: "error on exec active",
			args: args{
				commitBlue: mocktypes.NewCCIPOracle(t),
				execBlue:   mocktypes.NewCCIPOracle(t),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ccipDeployment{
				commit: activeCandidateDeployment{
					active: tt.args.commitBlue,
				},
				exec: activeCandidateDeployment{
					active: tt.args.execBlue,
				},
			}

			tt.expect(t, tt.args)
			defer tt.asserts(t, tt.args)
			err := c.StartActive()
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
		commitBlue *mocktypes.CCIPOracle
		execBlue   *mocktypes.CCIPOracle
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args)
		asserts func(t *testing.T, args args)
		wantErr bool
	}{
		{
			name: "no errors",
			args: args{
				commitBlue: mocktypes.NewCCIPOracle(t),
				execBlue:   mocktypes.NewCCIPOracle(t),
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
			name: "error on commit active",
			args: args{
				commitBlue: mocktypes.NewCCIPOracle(t),
				execBlue:   mocktypes.NewCCIPOracle(t),
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
			name: "error on exec active",
			args: args{
				commitBlue: mocktypes.NewCCIPOracle(t),
				execBlue:   mocktypes.NewCCIPOracle(t),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ccipDeployment{
				commit: activeCandidateDeployment{
					active: tt.args.commitBlue,
				},
				exec: activeCandidateDeployment{
					active: tt.args.execBlue,
				},
			}

			tt.expect(t, tt.args)
			defer tt.asserts(t, tt.args)
			err := c.CloseActive()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_ccipDeployment_HandleBlueGreen_PrevDeploymentNil(t *testing.T) {
	require.Error(t, (&ccipDeployment{}).TransitionDeployment(nil))
}

func Test_ccipDeployment_HandleBlueGreen(t *testing.T) {
	type args struct {
		commitBlue  *mocktypes.CCIPOracle
		commitGreen *mocktypes.CCIPOracle
		execBlue    *mocktypes.CCIPOracle
		execGreen   *mocktypes.CCIPOracle
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
			name: "promotion active to candidate",
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
			name: "new candidate deployment",
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
			name: "error on commit candidate start",
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
			name: "error on exec candidate start",
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
			name: "invalid active-candidate deployment transition commit: both prev and future deployment have candidate",
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
			name: "invalid active-candidate deployment transition exec: both prev and future exec deployment have candidate",
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
				commit: activeCandidateDeployment{
					active: tt.argsFutureDeployment.commitBlue,
				},
				exec: activeCandidateDeployment{
					active: tt.argsFutureDeployment.execBlue,
				},
			}
			if tt.argsFutureDeployment.commitGreen != nil {
				futDeployment.commit.candidate = tt.argsFutureDeployment.commitGreen
			}
			if tt.argsFutureDeployment.execGreen != nil {
				futDeployment.exec.candidate = tt.argsFutureDeployment.execGreen
			}

			prevDeployment := &ccipDeployment{
				commit: activeCandidateDeployment{
					active: tt.argsPrevDeployment.commitBlue,
				},
				exec: activeCandidateDeployment{
					active: tt.argsPrevDeployment.execBlue,
				},
			}
			if tt.argsPrevDeployment.commitGreen != nil {
				prevDeployment.commit.candidate = tt.argsPrevDeployment.commitGreen
			}
			if tt.argsPrevDeployment.execGreen != nil {
				prevDeployment.exec.candidate = tt.argsPrevDeployment.execGreen
			}

			tt.expect(t, tt.argsFutureDeployment, tt.argsPrevDeployment)
			defer tt.asserts(t, tt.argsFutureDeployment, tt.argsPrevDeployment)
			err := futDeployment.TransitionDeployment(prevDeployment)
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
			"prev deployment only active",
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
					{}, {},
				},
				prevDeployment: ccipDeployment{
					commit: activeCandidateDeployment{
						active: mocktypes.NewCCIPOracle(t),
					},
				},
			},
			true,
		},
		{
			"candidate -> active promotion",
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
					{},
				},
				prevDeployment: ccipDeployment{
					commit: activeCandidateDeployment{
						active:    mocktypes.NewCCIPOracle(t),
						candidate: mocktypes.NewCCIPOracle(t),
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNewCandidateInstance(tt.args.pluginType, tt.args.ocrConfigs, tt.args.prevDeployment)
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
			"prev deployment only active",
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
					{}, {},
				},
				prevDeployment: ccipDeployment{
					commit: activeCandidateDeployment{
						active: mocktypes.NewCCIPOracle(t),
					},
				},
			},
			false,
		},
		{
			"candidate -> active promotion",
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
					{},
				},
				prevDeployment: ccipDeployment{
					commit: activeCandidateDeployment{
						active:    mocktypes.NewCCIPOracle(t),
						candidate: mocktypes.NewCCIPOracle(t),
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
		commit activeCandidateDeployment
		exec   activeCandidateDeployment
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
			"commit candidate present",
			fields{
				commit: activeCandidateDeployment{
					active:    mocktypes.NewCCIPOracle(t),
					candidate: mocktypes.NewCCIPOracle(t),
				},
			},
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
			},
			true,
		},
		{
			"commit candidate not present",
			fields{
				commit: activeCandidateDeployment{
					active: mocktypes.NewCCIPOracle(t),
				},
			},
			args{
				pluginType: cctypes.PluginTypeCCIPCommit,
			},
			false,
		},
		{
			"exec candidate present",
			fields{
				exec: activeCandidateDeployment{
					active:    mocktypes.NewCCIPOracle(t),
					candidate: mocktypes.NewCCIPOracle(t),
				},
			},
			args{
				pluginType: cctypes.PluginTypeCCIPExec,
			},
			true,
		},
		{
			"exec candidate not present",
			fields{
				exec: activeCandidateDeployment{
					active: mocktypes.NewCCIPOracle(t),
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
			if tt.fields.commit.active != nil {
				c.commit.active = tt.fields.commit.active
			}
			if tt.fields.commit.candidate != nil {
				c.commit.candidate = tt.fields.commit.candidate
			}
			if tt.fields.exec.active != nil {
				c.exec.active = tt.fields.exec.active
			}
			if tt.fields.exec.candidate != nil {
				c.exec.candidate = tt.fields.exec.candidate
			}
			got := c.HasCandidateInstance(tt.args.pluginType)
			require.Equal(t, tt.want, got)
		})
	}
}
