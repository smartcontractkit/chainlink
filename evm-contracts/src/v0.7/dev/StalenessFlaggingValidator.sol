// SPDX-License-Identifier: MIT
pragma solidity 0.7.0;

import "./ConfirmedOwner.sol";
import "../vendor/SafeMathChainlink.sol";
import "../interfaces/FlagsInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";

contract StalenessFlaggingValidator is ConfirmedOwner {
  using SafeMathChainlink for uint256;

  FlagsInterface private s_flags;
  uint256 private s_threshold;

  event FlagsAddressUpdated(
    address indexed previous,
    address indexed current
  );
  event FlaggingThresholdUpdated(
    uint256 indexed previous,
    uint256 indexed current
  );

  constructor(address flagsAddress, uint256 flaggingThreshold) 
    public 
    ConfirmedOwner(msg.sender)
  {
    setFlagsAddress(flagsAddress);
    setThreshold(flaggingThreshold);
  }

  function update(address[] calldata aggregators) external {
    uint256 currentTimestamp = block.timestamp;

    for (uint256 i = 0; i < aggregators.length; i++) {
      address aggregator = aggregators[i];
      if (isStale(aggregator, currentTimestamp)) {
        s_flags.raiseFlag(aggregator);
      }
    }
  }

  function check(address[] memory aggregators) public returns (bool) {
    uint256 currentTimestamp = block.timestamp;

    for (uint256 i = 0; i < aggregators.length; i++) {
      if (isStale(aggregators[i], currentTimestamp)) {
        return true;
      }
    }
    return false;
  }

  function setThreshold(uint256 flaggingThreshold) public onlyOwner() {
    uint256 previous = s_threshold;
    if (previous != flaggingThreshold) {
      s_threshold = flaggingThreshold;
      emit FlaggingThresholdUpdated(previous, flaggingThreshold);
    }
  }

  /**
   * @notice updates the flagging contract address for raising flags
   * @param flagsAddress sets the address of the flags contract
   */
  function setFlagsAddress(address flagsAddress)
    public
    onlyOwner()
  {
    address previous = address(s_flags);
    if (previous != flagsAddress) {
      s_flags = FlagsInterface(flagsAddress);
      emit FlagsAddressUpdated(previous, flagsAddress);
    }
  }

  function isStale(address aggregator, uint256 currentTimestamp) 
    private 
    returns (bool stale)
  {
    (,,,uint updatedAt,) = AggregatorV3Interface(aggregator).latestRoundData();
    uint256 diff = currentTimestamp.sub(updatedAt);
    if (diff > s_threshold) {
      stale = true;
    }
  }
}