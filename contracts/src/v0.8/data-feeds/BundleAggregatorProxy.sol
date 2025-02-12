// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";
import {IBundleAggregator} from "./interfaces/IBundleAggregator.sol";
import {IBundleAggregatorProxy} from "./interfaces/IBundleAggregatorProxy.sol";

/// @title A trusted proxy for updating where current answers are read from
/// @notice This contract provides a consistent address for the
/// CurrentAnswerInterface but delegates where it reads from to the owner, who is
/// trusted to update it.
contract BundleAggregatorProxy is IBundleAggregatorProxy, ITypeAndVersion, ConfirmedOwner {
  string public constant override typeAndVersion = "BundleAggregatorProxy 1.0.0";

  IBundleAggregator private s_currentAggregator;
  IBundleAggregator private s_proposedAggregator;

  event AggregatorProposed(address indexed current, address indexed proposed);
  event AggregatorConfirmed(address indexed previous, address indexed latest);

  error AggregatorNotProposed(address aggregator);

  constructor(address aggregatorAddress, address owner) ConfirmedOwner(owner) {
    s_currentAggregator = IBundleAggregator(aggregatorAddress);
  }

  function latestBundle() external view returns (bytes memory bundle) {
    return s_currentAggregator.latestBundle();
  }

  function latestBundleTimestamp() external view returns (uint256 timestamp) {
    return s_currentAggregator.latestBundleTimestamp();
  }

  /// @notice returns the current aggregator address.
  function aggregator() external view returns (address) {
    return address(s_currentAggregator);
  }

  /// @notice represents the number of decimals the aggregator responses represent.
  function bundleDecimals() external view override returns (uint8[] memory decimals) {
    return s_currentAggregator.bundleDecimals();
  }

  /// @notice the version number representing the type of aggregator the proxy
  /// points to.
  function version() external view override returns (uint256 aggregatorVersion) {
    return s_currentAggregator.version();
  }

  /// @notice returns the description of the aggregator the proxy points to.
  function description() external view returns (string memory aggregatorDescription) {
    return s_currentAggregator.description();
  }

  /// @notice returns the current proposed aggregator
  function proposedAggregator() external view returns (address proposedAggregatorAddress) {
    return address(s_proposedAggregator);
  }

  /// @notice Allows the owner to propose a new address for the aggregator
  /// @param aggregatorAddress The new address for the aggregator contract
  function proposeAggregator(
    address aggregatorAddress
  ) external onlyOwner {
    s_proposedAggregator = IBundleAggregator(aggregatorAddress);
    emit AggregatorProposed(address(s_currentAggregator), aggregatorAddress);
  }

  /// @notice Allows the owner to confirm and change the address
  /// to the proposed aggregator
  /// @dev Reverts if the given address doesn't match what was previously proposed
  /// @param aggregatorAddress The new address for the aggregator contract
  function confirmAggregator(
    address aggregatorAddress
  ) external onlyOwner {
    if (aggregatorAddress != address(s_proposedAggregator)) {
      revert AggregatorNotProposed(aggregatorAddress);
    }
    address previousAggregator = address(s_currentAggregator);
    delete s_proposedAggregator;
    s_currentAggregator = IBundleAggregator(aggregatorAddress);
    emit AggregatorConfirmed(previousAggregator, aggregatorAddress);
  }
}
