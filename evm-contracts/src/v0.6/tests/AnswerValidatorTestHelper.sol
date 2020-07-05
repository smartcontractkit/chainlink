pragma solidity 0.6.6;

import '../dev/AnswerValidatorInterface.sol';

contract AnswerValidatorTestHelper is AnswerValidatorInterface {
  event Validated(
    int256 indexed previous,
    int256 indexed current
  );

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
    emit Validated(previous, current);
    return true;
  }
}

