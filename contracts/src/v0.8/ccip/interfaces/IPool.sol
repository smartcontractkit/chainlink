// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Pool} from "../libraries/Pool.sol";

import {IERC165} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

/// @notice Shared public interface for multiple V1 pool types.
/// Each pool type handles a different child token model e.g. lock/unlock, mint/burn.
interface IPoolV1 is IERC165 {
  /// @notice Lock tokens into the pool or burn the tokens.
  /// @param lockOrBurnIn Encoded data fields for the processing of tokens on the source chain.
  /// @return lockOrBurnOut Encoded data fields for the processing of tokens on the destination chain.
  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external returns (Pool.LockOrBurnOutV1 memory lockOrBurnOut);

  /// @notice Releases or mints tokens to the receiver address.
  /// @param releaseOrMintIn All data required to release or mint tokens.
  /// @return releaseOrMintOut The amount of tokens released or minted on the local chain, denominated
  /// in the local token's decimals.
  /// @dev The offramp asserts that the balanceOf of the receiver has been incremented by exactly the number
  /// of tokens that is returned in ReleaseOrMintOutV1.destinationAmount. If the amounts do not match, the tx reverts.
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external returns (Pool.ReleaseOrMintOutV1 memory);

  /// @notice Checks whether a remote chain is supported in the token pool.
  /// @param remoteChainSelector The selector of the remote chain.
  /// @return true if the given chain is a permissioned remote chain.
  function isSupportedChain(
    uint64 remoteChainSelector
  ) external view returns (bool);

  /// @notice Returns if the token pool supports the given token.
  /// @param token The address of the token.
  /// @return true if the token is supported by the pool.
  function isSupportedToken(
    address token
  ) external view returns (bool);
}
