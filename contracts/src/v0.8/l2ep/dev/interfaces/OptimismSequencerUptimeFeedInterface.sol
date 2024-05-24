// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// solhint-disable-next-line interface-starts-with-i
interface OptimismSequencerUptimeFeedInterface {
  function updateStatus(bool status, uint64 timestamp) external;
}
