pragma solidity 0.6.6;

import '../dev/AnswerValidatorInterface.sol';

contract AnswerValidatorGasGuzzler is AnswerValidatorInterface {
  function validate(
    uint256 previousRoundId,
    int256 previous,
    uint256 currentRoundId,
    int256 current
  )
    external
    override
    returns (bool)
  {
    while (true) {
    }
    return false;
  }
}

