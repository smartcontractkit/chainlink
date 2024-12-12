// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ISequencerUptimeFeed {
  function updateStatus(bool status, uint64 timestamp) external;
}
