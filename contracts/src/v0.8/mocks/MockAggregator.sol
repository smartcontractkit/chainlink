// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IOffchainAggregator} from "../HeartbeatRequester.sol";

contract MockAggregator is IOffchainAggregator {
  int256 public s_answer;
  bool public newRoundCalled;

  function setLatestAnswer(int256 answer) public {
    s_answer = answer;
  }

  function latestAnswer() public view returns (int256) {
    return s_answer;
  }

  function requestNewRound() external override returns (uint80) {
    // do we need the actual logic of requestNewRound?
    // or is dummy requestNewRound ok?
    newRoundCalled = true;
    return 1;
  }
}
