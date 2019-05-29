pragma solidity 0.4.24;

interface CurrentAnswerInterface {
  function currentAnswer() external returns (int256);
  function updatedHeight() external returns (uint256);
}
