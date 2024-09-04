// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";

interface IPriceRegistry {
  /// @notice Update the price for given tokens and gas prices for given chains.
  /// @param priceUpdates The price updates to apply.
  function updatePrices(Internal.PriceUpdates memory priceUpdates) external;

  /// @notice Get the `tokenPrice` for a given token.
  /// @param token The token to get the price for.
  /// @return tokenPrice The tokenPrice for the given token.
  function getTokenPrice(address token) external view returns (Internal.TimestampedPackedUint224 memory);

  /// @notice Get the `tokenPrice` for a given token, checks if the price is valid.
  /// @param token The token to get the price for.
  /// @return tokenPrice The tokenPrice for the given token if it exists and is valid.
  function getValidatedTokenPrice(address token) external view returns (uint224);

  /// @notice Get the `tokenPrice` for an array of tokens.
  /// @param tokens The tokens to get prices for.
  /// @return tokenPrices The tokenPrices for the given tokens.
  function getTokenPrices(address[] calldata tokens) external view returns (Internal.TimestampedPackedUint224[] memory);

  /// @notice Get an encoded `gasPrice` for a given destination chain ID.
  /// The 224-bit result encodes necessary gas price components.
  /// On L1 chains like Ethereum or Avax, the only component is the gas price.
  /// On Optimistic Rollups, there are two components - the L2 gas price, and L1 base fee for data availability.
  /// On future chains, there could be more or differing price components.
  /// PriceRegistry does not contain chain-specific logic to parse destination chain price components.
  /// @param destChainSelector The destination chain to get the price for.
  /// @return gasPrice The encoded gasPrice for the given destination chain ID.
  function getDestinationChainGasPrice(
    uint64 destChainSelector
  ) external view returns (Internal.TimestampedPackedUint224 memory);

  /// @notice Gets the fee token price and the gas price, both denominated in dollars.
  /// @param token The source token to get the price for.
  /// @param destChainSelector The destination chain to get the gas price for.
  /// @return tokenPrice The price of the feeToken in 1e18 dollars per base unit.
  /// @return gasPrice The price of gas in 1e18 dollars per base unit.
  function getTokenAndGasPrices(
    address token,
    uint64 destChainSelector
  ) external view returns (uint224 tokenPrice, uint224 gasPrice);

  /// @notice Convert a given token amount to target token amount.
  /// @param fromToken The given token address.
  /// @param fromTokenAmount The given token amount.
  /// @param toToken The target token address.
  /// @return toTokenAmount The target token amount.
  function convertTokenAmount(
    address fromToken,
    uint256 fromTokenAmount,
    address toToken
  ) external view returns (uint256 toTokenAmount);

  /// @notice Get the list of fee tokens.
  /// @return feeTokens The tokens set as fee tokens.
  function getFeeTokens() external view returns (address[] memory);

  /// @notice Validates the ccip message & returns the fee
  /// @param destChainSelector The destination chain selector.
  /// @param message The message to get quote for.
  /// @return feeTokenAmount The amount of fee token needed for the fee, in smallest denomination of the fee token.
  function getValidatedFee(
    uint64 destChainSelector,
    Client.EVM2AnyMessage calldata message
  ) external view returns (uint256 feeTokenAmount);

  /// @notice Converts the extraArgs to the latest version and returns the converted message fee in juels
  /// @param destChainSelector destination chain selector to process
  /// @param feeToken Fee token address used to pay for message fees
  /// @param feeTokenAmount Fee token amount
  /// @param extraArgs Message extra args that were passed in by the client
  /// @return msgFeeJuels message fee in juels
  /// @return isOutOfOrderExecution true if the message should be executed out of order
  /// @return convertedExtraArgs extra args converted to the latest family-specific args version
  function processMessageArgs(
    uint64 destChainSelector,
    address feeToken,
    uint256 feeTokenAmount,
    bytes memory extraArgs
  ) external view returns (uint256 msgFeeJuels, bool isOutOfOrderExecution, bytes memory convertedExtraArgs);

  /// @notice Validates pool return data
  /// @param destChainSelector Destination chain selector to which the token amounts are sent to
  /// @param rampTokenAmounts Token amounts with populated pool return data
  /// @param sourceTokenAmounts Token amounts originally sent in a Client.EVM2AnyMessage message
  /// @return destExecData Destination chain execution data
  function processPoolReturnData(
    uint64 destChainSelector,
    Internal.RampTokenAmount[] memory rampTokenAmounts,
    Client.EVMTokenAmount[] calldata sourceTokenAmounts
  ) external view returns (bytes[] memory);
}
