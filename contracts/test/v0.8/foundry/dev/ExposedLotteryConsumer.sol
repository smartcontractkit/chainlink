// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {LotteryConsumer} from "../../../../src/v0.8/dev/LotteryConsumer.sol";

contract ExposedLotteryConsumer is LotteryConsumer {
  constructor(address _vrfCoordinator) LotteryConsumer(_vrfCoordinator) {

  }

  function fulfillRandomWordsExternal(uint256 requestId, uint256[] memory randomWords) external {
    fulfillRandomWords(requestId, randomWords);
  }
}
