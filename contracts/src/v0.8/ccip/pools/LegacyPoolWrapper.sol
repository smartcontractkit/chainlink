// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IPoolPriorTo1_5} from "../interfaces/IPoolPriorTo1_5.sol";

import {Pool} from "../libraries/Pool.sol";
import {TokenPool} from "./TokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

abstract contract LegacyPoolWrapper is TokenPool {
  using SafeERC20 for IERC20;

  event LegacyPoolChanged(IPoolPriorTo1_5 oldPool, IPoolPriorTo1_5 newPool);

  /// @dev The previous pool, if there is any. This is a property to make the older 1.0-1.4 pools
  /// compatible with the current 1.5 pool. To achieve this, we set the previous pool address to the
  /// currently deployed legacy pool. Then we configure this new pool as onRamp and offRamp on the legacy pools.
  /// In the case of a 1.4 pool, this new pool contract has to be set to the Router as well, as it validates
  /// who can call it through the router calls. This contract will always return itself as the only allowed ramp.
  /// @dev Can be address(0), this would indicate that this pool is operating as a normal pool as opposed to
  /// a proxy pool.
  IPoolPriorTo1_5 internal s_previousPool;

  constructor(
    IERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) TokenPool(token, allowlist, rmnProxy, router) {}

  // ================================================================
  // │                      Legacy Fallbacks                        │
  // ================================================================
  // Legacy fallbacks for older token pools that do not implement the new interface.

  /// @notice Legacy fallback for the 1.4 token pools.
  function getOnRamp(uint64) external view returns (address onRampAddress) {
    return address(this);
  }

  /// @notice Return true if the given offRamp is a configured offRamp for the given source chain.
  function isOffRamp(uint64 sourceChainSelector, address offRamp) external view returns (bool) {
    return offRamp == address(this) || s_router.isOffRamp(sourceChainSelector, offRamp);
  }

  /// @notice Configures the legacy fallback option. If the previous pool is set, this pool will act as a proxy for
  /// the legacy pool.
  /// @param prevPool The address of the previous pool.
  function setPreviousPool(IPoolPriorTo1_5 prevPool) external onlyOwner {
    IPoolPriorTo1_5 oldPrevPool = s_previousPool;
    s_previousPool = prevPool;

    emit LegacyPoolChanged(oldPrevPool, prevPool);
  }

  /// @notice Returns the address of the previous pool.
  function getPreviousPool() external view returns (address) {
    return address(s_previousPool);
  }

  function _hasLegacyPool() internal view returns (bool) {
    return address(s_previousPool) != address(0);
  }

  function _lockOrBurnLegacy(Pool.LockOrBurnInV1 memory lockOrBurnIn) internal {
    i_token.safeTransfer(address(s_previousPool), lockOrBurnIn.amount);
    s_previousPool.lockOrBurn(
      lockOrBurnIn.originalSender, lockOrBurnIn.receiver, lockOrBurnIn.amount, lockOrBurnIn.remoteChainSelector, ""
    );
  }

  /// @notice This call converts the arguments from a >=1.5 pool call to those of a <1.5 pool call, and uses these
  /// to call the previous pool.
  /// @param releaseOrMintIn The 1.5 style release or mint arguments.
  /// @dev Overwrites the receiver so the previous pool sends the tokens to the sender of this call, which is the
  /// offRamp. This is due to the older pools sending funds directly to the receiver, while the new pools do a hop
  /// through the offRamp to ensure the correct tokens are sent.
  /// @dev Since extraData has never been used in LockRelease or MintBurn token pools, we can safely ignore it.
  function _releaseOrMintLegacy(Pool.ReleaseOrMintInV1 memory releaseOrMintIn) internal {
    s_previousPool.releaseOrMint(
      releaseOrMintIn.originalSender,
      releaseOrMintIn.receiver,
      releaseOrMintIn.amount,
      releaseOrMintIn.remoteChainSelector,
      ""
    );
  }
}
