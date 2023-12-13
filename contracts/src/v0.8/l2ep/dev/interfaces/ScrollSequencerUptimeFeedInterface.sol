// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

interface ScrollSequencerUptimeFeedInterface {
  function updateStatus(bool status, uint64 timestamp) external;
}
