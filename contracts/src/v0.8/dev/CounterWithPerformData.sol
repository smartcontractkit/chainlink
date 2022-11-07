// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../interfaces/AutomationCompatibleInterface.sol";

contract CounterWithPerformData is AutomationCompatibleInterface {
  /**
   * Public counter variable
   */
  uint256 public counter;
  bytes internal constant DATA_PADDING = "0xfffffffffffffffffffffffffffffffffffffffffff";

  /**
   * Use an interval in seconds and a timestamp to slow execution of Upkeep
   */
  uint256 public immutable interval;
  uint256 public lastTimeStamp;
  event Logger(string message, uint256 timestamp, uint256 blocknbr, bytes abc);

  constructor(uint256 updateInterval) {
    interval = updateInterval;
    lastTimeStamp = block.timestamp;

    counter = 0;
  }

  function checkUpkeep(bytes calldata checkData)
    external
    view
    override
    returns (bool upkeepNeeded, bytes memory performData)
  {
    upkeepNeeded = (block.timestamp - lastTimeStamp) > interval;
    bytes memory performData = bytes.concat(DATA_PADDING, checkData);
    return (upkeepNeeded, performData);
  }

  function performUpkeep(bytes calldata performData) external override {
    //We highly recommend revalidating the upkeep in the performUpkeep function

    if ((block.timestamp - lastTimeStamp) > interval) {
      lastTimeStamp = block.timestamp;
      counter = counter + 1;
      emit Logger("add 1", block.timestamp, block.number, performData);
    }
  }
}
