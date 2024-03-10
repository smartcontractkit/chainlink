// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../interfaces/AggregatorV3Interface.sol";

contract MockGASAggregator is AggregatorV3Interface {
    int256 public answer;
    constructor (int256 _answer) public {
        answer = _answer;
    }
    function decimals() external override view returns (uint8) {
        return 18;
    }
    function description() external override view returns (string memory) {
        return "MockGASAggregator";
    }
    function version() external override view returns (uint256) {
        return 1;
    }
    function getRoundData(uint80 _roundId) external override view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    ) {
        return (1, answer, block.timestamp, block.timestamp, 1);
    }
    function latestRoundData() external override view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    ) {
        return (1, answer, block.timestamp, block.timestamp, 1);
    }
}