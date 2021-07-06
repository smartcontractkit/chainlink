// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

interface AggregatorValidatorInterface {
  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  ) external returns (bool);
  function setGasConfiguration(
    uint32 maximumGasPrice,
    uint256 gasCostL2
  ) external;
  function setRefundableAddress(
    address refundableAddress
  ) external;
}