package deployment

import (
	"context"
	"encoding/binary"
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
}

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
		return nil, fmt.Errorf("data too short")
	}
	state := &UpgradeableLoaderState{}
	state.Type = binary.LittleEndian.Uint32(data[:4])

	switch state.Type {
	case 2: // Program
		if len(data) < 36 {
			return nil, fmt.Errorf("program data too short")
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
				return nil, fmt.Errorf("missing authority pubkey")
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

func GetProgramDataAddress(client *solRpc.Client, ctx context.Context, progPubkey solana.PublicKey) (solana.PublicKey, error) {
	resp, err := client.GetAccountInfo(ctx, progPubkey)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to fetch program account: %w", err)
	}
	if resp.Value == nil {
		return solana.PublicKey{}, fmt.Errorf("program account does not exist")
	}

	state, err := decodeUpgradeableLoaderState(resp.Value.Data.GetBinary())
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("decode error: %w", err)
	}
	if state.Program == nil {
		return solana.PublicKey{}, fmt.Errorf("account is not an upgradeable program")
	}
	return state.Program.ProgramData, nil
}

func GetUpgradeAuthority(client *solRpc.Client, ctx context.Context, progDataPubkey solana.PublicKey) (solana.PublicKey, bool, error) {
	resp, err := client.GetAccountInfo(ctx, progDataPubkey)
	if err != nil {
		return solana.PublicKey{}, false, fmt.Errorf("failed to fetch programdata account: %w", err)
	}
	if resp.Value == nil {
		return solana.PublicKey{}, false, fmt.Errorf("programdata account does not exist")
	}

	state, err := decodeUpgradeableLoaderState(resp.Value.Data.GetBinary())
	if err != nil {
		return solana.PublicKey{}, false, fmt.Errorf("decode error: %w", err)
	}
	if state.ProgramData == nil {
		return solana.PublicKey{}, false, fmt.Errorf("unexpected state: not programdata")
	}

	if state.ProgramData.AuthorityOption == 0 {
		// No authority – the program is immutable
		return solana.PublicKey{}, false, nil
	}
	return state.ProgramData.Authority, true, nil
}
