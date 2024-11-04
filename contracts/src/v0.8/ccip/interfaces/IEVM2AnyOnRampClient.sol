// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IPoolV1} from "./IPool.sol";

import {Client} from "../libraries/Client.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

interface IEVM2AnyOnRampClient {
  /// @notice Get the fee for a given ccip message.
  /// @param destChainSelector The destination chain selector.
  /// @param message The message to calculate the cost for.
  /// @return fee The calculated fee.
  function getFee(uint64 destChainSelector, Client.EVM2AnyMessage calldata message) external view returns (uint256 fee);

  /// @notice Get the pool for a specific token.
  /// @param destChainSelector The destination chain selector.
  /// @param sourceToken The source chain token to get the pool for.
  /// @return pool Token pool.
  function getPoolBySourceToken(uint64 destChainSelector, IERC20 sourceToken) external view returns (IPoolV1);

  /// @notice Gets a list of all supported source chain tokens.
  /// @param destChainSelector The destination chain selector.
  /// @return tokens The addresses of all tokens that this onRamp supports the given destination chain.
  function getSupportedTokens(
    uint64 destChainSelector
  ) external view returns (address[] memory tokens);

  /// @notice Send a message to the remote chain.
  /// @dev only callable by the Router.
  /// @dev approve() must have already been called on the token using the this ramp address as the spender.
  /// @dev if the contract is paused, this function will revert.
  /// @param destChainSelector The destination chain selector.
  /// @param message Message struct to send.
  /// @param feeTokenAmount Amount of fee tokens for payment.
  /// @param originalSender The original initiator of the CCIP request.
  /// @return messageId The message id.
  function forwardFromRouter(
    uint64 destChainSelector,
    Client.EVM2AnyMessage memory message,
    uint256 feeTokenAmount,
    address originalSender
  ) external returns (bytes32);
}
