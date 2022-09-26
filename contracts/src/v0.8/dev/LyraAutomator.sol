// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../interfaces/AutomationCompatibleInterface.sol";
import "../ConfirmedOwner.sol";
import "./OptionMarketInterface.sol";

contract LyraAutomator is AutomationCompatibleInterface, ConfirmedOwner {
  OptionMarketInterface public s_optionMarket;

  event BoardExpired(uint256 indexed boardId, uint256 blocknumber);
  event SettlementFailed(uint256 indexed boardId, bytes lowLevelData);

  constructor(OptionMarketInterface optionMarket) ConfirmedOwner(msg.sender) {
    s_optionMarket = optionMarket;
  }

  function setOptionMarket(OptionMarketInterface optionMarket) external onlyOwner {
    s_optionMarket = optionMarket;
  }

  function checkUpkeep(bytes calldata checkData)
    external
    override
    returns (bool upkeepNeeded, bytes memory performData)
  {
    uint256[] memory liveBoards = s_optionMarket.getLiveBoards();
    uint256 index = 0;

    for (uint256 i = 0; i < liveBoards.length; i++) {
      uint256 boardId = liveBoards[i];
      OptionBoard memory board = s_optionMarket.getOptionBoard(boardId);
      if (board.expiry < block.timestamp) {
        liveBoards[index++] = boardId;
      }
    }

    if (index > 0) {
      return (true, abi.encode(index, liveBoards));
    }

    return (false, "");
  }

  function performUpkeep(bytes calldata performData) external override {
    (uint256 index, uint256[] memory boardIds) = abi.decode(performData, (uint256, uint256[]));
    if (index == 0) {
      return;
    }

    for (uint256 i = 0; i < index; i++) {
      uint256 boardId = boardIds[i];

      try s_optionMarket.settleExpiredBoard(boardId) {
        emit BoardExpired(boardId, block.number);
      } catch (bytes memory lowLevelData) {
        emit SettlementFailed(boardId, lowLevelData);
      }
    }
  }
}
