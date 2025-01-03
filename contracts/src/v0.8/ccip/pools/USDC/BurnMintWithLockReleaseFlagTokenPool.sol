// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../../libraries/Pool.sol";
import {BurnMintTokenPool} from "../BurnMintTokenPool.sol";
import {LOCK_RELEASE_FLAG} from "./HybridLockReleaseUSDCTokenPool.sol";

/// @notice A standard BurnMintTokenPool with modified destPoolData so that the remote pool knows to release tokens
/// instead of minting. This enables interoperability with HybridLockReleaseUSDCTokenPool which uses
// the destPoolData to determine whether to mint or release tokens.
/// @dev The only difference between this contract and BurnMintTokenPool is the destPoolData returns the
/// abi-encoded LOCK_RELEASE_FLAG instead of an empty string.
contract BurnMintWithLockReleaseFlagTokenPool is BurnMintTokenPool {
  constructor(
    IBurnMintERC20 token,
    uint8 localTokenDecimals,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) BurnMintTokenPool(token, localTokenDecimals, allowlist, rmnProxy, router) {}

  /// @notice Burn the token in the pool
  /// @dev The _validateLockOrBurn check is an essential security check
  /// @dev Performs the exact same functionality as BurnMintTokenPool, but returns the LOCK_RELEASE_FLAG
  /// as the destPoolData to signal to the remote pool to release tokens instead of minting them.
  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external override returns (Pool.LockOrBurnOutV1 memory) {
    _validateLockOrBurn(lockOrBurnIn);

    _burn(lockOrBurnIn.amount);

    emit Burned(msg.sender, lockOrBurnIn.amount);

    // LOCK_RELEASE_FLAG = bytes4(keccak256("NO_CCTP_USE_LOCK_RELEASE"))
    return Pool.LockOrBurnOutV1({
      destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector),
      destPoolData: abi.encode(LOCK_RELEASE_FLAG)
    });
  }
}
