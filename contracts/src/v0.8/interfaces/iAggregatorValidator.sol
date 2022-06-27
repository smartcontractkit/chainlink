// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface iAggregatorValidator {
  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  ) external returns (bool);
}
