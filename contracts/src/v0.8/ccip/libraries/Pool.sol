// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice This library contains various token pool functions to aid constructing the return data.
library Pool {
  // The tag used to signal support for the pool v1 standard
  // bytes4(keccak256("CCIP_POOL_V1"))
  bytes4 public constant CCIP_POOL_V1 = 0xaff2afbf;

  // The number of bytes in the return data for a pool v1 releaseOrMint call.
  // This should match the size of the ReleaseOrMintOutV1 struct.
  uint16 public constant CCIP_POOL_V1_RET_BYTES = 32;

  // The default max number of bytes in the return data for a pool v1 lockOrBurn call.
  // This data can be used to send information to the destination chain token pool. Can be overwritten
  // in the TokenTransferFeeConfig.destBytesOverhead if more data is required.
  uint32 public constant CCIP_LOCK_OR_BURN_V1_RET_BYTES = 32;

  struct LockOrBurnInV1 {
    bytes receiver; //  The recipient of the tokens on the destination chain, abi encoded
    uint64 remoteChainSelector; // ─╮ The chain ID of the destination chain
    address originalSender; // ─────╯ The original sender of the tx on the source chain
    uint256 amount; //  The amount of tokens to lock or burn, denominated in the source token's decimals
    address localToken; //  The address on this chain of the token to lock or burn
  }

  struct LockOrBurnOutV1 {
    // The address of the destination token, abi encoded in the case of EVM chains
    // This value is UNTRUSTED as any pool owner can return whatever value they want.
    bytes destTokenAddress;
    // Optional pool data to be transferred to the destination chain. Be default this is capped at
    // CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
    // has to be set for the specific token.
    bytes destPoolData;
  }

  struct ReleaseOrMintInV1 {
    bytes originalSender; //          The original sender of the tx on the source chain
    uint64 remoteChainSelector; // ─╮ The chain ID of the source chain
    address receiver; // ───────────╯ The recipient of the tokens on the destination chain.
    uint256 amount; //                The amount of tokens to release or mint, denominated in the source token's decimals
    address localToken; //            The address on this chain of the token to release or mint
    /// @dev WARNING: sourcePoolAddress should be checked prior to any processing of funds. Make sure it matches the
    /// expected pool address for the given remoteChainSelector.
    bytes sourcePoolAddress; //       The address of the source pool, abi encoded in the case of EVM chains
    bytes sourcePoolData; //          The data received from the source pool to process the release or mint
    /// @dev WARNING: offchainTokenData is untrusted data.
    bytes offchainTokenData; //       The offchain data to process the release or mint
  }

  struct ReleaseOrMintOutV1 {
    // The number of tokens released or minted on the destination chain, denominated in the local token's decimals.
    // This value is expected to be equal to the ReleaseOrMintInV1.amount in the case where the source and destination
    // chain have the same number of decimals.
    uint256 destinationAmount;
  }
}
