// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import {ILogAutomation, Log} from "../interfaces/ILogAutomation.sol";

contract SimpleLogUpkeepCounter is ILogAutomation {

  event PerformingUpkeep(
    address indexed from,
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    uint256 timeToPerform
  );

  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  uint256 public timeToPerform;

  constructor() {
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
  }

  function checkLog(Log calldata log, bytes memory) external view override returns (bool, bytes memory) {
    return (true, abi.encode(log));
  }

  function performUpkeep(bytes calldata performData) external override {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    lastBlock = block.number;
    counter = counter + 1;
    previousPerformBlock = lastBlock;
    Log memory log = abi.decode(performData, (Log));
    timeToPerform = block.timestamp - log.timestamp;
    emit PerformingUpkeep(tx.origin, initialBlock, lastBlock, previousPerformBlock, counter, timeToPerform);
  }
}
