// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../shared/interfaces/AggregatorValidatorInterface.sol";

contract MockAggregatorValidator is AggregatorValidatorInterface {
  uint8 immutable id;

  constructor(uint8 id_) {
    id = id_;
  }

  event ValidateCalled(
    uint8 id,
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  );

  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  ) external override returns (bool) {
    emit ValidateCalled(id, previousRoundId, previousAnswer, currentRoundId, currentAnswer);
    return true;
  }
}
