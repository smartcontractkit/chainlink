// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IPool} from "./pools/IPool.sol";

import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.0/token/ERC20/IERC20.sol";

interface IEVM2AnyOnRamp {
  /// @notice Get the fee for a given ccip message
  /// @param message The message to calculate the cost for
  /// @return fee The calculated fee
  function getFee(Client.EVM2AnyMessage calldata message) external view returns (uint256 fee);

  /// @notice Get the pool for a specific token
  /// @param sourceToken The source chain token to get the pool for
  /// @return pool Token pool
  function getPoolBySourceToken(IERC20 sourceToken) external view returns (IPool);

  /// @notice Gets a list of all supported source chain tokens.
  /// @return tokens The addresses of all tokens that this onRamp supports for sending.
  function getSupportedTokens() external view returns (address[] memory tokens);

  /// @notice Gets the next sequence number to be used in the onRamp
  /// @return the next sequence number to be used
  function getExpectedNextSequenceNumber() external view returns (uint64);

  /// @notice Get the next nonce for a given sender
  /// @param sender The sender to get the nonce for
  /// @return nonce The next nonce for the sender
  function getSenderNonce(address sender) external view returns (uint64 nonce);

  /// @notice Adds and removed token pools.
  /// @param removes The tokens and pools to be removed
  /// @param adds The tokens and pools to be added.
  function applyPoolUpdates(Internal.PoolUpdate[] memory removes, Internal.PoolUpdate[] memory adds) external;

  /// @notice Send a message to the remote chain
  /// @dev only callable by the Router
  /// @dev approve() must have already been called on the token using the this ramp address as the spender.
  /// @dev if the contract is paused, this function will revert.
  /// @param message Message struct to send
  /// @param originalSender The original initiator of the CCIP request
  function forwardFromRouter(
    Client.EVM2AnyMessage memory message,
    uint256 feeTokenAmount,
    address originalSender
  ) external returns (bytes32);
}
