// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import {ILogAutomation, Log} from "../interfaces/ILogAutomation.sol";

struct CheckData {
  uint256 checkBurnAmount;
  uint256 performBurnAmount;
  bytes32 eventSig;
}

contract SimpleLogUpkeepCounter is ILogAutomation {
  event PerformingUpkeep(
    address indexed from,
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    uint256 timeToPerform
  );

  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup
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

  function _checkDataConfig(CheckData memory) external {}

  function checkLog(Log calldata log, bytes calldata checkData) external view override returns (bool, bytes memory) {
    (uint256 checkBurnAmount, uint256 performBurnAmount, bytes32 eventSig) = abi.decode(
      checkData,
      (uint256, uint256, bytes32)
    );
    uint256 startGas = gasleft();
    bytes32 dummyIndex = blockhash(block.number - 1);
    bool dummy;
    // burn gas
    if (checkBurnAmount > 0) {
      while (startGas - gasleft() < checkBurnAmount) {
        dummy = dummy && dummyMap[dummyIndex]; // arbitrary storage reads
        dummyIndex = keccak256(abi.encode(dummyIndex, address(this)));
      }
    }
    if (log.topics[1] == eventSig) {
      return (true, abi.encode(log, checkData));
    }
    return (false, abi.encode(log, checkData));
  }

  function performUpkeep(bytes calldata performData) external override {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    lastBlock = block.number;
    counter = counter + 1;
    previousPerformBlock = lastBlock;
    (Log memory log, bytes memory extraData) = abi.decode(performData, (Log, bytes));
    timeToPerform = block.timestamp - log.timestamp;
    (uint256 checkBurnAmount, uint256 performBurnAmount, bytes32 eventSig) = abi.decode(
      extraData,
      (uint256, uint256, bytes32)
    );
    uint256 startGas = gasleft();
    bytes32 dummyIndex = blockhash(block.number - 1);
    bool dummy;
    // burn gas
    if (performBurnAmount > 0) {
      while (startGas - gasleft() < performBurnAmount) {
        dummy = dummy && dummyMap[dummyIndex]; // arbitrary storage reads
        dummyIndex = keccak256(abi.encode(dummyIndex, address(this)));
      }
    }
    emit PerformingUpkeep(tx.origin, initialBlock, lastBlock, previousPerformBlock, counter, timeToPerform);
  }
}
