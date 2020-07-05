pragma solidity 0.6.6;

interface AnswerValidatorInterface {
  function validate(
    uint256 previousRoundId,
    int256 previous,
    uint256 currentRoundId,
    int256 current
  ) external virtual returns (bool);
}
