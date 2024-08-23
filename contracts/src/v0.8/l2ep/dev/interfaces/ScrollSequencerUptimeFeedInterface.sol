// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

// solhint-disable-next-line interface-starts-with-i
interface ScrollSequencerUptimeFeedInterface {
  function updateStatus(bool status, uint64 timestamp) external;
}
