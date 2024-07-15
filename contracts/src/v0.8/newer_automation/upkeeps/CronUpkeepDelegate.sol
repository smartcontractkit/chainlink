// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {Cron, Spec} from "../libraries/internal/Cron.sol";

/**
 * @title The CronUpkeepDelegate contract
 * @notice This contract serves as a delegate for all instances of CronUpkeep. Those contracts
 * delegate their checkUpkeep calls onto this contract. Utilizing this pattern reduces the size
 * of the CronUpkeep contracts.
 */
contract CronUpkeepDelegate {
  using EnumerableSet for EnumerableSet.UintSet;
  using Cron for Spec;

  address private s_owner; // from ConfirmedOwner
  address private s_delegate;
  uint256 private s_nextCronJobID;
  EnumerableSet.UintSet private s_activeCronJobIDs;
  mapping(uint256 => uint256) private s_lastRuns;
  mapping(uint256 => Spec) private s_specs;
  mapping(uint256 => address) private s_targets;
  mapping(uint256 => bytes) private s_handlers;

  /**
   * @notice Get the id of an eligible cron job
   * @return upkeepNeeded signals if upkeep is needed, performData is an abi encoding
   * of the id and "next tick" of the eligible cron job
   */
  function checkUpkeep(bytes calldata) external view returns (bool, bytes memory) {
    // DEV: start at a random spot in the list so that checks are
    // spread evenly among cron jobs
    uint256 numCrons = s_activeCronJobIDs.length();
    if (numCrons == 0) {
      return (false, bytes(""));
    }
    uint256 startIdx = block.number % numCrons;
    bool result;
    bytes memory payload;
    (result, payload) = checkInRange(startIdx, numCrons);
    if (result) {
      return (result, payload);
    }
    (result, payload) = checkInRange(0, startIdx);
    if (result) {
      return (result, payload);
    }
    return (false, bytes(""));
  }

  /**
   * @notice checks the cron jobs in a given range
   * @param start the starting id to check (inclusive)
   * @param end the ending id to check (exclusive)
   * @return upkeepNeeded signals if upkeep is needed, performData is an abi encoding
   * of the id and "next tick" of the eligible cron job
   */
  function checkInRange(uint256 start, uint256 end) private view returns (bool, bytes memory) {
    uint256 id;
    uint256 lastTick;
    for (uint256 idx = start; idx < end; idx++) {
      id = s_activeCronJobIDs.at(idx);
      lastTick = s_specs[id].lastTick();
      if (lastTick > s_lastRuns[id]) {
        return (true, abi.encode(id, lastTick, s_targets[id], s_handlers[id]));
      }
    }
  }
}
