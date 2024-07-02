package launcher

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types"
	mockcctypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func Test_createOracle(t *testing.T) {
	var p2pKeys []ragep2ptypes.PeerID
	for i := 0; i < 3; i++ {
		p2pKeys = append(p2pKeys, ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i+1))).PeerID()))
	}
	myP2PKey := p2pKeys[0]
	type args struct {
		p2pID         ragep2ptypes.PeerID
		oracleCreator *mockcctypes.OracleCreator
		pluginType    cctypes.PluginType
		ocrConfigs    []ccipreaderpkg.OCR3ConfigWithMeta
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args, oracleCreator *mockcctypes.OracleCreator)
		wantErr bool
	}{
		{
			"success, no bootstrap",
			args{
				myP2PKey,
				mockcctypes.NewOracleCreator(t),
				cctypes.PluginTypeCCIPCommit,
				[]ccipreaderpkg.OCR3ConfigWithMeta{
					{
						Config:       ccipreaderpkg.OCR3Config{},
						ConfigCount:  1,
						ConfigDigest: testutils.Random32Byte(),
					},
				},
			},
			func(t *testing.T, args args, oracleCreator *mockcctypes.OracleCreator) {
				oracleCreator.
					On("CreatePluginOracle", cctypes.PluginTypeCCIPCommit, cctypes.OCR3ConfigWithMeta(args.ocrConfigs[0])).
					Return(mockcctypes.NewCCIPOracle(t), nil)
			},
			false,
		},
		{
			"success, with bootstrap",
			args{
				myP2PKey,
				mockcctypes.NewOracleCreator(t),
				cctypes.PluginTypeCCIPCommit,
				[]ccipreaderpkg.OCR3ConfigWithMeta{
					{
						Config: ccipreaderpkg.OCR3Config{
							BootstrapP2PIds: [][32]byte{myP2PKey},
						},
						ConfigCount:  1,
						ConfigDigest: testutils.Random32Byte(),
					},
				},
			},
			func(t *testing.T, args args, oracleCreator *mockcctypes.OracleCreator) {
				oracleCreator.
					On("CreatePluginOracle", cctypes.PluginTypeCCIPCommit, cctypes.OCR3ConfigWithMeta(args.ocrConfigs[0])).
					Return(mockcctypes.NewCCIPOracle(t), nil)
				oracleCreator.
					On("CreateBootstrapOracle", cctypes.OCR3ConfigWithMeta(args.ocrConfigs[0])).
					Return(mockcctypes.NewCCIPOracle(t), nil)
			},
			false,
		},
		{
			"error creating plugin oracle",
			args{
				myP2PKey,
				mockcctypes.NewOracleCreator(t),
				cctypes.PluginTypeCCIPCommit,
				[]ccipreaderpkg.OCR3ConfigWithMeta{
					{
						Config:       ccipreaderpkg.OCR3Config{},
						ConfigCount:  1,
						ConfigDigest: testutils.Random32Byte(),
					},
				},
			},
			func(t *testing.T, args args, oracleCreator *mockcctypes.OracleCreator) {
				oracleCreator.
					On("CreatePluginOracle", cctypes.PluginTypeCCIPCommit, cctypes.OCR3ConfigWithMeta(args.ocrConfigs[0])).
					Return(nil, errors.New("error creating oracle"))
			},
			true,
		},
		{
			"error creating bootstrap oracle",
			args{
				myP2PKey,
				mockcctypes.NewOracleCreator(t),
				cctypes.PluginTypeCCIPCommit,
				[]ccipreaderpkg.OCR3ConfigWithMeta{
					{
						Config: ccipreaderpkg.OCR3Config{
							BootstrapP2PIds: [][32]byte{myP2PKey},
						},
						ConfigCount:  1,
						ConfigDigest: testutils.Random32Byte(),
					},
				},
			},
			func(t *testing.T, args args, oracleCreator *mockcctypes.OracleCreator) {
				oracleCreator.
					On("CreatePluginOracle", cctypes.PluginTypeCCIPCommit, cctypes.OCR3ConfigWithMeta(args.ocrConfigs[0])).
					Return(mockcctypes.NewCCIPOracle(t), nil)
				oracleCreator.
					On("CreateBootstrapOracle", cctypes.OCR3ConfigWithMeta(args.ocrConfigs[0])).
					Return(nil, errors.New("error creating oracle"))
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect(t, tt.args, tt.args.oracleCreator)
			_, _, err := createOracle(tt.args.p2pID, tt.args.oracleCreator, tt.args.pluginType, tt.args.ocrConfigs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_createDON(t *testing.T) {
	type args struct {
		lggr            logger.Logger
		p2pID           ragep2ptypes.PeerID
		homeChainReader *mockcctypes.HomeChainReader
		oracleCreator   *mockcctypes.OracleCreator
		don             kcr.CapabilitiesRegistryDONInfo
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args, oracleCreator *mockcctypes.OracleCreator, homeChainReader *mockcctypes.HomeChainReader)
		wantErr bool
	}{
		{
			"not a member of the DON",
			args{
				logger.TestLogger(t),
				ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID()),
				mockcctypes.NewHomeChainReader(t),
				mockcctypes.NewOracleCreator(t),
				kcr.CapabilitiesRegistryDONInfo{
					NodeP2PIds: [][32]byte{
						p2pkey.MustNewV2XXXTestingOnly(big.NewInt(2)).PeerID(),
					},
					Id: 2,
				},
			},
			func(t *testing.T, args args, oracleCreator *mockcctypes.OracleCreator, homeChainReader *mockcctypes.HomeChainReader) {
			},
			false,
		},
		{
			"success, no bootstrap",
			args{
				logger.TestLogger(t),
				ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID()),
				mockcctypes.NewHomeChainReader(t),
				mockcctypes.NewOracleCreator(t),
				kcr.CapabilitiesRegistryDONInfo{
					NodeP2PIds: [][32]byte{
						p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID(),
					},
					Id: 1,
				},
			},
			func(t *testing.T, args args, oracleCreator *mockcctypes.OracleCreator, homeChainReader *mockcctypes.HomeChainReader) {
				homeChainReader.
					On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPCommit)).
					Return([]ccipreaderpkg.OCR3ConfigWithMeta{{}}, nil)
				homeChainReader.
					On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPExec)).
					Return([]ccipreaderpkg.OCR3ConfigWithMeta{{}}, nil)
				oracleCreator.
					On("CreatePluginOracle", cctypes.PluginTypeCCIPCommit, mock.Anything).
					Return(mockcctypes.NewCCIPOracle(t), nil)
				oracleCreator.
					On("CreatePluginOracle", cctypes.PluginTypeCCIPExec, mock.Anything).
					Return(mockcctypes.NewCCIPOracle(t), nil)
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expect != nil {
				tt.expect(t, tt.args, tt.args.oracleCreator, tt.args.homeChainReader)
			}
			_, err := createDON(tt.args.lggr, tt.args.p2pID, tt.args.homeChainReader, tt.args.oracleCreator, tt.args.don)
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
		prevDeployment ccipDeployment
		ocrConfigs     []ccipreaderpkg.OCR3ConfigWithMeta
		oracleCreator  *mockcctypes.OracleCreator
		pluginType     cctypes.PluginType
	}
	tests := []struct {
		name    string
		args    args
		want    blueGreenDeployment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createFutureBlueGreenDeployment(tt.args.prevDeployment, tt.args.ocrConfigs, tt.args.oracleCreator, tt.args.pluginType)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFutureBlueGreenDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFutureBlueGreenDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateDON(t *testing.T) {
	type args struct {
		lggr            logger.Logger
		p2pID           ragep2ptypes.PeerID
		homeChainReader *mockcctypes.HomeChainReader
		oracleCreator   *mockcctypes.OracleCreator
		prevDeployment  ccipDeployment
		don             kcr.CapabilitiesRegistryDONInfo
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
			gotFutDeployment, err := updateDON(tt.args.lggr, tt.args.p2pID, tt.args.homeChainReader, tt.args.oracleCreator, tt.args.prevDeployment, tt.args.don)
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
		homeChainReader *mockcctypes.HomeChainReader
		oracleCreator   *mockcctypes.OracleCreator
		dons            map[registrysyncer.DonID]*ccipDeployment
		regState        registrysyncer.State
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
						commit: blueGreenDeployment{
							blue: newMock(t,
								func(t *testing.T) *mockcctypes.CCIPOracle { return mockcctypes.NewCCIPOracle(t) },
								func(m *mockcctypes.CCIPOracle) {
									m.On("Close").Return(nil)
								}),
						},
						exec: blueGreenDeployment{
							blue: newMock(t,
								func(t *testing.T) *mockcctypes.CCIPOracle { return mockcctypes.NewCCIPOracle(t) },
								func(m *mockcctypes.CCIPOracle) {
									m.On("Close").Return(nil)
								}),
						},
					},
				},
				regState: registrysyncer.State{
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
						},
					},
				},
			},
			args{
				diff: diffResult{
					removed: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
						},
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
				p2pID: ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID()),
				homeChainReader: newMock(t, func(t *testing.T) *mockcctypes.HomeChainReader {
					return mockcctypes.NewHomeChainReader(t)
				}, func(m *mockcctypes.HomeChainReader) {
					m.On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPCommit)).
						Return([]ccipreaderpkg.OCR3ConfigWithMeta{{}}, nil)
					m.On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPExec)).
						Return([]ccipreaderpkg.OCR3ConfigWithMeta{{}}, nil)
				}),
				oracleCreator: newMock(t, func(t *testing.T) *mockcctypes.OracleCreator {
					return mockcctypes.NewOracleCreator(t)
				}, func(m *mockcctypes.OracleCreator) {
					commitOracle := mockcctypes.NewCCIPOracle(t)
					commitOracle.On("Start").Return(nil)
					execOracle := mockcctypes.NewCCIPOracle(t)
					execOracle.On("Start").Return(nil)
					m.On("CreatePluginOracle", cctypes.PluginTypeCCIPCommit, mock.Anything).
						Return(commitOracle, nil)
					m.On("CreatePluginOracle", cctypes.PluginTypeCCIPExec, mock.Anything).
						Return(execOracle, nil)
				}),
				dons: map[registrysyncer.DonID]*ccipDeployment{},
				regState: registrysyncer.State{
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{},
				},
			},
			args{
				diff: diffResult{
					added: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							NodeP2PIds: [][32]byte{
								p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID(),
							},
						},
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
			"don updated new green instance success",
			fields{
				lggr:  logger.TestLogger(t),
				p2pID: ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID()),
				homeChainReader: newMock(t, func(t *testing.T) *mockcctypes.HomeChainReader {
					return mockcctypes.NewHomeChainReader(t)
				}, func(m *mockcctypes.HomeChainReader) {
					m.On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPCommit)).
						Return([]ccipreaderpkg.OCR3ConfigWithMeta{{}, {}}, nil)
					m.On("GetOCRConfigs", mock.Anything, uint32(1), uint8(cctypes.PluginTypeCCIPExec)).
						Return([]ccipreaderpkg.OCR3ConfigWithMeta{{}, {}}, nil)
				}),
				oracleCreator: newMock(t, func(t *testing.T) *mockcctypes.OracleCreator {
					return mockcctypes.NewOracleCreator(t)
				}, func(m *mockcctypes.OracleCreator) {
					commitOracle := mockcctypes.NewCCIPOracle(t)
					commitOracle.On("Start").Return(nil)
					execOracle := mockcctypes.NewCCIPOracle(t)
					execOracle.On("Start").Return(nil)
					m.On("CreatePluginOracle", cctypes.PluginTypeCCIPCommit, mock.Anything).
						Return(commitOracle, nil)
					m.On("CreatePluginOracle", cctypes.PluginTypeCCIPExec, mock.Anything).
						Return(execOracle, nil)
				}),
				dons: map[registrysyncer.DonID]*ccipDeployment{
					1: {
						commit: blueGreenDeployment{
							blue: newMock(t, func(t *testing.T) *mockcctypes.CCIPOracle {
								return mockcctypes.NewCCIPOracle(t)
							}, func(m *mockcctypes.CCIPOracle) {}),
						},
						exec: blueGreenDeployment{
							blue: newMock(t, func(t *testing.T) *mockcctypes.CCIPOracle {
								return mockcctypes.NewCCIPOracle(t)
							}, func(m *mockcctypes.CCIPOracle) {}),
						},
					},
				},
				regState: registrysyncer.State{
					IDsToDONs: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							NodeP2PIds: [][32]byte{
								p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID(),
							},
						},
					},
				},
			},
			args{
				diff: diffResult{
					updated: map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo{
						1: {
							Id: 1,
							NodeP2PIds: [][32]byte{
								p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID(),
								p2pkey.MustNewV2XXXTestingOnly(big.NewInt(2)).PeerID(), // new node in don
							},
						},
					},
				},
			},
			func(t *testing.T, l *launcher) {
				require.Len(t, l.dons, 1)
				require.Len(t, l.regState.IDsToDONs, 1)
				require.Len(t, l.regState.IDsToDONs[1].NodeP2PIds, 2)
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

func newMock[T any](t *testing.T, newer func(t *testing.T) T, expect func(m T)) T {
	o := newer(t)
	expect(o)
	return o
}
