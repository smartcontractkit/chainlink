package deployment

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
)

const (
	RouterProgramName               = "ccip_router"
	OffRampProgramName              = "ccip_offramp"
	FeeQuoterProgramName            = "fee_quoter"
	BurnMintTokenPoolProgramName    = "burnmint_token_pool"
	LockReleaseTokenPoolProgramName = "lockrelease_token_pool"
	AccessControllerProgramName     = "access_controller"
	TimelockProgramName             = "timelock"
	McmProgramName                  = "mcm"
	RMNRemoteProgramName            = "rmn_remote"
	ReceiverProgramName             = "test_ccip_receiver"
	KeystoneForwarderProgramName    = "keystone_forwarder"
	CCTPTokenPoolProgramName        = "cctp_token_pool"
	DataFeedsCacheProgramName       = "data_feeds_cache"
	BaseSignerRegistryProgramName   = "ccip_signer_registry"
)

// https://docs.google.com/document/d/1Fk76lOeyS2z2X6MokaNX_QTMFAn5wvSZvNXJluuNV1E/edit?tab=t.0#heading=h.uij286zaarkz
// https://docs.google.com/document/d/1nCNuam0ljOHiOW0DUeiZf4ntHf_1Bw94Zi7ThPGoKR4/edit?tab=t.0#heading=h.hju45z55bnqd
var SolanaProgramBytes = map[string]int{
	RouterProgramName:               5 * 1024 * 1024,
	OffRampProgramName:              1.5 * 1024 * 1024, // router should be redeployed but it does support upgrades if required (big fixes etc.)
	FeeQuoterProgramName:            5 * 1024 * 1024,
	BurnMintTokenPoolProgramName:    3 * 1024 * 1024,
	LockReleaseTokenPoolProgramName: 3 * 1024 * 1024,
	AccessControllerProgramName:     1 * 1024 * 1024,
	TimelockProgramName:             1 * 1024 * 1024,
	McmProgramName:                  1 * 1024 * 1024,
	RMNRemoteProgramName:            3 * 1024 * 1024,
	CCTPTokenPoolProgramName:        3 * 1024 * 1024,
	BaseSignerRegistryProgramName:   1 * 1024 * 1024,
}

// PROGRAM ID for Metaplex Metadata Program
var MplTokenMetadataID solana.PublicKey = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

// UpgradeableLoaderState mirrors the Rust enum in the Solana SDK.
type UpgradeableLoaderState struct {
	Type          uint32
	Program       *Program
	ProgramData   *ProgramData
	Uninitialized bool
}

// Program holds the address of the ProgramData account.
type Program struct {
	ProgramData solana.PublicKey
}

// ProgramData holds the optional UpgradeAuthority.
type ProgramData struct {
	Slot            uint64
	AuthorityOption uint32 // 0 = none, 1 = present
	Authority       solana.PublicKey
}

func decodeUpgradeableLoaderState(data []byte) (*UpgradeableLoaderState, error) {
	if len(data) < 4 {
		return nil, errors.New("data too short")
	}
	state := &UpgradeableLoaderState{}
	state.Type = binary.LittleEndian.Uint32(data[:4])

	switch state.Type {
	case 2: // Program
		if len(data) < 36 {
			return nil, errors.New("program data too short")
		}
		state.Program = &Program{
			ProgramData: solana.PublicKeyFromBytes(data[4:36]),
		}
	case 3: // ProgramData
		slot := binary.LittleEndian.Uint64(data[4:12])
		opt := data[12]
		var auth *solana.PublicKey
		if opt == 1 {
			if len(data) < 45 {
				return nil, errors.New("missing authority pubkey")
			}
			pk := solana.PublicKeyFromBytes(data[13:45])
			auth = &pk
		}
		state.ProgramData = &ProgramData{
			Slot:            slot,
			AuthorityOption: uint32(opt),
		}
		if state.ProgramData.AuthorityOption == 1 {
			state.ProgramData.Authority = *auth
		}
	default:
		// other variants (Uninitialized, Buffer) are not needed here
	}
	return state, nil
}

func getUpgradeableLoaderState(client *solRpc.Client, progPubkey solana.PublicKey) (*UpgradeableLoaderState, error) {
	resp, err := client.GetAccountInfo(context.Background(), progPubkey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch program account: %w", err)
	}
	if resp.Value == nil {
		return nil, errors.New("program account does not exist")
	}

	state, err := decodeUpgradeableLoaderState(resp.Value.Data.GetBinary())
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	return state, nil
}

func GetProgramDataAddress(client *solRpc.Client, progPubkey solana.PublicKey) (solana.PublicKey, error) {
	state, err := getUpgradeableLoaderState(client, progPubkey)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to get program data address for program %s: %w", progPubkey.String(), err)
	}

	if state.Program == nil {
		return solana.PublicKey{}, errors.New("account is not an upgradeable program")
	}
	return state.Program.ProgramData, nil
}

func GetUpgradeAuthority(client *solRpc.Client, progDataPubkey solana.PublicKey) (solana.PublicKey, bool, error) {
	state, err := getUpgradeableLoaderState(client, progDataPubkey)
	if err != nil {
		return solana.PublicKey{}, false, fmt.Errorf("failed to get upgrade authority for program data %s: %w", progDataPubkey.String(), err)
	}

	if state.ProgramData == nil {
		return solana.PublicKey{}, false, errors.New("unexpected state: not programdata")
	}

	if state.ProgramData.AuthorityOption == 0 {
		// No authority – the program is immutable
		return solana.PublicKey{}, false, nil
	}
	return state.ProgramData.Authority, true, nil
}

func FindMplTokenMetadataPDA(mint solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("metadata"),
		MplTokenMetadataID.Bytes(),
		mint.Bytes(),
	}
	pda, _, err := solana.FindProgramAddress(seeds, MplTokenMetadataID)
	return pda, err
}
