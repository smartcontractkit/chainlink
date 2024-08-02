// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

struct OptionBoard {
  // board identifier
  uint256 id;
  // expiry of all strikes belonging to board
  uint256 expiry;
  // volatility component specific to board (boardIv * skew = vol of strike)
  uint256 iv;
  // admin settable flag blocking all trading on this board
  bool frozen;
  // list of all strikes belonging to this board
  uint256[] strikeIds;
}

interface OptionMarketInterface {
  function getLiveBoards() external view returns (uint256[] memory _liveBoards);

  function getOptionBoard(uint256 boardId) external view returns (OptionBoard memory);

  function settleExpiredBoard(uint256 boardId) external;
}
