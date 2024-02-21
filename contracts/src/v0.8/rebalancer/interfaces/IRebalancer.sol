// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IBridgeAdapter} from "./IBridge.sol";

interface IRebalancer {
  struct SendLiquidityParams {
    /// @notice The amount of tokens to be sent to the remote chain.
    uint256 amount;
    /// @notice The amount of native that should be sent by the rebalancer in the sendERC20 call.
    /// @notice This is used to pay for the bridge fees.
    uint256 nativeBridgeFee;
    /// @notice The selector of the remote chain.
    uint64 remoteChainSelector;
    /// @notice The bridge data that should be passed to the sendERC20 call.
    bytes bridgeData;
  }

  struct ReceiveLiquidityParams {
    uint256 amount;
    uint64 remoteChainSelector;
    bytes bridgeData;
  }

  struct LiquidityInstructions {
    SendLiquidityParams[] sendLiquidityParams;
    ReceiveLiquidityParams[] receiveLiquidityParams;
  }

  struct CrossChainRebalancerArgs {
    address remoteRebalancer;
    IBridgeAdapter localBridge;
    address remoteToken;
    uint64 remoteChainSelector;
    bool enabled;
  }

  /// @notice Returns the current liquidity in the liquidity container.
  function getLiquidity() external view returns (uint256 currentLiquidity);

  function getAllCrossChainRebalancers() external view returns (CrossChainRebalancerArgs[] memory);
}
