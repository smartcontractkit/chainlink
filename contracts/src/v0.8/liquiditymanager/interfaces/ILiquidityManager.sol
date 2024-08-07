// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IBridgeAdapter} from "./IBridge.sol";

interface ILiquidityManager {
  /// @notice Parameters for sending liquidity to a remote chain.
  /// @param amount The amount of tokens to be sent to the remote chain.
  /// @param nativeBridgeFee The amount of native that should be sent by the liquiditymanager in the sendERC20 call.
  ///        Used to pay for the bridge fees.
  /// @param remoteChainSelector The selector of the remote chain.
  /// @param bridgeData The bridge data that should be passed to the sendERC20 call.
  struct SendLiquidityParams {
    uint256 amount;
    uint256 nativeBridgeFee;
    uint64 remoteChainSelector;
    bytes bridgeData;
  }

  /// @notice Parameters for receiving liquidity from a remote chain.
  /// @param amount The amount of tokens to be received from the remote chain.
  /// @param remoteChainSelector The selector of the remote chain.
  /// @param bridgeData The bridge data that should be passed to the finalizeWithdrawERC20 call.
  /// @param shouldWrapNative Whether the received native token should be wrapped into wrapped native.
  ///        This is needed for when the bridge being used doesn't bridge wrapped native but native directly.
  struct ReceiveLiquidityParams {
    uint256 amount;
    uint64 remoteChainSelector;
    bool shouldWrapNative;
    bytes bridgeData;
  }

  /// @notice Instructions for the rebalancer on what to do with the available liquidity.
  /// @param sendLiquidityParams The parameters for sending liquidity to a remote chain.
  /// @param receiveLiquidityParams The parameters for receiving liquidity from a remote chain.
  struct LiquidityInstructions {
    SendLiquidityParams[] sendLiquidityParams;
    ReceiveLiquidityParams[] receiveLiquidityParams;
  }

  /// @notice Parameters for adding a cross-chain rebalancer.
  /// @param remoteRebalancer The address of the remote rebalancer.
  /// @param localBridge The local bridge adapter address.
  /// @param remoteToken The address of the remote token.
  /// @param remoteChainSelector The selector of the remote chain.
  /// @param enabled Whether the rebalancer is enabled.
  struct CrossChainRebalancerArgs {
    address remoteRebalancer;
    IBridgeAdapter localBridge;
    address remoteToken;
    uint64 remoteChainSelector;
    bool enabled;
  }

  /// @notice Returns the current liquidity in the liquidity container.
  /// @return currentLiquidity The current liquidity in the liquidity container.
  function getLiquidity() external view returns (uint256 currentLiquidity);

  /// @notice Returns all the cross-chain rebalancers.
  /// @return All the cross-chain rebalancers.
  function getAllCrossChainRebalancers() external view returns (CrossChainRebalancerArgs[] memory);
}
