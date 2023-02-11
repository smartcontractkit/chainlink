// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockComp {
  mapping(address => uint96) public s_priorVotesMapping;

  function getPriorVotes(
    address account,
    uint256 /* blockNumber */
  ) external view returns (uint96) {
    return s_priorVotesMapping[account];
  }

  function setPriorVotes(address account, uint96 votes) external {
    s_priorVotesMapping[account] = votes;
  }
}
