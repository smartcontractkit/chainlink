// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../shared/interfaces/AggregatorV3Interface.sol";

contract VRFMockETHLINKAggregator is AggregatorV3Interface {
  int256 public answer;
  uint256 private blockTimestampDeduction = 0;

  constructor(int256 _answer) {
    answer = _answer;
  }

  function decimals() external pure override returns (uint8) {
    return 18;
  }

  function description() external pure override returns (string memory) {
    return "VRFMockETHLINKAggregator";
  }

  function version() external pure override returns (uint256) {
    return 1;
  }

  function getRoundData(
    uint80 /*_roundId*/
  )
    external
    view
    override
    returns (uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return (1, answer, getDeductedBlockTimestamp(), getDeductedBlockTimestamp(), 1);
  }

  function latestRoundData()
    external
    view
    override
    returns (uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return (1, answer, getDeductedBlockTimestamp(), getDeductedBlockTimestamp(), 1);
  }

  function getDeductedBlockTimestamp() internal view returns (uint256) {
    return block.timestamp - blockTimestampDeduction;
  }

  function setBlockTimestampDeduction(uint256 _blockTimestampDeduction) external {
    blockTimestampDeduction = _blockTimestampDeduction;
  }
}
