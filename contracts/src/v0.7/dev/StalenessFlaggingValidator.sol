// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../ConfirmedOwner.sol";
import "../vendor/SafeMathChainlink.sol";
import "../interfaces/FlagsInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../interfaces/KeeperCompatibleInterface.sol";

contract StalenessFlaggingValidator is ConfirmedOwner, KeeperCompatibleInterface {
  using SafeMathChainlink for uint256;

  FlagsInterface private s_flags;
  mapping(address => uint256) private s_thresholdSeconds;

  event FlagsAddressUpdated(
    address indexed previous,
    address indexed current
  );
  event FlaggingThresholdUpdated(
    address indexed aggregator,
    uint256 indexed previous,
    uint256 indexed current
  );

  /**
   * @notice Create a new StalenessFlaggingValidator
   * @param flagsAddress Address of the flag contract
   * @dev Ensure that this contract has sufficient write permissions
   * on the flag contract
   */
  constructor(
    address flagsAddress
  )
    ConfirmedOwner(msg.sender)
  {
    setFlagsAddress(flagsAddress);
  }

  /**
   * @notice Updates the flagging contract address for raising flags
   * @param flagsAddress sets the address of the flags contract
   */
  function setFlagsAddress(
    address flagsAddress
  )
    public
    onlyOwner()
  {
    address previous = address(s_flags);
    if (previous != flagsAddress) {
      s_flags = FlagsInterface(flagsAddress);
      emit FlagsAddressUpdated(previous, flagsAddress);
    }
  }

  /**
   * @notice Set the threshold limits for each aggregator
   * @dev parameters must be same length
   * @param aggregators address[] memory
   * @param flaggingThresholds uint256[] memory
   */
  function setThresholds(
    address[] memory aggregators,
    uint256[] memory flaggingThresholds
  )
    public 
    onlyOwner()
  {
    require(aggregators.length == flaggingThresholds.length, "Different sized arrays");
    for (uint256 i = 0; i < aggregators.length; i++) {
      address aggregator = aggregators[i];
      uint256 previousThreshold = s_thresholdSeconds[aggregator];
      uint256 newThreshold = flaggingThresholds[i];
      if (previousThreshold != newThreshold) {
        s_thresholdSeconds[aggregator] = newThreshold;
        emit FlaggingThresholdUpdated(aggregator, previousThreshold, newThreshold);
      }
    }
  }

  /**
   * @notice Check for staleness in an array of aggregators
   * @dev If any of the aggregators are stale, this function will return true,
   * otherwise false
   * @param aggregators address[] memory
   * @return address[] memory stale aggregators
   */
  function check(
    address[] memory aggregators
  )
    public
    view
    returns (
      address[] memory
    )
  {
    uint256 currentTimestamp = block.timestamp;
    address[] memory staleAggregators = new address[](aggregators.length);
    uint256 staleCount = 0;
    for (uint256 i = 0; i < aggregators.length; i++) {
      address aggregator = aggregators[i];
      if (isStale(aggregator, currentTimestamp)) {
        staleAggregators[staleCount] = aggregator;
        staleCount++;
      }
    }

    if (aggregators.length != staleCount) {
      assembly {
        mstore(staleAggregators, staleCount)
      }
    }
    return staleAggregators;
  }

  /**
   * @notice Check for staleness in an array of aggregators, raise a flag
   * on the flags contract for each aggregator that is stale
   * @dev This contract must have write permissions on the flags contract
   * @param aggregators address[] memory
   * @return address[] memory stale aggregators
   */
  function update(
    address[] memory aggregators
  )
    public
    returns (
      address[] memory
    )
  {
    address[] memory staleAggregators = check(aggregators);
    s_flags.raiseFlags(staleAggregators);
    return staleAggregators;
  }

  /**
   * @notice Check for staleness in an array of aggregators
   * @dev Overriding KeeperInterface
   * @param data bytes encoded address array
   * @return needsUpkeep bool indicating whether upkeep needs to be performed
   * @return staleAggregators bytes encoded address array of stale aggregator addresses
   */
  function checkUpkeep(
    bytes calldata data
  )
    external
    view
    override
    returns (
      bool,
      bytes memory
    )
  {
    address[] memory staleAggregators = check(abi.decode(data, (address[])));
    bool needsUpkeep = (staleAggregators.length > 0);
    return (needsUpkeep, abi.encode(staleAggregators));
  }

  /**
   * @notice Check for staleness in an array of aggregators, raise a flag
   * on the flags contract for each aggregator that is stale
   * @dev Overriding KeeperInterface
   * @param data bytes encoded address array
   */
  function performUpkeep(
    bytes calldata data
  )
    external
    override
  {
    update(abi.decode(data, (address[])));
  }

  /**
   * @notice Get the threshold of an aggregator
   * @param aggregator address
   * @return uint256
   */
  function threshold(
    address aggregator
  )
    external
    view
    returns (uint256)
  {
    return s_thresholdSeconds[aggregator];
  }

  /**
   * @notice Get the flags address
   * @return address
   */
  function flags()
    external
    view
    returns (
      address
    )
  {
    return address(s_flags);
  }

  /**
   * @notice Check if an aggregator is stale.
   * @dev Staleness is where an aggregator's `updatedAt` field is older
   * than the threshold set for it in this contract
   * @param aggregator address
   * @param currentTimestamp uint256
   * @return stale bool
   */
  function isStale(
    address aggregator,
    uint256 currentTimestamp
  ) 
    private
    view
    returns (
      bool stale
    )
  {
    if (s_thresholdSeconds[aggregator] == 0) {
      return false;
    }
    (,,,uint updatedAt,) = AggregatorV3Interface(aggregator).latestRoundData();
    uint256 diff = currentTimestamp.sub(updatedAt);
    if (diff > s_thresholdSeconds[aggregator]) {
      return true;
    }
    return false;
  }
}
