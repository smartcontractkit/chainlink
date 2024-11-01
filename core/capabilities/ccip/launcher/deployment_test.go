package launcher

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"

	"github.com/stretchr/testify/require"

	mocktypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types/mocks"
)

func Test_ccipDeployment_Transitions(t *testing.T) {
	// we use a pointer to the oracle here for mock assertions
	type args struct {
		prevDeployment map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle
		currDeployment map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle
	}
	assertions := func(t *testing.T, args args) {
		for i, _ := range args.prevDeployment {
			args.prevDeployment[i].AssertExpectations(t)
		}
		for i, _ := range args.currDeployment {
			args.currDeployment[i].AssertExpectations(t)
		}
	}
	tests := []struct {
		name     string
		makeArgs func(t *testing.T) args
		expect   func(t *testing.T, args args)
		wantErr  bool
	}{
		{
			name: "all plugins are new",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 4 {
					prevP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}

				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 4 {
					currP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}
				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect: func(t *testing.T, args args) {
				for _, plugin := range args.prevDeployment {
					plugin.On("Close").Return(nil).Once()
				}
				for _, plugin := range args.currDeployment {
					plugin.On("Start").Return(nil).Once()
				}
			},
			wantErr: false,
		},
		{
			name: "no configs -> candidates",
			makeArgs: func(t *testing.T) args {
				prev := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)

				curr := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 2 {
					curr[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}
				return args{prevDeployment: prev, currDeployment: curr}
			},
			expect: func(t *testing.T, args args) {
				// When we are creating candidates, they should be started
				for _, plugin := range args.currDeployment {
					plugin.On("Start").Return(nil).Once()
				}
			},
			wantErr: false,
		},
		{
			name: "candidates -> active",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 2 {
					prevP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}

				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for digest, oracle := range prevP {
					currP[digest] = oracle
				}
				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect: func(t *testing.T, args args) {
				// if candidates are being promoted, there should be nothing to start or stop
				for _, plugin := range args.currDeployment {
					plugin.AssertNotCalled(t, "Start")
					plugin.AssertNotCalled(t, "Close")
				}
			},
			wantErr: false,
		},
		{
			name: "active -> active+candidates",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 2 {
					prevP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}

				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for digest, oracle := range prevP {
					currP[digest] = oracle
				}
				for range 2 {
					currP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}
				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect: func(t *testing.T, args args) {
				for digest, plugin := range args.currDeployment {
					// if it previously existed, there should be noop
					// if it's a new instance, it should be started
					if _, ok := args.prevDeployment[digest]; ok {
						plugin.AssertNotCalled(t, "Start")
						plugin.AssertNotCalled(t, "Close")
					} else {
						plugin.On("Start").Return(nil).Once()
					}
				}
			},
			wantErr: false,
		},
		{
			name: "active+candidate -> active",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 4 {
					prevP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}

				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				// copy two digests over
				i := 2
				for digest, oracle := range prevP {
					if i == 0 {
						continue
					}
					currP[digest] = oracle
					i--
				}
				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect: func(t *testing.T, args args) {
				for digest, plugin := range args.prevDeployment {

					if _, ok := args.currDeployment[digest]; !ok {
						// if the instance is no longer present, it should have been deleted
						plugin.On("Close").Return(nil).Once()
					} else {
						// otherwise, it should have been left alone
						plugin.AssertNotCalled(t, "Close")
						plugin.AssertNotCalled(t, "Start")
					}
				}
			},
			wantErr: false,
		},
		{
			name: "candidate -> different candidate",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 2 {
					prevP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}

				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 2 {
					currP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}

				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect: func(t *testing.T, args args) {
				for _, plugin := range args.prevDeployment {
					plugin.On("Close").Return(nil).Once()
				}
				for _, plugin := range args.currDeployment {
					plugin.On("Start").Return(nil).Once()
				}
			},
			wantErr: false,
		},
		{
			name: "close all instances",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 4 {
					prevP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}

				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)

				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect: func(t *testing.T, args args) {
				for _, plugin := range args.prevDeployment {
					plugin.On("Close").Return(nil).Once()
				}
			},
			wantErr: false,
		},
		{
			name: "start all instances",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)

				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 4 {
					currP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}

				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect: func(t *testing.T, args args) {
				for _, plugin := range args.currDeployment {
					plugin.On("Start").Return(nil).Once()
				}
			},
			wantErr: false,
		},
		{
			name: "should handle nil to nil",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect: func(t *testing.T, args args) {
				for _, plugin := range args.currDeployment {
					plugin.On("Start").Return(nil).Once()
				}
			},
			wantErr: false,
		},
		{
			name: "should throw error if there are more than 5 instances",
			makeArgs: func(t *testing.T) args {
				prevP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				currP := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 5 {
					currP[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}
				return args{prevDeployment: prevP, currDeployment: currP}
			},
			expect:  func(t *testing.T, args args) {},
			wantErr: true,
		},
		{
			name: "candidate -> init",
			makeArgs: func(t *testing.T) args {
				prev := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)
				for range 2 {
					prev[utils.RandomBytes32()] = mocktypes.NewCCIPOracle(t)
				}
				curr := make(map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle)

				return args{prevDeployment: prev, currDeployment: curr}
			},
			expect: func(t *testing.T, args args) {
				// When we are creating candidates, they should be started
				for _, plugin := range args.prevDeployment {
					plugin.On("Close").Return(nil).Once()
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.makeArgs(t)
			prev := make(ccipPlugins)
			for digest, oracle := range args.prevDeployment {
				prev[digest] = oracle
			}
			curr := make(ccipPlugins)
			for digest, oracle := range args.currDeployment {
				curr[digest] = oracle
			}
			tt.expect(t, args)
			defer assertions(t, args)
			err := curr.Transition(prev)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

//
//func Test_ccipDeployment_Close(t *testing.T) {
//	type args struct {
//		commitBlue  *mocktypes.CCIPOracle
//		commitGreen *mocktypes.CCIPOracle
//		execBlue    *mocktypes.CCIPOracle
//		execGreen   *mocktypes.CCIPOracle
//	}
//	tests := []struct {
//		name    string
//		args    args
//		expect  func(t *testing.T, args args)
//		asserts func(t *testing.T, args args)
//		wantErr bool
//	}{
//		{
//			name: "no errors, active only",
//			args: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: nil,
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   nil,
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Close").Return(nil).Once()
//				args.execBlue.On("Close").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//			},
//			wantErr: false,
//		},
//		{
//			name: "no errors, active and candidate",
//			args: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: mocktypes.NewCCIPOracle(t),
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Close").Return(nil).Once()
//				args.commitGreen.On("Close").Return(nil).Once()
//				args.execBlue.On("Close").Return(nil).Once()
//				args.execGreen.On("Close").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.commitGreen.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//				args.execGreen.AssertExpectations(t)
//			},
//			wantErr: false,
//		},
//		{
//			name: "error on commit active",
//			args: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: nil,
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   nil,
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Close").Return(errors.New("failed")).Once()
//				args.execBlue.On("Close").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &ccipDeployment{
//				commit: activeCandidateDeployment{
//					active: tt.args.commitBlue,
//				},
//				exec: activeCandidateDeployment{
//					active: tt.args.execBlue,
//				},
//			}
//			if tt.args.commitGreen != nil {
//				c.commit.candidate = tt.args.commitGreen
//			}
//
//			if tt.args.execGreen != nil {
//				c.exec.candidate = tt.args.execGreen
//			}
//
//			tt.expect(t, tt.args)
//			defer tt.asserts(t, tt.args)
//			err := c.Close()
//			if tt.wantErr {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}
//
//func Test_ccipDeployment_StartBlue(t *testing.T) {
//	type args struct {
//		commitBlue *mocktypes.CCIPOracle
//		execBlue   *mocktypes.CCIPOracle
//	}
//	tests := []struct {
//		name    string
//		args    args
//		expect  func(t *testing.T, args args)
//		asserts func(t *testing.T, args args)
//		wantErr bool
//	}{
//		{
//			name: "no errors",
//			args: args{
//				commitBlue: mocktypes.NewCCIPOracle(t),
//				execBlue:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Start").Return(nil).Once()
//				args.execBlue.On("Start").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//			},
//			wantErr: false,
//		},
//		{
//			name: "error on commit active",
//			args: args{
//				commitBlue: mocktypes.NewCCIPOracle(t),
//				execBlue:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Start").Return(errors.New("failed")).Once()
//				args.execBlue.On("Start").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//			},
//			wantErr: true,
//		},
//		{
//			name: "error on exec active",
//			args: args{
//				commitBlue: mocktypes.NewCCIPOracle(t),
//				execBlue:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Start").Return(nil).Once()
//				args.execBlue.On("Start").Return(errors.New("failed")).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &ccipDeployment{
//				commit: activeCandidateDeployment{
//					active: tt.args.commitBlue,
//				},
//				exec: activeCandidateDeployment{
//					active: tt.args.execBlue,
//				},
//			}
//
//			tt.expect(t, tt.args)
//			defer tt.asserts(t, tt.args)
//			err := c.StartActive()
//			if tt.wantErr {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}
//
//func Test_ccipDeployment_CloseBlue(t *testing.T) {
//	type args struct {
//		commitBlue *mocktypes.CCIPOracle
//		execBlue   *mocktypes.CCIPOracle
//	}
//	tests := []struct {
//		name    string
//		args    args
//		expect  func(t *testing.T, args args)
//		asserts func(t *testing.T, args args)
//		wantErr bool
//	}{
//		{
//			name: "no errors",
//			args: args{
//				commitBlue: mocktypes.NewCCIPOracle(t),
//				execBlue:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Close").Return(nil).Once()
//				args.execBlue.On("Close").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//			},
//			wantErr: false,
//		},
//		{
//			name: "error on commit active",
//			args: args{
//				commitBlue: mocktypes.NewCCIPOracle(t),
//				execBlue:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Close").Return(errors.New("failed")).Once()
//				args.execBlue.On("Close").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//			},
//			wantErr: true,
//		},
//		{
//			name: "error on exec active",
//			args: args{
//				commitBlue: mocktypes.NewCCIPOracle(t),
//				execBlue:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args) {
//				args.commitBlue.On("Close").Return(nil).Once()
//				args.execBlue.On("Close").Return(errors.New("failed")).Once()
//			},
//			asserts: func(t *testing.T, args args) {
//				args.commitBlue.AssertExpectations(t)
//				args.execBlue.AssertExpectations(t)
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &ccipDeployment{
//				commit: activeCandidateDeployment{
//					active: tt.args.commitBlue,
//				},
//				exec: activeCandidateDeployment{
//					active: tt.args.execBlue,
//				},
//			}
//
//			tt.expect(t, tt.args)
//			defer tt.asserts(t, tt.args)
//			err := c.CloseActive()
//			if tt.wantErr {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}
//
//func Test_ccipDeployment_HandleBlueGreen_PrevDeploymentNil(t *testing.T) {
//	require.Error(t, (&ccipDeployment{}).TransitionDeployment(nil))
//}
//
//func Test_ccipDeployment_HandleBlueGreen(t *testing.T) {
//	type args struct {
//		commitBlue  *mocktypes.CCIPOracle
//		commitGreen *mocktypes.CCIPOracle
//		execBlue    *mocktypes.CCIPOracle
//		execGreen   *mocktypes.CCIPOracle
//	}
//	tests := []struct {
//		name                 string
//		argsPrevDeployment   args
//		argsFutureDeployment args
//		expect               func(t *testing.T, args args, argsPrevDeployment args)
//		asserts              func(t *testing.T, args args, argsPrevDeployment args)
//		wantErr              bool
//	}{
//		{
//			name: "promotion active to candidate",
//			argsPrevDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: mocktypes.NewCCIPOracle(t),
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			argsFutureDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: nil,
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   nil,
//			},
//			expect: func(t *testing.T, args args, argsPrevDeployment args) {
//				argsPrevDeployment.commitBlue.On("Close").Return(nil).Once()
//				argsPrevDeployment.execBlue.On("Close").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
//				argsPrevDeployment.commitBlue.AssertExpectations(t)
//				argsPrevDeployment.execBlue.AssertExpectations(t)
//			},
//			wantErr: false,
//		},
//		{
//			name: "new candidate deployment",
//			argsPrevDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: nil,
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   nil,
//			},
//			argsFutureDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: mocktypes.NewCCIPOracle(t),
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args, argsPrevDeployment args) {
//				args.commitGreen.On("Start").Return(nil).Once()
//				args.execGreen.On("Start").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
//				args.commitGreen.AssertExpectations(t)
//				args.execGreen.AssertExpectations(t)
//			},
//			wantErr: false,
//		},
//		{
//			name: "error on commit candidate start",
//			argsPrevDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: nil,
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   nil,
//			},
//			argsFutureDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: mocktypes.NewCCIPOracle(t),
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args, argsPrevDeployment args) {
//				args.commitGreen.On("Start").Return(errors.New("failed")).Once()
//				args.execGreen.On("Start").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
//				args.commitGreen.AssertExpectations(t)
//				args.execGreen.AssertExpectations(t)
//			},
//			wantErr: true,
//		},
//		{
//			name: "error on exec candidate start",
//			argsPrevDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: nil,
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   nil,
//			},
//			argsFutureDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: mocktypes.NewCCIPOracle(t),
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args, argsPrevDeployment args) {
//				args.commitGreen.On("Start").Return(nil).Once()
//				args.execGreen.On("Start").Return(errors.New("failed")).Once()
//			},
//			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
//				args.commitGreen.AssertExpectations(t)
//				args.execGreen.AssertExpectations(t)
//			},
//			wantErr: true,
//		},
//		{
//			name: "invalid active-candidate deployment transition commit: both prev and future deployment have candidate",
//			argsPrevDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: mocktypes.NewCCIPOracle(t),
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			argsFutureDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: mocktypes.NewCCIPOracle(t),
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			expect:  func(t *testing.T, args args, argsPrevDeployment args) {},
//			asserts: func(t *testing.T, args args, argsPrevDeployment args) {},
//			wantErr: true,
//		},
//		{
//			name: "invalid active-candidate deployment transition exec: both prev and future exec deployment have candidate",
//			argsPrevDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: nil,
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			argsFutureDeployment: args{
//				commitBlue:  mocktypes.NewCCIPOracle(t),
//				commitGreen: mocktypes.NewCCIPOracle(t),
//				execBlue:    mocktypes.NewCCIPOracle(t),
//				execGreen:   mocktypes.NewCCIPOracle(t),
//			},
//			expect: func(t *testing.T, args args, argsPrevDeployment args) {
//				args.commitGreen.On("Start").Return(nil).Once()
//			},
//			asserts: func(t *testing.T, args args, argsPrevDeployment args) {
//				args.commitGreen.AssertExpectations(t)
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			futDeployment := &ccipDeployment{
//				commit: activeCandidateDeployment{
//					active: tt.argsFutureDeployment.commitBlue,
//				},
//				exec: activeCandidateDeployment{
//					active: tt.argsFutureDeployment.execBlue,
//				},
//			}
//			if tt.argsFutureDeployment.commitGreen != nil {
//				futDeployment.commit.candidate = tt.argsFutureDeployment.commitGreen
//			}
//			if tt.argsFutureDeployment.execGreen != nil {
//				futDeployment.exec.candidate = tt.argsFutureDeployment.execGreen
//			}
//
//			prevDeployment := &ccipDeployment{
//				commit: activeCandidateDeployment{
//					active: tt.argsPrevDeployment.commitBlue,
//				},
//				exec: activeCandidateDeployment{
//					active: tt.argsPrevDeployment.execBlue,
//				},
//			}
//			if tt.argsPrevDeployment.commitGreen != nil {
//				prevDeployment.commit.candidate = tt.argsPrevDeployment.commitGreen
//			}
//			if tt.argsPrevDeployment.execGreen != nil {
//				prevDeployment.exec.candidate = tt.argsPrevDeployment.execGreen
//			}
//
//			tt.expect(t, tt.argsFutureDeployment, tt.argsPrevDeployment)
//			defer tt.asserts(t, tt.argsFutureDeployment, tt.argsPrevDeployment)
//			err := futDeployment.TransitionDeployment(prevDeployment)
//			if tt.wantErr {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}
//
//func Test_isNewGreenInstance(t *testing.T) {
//	type args struct {
//		pluginType     cctypes.PluginType
//		ocrConfigs     []ccipreaderpkg.OCR3ConfigWithMeta
//		prevDeployment ccipDeployment
//	}
//	tests := []struct {
//		name string
//		args args
//		want bool
//	}{
//		{
//			"prev deployment only active",
//			args{
//				pluginType: cctypes.PluginTypeCCIPCommit,
//				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
//					{}, {},
//				},
//				prevDeployment: ccipDeployment{
//					commit: activeCandidateDeployment{
//						active: mocktypes.NewCCIPOracle(t),
//					},
//				},
//			},
//			true,
//		},
//		{
//			"candidate -> active promotion",
//			args{
//				pluginType: cctypes.PluginTypeCCIPCommit,
//				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
//					{},
//				},
//				prevDeployment: ccipDeployment{
//					commit: activeCandidateDeployment{
//						active:    mocktypes.NewCCIPOracle(t),
//						candidate: mocktypes.NewCCIPOracle(t),
//					},
//				},
//			},
//			false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := isNewCandidateInstance(tt.args.pluginType, tt.args.ocrConfigs, tt.args.prevDeployment)
//			require.Equal(t, tt.want, got)
//		})
//	}
//}
//
//func Test_isPromotion(t *testing.T) {
//	type args struct {
//		pluginType     cctypes.PluginType
//		ocrConfigs     []ccipreaderpkg.OCR3ConfigWithMeta
//		prevDeployment ccipDeployment
//	}
//	tests := []struct {
//		name string
//		args args
//		want bool
//	}{
//		{
//			"prev deployment only active",
//			args{
//				pluginType: cctypes.PluginTypeCCIPCommit,
//				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
//					{}, {},
//				},
//				prevDeployment: ccipDeployment{
//					commit: activeCandidateDeployment{
//						active: mocktypes.NewCCIPOracle(t),
//					},
//				},
//			},
//			false,
//		},
//		{
//			"candidate -> active promotion",
//			args{
//				pluginType: cctypes.PluginTypeCCIPCommit,
//				ocrConfigs: []ccipreaderpkg.OCR3ConfigWithMeta{
//					{},
//				},
//				prevDeployment: ccipDeployment{
//					commit: activeCandidateDeployment{
//						active:    mocktypes.NewCCIPOracle(t),
//						candidate: mocktypes.NewCCIPOracle(t),
//					},
//				},
//			},
//			true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := isPromotion(tt.args.pluginType, tt.args.ocrConfigs, tt.args.prevDeployment); got != tt.want {
//				t.Errorf("isPromotion() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_ccipDeployment_HasGreenInstance(t *testing.T) {
//	type fields struct {
//		commit activeCandidateDeployment
//		exec   activeCandidateDeployment
//	}
//	type args struct {
//		pluginType cctypes.PluginType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		{
//			"commit candidate present",
//			fields{
//				commit: activeCandidateDeployment{
//					active:    mocktypes.NewCCIPOracle(t),
//					candidate: mocktypes.NewCCIPOracle(t),
//				},
//			},
//			args{
//				pluginType: cctypes.PluginTypeCCIPCommit,
//			},
//			true,
//		},
//		{
//			"commit candidate not present",
//			fields{
//				commit: activeCandidateDeployment{
//					active: mocktypes.NewCCIPOracle(t),
//				},
//			},
//			args{
//				pluginType: cctypes.PluginTypeCCIPCommit,
//			},
//			false,
//		},
//		{
//			"exec candidate present",
//			fields{
//				exec: activeCandidateDeployment{
//					active:    mocktypes.NewCCIPOracle(t),
//					candidate: mocktypes.NewCCIPOracle(t),
//				},
//			},
//			args{
//				pluginType: cctypes.PluginTypeCCIPExec,
//			},
//			true,
//		},
//		{
//			"exec candidate not present",
//			fields{
//				exec: activeCandidateDeployment{
//					active: mocktypes.NewCCIPOracle(t),
//				},
//			},
//			args{
//				pluginType: cctypes.PluginTypeCCIPExec,
//			},
//			false,
//		},
//		{
//			"invalid plugin type",
//			fields{},
//			args{
//				pluginType: cctypes.PluginType(100),
//			},
//			false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &ccipDeployment{}
//			if tt.fields.commit.active != nil {
//				c.commit.active = tt.fields.commit.active
//			}
//			if tt.fields.commit.candidate != nil {
//				c.commit.candidate = tt.fields.commit.candidate
//			}
//			if tt.fields.exec.active != nil {
//				c.exec.active = tt.fields.exec.active
//			}
//			if tt.fields.exec.candidate != nil {
//				c.exec.candidate = tt.fields.exec.candidate
//			}
//			got := c.HasCandidateInstance(tt.args.pluginType)
//			require.Equal(t, tt.want, got)
//		})
//	}
//}
