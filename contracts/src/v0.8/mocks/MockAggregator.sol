// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockAggregator {

  int256 public s_answer;

  function setLatestAnswer(
    int256 answer
  )
    public
  {
    s_answer = answer;
  }
  
  function latestAnswer()
    public
    view
    returns(
      int256
    )
  {
    return s_answer;
  }
}