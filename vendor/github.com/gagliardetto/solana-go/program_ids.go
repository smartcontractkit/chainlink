// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solana

var (
	// Create new accounts, allocate account data, assign accounts to owning programs,
	// transfer lamports from System Program owned accounts and pay transacation fees.
	SystemProgramID = MustPublicKeyFromBase58("11111111111111111111111111111111")

	// Add configuration data to the chain and the list of public keys that are permitted to modify it.
	ConfigProgramID = MustPublicKeyFromBase58("Config1111111111111111111111111111111111111")

	// Create and manage accounts representing stake and rewards for delegations to validators.
	StakeProgramID = MustPublicKeyFromBase58("Stake11111111111111111111111111111111111111")

	// Create and manage accounts that track validator voting state and rewards.
	VoteProgramID = MustPublicKeyFromBase58("Vote111111111111111111111111111111111111111")

	BPFLoaderDeprecatedProgramID = MustPublicKeyFromBase58("BPFLoader1111111111111111111111111111111111")
	// Deploys, upgrades, and executes programs on the chain.
	BPFLoaderProgramID            = MustPublicKeyFromBase58("BPFLoader2111111111111111111111111111111111")
	BPFLoaderUpgradeableProgramID = MustPublicKeyFromBase58("BPFLoaderUpgradeab1e11111111111111111111111")

	// Verify secp256k1 public key recovery operations (ecrecover).
	Secp256k1ProgramID = MustPublicKeyFromBase58("KeccakSecp256k11111111111111111111111111111")

	FeatureProgramID = MustPublicKeyFromBase58("Feature111111111111111111111111111111111111")
)

// SPL:
var (
	// A Token program on the Solana blockchain.
	// This program defines a common implementation for Fungible and Non Fungible tokens.
	TokenProgramID = MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

	// A Uniswap-like exchange for the Token program on the Solana blockchain,
	// implementing multiple automated market maker (AMM) curves.
	TokenSwapProgramID = MustPublicKeyFromBase58("SwaPpA9LAaLfeLi3a68M4DjnLqgtticKg6CnyNwgAC8")
	TokenSwapFeeOwner  = MustPublicKeyFromBase58("HfoTxFR1Tm6kGmWgYWD6J7YHVy1UwqSULUGVLXkJqaKN")

	// A lending protocol for the Token program on the Solana blockchain inspired by Aave and Compound.
	TokenLendingProgramID = MustPublicKeyFromBase58("LendZqTs8gn5CTSJU1jWKhKuVpjJGom45nnwPb2AMTi")

	// This program defines the convention and provides the mechanism for mapping
	// the user's wallet address to the associated token accounts they hold.
	SPLAssociatedTokenAccountProgramID = MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")

	// The Memo program is a simple program that validates a string of UTF-8 encoded characters
	// and verifies that any accounts provided are signers of the transaction.
	// The program also logs the memo, as well as any verified signer addresses,
	// to the transaction log, so that anyone can easily observe memos
	// and know they were approved by zero or more addresses
	// by inspecting the transaction log from a trusted provider.
	MemoProgramID = MustPublicKeyFromBase58("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")
)

var (
	// The Mint for native SOL Token accounts
	SolMint    = MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	WrappedSol = SolMint
)

var (
	TokenMetadataProgramID = MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
)
