// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../libraries/Pool.sol";
import {LegacyPoolWrapper} from "./LegacyPoolWrapper.sol";

import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

contract BurnWithFromMintTokenPoolAndProxy is ITypeAndVersion, LegacyPoolWrapper {
  using SafeERC20 for IBurnMintERC20;

  string public constant override typeAndVersion = "BurnWithFromMintTokenPoolAndProxy 1.5.0";

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) LegacyPoolWrapper(token, allowlist, rmnProxy, router) {
    // Some tokens allow burning from the sender without approval, but not all do.
    // To be safe, we approve the pool to burn from the pool.
    token.safeIncreaseAllowance(address(this), type(uint256).max);
  }

  /// @notice Burn the token in the pool
  /// @dev The _validateLockOrBurn check is an essential security check
  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external virtual override returns (Pool.LockOrBurnOutV1 memory) {
    _validateLockOrBurn(lockOrBurnIn);

    if (!_hasLegacyPool()) {
      IBurnMintERC20(address(i_token)).burnFrom(address(this), lockOrBurnIn.amount);
    } else {
      _lockOrBurnLegacy(lockOrBurnIn);
    }

    emit Burned(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  /// @notice Mint tokens from the pool to the recipient
  /// @dev The _validateReleaseOrMint check is an essential security check
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    _validateReleaseOrMint(releaseOrMintIn);

    if (!_hasLegacyPool()) {
      IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, releaseOrMintIn.amount);
    } else {
      _releaseOrMintLegacy(releaseOrMintIn);
    }

    emit Minted(msg.sender, releaseOrMintIn.receiver, releaseOrMintIn.amount);

    return Pool.ReleaseOrMintOutV1({destinationAmount: releaseOrMintIn.amount});
  }
}
