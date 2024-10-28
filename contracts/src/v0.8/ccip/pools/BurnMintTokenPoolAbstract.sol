// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../libraries/Pool.sol";
import {TokenPool} from "./TokenPool.sol";

abstract contract BurnMintTokenPoolAbstract is TokenPool {
  /// @notice Contains the specific burn call for a pool.
  /// @dev overriding this method allows us to create pools with different burn signatures
  /// without duplicating the underlying logic.
  function _burn(
    uint256 amount
  ) internal virtual;

  /// @notice Burn the token in the pool
  /// @dev The _validateLockOrBurn check is an essential security check
  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external virtual override returns (Pool.LockOrBurnOutV1 memory) {
    _validateLockOrBurn(lockOrBurnIn);

    _burn(lockOrBurnIn.amount);

    emit Burned(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  /// @notice Mint tokens from the pool to the recipient
  /// @dev The _validateReleaseOrMint check is an essential security check
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    _validateReleaseOrMint(releaseOrMintIn);

    // Mint to the receiver
    IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, releaseOrMintIn.amount);

    emit Minted(msg.sender, releaseOrMintIn.receiver, releaseOrMintIn.amount);

    return Pool.ReleaseOrMintOutV1({destinationAmount: releaseOrMintIn.amount});
  }
}
