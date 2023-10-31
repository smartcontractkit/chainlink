// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title Chainlink Functions DON billing interface.
interface IFunctionsBilling {
  /// @notice Return the current conversion from WEI of ETH to LINK from the configured Chainlink data feed
  /// @return weiPerUnitLink - The amount of WEI in one LINK
  function getWeiPerUnitLink() external view returns (uint256);

  /// @notice Determine the fee that will be split between Node Operators for servicing a request
  /// @param requestCBOR - CBOR encoded Chainlink Functions request data, use FunctionsRequest library to encode a request
  /// @return fee - Cost in Juels (1e18) of LINK
  function getDONFee(bytes memory requestCBOR) external view returns (uint72);

  /// @notice Determine the fee that will be paid to the Router owner for operating the network
  /// @return fee - Cost in Juels (1e18) of LINK
  function getAdminFee() external view returns (uint72);

  /// @notice Estimate the total cost that will be charged to a subscription to make a request: transmitter gas re-reimbursement, plus DON fee, plus Registry fee
  /// @param - subscriptionId An identifier of the billing account
  /// @param - data Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
  /// @param - callbackGasLimit Gas limit for the fulfillment callback
  /// @param - gasPriceWei The blockchain's gas price to estimate with
  /// @return - billedCost Cost in Juels (1e18) of LINK
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 callbackGasLimit,
    uint256 gasPriceWei
  ) external view returns (uint96);

  /// @notice Remove a request commitment that the Router has determined to be stale
  /// @param requestId - The request ID to remove
  function deleteCommitment(bytes32 requestId) external;

  /// @notice Oracle withdraw LINK earned through fulfilling requests
  /// @notice If amount is 0 the full balance will be withdrawn
  /// @param recipient where to send the funds
  /// @param amount amount to withdraw
  function oracleWithdraw(address recipient, uint96 amount) external;

  /// @notice Withdraw all LINK earned by Oracles through fulfilling requests
  function oracleWithdrawAll() external;
}
