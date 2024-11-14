// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {AggregatorValidatorInterface} from "../../../shared/interfaces/AggregatorValidatorInterface.sol";

contract AggregatorValidator is AggregatorValidatorInterface {
  uint256 internal s_previousRoundId;
  int256 internal s_previousAnswer;
  uint256 internal s_currentRoundId;
  int256 internal s_currentAnswer;

  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  ) external override returns (bool) {
    s_previousRoundId = previousRoundId;
    s_previousAnswer = previousAnswer;
    s_currentRoundId = currentRoundId;
    s_currentAnswer = currentAnswer;

    return true;
  }

  function getLatestValidatedValues() public view returns (uint256, int256, uint256, int256) {
    return (s_previousRoundId, s_previousAnswer, s_currentRoundId, s_currentAnswer);
  }
}
