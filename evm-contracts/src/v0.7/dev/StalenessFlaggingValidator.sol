// SPDX-License-Identifier: MIT
pragma solidity 0.7.0;

import "./ConfirmedOwner.sol";
import "../vendor/SafeMathChainlink.sol";
import "../interfaces/FlagsInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";

contract StalenessFlaggingValidator is ConfirmedOwner {
  using SafeMathChainlink for uint256;

  FlagsInterface private s_flags;
  mapping(address => uint256) s_thresholds;

  event FlagsAddressUpdated(
    address indexed previous,
    address indexed current
  );
  event FlaggingThresholdUpdated(
    uint256 indexed previous,
    uint256 indexed current
  );

  constructor(address flagsAddress) 
    public 
    ConfirmedOwner(msg.sender)
  {
    setFlagsAddress(flagsAddress);
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

  function setThresholds(address[] memory aggregators, uint256[] memory flaggingThresholds)
    public 
    onlyOwner()
  {
    require(aggregators.length == flaggingThresholds.length, "Different sized arrays");
    for (uint256 i = 0; i < aggregators.length; i++) {
      address aggregator = aggregators[i];
      uint256 previousThreshold = s_thresholds[aggregator];
      uint newThreshold = flaggingThresholds[i];
      if (previousThreshold != newThreshold) {
        s_thresholds[aggregator] = newThreshold;
        emit FlaggingThresholdUpdated(previousThreshold, newThreshold);
      }
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

  function threshold(address aggregator) public view returns (uint256) {
    return s_thresholds[aggregator];
  }

  function flags() public view returns (address) {
    return address(s_flags);
  }

  function isStale(address aggregator, uint256 currentTimestamp) 
    private 
    returns (bool stale)
  {
    if (s_thresholds[aggregator] == 0) {
      return false;
    }
    (,,,uint updatedAt,) = AggregatorV3Interface(aggregator).latestRoundData();
    uint256 diff = currentTimestamp.sub(updatedAt);
    if (diff > s_thresholds[aggregator]) {
      stale = true;
    }
  }
}