// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

/// @title Sorted Set Validation Utility
/// @notice Provides utility functions for validating sorted sets and their subset relationships.
/// @dev This library is used to ensure that arrays of bytes32 are sorted sets and to check subset relations.
library SortedSetValidationUtil {
  /// @dev Error to be thrown when an operation is attempted on an empty set.
  error EmptySet();
  /// @dev Error to be thrown when the set is not in ascending unique order.
  error NotASortedSet(bytes32[] set);
  /// @dev Error to be thrown when the first array is not a subset of the second array.
  error NotASubset(bytes32[] subset, bytes32[] superset);

  /// @notice Checks if `subset` is a valid and unique subset of `superset`.
  /// NOTE: Empty set is not considered a valid subset of superset for our use case.
  /// @dev Both arrays must be valid sets (unique, sorted in ascending order) and `subset` must be entirely contained within `superset`.
  /// @param subset The array of bytes32 to validate as a subset.
  /// @param superset The array of bytes32 in which subset is checked against.
  /// @custom:revert EmptySet If either `subset` or `superset` is empty.
  /// @custom:revert NotASubset If `subset` is not a subset of `superset`.
  function _checkIsValidUniqueSubset(bytes32[] memory subset, bytes32[] memory superset) internal pure {
    if (subset.length == 0 || superset.length == 0) {
      revert EmptySet();
    }

    _checkIsValidSet(subset);
    _checkIsValidSet(superset);

    uint256 i = 0; // Pointer for 'subset'
    uint256 j = 0; // Pointer for 'superset'

    while (i < subset.length && j < superset.length) {
      if (subset[i] > superset[j]) {
        ++j; // Move the pointer in 'superset' to find a match
      } else if (subset[i] == superset[j]) {
        ++i; // Found a match, move the pointer in 'subset'
        ++j; // Also move in 'superset' to continue checking
      } else {
        revert NotASubset(subset, superset);
      }
    }

    if (i < subset.length) {
      revert NotASubset(subset, superset);
    }
  }

  /// @notice Validates that a given set is sorted and has unique elements.
  /// @dev Iterates through the array to check that each element is greater than the previous.
  /// @param set The array of bytes32 to validate.
  /// @custom:revert NotASortedSet If any element is not greater than its predecessor or if any two consecutive elements are equal.
  function _checkIsValidSet(bytes32[] memory set) private pure {
    for (uint256 i = 1; i < set.length; ++i) {
      if (set[i] <= set[i - 1]) {
        revert NotASortedSet(set);
      }
    }
  }
}
