pragma solidity 0.6.6;

import '../dev/AnswerValidatorInterface.sol';

contract AnswerValidatorTestHelper is AnswerValidatorInterface {
  event Validated(
    int256 indexed previous,
    int256 indexed current
  );

  function validate(
    int256 previous,
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

