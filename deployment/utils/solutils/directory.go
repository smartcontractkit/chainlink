package solutils

// Program names
const (
	// CCIP Programs
	ProgCCIPRouter             = "ccip_router"
	ProgTestTokenPool          = "test_token_pool"
	ProgBurnMintTokenPool      = "burnmint_token_pool"
	ProgLockReleaseTokenPool   = "lockrelease_token_pool"
	ProgFeeQuoter              = "fee_quoter"
	ProgTestCCIPReceiver       = "test_ccip_receiver"
	ProgCCIPOfframp            = "ccip_offramp"
	ProgExternalProgramCPIStub = "external_program_cpi_stub"
	ProgRMNRemote              = "rmn_remote"
	ProgCCTPTokenPool          = "cctp_token_pool"
	ProgBaseSignerRegistry     = "ccip_signer_registry" // The base implementation of the signer registry program. This is not actually deployed as a program, so it is not included in the CCIP program names list.

	// Data Feeds Programs
	ProgDataFeedsCache = "data_feeds_cache"

	// Keystone Programs
	ProgKeystoneForwarder = "keystone_forwarder"

	// MCMS Programs
	ProgMCM              = "mcm"
	ProgTimelock         = "timelock"
	ProgAccessController = "access_controller"
)

// Program names grouped by their usage.
var (
	CCIPProgramNames = []string{
		ProgCCIPRouter,
		ProgTestTokenPool,
		ProgBurnMintTokenPool,
		ProgLockReleaseTokenPool,
		ProgFeeQuoter,
		ProgTestCCIPReceiver,
		ProgCCIPOfframp,
		ProgExternalProgramCPIStub,
		ProgRMNRemote,
		ProgCCTPTokenPool,
	}
	DataFeedsProgramNames = []string{ProgDataFeedsCache}
	KeystoneProgramNames  = []string{ProgKeystoneForwarder}
	MCMSProgramNames      = []string{ProgMCM, ProgTimelock, ProgAccessController}
)

// Repositories that contain the program artifacts.
const (
	repoCCIP   = "chainlink-ccip"
	repoSolana = "chainlink-solana"
)

// Programs maps program names to it's information.
//
// This is the source of truth for the program IDs and repositories.
var Directory = directory{
	// CCIP Programs
	ProgCCIPRouter:             {ID: "Ccip842gzYHhvdDkSyi2YVCoAWPbYJoApMFzSxQroE9C", Repo: repoCCIP, ProgramBufferBytes: 5 * 1024 * 1024},
	ProgTestTokenPool:          {ID: "JuCcZ4smxAYv9QHJ36jshA7pA3FuQ3vQeWLUeAtZduJ", Repo: repoCCIP},
	ProgBurnMintTokenPool:      {ID: "41FGToCmdaWa1dgZLKFAjvmx6e6AjVTX7SVRibvsMGVB", Repo: repoCCIP, ProgramBufferBytes: 3 * 1024 * 1024},
	ProgLockReleaseTokenPool:   {ID: "8eqh8wppT9c5rw4ERqNCffvU6cNFJWff9WmkcYtmGiqC", Repo: repoCCIP, ProgramBufferBytes: 3 * 1024 * 1024},
	ProgFeeQuoter:              {ID: "FeeQPGkKDeRV1MgoYfMH6L8o3KeuYjwUZrgn4LRKfjHi", Repo: repoCCIP, ProgramBufferBytes: 5 * 1024 * 1024},
	ProgTestCCIPReceiver:       {ID: "EvhgrPhTDt4LcSPS2kfJgH6T6XWZ6wT3X9ncDGLT1vui", Repo: repoCCIP},
	ProgCCIPOfframp:            {ID: "offqSMQWgQud6WJz694LRzkeN5kMYpCHTpXQr3Rkcjm", Repo: repoCCIP, ProgramBufferBytes: 1.5 * 1024 * 1024},
	ProgExternalProgramCPIStub: {ID: "2zZwzyptLqwFJFEFxjPvrdhiGpH9pJ3MfrrmZX6NTKxm", Repo: repoCCIP},
	ProgRMNRemote:              {ID: "RmnXLft1mSEwDgMKu2okYuHkiazxntFFcZFrrcXxYg7", Repo: repoCCIP, ProgramBufferBytes: 3 * 1024 * 1024},
	ProgCCTPTokenPool:          {ID: "CCiTPESGEevd7TBU8EGBKrcxuRq7jx3YtW6tPidnscaZ", Repo: repoCCIP, ProgramBufferBytes: 3 * 1024 * 1024},
	ProgBaseSignerRegistry:     {ID: "S1GN4jus9XzKVVnoHqfkjo1GN8bX46gjXZQwsdGBPHE", Repo: "", ProgramBufferBytes: 1 * 1024 * 1024},

	// MCMS Programs
	ProgMCM:              {ID: "5vNJx78mz7KVMjhuipyr9jKBKcMrKYGdjGkgE4LUmjKk", Repo: repoCCIP, ProgramBufferBytes: 1 * 1024 * 1024},
	ProgTimelock:         {ID: "DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA", Repo: repoCCIP, ProgramBufferBytes: 1 * 1024 * 1024},
	ProgAccessController: {ID: "6KsN58MTnRQ8FfPaXHiFPPFGDRioikj9CdPvPxZJdCjb", Repo: repoCCIP, ProgramBufferBytes: 1 * 1024 * 1024},

	// Keystone Programs
	ProgKeystoneForwarder: {ID: "whV7Q5pi17hPPyaPksToDw1nMx6Lh8qmNWKFaLRQ4wz", Repo: repoSolana},

	// Data Feeds Programs
	ProgDataFeedsCache: {ID: "3kX63udXtYcsdj2737Wi2KGd2PhqiKPgAFAxstrjtRUa", Repo: repoSolana},
}

// GetProgramID returns the program ID for a given program name.
//
// Returns the program ID for the given program name or an empty string if the program is not
// found.
func GetProgramID(name string) string {
	info, ok := Directory[name]
	if !ok {
		return ""
	}

	return info.ID
}

// GetProgramBufferBytes returns the size of the program buffer in bytes for the given program name.
//
// Returns 0 if the program is not found or if the program is not upgradable.
func GetProgramBufferBytes(name string) int {
	info, ok := Directory[name]
	if !ok {
		return 0
	}

	return info.ProgramBufferBytes
}

// programInfo contains the information about a program.
type programInfo struct {
	// ID is the program ID of the program.
	ID string

	// Repo is the repository name of where the program is located.
	Repo string

	// ProgramBufferBytes is the size of the program buffer in bytes. Used for upgrades.
	// Can be left blank if the program is not upgradable.
	//
	// https://docs.google.com/document/d/1Fk76lOeyS2z2X6MokaNX_QTMFAn5wvSZvNXJluuNV1E/edit?tab=t.0#heading=h.uij286zaarkz
	// https://docs.google.com/document/d/1nCNuam0ljOHiOW0DUeiZf4ntHf_1Bw94Zi7ThPGoKR4/edit?tab=t.0#heading=h.hju45z55bnqd
	ProgramBufferBytes int
}

// directory maps the program name to the program information.
type directory map[string]programInfo
