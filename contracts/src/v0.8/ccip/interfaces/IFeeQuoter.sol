// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {IPriceRegistry} from "./IPriceRegistry.sol";

interface IFeeQuoter is IPriceRegistry {
  /// @notice Validates the ccip message & returns the fee
  /// @param destChainSelector The destination chain selector.
  /// @param message The message to get quote for.
  /// @return feeTokenAmount The amount of fee token needed for the fee, in smallest denomination of the fee token.
  function getValidatedFee(
    uint64 destChainSelector,
    Client.EVM2AnyMessage calldata message
  ) external view returns (uint256 feeTokenAmount);

  /// @notice Converts the extraArgs to the latest version and returns the converted message fee in juels
  /// @notice Validates pool return data
  /// @param destChainSelector destination chain selector to process
  /// @param feeToken Fee token address used to pay for message fees
  /// @param feeTokenAmount Fee token amount
  /// @param extraArgs Message extra args that were passed in by the client
  /// @param onRampTokenTransfers Token amounts with populated pool return data
  /// @param sourceTokenAmounts Token amounts originally sent in a Client.EVM2AnyMessage message
  /// @return msgFeeJuels message fee in juels
  /// @return isOutOfOrderExecution true if the message should be executed out of order
  /// @return convertedExtraArgs extra args converted to the latest family-specific args version
  /// @return destExecDataPerToken Destination chain execution data
  function processMessageArgs(
    uint64 destChainSelector,
    address feeToken,
    uint256 feeTokenAmount,
    bytes memory extraArgs,
    Internal.EVM2AnyTokenTransfer[] calldata onRampTokenTransfers,
    Client.EVMTokenAmount[] calldata sourceTokenAmounts
  )
    external
    view
    returns (
      uint256 msgFeeJuels,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      bytes[] memory destExecDataPerToken
    );
}
