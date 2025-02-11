package launcher

import (
	"testing"

	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	mocktypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types/mocks"
)

func Test_ccipDeployment_Transitions(t *testing.T) {
	// we use a pointer to the oracle here for mock assertions
	type args struct {
		prevDeployment map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle
		currDeployment map[ocrtypes.ConfigDigest]*mocktypes.CCIPOracle
	}
	assertions := func(t *testing.T, args args) {
		for i := range args.prevDeployment {
			args.prevDeployment[i].AssertExpectations(t)
		}
		for i := range args.currDeployment {
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
			prev := make(pluginRegistry)
			for digest, oracle := range args.prevDeployment {
				prev[digest] = oracle
			}
			curr := make(pluginRegistry)
			for digest, oracle := range args.currDeployment {
				curr[digest] = oracle
			}
			tt.expect(t, args)
			defer assertions(t, args)
			err := curr.TransitionFrom(prev)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
