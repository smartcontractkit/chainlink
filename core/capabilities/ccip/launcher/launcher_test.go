package launcher

import (
	"math/big"
	"reflect"
	"testing"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types/mocks"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func Test_createDON(t *testing.T) {
	type args struct {
		lggr            logger.Logger
		p2pID           ragep2ptypes.PeerID
		homeChainReader *mocks.HomeChainReader
		oracleCreator   *mocks.OracleCreator
		don             registrysyncer.DON
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args, oracleCreator *mocks.OracleCreator, homeChainReader *mocks.HomeChainReader)
		wantErr bool
	}{
		{
			"not a member of the DON and not a bootstrap node",
			args{
				logger.TestLogger(t),
				p2pID1,
				mocks.NewHomeChainReader(t),
				mocks.NewOracleCreator(t),
				registrysyncer.DON{
					DON:                      getDON(2, []ragep2ptypes.PeerID{p2pID3, p2pID4}, 0),
					CapabilityConfigurations: defaultCapCfgs,
				},
			},
			func(t *testing.T, args args, oracleCreator *mocks.OracleCreator, homeChainReader *mocks.HomeChainReader) {
				homeChainReader.
					On("GetOCRConfigs", mock.Anything, uint32(2), uint8(cctypes.PluginTypeCCIPCommit)).
					Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
						Config: ccipreaderpkg.OCR3Config{
							PluginType: uint8(cctypes.PluginTypeCCIPCommit),
							Nodes:      getOCR3Nodes(3, 4),
						},
					}}, nil)
				homeChainReader.
					On("GetOCRConfigs", mock.Anything, uint32(2), uint8(cctypes.PluginTypeCCIPExec)).
					Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
						Config: ccipreaderpkg.OCR3Config{
							PluginType: uint8(cctypes.PluginTypeCCIPExec),
							Nodes:      getOCR3Nodes(3, 4),
						},
					}}, nil)
				oracleCreator.EXPECT().Type().Return(cctypes.OracleTypePlugin).Once()
			},
			false,
		},
		{
			"not a member of the DON but a running a bootstrap oracle creator",
			args{
				logger.TestLogger(t),
				ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID()),
				mocks.NewHomeChainReader(t),
				mocks.NewOracleCreator(t),
				registrysyncer.DON{
					DON:                      getDON(2, []ragep2ptypes.PeerID{p2pID3, p2pID4}, 0),
					CapabilityConfigurations: defaultCapCfgs,
				},
			},
			func(t *testing.T, args args, oracleCreator *mocks.OracleCreator, homeChainReader *mocks.HomeChainReader) {
				homeChainReader.
					On("GetOCRConfigs", mock.Anything, uint32(2), uint8(cctypes.PluginTypeCCIPCommit)).
					Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
						Config: ccipreaderpkg.OCR3Config{
							PluginType: uint8(cctypes.PluginTypeCCIPCommit),
							Nodes:      getOCR3Nodes(3, 4),
						},
					}}, nil)
				homeChainReader.
					On("GetOCRConfigs", mock.Anything, uint32(2), uint8(cctypes.PluginTypeCCIPExec)).
					Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
						Config: ccipreaderpkg.OCR3Config{
							PluginType: uint8(cctypes.PluginTypeCCIPExec),
							Nodes:      getOCR3Nodes(3, 4),
						},
					}}, nil)
				oracleCreator.EXPECT().Type().Return(cctypes.OracleTypeBootstrap).Once()
				oracleCreator.EXPECT().Create(mock.Anything, mock.Anything).Return(mocks.NewCCIPOracle(t), nil).Twice()
			},
			false,
		},
		{
			"success",
			args{
				logger.TestLogger(t),
				p2pID1,
				mocks.NewHomeChainReader(t),
				mocks.NewOracleCreator(t),
				defaultRegistryDon,
			},
			func(t *testing.T, args args, oracleCreator *mocks.OracleCreator, homeChainReader *mocks.HomeChainReader) {
				homeChainReader.
					On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPCommit)).
					Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
						Config: ccipreaderpkg.OCR3Config{
							PluginType: uint8(cctypes.PluginTypeCCIPCommit),
							Nodes:      getOCR3Nodes(3, 4),
						},
					}}, nil)
				homeChainReader.
					On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPExec)).
					Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
						Config: ccipreaderpkg.OCR3Config{
							PluginType: uint8(cctypes.PluginTypeCCIPExec),
							Nodes:      getOCR3Nodes(3, 4),
						},
					}}, nil)

				oracleCreator.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cfg cctypes.OCR3ConfigWithMeta) bool {
					return cfg.Config.PluginType == uint8(cctypes.PluginTypeCCIPCommit)
				})).
					Return(mocks.NewCCIPOracle(t), nil)
				oracleCreator.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cfg cctypes.OCR3ConfigWithMeta) bool {
					return cfg.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec)
				})).
					Return(mocks.NewCCIPOracle(t), nil)
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expect != nil {
				tt.expect(t, tt.args, tt.args.oracleCreator, tt.args.homeChainReader)
			}

			_, err := createDON(tt.args.lggr, tt.args.p2pID, tt.args.homeChainReader, tt.args.don, tt.args.oracleCreator)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_createFutureBlueGreenDeployment(t *testing.T) {
	type args struct {
		donID          uint32
		prevDeployment ccipDeployment
		ocrConfigs     []ccipreaderpkg.OCR3ConfigWithMeta
		oracleCreator  *mocks.OracleCreator
		pluginType     cctypes.PluginType
	}
	tests := []struct {
		name    string
		args    args
		want    activeCandidateDeployment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createFutureActiveCandidateDeployment(tt.args.donID, tt.args.prevDeployment, tt.args.ocrConfigs, tt.args.oracleCreator, tt.args.pluginType)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFutureActiveCandidateDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFutureActiveCandidateDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateDON(t *testing.T) {
	type args struct {
		lggr            logger.Logger
		p2pID           ragep2ptypes.PeerID
		homeChainReader *mocks.HomeChainReader
		oracleCreator   *mocks.OracleCreator
		prevDeployment  ccipDeployment
		don             registrysyncer.DON
	}
	tests := []struct {
		name              string
		args              args
		wantFutDeployment *ccipDeployment
		wantErr           bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFutDeployment, err := updateDON(tt.args.lggr, tt.args.p2pID, tt.args.homeChainReader, tt.args.prevDeployment, tt.args.don, tt.args.oracleCreator)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateDON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFutDeployment, tt.wantFutDeployment) {
				t.Errorf("updateDON() = %v, want %v", gotFutDeployment, tt.wantFutDeployment)
			}
		})
	}
}

func Test_launcher_processDiff(t *testing.T) {
	type fields struct {
		lggr            logger.Logger
		p2pID           ragep2ptypes.PeerID
		homeChainReader *mocks.HomeChainReader
		oracleCreator   *mocks.OracleCreator
		dons            map[registrysyncer.DonID]*ccipDeployment
		regState        registrysyncer.LocalRegistry
	}
	type args struct {
		diff diffResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		assert  func(t *testing.T, l *launcher)
		wantErr bool
	}{
		{
			"don removed success",
			fields{
				dons: map[registrysyncer.DonID]*ccipDeployment{
					1: {
						commit: activeCandidateDeployment{
							active: newMock(t,
								func(t *testing.T) *mocks.CCIPOracle { return mocks.NewCCIPOracle(t) },
								func(m *mocks.CCIPOracle) {
									m.On("Close").Return(nil)
								}),
						},
						exec: activeCandidateDeployment{
							active: newMock(t,
								func(t *testing.T) *mocks.CCIPOracle { return mocks.NewCCIPOracle(t) },
								func(m *mocks.CCIPOracle) {
									m.On("Close").Return(nil)
								}),
						},
					},
				},
				regState: registrysyncer.LocalRegistry{
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
				},
			},
			args{
				diff: diffResult{
					removed: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
				},
			},
			func(t *testing.T, l *launcher) {
				require.Len(t, l.dons, 0)
				require.Len(t, l.regState.IDsToDONs, 0)
			},
			false,
		},
		{
			"don added success",
			fields{
				lggr:  logger.TestLogger(t),
				p2pID: p2pID1,
				homeChainReader: newMock(t, func(t *testing.T) *mocks.HomeChainReader {
					return mocks.NewHomeChainReader(t)
				}, func(m *mocks.HomeChainReader) {
					m.On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPCommit)).
						Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
							Config: ccipreaderpkg.OCR3Config{
								PluginType: uint8(cctypes.PluginTypeCCIPCommit),
							},
						}}, nil)
					m.On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPExec)).
						Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
							Config: ccipreaderpkg.OCR3Config{
								PluginType: uint8(cctypes.PluginTypeCCIPExec),
							},
						}}, nil)
				}),
				oracleCreator: newMock(t, func(t *testing.T) *mocks.OracleCreator {
					return mocks.NewOracleCreator(t)
				}, func(m *mocks.OracleCreator) {
					commitOracle := mocks.NewCCIPOracle(t)
					commitOracle.On("Start").Return(nil)
					execOracle := mocks.NewCCIPOracle(t)
					execOracle.On("Start").Return(nil)
					m.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cfg cctypes.OCR3ConfigWithMeta) bool {
						return cfg.Config.PluginType == uint8(cctypes.PluginTypeCCIPCommit)
					})).
						Return(commitOracle, nil)
					m.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cfg cctypes.OCR3ConfigWithMeta) bool {
						return cfg.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec)
					})).
						Return(execOracle, nil)
				}),
				dons: map[registrysyncer.DonID]*ccipDeployment{},
				regState: registrysyncer.LocalRegistry{
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{},
				},
			},
			args{
				diff: diffResult{
					added: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
				},
			},
			func(t *testing.T, l *launcher) {
				require.Len(t, l.dons, 1)
				require.Len(t, l.regState.IDsToDONs, 1)
			},
			false,
		},
		{
			"don updated new candidate instance success",
			fields{
				lggr:  logger.TestLogger(t),
				p2pID: p2pID1,
				homeChainReader: newMock(t, func(t *testing.T) *mocks.HomeChainReader {
					return mocks.NewHomeChainReader(t)
				}, func(m *mocks.HomeChainReader) {
					m.On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPCommit)).
						Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
							Config: ccipreaderpkg.OCR3Config{
								PluginType: uint8(cctypes.PluginTypeCCIPCommit),
							},
						}, {
							Config: ccipreaderpkg.OCR3Config{
								PluginType: uint8(cctypes.PluginTypeCCIPCommit),
							},
						}}, nil)
					m.On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPExec)).
						Return([]ccipreaderpkg.OCR3ConfigWithMeta{{
							Config: ccipreaderpkg.OCR3Config{
								PluginType: uint8(cctypes.PluginTypeCCIPExec),
							},
						}, {
							Config: ccipreaderpkg.OCR3Config{
								PluginType: uint8(cctypes.PluginTypeCCIPExec),
							},
						}}, nil)
				}),
				oracleCreator: newMock(t, func(t *testing.T) *mocks.OracleCreator {
					return mocks.NewOracleCreator(t)
				}, func(m *mocks.OracleCreator) {
					commitOracle := mocks.NewCCIPOracle(t)
					commitOracle.On("Start").Return(nil)
					execOracle := mocks.NewCCIPOracle(t)
					execOracle.On("Start").Return(nil)
					m.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cfg cctypes.OCR3ConfigWithMeta) bool {
						return cfg.Config.PluginType == uint8(cctypes.PluginTypeCCIPCommit)
					})).
						Return(commitOracle, nil)
					m.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cfg cctypes.OCR3ConfigWithMeta) bool {
						return cfg.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec)
					})).
						Return(execOracle, nil)
				}),
				dons: map[registrysyncer.DonID]*ccipDeployment{
					1: {
						commit: activeCandidateDeployment{
							active: newMock(t, func(t *testing.T) *mocks.CCIPOracle {
								return mocks.NewCCIPOracle(t)
							}, func(m *mocks.CCIPOracle) {}),
						},
						exec: activeCandidateDeployment{
							active: newMock(t, func(t *testing.T) *mocks.CCIPOracle {
								return mocks.NewCCIPOracle(t)
							}, func(m *mocks.CCIPOracle) {}),
						},
					},
				},
				regState: registrysyncer.LocalRegistry{
					IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
						1: defaultRegistryDon,
					},
				},
			},
			args{
				diff: diffResult{
					updated: map[registrysyncer.DonID]registrysyncer.DON{
						1: {
							// new Node in Don: p2pID2
							DON:                      getDON(1, []ragep2ptypes.PeerID{p2pID1, p2pID2}, 0),
							CapabilityConfigurations: defaultCapCfgs,
						},
					},
				},
			},
			func(t *testing.T, l *launcher) {
				require.Len(t, l.dons, 1)
				require.Len(t, l.regState.IDsToDONs, 1)
				require.Len(t, l.regState.IDsToDONs[1].Members, 2)
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &launcher{
				dons:            tt.fields.dons,
				regState:        tt.fields.regState,
				p2pID:           tt.fields.p2pID,
				lggr:            tt.fields.lggr,
				homeChainReader: tt.fields.homeChainReader,
				oracleCreator:   tt.fields.oracleCreator,
			}
			err := l.processDiff(tt.args.diff)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			tt.assert(t, l)
		})
	}
}

func getOCR3Nodes(p2pIDs ...int64) []ccipreaderpkg.OCR3Node {
	nodes := make([]ccipreaderpkg.OCR3Node, len(p2pIDs))
	for i, p2pID := range p2pIDs {
		nodes[i] = ccipreaderpkg.OCR3Node{P2pID: p2pkey.MustNewV2XXXTestingOnly(big.NewInt(p2pID)).PeerID()}
	}
	return nodes
}
func newMock[T any](t *testing.T, newer func(t *testing.T) T, expect func(m T)) T {
	o := newer(t)
	expect(o)
	return o
}
