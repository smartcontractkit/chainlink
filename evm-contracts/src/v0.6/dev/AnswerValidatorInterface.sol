pragma solidity 0.6.6;

interface AnswerValidatorInterface {
  function validate(
    int256 previous,
    int256 current
  ) external virtual returns (bool);
}
